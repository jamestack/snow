package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"snow/pb"
)

// 本地节点
type localNode struct {
	iNode   interface{}
}

// 远程同伴节点
type peerNode struct {
	peerAddr string
	path string
}

// -------------- iNode内部实现 ------------
type Node struct {
	*Cluster
	addr string
	serviceName string
	nodeName string
	iNode   interface{}
}

// 是否为本地节点
func (i *Node) IsLocal() bool {
	return i.addr == i.Cluster.peerAddr
}

// 是否为远程节点
func (i *Node) IsRemote() bool {
	return i.addr != i.Cluster.peerAddr
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
