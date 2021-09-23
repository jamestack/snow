package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc"
)

// 任何一个包含*Node的struct都可以视为合法的节点对象
type INode interface{}

type myErr struct {
	S     string
	IsNil bool
}

func (e *myErr) Error() string {
	return e.S
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func isErr(obj reflect.Type) bool {
	return obj.Implements(errorInterface)
}

// -------------- iNode内部实现 ------------
type Node struct {
	*Cluster
	serviceName string
	nodeName    string
	iNode       interface{}
	mTime       int64 // 挂载时间
}

// 是否为本地节点
func (i *Node) IsLocal() bool {
	_, ok := i.localNodes.Load(i.serviceName + "/" + i.nodeName)
	return ok
}

// 是否为远程节点
func (i *Node) IsRemote() bool {
	_, ok := i.localNodes.Load(i.serviceName + "/" + i.nodeName)
	return !ok
}

// 节点名
func (i *Node) Name() string {
	return i.serviceName + "/" + i.nodeName
}

// 挂载时间
func (i *Node) MountTime() int64 {
	return i.mTime
}

// 取消挂载
func (i *Node) UnMount() error {
	return i.Cluster.UnMount(i.serviceName + "/" + i.nodeName)
}

type localCallback struct {
	lock     *sync.Mutex
	res      []reflect.Value
	thenErr  error
	isDone   bool
	callback CallbackFn
}

type rpcCallback struct {
	cs  grpc.ClientStream
	ack *pb.CallAck
}

type CallBack struct {
	local *localCallback
	rpc   *rpcCallback
	err   error
}

func (c *CallBack) Error() error {
	return c.err
}

type CallbackFn interface{}
type Any interface{}

func (c *CallBack) Then(callback CallbackFn) error {
	if callback == nil {
		return errors.New("callback is nil")
	}

	if c.rpc != nil {
		return c.thenRpc(callback)
	}

	c.local.lock.Lock()
	if c.err != nil {
		c.local.lock.Unlock()
		return c.err
	}

	isRun := false
	if c.local.callback == nil {
		c.local.callback = callback
		if c.local.isDone {
			isRun = true
		}
	} else {
		c.local.lock.Unlock()
		return errors.New("repeat Then()")
	}
	c.local.lock.Unlock()

	if isRun {
		return c.then()
	}

	return c.local.thenErr
}

func (c *CallBack) then() (err error) {
	cb := reflect.ValueOf(c.local.callback)
	if len(c.local.res) != cb.Type().NumIn() {
		return errors.New("cb args length not match")
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			printStack(pErr)
			err = fmt.Errorf("%v", pErr)
			return
		}
	}()

	cb.Call(c.local.res)

	return nil
}

func (c *CallBack) thenRpc(callback CallbackFn) (err error) {
	cb := reflect.ValueOf(callback)
	cbt := cb.Type()
	cbn := cbt.NumIn()

	var rpcRes *pb.CallAck

	if c.rpc.ack == nil {
		rpcRes = &pb.CallAck{}
		err = c.rpc.cs.RecvMsg(rpcRes)
		if err != nil {
			return err
		}
	} else {
		rpcRes = c.rpc.ack
	}

	if len(rpcRes.Args) != cbn {
		return errors.New("remote return length not match callback func args")
	}
	cin := make([]reflect.Value, cbn)
	reader := bytes.NewReader(nil)
	decoder := gob.NewDecoder(reader)

	for i := 0; i < cbn; i++ {
		it := cbt.In(i)
		if isErr(it) {
			it = reflect.TypeOf(&myErr{})
		}
		isPtr := false
		if it.Kind() == reflect.Ptr {
			it = it.Elem()
			isPtr = true
		}
		nw := reflect.New(it)
		reader.Reset(rpcRes.Args[i])
		err := decoder.DecodeValue(nw)
		if err != nil {
			return fmt.Errorf("remote args[%d] decode err: %v", i, err)
		}

		if e, ok := nw.Interface().(*myErr); ok && e.IsNil {
			nw = nilValue
		}

		if isPtr {
			cin[i] = nw
		} else {
			cin[i] = nw.Elem()
		}
	}

	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("node.thenRpc() panic: %v", errPanic)
		}
	}()
	cb.Call(cin)
	return
}

// 执行方法调用
func (i *Node) CallAsync(method string, args ...Any) (callback *CallBack) {
	// 本地调用
	if i.IsLocal() {
		return i.call(method, true, args...)
	} else {
		return i.callRpc(method, true, args...)
	}
}

// 执行方法调用
func (i *Node) Call(method string, args ...Any) (callback *CallBack) {
	// 本地调用
	if i.IsLocal() {
		return i.call(method, false, args...)
	} else {
		return i.callRpc(method, false, args...)
	}
}

