package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc"
)

func init() {
	// gob.Register(myErr{})
}

type INode interface {
	rPC() any
	IsLocal() bool
	IsRemote() bool
	ServiceName() string
	NodeName() string
	Name() string
	MountTime() int64
	UnMount() error
	Cluster() *Cluster
}

// type myErr struct {
// 	S     string
// 	IsNil bool
// }

// func (e *myErr) Error() string {
// 	return e.S
// }

// var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// func isErr(obj reflect.Type) bool {
// 	return obj.Implements(errorInterface)
// }

type Node[T any] struct {
	isLocal     bool
	RPC         *T
	cluster     *Cluster // 所属集群
	serviceName string   // 服务名
	nodeName    string   // 节点名
	mTime       int64    // 挂载时间
}

func (n *Node[T]) IsLocal() bool {
	return n.isLocal
}

func (n *Node[T]) IsRemote() bool {
	return !n.isLocal
}

func (n *Node[T]) ServiceName() string {
	return n.serviceName
}

func (n *Node[T]) NodeName() string {
	return n.nodeName
}

func (n *Node[T]) Name() string {
	return n.serviceName + "/" + n.nodeName
}

func (n *Node[T]) rPC() any {
	return n.RPC
}

// 挂载时间
func (i *Node[T]) MountTime() int64 {
	return i.mTime
}

// 取消挂载
func (i *Node[T]) UnMount() error {
	return UnMount[T](i.cluster, i.Name())
}

func (i *Node[T]) Cluster() *Cluster {
	return i.cluster
}

// // 由于reflect.Call()无法传入nil值(会报Zero Value错误)，并且就算成功Call((*myErr)(nil))成功，也会出现err == nil会为true的情况（而该又确实是nil,调用时也会触发nil point panic)
// // 所以这儿使用了发现的@hack写法，希望寻求更好的reflect.Call(nil)写法。
var nilValue = reflect.Zero(reflect.ValueOf(struct {
	Err error
}{}).Field(0).Type())

var unaryStreamDesc = &grpc.StreamDesc{ServerStreams: false, ClientStreams: false}

func Response[Res any, NodeType any](node *Node[NodeType], fn func() Res, args ...any) Res {
	if node.IsLocal() {
		return fn()
	}

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Response runtime.Caller(1) !ok")
	}

	method_path := runtime.FuncForPC(pc).Name()
	method_name := method_path[strings.LastIndex(method_path, ".")+1:]

	mv := reflect.ValueOf(node.RPC).MethodByName(method_name)
	if mv.IsZero() {
		panic("method " + method_name + "not found")
	}

	req := &pb.CallReq{
		ServiceName: node.Name(),
		Method:      method_name,
		Args:        nil,
	}

	if len(args) > 0 {
		req.Args = make([][]byte, len(args))
		buff := bytes.NewBuffer([]byte{})
		encoder := gob.NewEncoder(buff)
		for i, v := range args {
			// var err error
			// if isErr(reflect.TypeOf(v)) {
			// 	fmt.Println("wwdferfe")
			// 	err = encoder.Encode(&myErr{S: v.(error).Error(), IsNil: reflect.ValueOf(v).IsNil() || v.(error) == nil})
			// } else {
			// 	err = encoder.Encode(v)
			// }
			err := encoder.Encode(v)
			if err != nil {
				panic(err)
			}
			req.Args[i] = buff.Bytes()
			buff.Reset()
		}
	}

	nodeInfo, err := node.cluster.mountProcessor.Find(node.ServiceName(), node.NodeName())
	if err != nil {
		panic(err)
	}

	conn, err := node.cluster.getGrpcConn(nodeInfo.Address)
	if err != nil {
		panic(err)
	}

	cs, err := conn.NewStream(context.Background(), unaryStreamDesc, "/pb.PeerRpc/Call")
	if err != nil {
		panic(err)
	}

	err = cs.SendMsg(req)
	if err != nil {
		panic(err)
	}

	rpcRes := &pb.CallAck{}
	err = cs.RecvMsg(rpcRes)
	if err != nil {
		panic(err)
	}

	var res Res

	if len(rpcRes.Args) > 0 {
		reader := bytes.NewReader(nil)
		decoder := gob.NewDecoder(reader)
		reader.Reset(rpcRes.Args[0])
		err = decoder.Decode(&res)
		if err != nil {
			// if strings.Contains(err.Error(), "myErr") {
			// 	tmpErr := myErr{}
			// 	fmt.Println(err)
			// 	reader.Reset(rpcRes.Args[0])
			// 	err = decoder.Decode(&tmpErr)
			// 	fmt.Println("tmpErr:", tmpErr, err)
			// } else {
			// 	panic(err)
			// }
			panic(err)
		}
	}

	return res
}

func StreamResponse[NodeType any](node *Node[NodeType], fn func(stream *Stream), args ...any) (*Stream, error) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("Response runtime.Caller(1) !ok")
	}

	method_path := runtime.FuncForPC(pc).Name()
	method_name := method_path[strings.LastIndex(method_path, ".")+1:]

	mv := reflect.ValueOf(node.RPC).MethodByName(method_name)
	if mv.IsZero() {
		return nil, errors.New("method " + method_name + "not found")
	}

	if node.IsLocal() {
		exStream := node.cluster.popStream(node.Name() + "." + method_name)
		if exStream != nil {
			stream := &Stream{
				rpcServer: exStream,
			}
			fn(stream)
			return stream, nil
		} else {
			x := make(chan []byte)
			y := make(chan []byte)
			stream := &Stream{
				read:  &x,
				write: &y,
			}

			go fn(stream)
			stream = &Stream{
				read:  &y,
				write: &x,
			}
			return stream, nil
		}
	}

	nodeInfo, err := node.cluster.mountProcessor.Find(node.serviceName, node.nodeName)
	if err != nil {
		return nil, err
	}
	rpc, err := node.cluster.getRpcClient(nodeInfo.Address)
	if err != nil {
		return nil, err
	}
	rpcStream, err := rpc.Stream(context.TODO())
	if err != nil {
		return nil, err
	}
	stream := &Stream{
		rpcClient: rpcStream,
	}
	req := &pb.CallReq{
		ServiceName: node.serviceName + "/" + node.nodeName,
		Method:      method_name,
		Args:        [][]byte{},
	}

	var vl []reflect.Value
	for _, ai := range args {
		at := reflect.ValueOf(ai)

		// if ai == nil || isErr(at.Type()) {
		// 	if ai == nil {
		// 		at = reflect.ValueOf(&myErr{S: "", IsNil: true})
		// 	} else {
		// 		at = reflect.ValueOf(&myErr{S: ai.(error).Error(), IsNil: false})
		// 	}
		// }

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
