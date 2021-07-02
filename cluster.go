package snow

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"reflect"
	"github.com/jamestack/snow/mount_processor"
	"github.com/jamestack/snow/pb"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 集群抽象
type Cluster struct {
	rpcLock sync.Mutex
	findLock sync.Mutex
	findAllLock sync.Mutex
	eventPool *GoPool
	// 挂载点处理器
	mountProcessor mount_processor.IMountProcessor
	// 挂载点缓存
	findCache sync.Map // map[string]*mount_processor.Node
	findAllCache sync.Map // map[string]*mount_processor.Service
	// 本地挂载点
	localNodes sync.Map // map[string]interface{}
	// 监听地址
	listenAddr string
	// 对外暴露地址
	peerAddr string
	// gRpc连接池
	gRpcClients sync.Map   // map[string]*ring.Ring
	server *grpc.Server
}

// 初始化连接
func (c *Cluster) initGRpc(addr string) (conn *ring.Ring,err error) {
	c.rpcLock.Lock()
	defer c.rpcLock.Unlock()

	if exRing,ok := c.gRpcClients.Load(addr); ok {
		return exRing.(*ring.Ring), nil
	}

	var MAX_CONN int = 3
	var r = ring.New(MAX_CONN)
	for i := 0; i < MAX_CONN; i++ {
		ctx,_ := context.WithTimeout(context.TODO(), 5*time.Second)
		conn,err := grpc.DialContext(ctx ,addr, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("initGRpc err: %v", err)
		}
		r.Value = pb.NewPeerRpcClient(conn)
		r = r.Next()
	}
	return r, nil
}

// 获取连接
func (c *Cluster)getRpcClient(addr string) (client pb.PeerRpcClient,err error) {
	var r *ring.Ring

	exRing,ok := c.gRpcClients.Load(addr)
	if !ok {
		r, err = c.initGRpc(addr)
		if err != nil {
			return nil, err
		}
	}else {
		r = exRing.(*ring.Ring)
	}

	client = r.Value.(pb.PeerRpcClient)
	c.gRpcClients.Store(addr, r.Next())

	return client, nil
}

// 初始化集群
// @listenAddr  rpc监听地址
// &peerAddr    rpc对外访问的地址
func NewCluster(listenAddr string, peerAddr string, mountProcessor mount_processor.IMountProcessor) *Cluster {
	c := &Cluster{
		mountProcessor: mountProcessor,
		listenAddr:     listenAddr,
		peerAddr:       peerAddr,
		eventPool: NewGoPool(uint32(runtime.NumCPU()) * 4),
	}

	if mountProcessor == nil {
		panic("mountProcessor not found")
		return nil
	}

	return c
}

// 本地集群
func NewClusterWithLocal() *Cluster {
	mountProcessor := &mount_processor.LocalProcessor{}
	return NewCluster("", "", mountProcessor)
}

// 支持Consul挂载点的集群
// @listenAddr  rpc监听地址
// &peerAddr    rpc对外访问的地址
func NewClusterWithConsul(listenAddr string, peerAddr string, client ...*api.Client) *Cluster {
	mountProcessor := &mount_processor.ConsulProcessor{}
	if len(client) != 0 {
		mountProcessor.Client = client[0]
	}
	return NewCluster(listenAddr, peerAddr, mountProcessor)
}

