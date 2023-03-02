package snow

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc"
)

// 集群抽象
type Cluster struct {
	rpcLock     sync.Mutex
	findLock    sync.Mutex
	findAllLock sync.Mutex
	streamLock  sync.Mutex
	eventPool   *GoPool
	// 挂载点处理器
	mountProcessor IMountProcessor
	// 本地挂载点
	localNodes sync.Map // map[string]interface{}
	// 监听地址
	listenAddr string
	// 对外暴露地址
	peerAddr string
	// gRpc连接池
	gRpcClients sync.Map // map[string]*ring.Ring
	server      *grpc.Server
	// 唯一验证字符串
	key     string
	closeCh chan os.Signal

	rpcStream map[string][]pb.PeerRpc_StreamServer // map[string][]*rpcReq
}

func (c *Cluster) pushStream(key string, req pb.PeerRpc_StreamServer) {
	c.streamLock.Lock()
	defer c.streamLock.Unlock()

	exList, ok := c.rpcStream[key]
	if ok {
		c.rpcStream[key] = append(exList, req)
	} else {
		c.rpcStream[key] = []pb.PeerRpc_StreamServer{req}
	}
}

func (c *Cluster) popStream(key string) (res pb.PeerRpc_StreamServer) {
	c.streamLock.Lock()
	defer c.streamLock.Unlock()

	exList, ok := c.rpcStream[key]
	if ok {
		if len(exList) > 0 {
			res = exList[0]
			c.rpcStream[key] = exList[1:]
		}
		return res
	} else {
		return nil
	}
}

// 初始化连接
func (c *Cluster) initGRpc(addr string) (conn *ring.Ring, err error) {
	c.rpcLock.Lock()
	defer c.rpcLock.Unlock()

	c.rpcStream = make(map[string][]pb.PeerRpc_StreamServer)

	if exRing, ok := c.gRpcClients.Load(addr); ok {
		return exRing.(*ring.Ring), nil
	}

	var MAX_CONN int = 3
	var r = ring.New(MAX_CONN)
	for i := 0; i < MAX_CONN; i++ {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		_ = cancel // 避免静态检测警告
		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("initGRpc err: %v", err)
		}
		// r.Value = pb.NewPeerRpcClient(conn)
		r.Value = conn
		r = r.Next()
	}
	return r, nil
}

func (c *Cluster) getGrpcConn(addr string) (conn *grpc.ClientConn, err error) {
	var r *ring.Ring
	exRing, ok := c.gRpcClients.Load(addr)
	if !ok {
		r, err = c.initGRpc(addr)
		if err != nil {
			return nil, err
		}
	} else {
		r = exRing.(*ring.Ring)
	}

	conn = r.Value.(*grpc.ClientConn)
	c.gRpcClients.Store(addr, r.Next())
	return conn, nil
}

// 获取连接
func (c *Cluster) getRpcClient(addr string) (client pb.PeerRpcClient, err error) {
	conn, err := c.getGrpcConn(addr)
	if err != nil {
		return nil, err
	}
	client = pb.NewPeerRpcClient(conn)
	return client, nil
}

func (c *Cluster) ClusterKey() string {
	return c.key
}

// 初始化集群
// @listenAddr  rpc监听地址
// &peerAddr    rpc对外访问的地址
func NewClusterWithMountProcessor(listenAddr string, peerAddr string, mountProcessor IMountProcessor) *Cluster {
	c := &Cluster{
		mountProcessor: mountProcessor,
		listenAddr:     listenAddr,
		peerAddr:       peerAddr,
		eventPool:      NewGoPool(uint32(runtime.NumCPU()) * 4),
		server:         grpc.NewServer(),
		key:            RandStr(24),
	}

	if mountProcessor == nil {
		panic("mountProcessor not found")
	}

	return c
}

func (c *Cluster) GetPeerAddr() string {
	return c.peerAddr
}

// 本地集群
func NewClusterWithLocal() *Cluster {
	mountProcessor := &LocalProcessor{}
	return NewClusterWithMountProcessor("", "", mountProcessor)
}

// 支持Consul挂载点的集群
// @listenAddr  rpc监听地址
// &peerAddr    rpc对外访问的地址
func NewClusterWithConsul(listenAddr string, peerAddr string, client ...*api.Client) *Cluster {
	mountProcessor := &ConsulProcessor{}
	if len(client) != 0 {
		mountProcessor.Client = client[0]
	}
	return NewClusterWithMountProcessor(listenAddr, peerAddr, mountProcessor)
}

