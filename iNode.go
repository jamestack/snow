package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"github.com/jamestack/snow/pb"
)

type myErr struct {
	S string
	IsNil bool
}

func (e *myErr) Error() string {
	return e.S
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
func isErr(obj reflect.Type) bool {
	if obj.Implements(errorInterface) {
		return true
	}
	return false
}

// -------------- iNode内部实现 ------------
type Node struct {
	*Cluster
	addr string
	serviceName string
	nodeName string
	iNode   interface{}
	mTime int64  // 挂载时间
}

// 是否为本地节点
func (i *Node) IsLocal() bool {
	_,ok := i.localNodes.Load(i.serviceName+"/"+i.nodeName)
	return ok
}

// 是否为远程节点
func (i *Node) IsRemote() bool {
	_,ok := i.localNodes.Load(i.serviceName+"/"+i.nodeName)
	return !ok
}

// 节点名
func (i *Node) Name() string {
	return i.serviceName+"/"+i.nodeName
}

// 挂载时间
func (i *Node) MountTime() (int64, error) {
	if i.IsLocal() {
		return i.mTime, nil
	}else {
		var rpc pb.PeerRpcClient
		rpc,err := i.Cluster.getRpcClient(i.addr)
		if err != nil {
			return 0, err
		}
		res,err := rpc.MountTime(context.TODO(), &pb.NodeName{
			Str: i.serviceName + "/" + i.nodeName,
		})
		if err != nil {
			return 0, err
		}
		return res.Unix, nil
	}
}

// 取消挂载
func (i *Node) UnMount() error {
	return i.Cluster.UnMount(i.serviceName+"/"+i.nodeName)
}

// 执行方法调用
func (i *Node) Call(method string, args ...interface{}) (err error) {
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("node.Call() panic: %v", errPanic)
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
			printStack(pErr)
			err = fmt.Errorf("%v", pErr)
			return
		}
	}()

	// method call
	params := make([]reflect.Value, len(args))
	hasStream := false
	for i, v := range args {
		if _,ok := v.(*Stream); ok {
			hasStream = true
		}
		params[i] = reflect.ValueOf(v)
	}

	fn := func() []reflect.Value {
		if i,ok := i.iNode.(HookCall);ok {
			return i.OnCall(method, methodValue.Call, params)
		}
		return methodValue.Call(params)
	}

	if hasStream {
		// stream不执行回调函数
		i.eventPool.Go(func() {
			defer checkPanic()
			_ = fn()
		})
		return
	}else {
		res = fn()
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

// 由于reflect.Call()无法传入nil值(会报Zero Value错误)，并且就算成功Call((*myErr)(nil))成功，也会出现err == nil会为true的情况（而该又确实是nil,调用时也会触发nil point panic)
// 所以这儿使用了发现的@hack写法，希望寻求更好的reflect.Call(nil)写法。
var nilValue = reflect.Zero(reflect.ValueOf(struct {
	Err error
}{}).Field(0).Type())

// rpc远程调用
func (i *Node) rpcCall(method string, args ...interface{}) (err error) {
	// 远程调用
	var rpc pb.PeerRpcClient
	rpc,err = i.Cluster.getRpcClient(i.addr)
	if err != nil {
		return err
	}
	var res *pb.CallAck
	req := &pb.CallReq{
		ServiceName:i.serviceName,
		NodeName: i.nodeName,
		Method: method,
		Args:   nil,
	}

	var vl []reflect.Value
	var cb *reflect.Value
	for _, ai := range args {
		at := reflect.ValueOf(ai)
		if ai == nil || isErr(at.Type()) {
			if ai == nil {
				at = reflect.ValueOf(&myErr{S: "", IsNil: true})
			}else {
				at = reflect.ValueOf(&myErr{S: ai.(error).Error(), IsNil: false})
			}
		}

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
		if isErr(it) {
			it = reflect.TypeOf(&myErr{})
		}
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

		if e,ok := nw.Interface().(*myErr);ok && e.IsNil == true {
			nw = nilValue
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
		rpc,err := i.Cluster.getRpcClient(i.addr)
		if err != nil {
			return nil, err
		}
		rpcStream,err := rpc.Stream(context.TODO())
		if err != nil {
			return nil, err
		}
		stream.rpcClient = rpcStream
		req := &pb.CallReq{
			ServiceName: i.serviceName,
			NodeName: i.nodeName,
			Method: method,
			Args: [][]byte{},
		}


		var vl []reflect.Value
		for _, ai := range args {
			at := reflect.ValueOf(ai)

			if ai == nil || isErr(at.Type()) {
				if ai == nil {
					at = reflect.ValueOf(&myErr{S: "", IsNil: true})
				}else {
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