// 开始监听本
func (c *Cluster) Serve() (done chan os.Signal, err error){
	ch := make(chan os.Signal)
	done = make(chan os.Signal)
	// 初始化挂载点处理器
	err = c.mountProcessor.Init()
	if err != nil {
		return nil, err
	}
	// 执行异步任务
	c.eventPool.Go(func() {
		defer checkPanic()

		if c.mountProcessor == nil {
			return
		}

		for {
			// 每隔30秒检测一次本地挂载点
			<-time.After(30*time.Second)
			// 如果挂载点失效则尝试重新挂载
			c.localNodes.Range(func(key, value interface{}) bool {
				node := value.(*Node)
				_,err := c.mountProcessor.Find(node.serviceName, node.nodeName)
				if err != nil {
					err = c.mountProcessor.MountNode(node.serviceName, node.nodeName, c.peerAddr)
					if err != nil {
						fmt.Println("[SNOW] auto mount node "+node.Name()+" fail,err = "+err.Error())
					}
				}
				return true
			})
		}
	})

	// 监听进程退出信号
	signal.Notify(ch, os.Interrupt, os.Kill)
	c.eventPool.Go(func() {
		defer checkPanic()
		sig := <-ch

		fmt.Println("[Snow] Start UnMount All Local Nodes.")
		// 关闭取消挂载所有节点
		c.localNodes.Range(func(key, value interface{}) bool {
			node := value.(*Node)
			fmt.Println("[Snow] Start UnMount " + key.(string))
			_ = node.UnMount()
			fmt.Println("[Snow] End UnMount " + key.(string))
			return true
		})

		done <- sig
	})

	// 不监听
	if c.listenAddr == "" {
		fmt.Println("[Snow] Cluster Serve Success")
		done = ch
		return  done,nil
	}

	// 开始监听grpc
	listener,err := net.Listen("tcp", c.listenAddr)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	rpc := &PeerRpc{cluster:c}
	// 用于支持集群间通信
	pb.RegisterPeerRpcServer(server, rpc)
	// 用于支持Consul的gRpc健康检查
	pb.RegisterHealthServer(server, rpc)
	// 更新状态
	c.server = server
	// 监听rpc端口
	c.eventPool.Go(func() {
		defer checkPanic()
		err = server.Serve(listener)
		fmt.Println("[Snow] GRpc Serve() err:", err)
		ch <- os.Interrupt
	})

	fmt.Println("[Snow] Cluster Serve Success")
	return done, nil
}

// 查找某节点
func (c *Cluster) Find(name string) (*Node, error) {
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
		return nil, errors.New("mount name not validate: ["+name+"]")
	}

	key := serviceName +"/"+ nodeName
	// 本地节点
	if node,ok := c.localNodes.Load(key);ok {
		return node.(*Node), nil
	}

	// 远程节点
	node,err := c.find(name, serviceName, nodeName)
	if err != nil {
		return nil, err
	}

	return &Node{
		Cluster: c,
		addr: node.Address,
		serviceName: serviceName,
		nodeName: nodeName,
		iNode: nil,
	}, nil
}

func (c *Cluster) find(name string, serviceName string, nodeName string) (*mount_processor.Node, error) {
	if node,ok := c.findCache.Load(name);ok {
		node := node.(*mount_processor.Node)
		if node.CreateTime >= time.Now().Unix() - 5 {
			return node, nil
		}
	}

	c.findLock.Lock()
	defer c.findLock.Unlock()

	if node,ok := c.findCache.Load(name);ok {
		node := node.(*mount_processor.Node)
		if node.CreateTime >= time.Now().Unix() - 5 {
			return node, nil
		}
	}

	res,err := c.mountProcessor.Find(serviceName, nodeName)
	if err == nil {
		res.CreateTime = time.Now().Unix()
		c.findCache.Store(name, res)
	}
	return res,err
}

// 查询节点列表
func (c *Cluster) FindAll(serviceName string) ([]*Node, error) {
	// 远程节点
	service,err := c.findAll(serviceName)
	if err != nil {
		return nil, err
	}

	list := []*Node{}
	for _,node := range service.Nodes {
		var item *Node

		iNode,ok := c.localNodes.Load(serviceName + "/" + node.NodeName)
		if ok {
			item = iNode.(*Node)
		}else {
			item = &Node{
				Cluster: c,
				addr: node.Address,
				serviceName: serviceName,
				nodeName: node.NodeName,
				iNode: nil,
			}
		}

		list = append(list, item)
	}

	return list, nil
}