func NewClusterMaster(listenAddr string, peerAddr string) *Cluster {
	return NewClusterWithMountProcessor(listenAddr, peerAddr, &ClusterMountProcessorMaster{})
}

func NewClusterSlave(listenAddr string, peerAddr string, masterAddr string) *Cluster {
	return NewClusterWithMountProcessor(listenAddr, peerAddr, &ClusterMountProcessorSlave{
		MasterAddr: masterAddr,
	})
}

func (c *Cluster) GrpcServer() *grpc.Server {
	return c.server
}

// 开始监听本
func (c *Cluster) Serve() (done chan os.Signal, err error) {
	done = make(chan os.Signal, 4)

	var listener net.Listener
	// 监听grpc
	if c.listenAddr != "" {
		// 开始监听grpc
		listener, err = net.Listen("tcp", c.listenAddr)
		if err != nil {
			return nil, err
		}

		rpc := &PeerRpc{cluster: c}
		// 用于支持集群间通信
		pb.RegisterPeerRpcServer(c.server, rpc)
		// 用于支持Consul的gRpc健康检查
		pb.RegisterHealthServer(c.server, rpc)

		// 监听rpc端口
		c.eventPool.Go(func() {
			defer checkPanic()
			err = c.server.Serve(listener)
			fmt.Println("[Snow] GRpc Serve() err:", err)
			done <- os.Interrupt
		})
	}

	// 初始化挂载点处理器
	err = c.mountProcessor.Init(c)
	if err != nil {
		if listener != nil {
			_ = listener.Close()
		}
		return nil, err
	}

	fmt.Println("[Snow] Cluster Serve Success")

	// 监听进程退出信号
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	c.eventPool.Go(func() {
		defer checkPanic()
		sig := <-done

		fmt.Println("[Snow] Start UnMount All Local Nodes.")
		// 关闭取消挂载所有节点
		c.localNodes.Range(func(key, value interface{}) bool {
			// node := value.(*Node)
			fmt.Println("[Snow] Start UnMount " + key.(string))
			// err := node.UnMount()
			fmt.Println("[Snow] End UnMount", key.(string), " err:", err)
			return true
		})

		done <- sig
	})

	return done, nil
}

// 查找某节点
func Find[T any](c *Cluster, name string) (*Node[T], error) {
	var serviceName string
	var nodeName string
	name = strings.Trim(name, "/")
	list := strings.Split(name, "/")
	switch len(list) {
	case 1:
		serviceName = list[0]
		nodeName = "Master"
	case 2:
		serviceName = list[0]
		nodeName = list[1]
	default:
		return nil, errors.New("mount name not validate: [" + name + "]")
	}

	key := serviceName + "/" + nodeName
	// 本地节点
	if node, ok := c.localNodes.Load(key); ok {
		return node.(*Node[T]), nil
	}

	// 远程节点
	node, err := c.mountProcessor.Find(serviceName, nodeName)
	if err != nil {
		return nil, err
	}

	var t T
	n := &Node[T]{
		isLocal:     false,
		cluster:     c,
		serviceName: serviceName,
		nodeName:    nodeName,
		RPC:         &t,
		mTime:       node.CreateTime,
	}

	_, ok := reflect.TypeOf(&t).Elem().FieldByName("Node")
	if ok {
		fmt.Println("no zero")
		vt := reflect.ValueOf(&t).Elem().FieldByName("Node")
		vt.Set(reflect.ValueOf(n))
	}
	return n, nil
}

func (c *Cluster) Find(name string) (INode, error) {
	var serviceName string
	var nodeName string
	name = strings.Trim(name, "/")
	list := strings.Split(name, "/")
	switch len(list) {
	case 1:
		serviceName = list[0]
		nodeName = "Master"
	case 2:
		serviceName = list[0]
		nodeName = list[1]
	default:
		return nil, errors.New("mount name not validate: [" + name + "]")
	}

	key := serviceName + "/" + nodeName
	// 本地节点
	if node, ok := c.localNodes.Load(key); ok {
		return node.(INode), nil
	}

	// 远程节点
	node, err := c.mountProcessor.Find(serviceName, nodeName)
	if err != nil {
		return nil, err
	}

	return &Node[any]{
		isLocal:     false,
		cluster:     c,
		serviceName: serviceName,
		nodeName:    nodeName,
		RPC:         nil,
		mTime:       node.CreateTime,
	}, nil
}

