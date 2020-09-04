package snow

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"reflect"
	"snow/mount_processor"
	"snow/pb"
	"strings"
	"sync"
	"time"
)

// 集群抽象
type Cluster struct {
	rpcLock sync.Mutex
	findLock sync.Mutex
	findAllLock sync.Mutex
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
func NewCluster(listenAddr string, peerAddr string, mountProcessor mount_processor.IMountProcessor) *Cluster {
	c := &Cluster{
		mountProcessor: mountProcessor,
		listenAddr:     listenAddr,
		peerAddr:       peerAddr,
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
func NewClusterWithConsul(listenAddr string, peerAddr string) *Cluster {
	mountProcessor := &mount_processor.ConsulProcessor{}
	_ = mountProcessor.Init()
	return NewCluster(listenAddr, peerAddr, mountProcessor)
}

// 开始监听本
func (c *Cluster) Serve() (done chan bool, err error){
	ch := make(chan bool)
	// 不监听
	if c.listenAddr == "" {
		return  ch,nil
	}

	listener,err := net.Listen("tcp", c.listenAddr)
	if err != nil {
		return ch, err
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
	go func() {
		err = server.Serve(listener)
		ch <- true
	}()

	return ch, nil
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

		if node.Address == c.peerAddr {
			iNode,ok := c.localNodes.Load(serviceName + "/" + node.NodeName)
			if ok {
				item = iNode.(*Node)
			}else {
				continue
			}
		} else {
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
	}
	iNodefield.Set(reflect.ValueOf(newNode))

	// 写入本地挂载点
	c.localNodes.Store(key, newNode)

	if i,ok := iNode.(HookMount);ok {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%T OnMount() panic: %v\n", iNode, err)
				}
			}()
			i.OnMount()
		}()
	}
	return newNode, nil
}

// 取消挂载某节点
func (c *Cluster) UnMount(name string) error {
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
		return errors.New("mount name not validate: ["+name+"]")
	}

	key := serviceName +"/"+ nodeName

	// 检测是否为本地挂载点，拒绝取消挂载远程节点
	node,ok := c.localNodes.Load(key)
	if !ok {
		return errors.New("not found local node")
	}

	// 同步挂载到远程挂载点
	err := c.mountProcessor.UnMountNode(serviceName, nodeName)
	if err != nil {
		return err
	}

	// 写入本地挂载点
	c.localNodes.Delete(key)

	iNode := node.(*Node).iNode
	if i,ok := iNode.(HookUnMount);ok {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%T OnUnMount() panic: %v\n", iNode, err)
				}
			}()
			i.OnUnMount()
		}()
	}

	return nil
}