func (c *Cluster) findAll(serviceName string) (*mount_processor.Service, error) {
	if node,ok := c.findAllCache.Load(serviceName);ok {
		node := node.(*mount_processor.Service)
		if node.CreateTime >= time.Now().Unix() - 5 {
			return node, nil
		}
	}

	c.findAllLock.Lock()
	defer c.findAllLock.Unlock()

	if node,ok := c.findAllCache.Load(serviceName);ok {
		node := node.(*mount_processor.Service)
		if node.CreateTime >= time.Now().Unix() - 5 {
			return node, nil
		}
	}

	res,err :=  c.mountProcessor.FindAll(serviceName)
	if err == nil {
		res.CreateTime = time.Now().Unix()
		c.findAllCache.Store(serviceName, res)
	}

	return res, err
}

// 查询所有本机节点
func (c *Cluster) FindLocalAll() []*Node {
	list := []*Node{}
	c.localNodes.Range(func(_, value interface{}) bool {
		list = append(list, value.(*Node))
		return true
	})
	return list
}

// 挂载某节点
func (c *Cluster) Mount(name string, iNode interface{}) (*Node, error) {
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
		return nil, errors.New("mount name not validate: ["+name+"]")
	}

	key := serviceName + "/" + nodeName
	// 检测重复挂载
	if _,ok := c.localNodes.Load(key);ok {
		return nil, errors.New("repeated mount")
	}

	// 挂载本地节点
	vf := reflect.ValueOf(iNode)
	if vf.Kind() != reflect.Ptr {
		return nil, errors.New("must pointer.")
	}

	if field,ok := reflect.TypeOf(iNode).Elem().FieldByName("Node"); !ok || !field.Anonymous {
		return nil, errors.New("not found Anonymous Field *snow.Node")
	}

	iNodefield := vf.Elem().FieldByName("Node")
	if iNodefield.IsNil() == false {
		return nil, errors.New("*snow.Node must nil")
	}

	// 同步挂载到远程挂载点
	err := c.mountProcessor.MountNode(serviceName, nodeName, c.peerAddr)
	if err != nil {
		return nil, err
	}

	newNode := &Node{
		Cluster: c,
		addr: c.peerAddr,
		serviceName: serviceName,
		nodeName: nodeName,
		iNode:iNode,
		mTime: time.Now().Unix(),
	}
	iNodefield.Set(reflect.ValueOf(newNode))

	// 写入本地挂载点
	c.localNodes.Store(key, newNode)

	if i,ok := iNode.(HookMount);ok {
		c.eventPool.Go(func() {
			defer checkPanic()
			i.OnMount()
		})
	}

	fmt.Println("[Snow] Mount Node", key, "Success")
	return newNode, nil
}

// 挂载一个随机子节点
func (c *Cluster) MountRandNode(serviceName string, iNode interface{}) (*Node, error) {
	return c.Mount(serviceName+"/"+RandStr(6), iNode)
}

// 取消挂载某节点
func (c *Cluster) UnMount(name string) error {
	// 检测是否为本地挂载点，拒绝取消挂载远程节点
	node,err := c.Find(name)
	if err != nil {
		return err
	}

	if node.IsRemote() {
		rpc,err := c.getRpcClient(node.addr)
		if err != nil {
			return err
		}
		_,err = rpc.UnMount(context.TODO(), &pb.NodeName{Str: name})

		return err
	}

	// 同步挂载到远程挂载点
	err = c.mountProcessor.UnMountNode(node.serviceName, node.nodeName)
	if err != nil {
		return err
	}

	// 写入本地挂载点
	c.localNodes.Delete(node.serviceName + "/" + node.nodeName)

	if i,ok := node.iNode.(HookUnMount);ok {
		func() {
			defer checkPanic()
			i.OnUnMount()
		}()
	}

	return nil
}