// 查询节点列表
func (c *Cluster) FindAll(serviceName string) ([]INode, error) {
	// 远程节点
	service, err := c.mountProcessor.FindAll(serviceName)
	if err != nil {
		return nil, err
	}

	list := []INode{}
	for _, node := range service.Nodes {
		var item INode

		iNode, ok := c.localNodes.Load(serviceName + "/" + node.NodeName)
		if ok {
			item = iNode.(INode)
		} else {
			item = &Node[any]{
				isLocal:     false,
				cluster:     c,
				serviceName: serviceName,
				nodeName:    node.NodeName,
				// iNode:       nil,
				mTime: node.CreateTime,
			}
		}

		list = append(list, item)
	}

	return list, nil
}

// 查询所有本机节点
func (c *Cluster) FindLocalAll() []INode {
	list := []INode{}
	c.localNodes.Range(func(_, value interface{}) bool {
		list = append(list, value.(INode))
		return true
	})
	return list
}

// 挂载某节点
func Mount[T any](c *Cluster, name string, iNode *T) (*Node[T], error) {
	var serviceName string
	var nodeName string
	name = strings.Trim(name, "/")

	list := strings.Split(name, "/")
	switch len(list) {
	case 1:
		serviceName = list[0]
		nodeName = "Master"
	case 2:
		serviceName = list[0]
		nodeName = list[1]
	default:
		return nil, errors.New("mount name not validate: [" + name + "]")
	}

	key := serviceName + "/" + nodeName
	// 检测重复挂载
	if _, ok := c.localNodes.Load(key); ok {
		return nil, errors.New("repeated mount")
	}

	// 挂载本地节点
	vf := reflect.ValueOf(iNode)
	if vf.Kind() != reflect.Ptr {
		return nil, errors.New("must pointer")
	}

	newNode := &Node[T]{
		isLocal:     true,
		cluster:     c,
		serviceName: serviceName,
		nodeName:    nodeName,
		RPC:         iNode,
		mTime:       time.Now().Unix(),
	}

	if _, ok := reflect.TypeOf(iNode).Elem().FieldByName("Node"); ok {
		iNodefield := vf.Elem().FieldByName("Node")
		if iNodefield.IsNil() {
			iNodefield.Set(reflect.ValueOf(newNode))
		}
	}

	// 同步挂载到远程挂载点
	err := c.mountProcessor.MountNode(serviceName, nodeName, c.peerAddr, newNode.mTime)
	if err != nil {
		return nil, err
	}

	// 写入本地挂载点
	c.localNodes.Store(key, newNode)

	fmt.Println("[Snow] Mount Node", key, "Success")

	if i, ok := any(iNode).(HookMount); ok {
		c.eventPool.Go(func() {
			defer checkPanic()
			i.OnMount()
		})
	}

	return newNode, nil
}

// 挂载一个随机子节点
func MountRandNode[T any](c *Cluster, serviceName string, iNode *T) (*Node[T], error) {
	return Mount(c, serviceName+"/"+RandStr(6), iNode)
}

// 挂载一个随机子节点
// func (c *Cluster) MountRandNode(serviceName string, iNode INode) (*Node[T], error) {
// 	return c.Mount(serviceName+"/"+RandStr(6), iNode)
// }

// 取消挂载某节点
func UnMount[T any](c *Cluster, name string) error {
	// 检测是否为本地挂载点，拒绝取消挂载远程节点
	node, err := Find[T](c, name)
	if err != nil {
		return err
	}

	if node.IsRemote() {
		nodeInfo, err := c.mountProcessor.Find(node.serviceName, node.nodeName)
		if err != nil {
			return err
		}
		rpc, err := c.getRpcClient(nodeInfo.Address)
		if err != nil {
			return err
		}
		_, err = rpc.UnMount(context.TODO(), &pb.NodeName{Str: name})

		return err
	}

	// 同步挂载到远程挂载点
	err = c.mountProcessor.UnMountNode(node.serviceName, node.nodeName)
	if err != nil {
		return err
	}

	// 写入本地挂载点
	c.localNodes.Delete(node.serviceName + "/" + node.nodeName)

	if i, ok := any(node.RPC).(HookUnMount); ok {
		func() {
			defer checkPanic()
			i.OnUnMount()
		}()
	}

	fmt.Println("[Snow] UnMount Node " + name + " Success")

	return nil
}
