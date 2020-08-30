package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"snow/pb"
)

// 本地节点
type localNode struct {
	lock    sync.Mutex
	//parent *Node
	name 	string
	pNode   *Node   // 父节点
	iNode   interface{}
	childes map[string]*Node // 子节点
}

// 远程同伴节点
type peerNode struct {
	peerAddr string
	path string
}

// -------------- iNode内部实现 ------------
type Node struct {
	*localNode
	*peerNode
}

// 是否为本地节点
func (i *Node) IsLocal() bool {
	return i.localNode != nil
}

// 是否为远程节点
func (i *Node) IsRemote() bool {
	return i.peerNode != nil
}

// 初始化一个节点
func newNode(name string, node interface{}) *Node {
	if peerNode,ok := node.(*peerNode); !ok {
		return &Node{
			localNode: &localNode{
				name: name,
				iNode:  node,
			},
		}
	}else {
		return &Node{
			peerNode: peerNode,
		}
	}
}

func (i *Node) getChild(name string) *Node {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.childes == nil {
		return nil
	}

	return i.childes[name]
}

func (i *Node) addChild(name string, node *Node) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.childes == nil {
		i.childes = make(map[string]*Node)
	}

	i.childes[name] = node
}

func (i *Node) removeChild(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.childes == nil {
		return
	}

	delete(i.childes, name)
}

func (i *Node) RangeChild(fn func(name string, node *Node) bool) {
	i.lock.Lock()
	for k,v := range i.childes {
		i.lock.Unlock()
		ctl, err := func() (ctl bool, err error){
			defer func() {
				errPanic := recover()
				if errPanic != nil {
					err = fmt.Errorf("%v", errPanic)
				}
			}()
			ctl = fn(k,v)
			return
		}()
		if err != nil {
			return
		}
		if !ctl {
			break
		}

		i.lock.Lock()
	}
	i.lock.Unlock()
}

// 挂载子节点
func (i *Node) Mount(name string, node interface{}) (*Node, error) {
	if i == nil {
		return nil, errors.New("node is nil")
	}

	if !i.IsLocal() {
		return nil, errors.New("must mount local node")
	}
	if name == "" {
		return nil, errors.New("name is empty")
	}

	if strings.Contains(name, "/") {
		return nil, errors.New("can't contains \"/\"")
	}

	var iNodefield reflect.Value
	isPeerNode := false
	//var peerNodes *peerNode
	if _,isPeerNode = node.(*peerNode); !isPeerNode {
		// 挂载本地节点
		vf := reflect.ValueOf(node)
		if vf.Kind() != reflect.Ptr {
			return nil, errors.New("must pointer.")
		}

		if field,ok := reflect.TypeOf(node).Elem().FieldByName("Node"); !ok || !field.Anonymous {
			return nil, errors.New("not found Anonymous Field *snow.Node")
		}

		iNodefield = vf.Elem().FieldByName("Node")
		if iNodefield.IsNil() == false {
			return nil, errors.New("*snow.Node must nil")
		}
	}

	if exNode := i.getChild(name); exNode != nil {
		return nil, errors.New("mount point is already exist")
	}

	// 通知Master
	if i.IsRoot() && !isPeerNode {
		root := i.iNode.(*rootNode)
		if !root.isMaster {
			master,_ := root.getMasterRpc()
			_, err := master.Mount(context.TODO(), &pb.MountReq{
				Name: name,
			})
			if err != nil{
				return nil, err
			}
		}else {
			err := root.addSyncLog(false, name, true, root.peerAddr)
			if err != nil{
				return nil, err
			}
		}
	}

	newNode := newNode(name,node)
	if newNode.IsLocal() {
		newNode.pNode = i
	}

	if !isPeerNode {
		//newNode.parent = i
		iNodefield.Set(reflect.ValueOf(newNode))
	}

	i.addChild(name, newNode)

	// 执行回调
	if n,ok := node.(HookMount); ok {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("mount hook panic:", err)
				}
			}()
			n.OnMount()
		}()
	}

	return newNode, nil
}