// 本地执行
func (i *Node) call(method string, async bool, args ...Any) (callback *CallBack) {
	callback = &CallBack{
		local: &localCallback{
			lock: &sync.Mutex{},
		},
	}

	if i == nil {
		callback.err = errors.New("iNode is nil")
		return
	}

	methodValue := reflect.ValueOf(i.iNode).MethodByName(method)

	if methodValue.Kind() != reflect.Func {
		callback.err = errors.New("not found method " + method)
		return
	}

	if methodValue.Type().NumIn() != len(args) {
		callback.err = errors.New("method args length not match")
		return
	}

	// method call
	params := make([]reflect.Value, len(args))
	for i, v := range args {
		params[i] = reflect.ValueOf(v)
	}

	fn := func() []reflect.Value {
		if i, ok := i.iNode.(HookCall); ok {
			return i.OnCall(method, methodValue.Call, params)
		}
		return methodValue.Call(params)
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			printStack(pErr)
			callback.err = fmt.Errorf("%v", pErr)
			return
		}
	}()
	if async {
		_, callback.err = i.eventPool.Go(func() {
			if pErr := recover(); pErr != nil {
				printStack(pErr)
				callback.local.thenErr = fmt.Errorf("%v", pErr)
				callback.local.isDone = true
				return
			}

			callback.local.res = fn()
			callback.local.isDone = true
			isRun := false
			callback.local.lock.Lock()
			if callback.local.callback != nil {
				isRun = true
			}
			callback.local.lock.Unlock()
			if isRun {
				callback.local.thenErr = callback.then()
			}
		})
	} else {
		callback.local.res = fn()
		callback.local.isDone = true
	}
	return
}

// 由于reflect.Call()无法传入nil值(会报Zero Value错误)，并且就算成功Call((*myErr)(nil))成功，也会出现err == nil会为true的情况（而该又确实是nil,调用时也会触发nil point panic)
// 所以这儿使用了发现的@hack写法，希望寻求更好的reflect.Call(nil)写法。
var nilValue = reflect.Zero(reflect.ValueOf(struct {
	Err error
}{}).Field(0).Type())

var unaryStreamDesc = &grpc.StreamDesc{ServerStreams: false, ClientStreams: false}

func (i *Node) callRpc(method string, async bool, args ...Any) (callback *CallBack) {
	callback = &CallBack{
		rpc: &rpcCallback{},
	}
	if i == nil {
		callback.err = errors.New("iNode is nil")
		return
	}

	var vl []reflect.Value
	for _, ai := range args {
		at := reflect.ValueOf(ai)
		if ai == nil || isErr(at.Type()) {
			if ai == nil {
				at = reflect.ValueOf(&myErr{S: "", IsNil: true})
			} else {
				at = reflect.ValueOf(&myErr{S: ai.(error).Error(), IsNil: false})
			}
		}

		if at.Kind() == reflect.Ptr {
			at = at.Elem()
		}
		vl = append(vl, at)
	}

	req := &pb.CallReq{
		ServiceName: i.serviceName,
		NodeName:    i.nodeName,
		Method:      method,
		Args:        nil,
	}

	req.Args = make([][]byte, len(vl))
	buff := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buff)
	for i, v := range vl {
		err := encoder.EncodeValue(v)
		if err != nil {
			callback.err = fmt.Errorf("call args[%d] encode err:%v", i, err)
			return
		}
		req.Args[i] = buff.Bytes()
		buff.Reset()
	}

	nodeInfo, err := i.mountProcessor.Find(i.serviceName, i.nodeName)
	if err != nil {
		callback.err = err
		return
	}

	conn, err := i.Cluster.getGrpcConn(nodeInfo.Address)
	if err != nil {
		callback.err = err
		return
	}

	cs, err := conn.NewStream(context.Background(), unaryStreamDesc, "/pb.PeerRpc/Call")
	if err != nil {
		callback.err = err
		return
	}

	callback.rpc.cs = cs
	callback.err = cs.SendMsg(req)

	if callback.err == nil && async == false {
		callback.rpc.ack = &pb.CallAck{}
		callback.err = cs.RecvMsg(callback.rpc.ack)
	}
	return
}

func (i *Node) Stream(method string, args ...Any) (stream *Stream, err error) {
	x := make(chan []byte)
	y := make(chan []byte)
	stream = &Stream{
		read:  &x,
		write: &y,
	}
	if i.IsLocal() {
		err = i.call(method, false, append(args, &Stream{
			read:  &y,
			write: &x,
		})...).Error()
		return stream, err
	} else {
		nodeInfo, err := i.mountProcessor.Find(i.serviceName, i.nodeName)
		if err != nil {
			return nil, err
		}
		rpc, err := i.Cluster.getRpcClient(nodeInfo.Address)
		if err != nil {
			return nil, err
		}
		rpcStream, err := rpc.Stream(context.TODO())
		if err != nil {
			return nil, err
		}
		stream.rpcClient = rpcStream
		req := &pb.CallReq{
			ServiceName: i.serviceName,
			NodeName:    i.nodeName,
			Method:      method,
			Args:        [][]byte{},
		}

		var vl []reflect.Value
		for _, ai := range args {
			at := reflect.ValueOf(ai)

			if ai == nil || isErr(at.Type()) {
				if ai == nil {
					at = reflect.ValueOf(&myErr{S: "", IsNil: true})
				} else {
					at = reflect.ValueOf(&myErr{S: ai.(error).Error(), IsNil: false})
				}
			}

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
		for i, v := range vl {
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