func (i *Node) UnMountChild(name string) error {
	if i == nil {
		return errors.New("node is nil")
	}

	if i.IsRemote() {
		return errors.New("peerNode UnMount deny")
	}

	exChild := i.getChild(name)
	if exChild == nil {
		return errors.New("child not exist")
	}

	i.removeChild(name)

	if i.IsRoot() && exChild.IsLocal() {
		root := i.iNode.(*rootNode)
		if !i.IsMaster() {
			master,_ := root.getMasterRpc()
			_, err := master.UnMount(context.TODO(), &pb.MountReq{
				Name:                 name,
			})
			return err
		}else {
			return root.addSyncLog(false, name, false, root.peerAddr)
		}
	}

	return nil
}

func (i *Node) UnMount() error {
	if !i.IsLocal() {
		return errors.New("can't unmount remote node")
	}
	if i.IsRoot() {
		return errors.New("can't unmount root node")
	}
	return root.pNode.UnMountChild(i.name)
}

// 是否为根节点
func (r *Node) IsRoot() bool {
	if r.IsLocal() {
		_,ok := r.iNode.(*rootNode)
		return ok
	}else {
		// 不可能拿到远程节点的root
		return false
	}
}

// 是否为Master节点
func (r *Node) IsMaster() bool {
	root := root.iNode.(*rootNode)
	if r.IsLocal() {
		return root.isMaster
	}else {
		return root.peerNode.peerAddr == root.masterAddr
	}
}

func (node *Node) Find(name string) *Node {
	if node == nil || name == "" {
		return nil
	}

	path := strings.Split(name, "/")

	if strings.HasPrefix(name, "/") {
		// 转化为相对路径
		return root.Find(name[1:])
	}else {
		// 一级一级查
		if node.IsRoot() {
			child := node.getChild(path[0])
			if child == nil {
				return nil
			}
			if child.IsLocal() {
				if len(path) == 1 {
					return child
				}else {
					return child.Find(strings.Join(path[1:], "/"))
				}
			}else {
				return &Node{
					peerNode:  &peerNode{
						peerAddr: child.peerAddr,
						path:     "/"+name,
					},
				}
			}
		}else {
			if node.IsLocal() {
				child := node.getChild(path[0])
				if child == nil {
					return nil
				}
				if child.IsLocal() {
					if len(path) == 1 {
						return child
					}else {
						return child.Find(strings.Join(path[1:], "/"))
					}
				}
			}else {
				return &Node{
					peerNode:  &peerNode{
						peerAddr: node.peerAddr,
						path:     node.path+"/"+name,
					},
				}
			}
		}
	}

	return nil
}

// 执行方法调用
func (i *Node) Call(method string, args ...interface{}) (err error) {
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			err = fmt.Errorf("panic: %v", errPanic)
		}
	}()

	if i == nil {
		return errors.New("nil node")
	}
	// 本地调用
	if i.IsLocal() {
		return i.call(method, args...)
	}else {
		return i.rpcCall(method, args...)
	}
}

// 本地执行
func (i *Node) call(method string, args ...interface{}) (err error) {
	methodValue := reflect.ValueOf(i.iNode).MethodByName(method)
	methodType := methodValue.Type()
	if methodValue.Kind() != reflect.Func {

		return errors.New("not found method " + method)
	}

	var cb reflect.Value
	var cbType reflect.Type
	var hasCb bool

	var res []reflect.Value

	// cb validate
	if len(args) > 0 {
		cb = reflect.ValueOf(args[len(args)-1])
		if cb.Kind() == reflect.Func {
			hasCb = true
			args = args[:len(args)-1]
			cbType = cb.Type()
			if cbArgs := cbType.NumIn(); cbArgs > 0 && cbArgs != methodType.NumOut() {
				return errors.New("cb args length not match")
			}
		}
	}

	if methodType.NumIn() != len(args) {
		return errors.New("method args length not match")
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("%v", pErr)
			return
		}
	}()

	// method call
	if len(args) > 0 {
		params := make([]reflect.Value, len(args))
		hasStream := false
		for i, v := range args {
			if _,ok := v.(*Stream); ok {
				hasStream = true
			}
			params[i] = reflect.ValueOf(v)
		}

		if hasStream {
			// stream不执行回调函数
			go methodValue.Call(params)
			return
		}
		res = methodValue.Call(params)
	}else {
		res = methodValue.Call([]reflect.Value{})
	}

	// cb call
	if hasCb {
		if len(res) != cbType.NumIn() {
			return errors.New("cb args length not match")
		}
		cb.Call(res)
	}

	return nil
}

// rpc远程调用
func (i *Node) rpcCall(method string, args ...interface{}) (err error) {
	// 远程调用
	var rpc pb.PeerRpcClient
	rpc,err = root.iNode.(*rootNode).getPeerRpc(i.peerAddr)
	if err != nil {
		return err
	}
	var res *pb.CallAck
	req := &pb.CallReq{
		Path:   i.path,
		Method: method,
		Args:   nil,
	}

	var vl []reflect.Value
	var cb *reflect.Value
	for _, ai := range args {
		at := reflect.ValueOf(ai)
		if at.Kind() == reflect.Ptr {
			at = at.Elem()
		}
		if at.Kind() == reflect.Func {
			cb = &at
			continue
		}
		vl = append(vl, at)
	}

	req.Args = make([][]byte, len(vl))
	buff := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buff)
	for i,v := range vl {
		err := encoder.EncodeValue(v)
		if err != nil {
			return fmt.Errorf("call args[%d] encode err:%v", i, err)
		}
		req.Args[i] = buff.Bytes()
		buff.Reset()
	}

	res,err = rpc.Call(context.TODO(), req)
	if err != nil {
		return err
	}

	// callback
	if cb == nil {
		return err
	}

	cbt := cb.Type()
	cbn := cbt.NumIn()
	if len(res.Args) != cbn {
		return errors.New("remote return length not match callback func args")
	}
	cin := make([]reflect.Value, cbn)
	reader := bytes.NewReader(nil)
	decoder :=  gob.NewDecoder(reader)
	for i:=0;i<cbn;i++ {
		it := cbt.In(i)
		isPtr := false
		if it.Kind() == reflect.Ptr {
			it = it.Elem()
			isPtr = true
		}
		nw := reflect.New(it)
		reader.Reset(res.Args[i])
		err := decoder.DecodeValue(nw)
		if err != nil {
			return fmt.Errorf("remote args[%d] decode err: %v", i, err)
		}
		if isPtr {
			cin[i] = nw
		}else {
			cin[i] = nw.Elem()
		}
	}

	cb.Call(cin)

	return nil
}

func (i *Node) stream(method string, args ...interface{}) (err error) {
	return i.call(method, args...)
}

func (i *Node) Stream(method string, args ...interface{}) (stream *Stream, err error) {
	x := make(chan []byte, 0)
	y := make(chan []byte, 0)
	stream = &Stream{
		read:  &x,
		write: &y,
	}
	if i.IsLocal() {
		return stream, i.stream(method, append(args, &Stream{
			read:  &y,
			write: &x,
		})...)
	}else {
		root := root.iNode.(*rootNode)
		rpc,err := root.getPeerRpc(i.peerAddr)
		if err != nil {
			return nil, err
		}
		rpcStream,err := rpc.Stream(context.TODO())
		if err != nil {
			return nil, err
		}
		stream.rpcClient = rpcStream
		req := &pb.CallReq{
			Path: i.path,
			Method: method,
			Args: [][]byte{},
		}


		var vl []reflect.Value
		for _, ai := range args {
			at := reflect.ValueOf(ai)
			if at.Kind() == reflect.Ptr {
				at = at.Elem()
			}
			if at.Kind() == reflect.Func {
				continue
			}
			vl = append(vl, at)
		}
		req.Args = make([][]byte, len(vl))
		buff := bytes.NewBuffer([]byte{})
		encoder := gob.NewEncoder(buff)
		for i,v := range vl {
			err := encoder.EncodeValue(v)
			if err != nil {
				return nil, fmt.Errorf("call args[%d] encode err:%v", i, err)
			}
			req.Args[i] = buff.Bytes()
			buff.Reset()
		}

		err = rpcStream.Send(&pb.StreamMsg{
			StreamType: &pb.StreamMsg_Req{
				Req: req,
			},
		})
		if err != nil {
			return nil, err
		}
		return stream, nil
	}
}
