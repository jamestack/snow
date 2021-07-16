package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jamestack/snow/pb"
)

//// ======================== 主节点 ===========================
//type MasterRpc struct {
//	peerAddr string
//}
//
//// 心跳
//func (m *MasterRpc)Ping(context.Context, *pb.Empty) (ack *pb.Empty,err error) {
//	ack = emptyMsg
//	if m.peerAddr == "" {
//		err = errors.New("not register")
//		return
//	}
//
//	//root.iNode.(*rootNode).lastPing[m.peerAddr] = time.Now()
//
//	return
//}
//
//// 注册
//func (m *MasterRpc)Register(ctx context.Context, req *pb.RegisterReq) (ack *pb.RegisterAck, err error)  {
//	ack = &pb.RegisterAck{
//		MasterKey: 0,
//	}
//
//	if m.peerAddr != "" && m.peerAddr != req.PeerNode {
//		err = errors.New("is already register")
//		return
//	}
//
//	//root := root.iNode.(*rootNode)
//	//
//	//ack.MasterKey = root.masterKey
//	//if req.MasterKey != 0 && req.MasterKey != root.masterKey {
//	//	err = errors.New("register key is already expired")
//	//	return
//	//}
//	//
//	//m.peerAddr = req.PeerNode
//	//root.lastPing[m.peerAddr] = time.Now()
//
//	return
//}
//
//// 挂载节点
//func (m *MasterRpc)Mount(ctx context.Context, req *pb.MountReq) (ack *pb.Empty,err error) {
//	ack = emptyMsg
//	if m.peerAddr == "" {
//		err = errors.New("not register")
//		return
//	}
//
//	//err = root.iNode.(*rootNode).addSyncLog(true, req.Name, true, m.peerAddr)
//	return
//}
//// 移除节点
//func (m *MasterRpc)UnMount(ctx context.Context, req *pb.MountReq) (ack *pb.Empty,err error) {
//	ack = emptyMsg
//	if m.peerAddr == "" {
//		err = errors.New("not register")
//		return
//	}
//
//	//err = root.iNode.(*rootNode).addSyncLog(true, req.Name, false, m.peerAddr)
//	return
//}
//// 移除所有节点
//func (m *MasterRpc)UnMountAll(ctx context.Context,req *pb.Empty) (ack *pb.Empty,err error) {
//	ack = emptyMsg
//	if m.peerAddr == "" {
//		err = errors.New("not register")
//		return
//	}
//	//
//	//root.RangeChild(func(name string, node *Node) bool {
//	//	if !node.IsLocal() && node.peerNode.peerAddr == m.peerAddr {
//	//		_ = root.iNode.(*rootNode).addSyncLog(true, name, false, m.peerAddr)
//	//	}
//	//	return true
//	//})
//
//	return
//}
// 同步挂载点
//func (m *MasterRpc)Sync(req *pb.SyncReq, stream pb.MasterRpc_SyncServer) (err error) {
//	lastIndex := 0
//	for {
//		logs := root.iNode.(*rootNode).syncLogs
//		length := len(logs)
//		var list []*SyncLog
//		if length-lastIndex > 0 {
//			list = make([]*SyncLog, length-lastIndex)
//			copy(list, logs[lastIndex:length])
//			lastIndex = length
//		}
//
//		// 去除重复
//		for i:=0;i<len(list)-1;i++ {
//			li := (list)[i]
//			if li == nil {
//				continue
//			}
//			if li.Id <= req.Id {
//				continue
//			}
//			if li.IsAdd == true {
//				for j:=1;j<len(list);j++ {
//					lj := (list)[j]
//					if lj == nil {
//						continue
//					}
//					if lj.Name == li.Name && lj.IsAdd == false {
//						(list)[i] = nil
//						(list)[j] = nil
//						break
//					}
//				}
//			}
//		}
//
//		// Send
//		for _, item := range list {
//			if item == nil {
//				continue
//			}
//			if item.Id <= req.Id {
//				continue
//			}
//			err = stream.Send(&pb.MountLogItem{
//				Id:                   item.Id,
//				IsAdd:                item.IsAdd,
//				Name:                 item.Name,
//				PeerAddr:             item.PeerAddr,
//			})
//			if err != nil {
//				return err
//			}
//		}
//
//		// 检测更新的频率为100毫秒
//		<-time.After(100 * time.Millisecond)
//	}
//
//}

// ============================ 子节点 ===========================
type PeerRpc struct {
	cluster *Cluster
}

// 远程调用
func (p *PeerRpc) Call(ctx context.Context, req *pb.CallReq) (ack *pb.CallAck, err error) {
	ack = &pb.CallAck{}
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("PeerRpc.Call() panic: %v", errPanic)
		}
	}()

	node, err := p.cluster.Find(req.ServiceName + "/" + req.NodeName)
	if err != nil {
		err = errors.New("node not found: " + err.Error())
		return
	}

	if !node.IsLocal() {
		err = errors.New("not found node")
		return
	}
	vf := reflect.ValueOf(node.iNode)
	vm := vf.MethodByName(req.Method)

	vmt := vm.Type()

	fn := vmt.NumIn()
	if fn != len(req.Args) {
		err = errors.New("call args length not match")
		return
	}
	fin := make([]reflect.Value, fn)
	reader := bytes.NewReader(nil)
	decoder := gob.NewDecoder(reader)
	for i := 0; i < fn; i++ {
		tp := vmt.In(i)
		if isErr(tp) {
			tp = reflect.TypeOf(&myErr{})
		}
		isPtr := false
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
			isPtr = true
		}
		nw := reflect.New(tp)

		reader.Reset(req.Args[i])
		err = decoder.DecodeValue(nw)
		if err != nil {
			err = fmt.Errorf("call args[%d] decode err:%v", i, err)
			return
		}

		if isPtr {
			fin[i] = nw
			if e := nw.Interface().(*myErr); e.IsNil {
				fin[i] = nilValue
			}
		} else {
			fin[i] = nw.Elem()
		}
	}

	var fout []reflect.Value
	if i, ok := node.iNode.(HookCall); ok {
		fout = i.OnCall(req.Method, vm.Call, fin)
	} else {
		fout = vm.Call(fin)
	}

	ack = &pb.CallAck{
		Args: make([][]byte, len(fout)),
	}
	buff := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buff)
	for i, v := range fout {
		if isErr(v.Type()) {
			if v.IsNil() || v.Interface().(*myErr) == nil {
				v = reflect.ValueOf(&myErr{S: "", IsNil: true})
			} else {
				v = reflect.ValueOf(&myErr{S: v.Interface().(error).Error(), IsNil: false})
			}
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		err = encoder.EncodeValue(v)
		if err != nil {
			err = fmt.Errorf("call res[%d] encode err:%v", i, err)
			return
		}
		ack.Args[i] = buff.Bytes()
		buff.Reset()
	}

	return
}

// 流抽象
func (p *PeerRpc) Stream(stream pb.PeerRpc_StreamServer) (err error) {
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("PeerRpc.Stream() panic: %v", errPanic)
		}
	}()

	streamReq, err := stream.Recv()
	if err != nil {
		return err
	}
	req := streamReq.GetReq()
	if req == nil {
		return errors.New("StreamReq not CallReq")
	}

	node, err := p.cluster.Find(req.ServiceName + "/" + req.NodeName)
	if err != nil {
		err = errors.New("node not found: " + err.Error())
		return
	}

	if !node.IsLocal() {
		err = errors.New("not found node")
		return
	}
	vf := reflect.ValueOf(node.iNode)
	vm := vf.MethodByName(req.Method)

	vmt := vm.Type()

	fn := vmt.NumIn()
	if fn != len(req.Args)+1 {
		err = errors.New("call args length not match")
		return
	}
	fin := make([]reflect.Value, fn)
	reader := bytes.NewReader(nil)
	decoder := gob.NewDecoder(reader)

	for i := 0; i < fn-1; i++ {
		tp := vmt.In(i)
		if isErr(tp) {
			tp = reflect.TypeOf(&myErr{})
		}
		isPtr := false
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
			isPtr = true
		}
		nw := reflect.New(tp)

		reader.Reset(req.Args[i])
		err = decoder.DecodeValue(nw)
		if err != nil {
			err = fmt.Errorf("call args[%d] decode err:%v", i, err)
			return
		}

		if isPtr {
			fin[i] = nw
			if e := nw.Interface().(*myErr); e.IsNil {
				fin[i] = nilValue
			}
		} else {
			fin[i] = nw.Elem()
		}
	}

	// 注入rpcStream
	fin[len(fin)-1] = reflect.ValueOf(&Stream{
		rpcServer: stream,
	})

	if i, ok := node.iNode.(HookCall); ok {
		_ = i.OnCall(req.Method, vm.Call, fin)
	} else {
		_ = vm.Call(fin)
	}
	return nil
}

type Stream struct {
	rpcClient pb.PeerRpc_StreamClient
	rpcServer pb.PeerRpc_StreamServer
	write     *chan []byte
	read      *chan []byte
}

func (s *Stream) Read(p []byte) (n int, err error) {
	data, err := s.ReadData()
	if err != nil {
		return 0, err
	}

	copy(p, data)
	n = len(data)
	return
}

func (s *Stream) ReadData() (data []byte, err error) {
	switch {
	case s.rpcClient != nil:
		res, err := s.rpcClient.Recv()
		if err != nil {
			return nil, err
		}
		b := res.GetData()
		if b == nil {
			err = s.rpcClient.CloseSend()
			if err != nil {
				return nil, err
			} else {
				return s.ReadData()
			}
		}
		return b.Bytes, nil
	case s.rpcServer != nil:
		res, err := s.rpcServer.Recv()
		if err != nil {
			return nil, err
		}
		b := res.GetData()
		if b == nil {
			return nil, errors.New("not []byte")
		}
		return b.Bytes, nil
	case s.read != nil:
		var ok bool
		data, ok = <-*s.read
		if !ok {
			err = errors.New("stream is closed")
		}
	default:
		err = errors.New("stream not init")
	}
	return
}

func (s *Stream) Write(data []byte) (n int, err error) {
	switch {
	case s.rpcClient != nil:
		err = s.rpcClient.Send(&pb.StreamMsg{
			StreamType: &pb.StreamMsg_Data{
				Data: &pb.Bytes{
					Bytes: data,
				},
			},
		})
	case s.rpcServer != nil:
		err = s.rpcServer.Send(&pb.StreamMsg{
			StreamType: &pb.StreamMsg_Data{
				Data: &pb.Bytes{
					Bytes: data,
				},
			},
		})
	case s.write != nil:
		defer func() {
			if recover() != nil {
				err = errors.New("stream is closed")
			}
		}()
		*s.write <- data
	default:
		err = errors.New("stream not init")
	}

	if err == nil {
		n = len(data)
	}
	return
}

func (s *Stream) Close() (err error) {
	switch {
	case s.rpcClient != nil:
		err = s.rpcClient.CloseSend()
	case s.rpcServer != nil:
		err = s.rpcServer.Send(&pb.StreamMsg{
			StreamType: &pb.StreamMsg_Req{Req: &pb.CallReq{}},
		})
		s.rpcServer = nil
	case *s.read != nil || *s.write != nil:
		defer func() {
			if recover() != nil {
				err = errors.New("chan is closed")
			}
		}()

		close(*s.read)
		close(*s.write)
	default:
		err = errors.New("stream not init")
	}
	return
}

// ======================================== consul healthy check =======================================
func (p *PeerRpc) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	status := pb.HealthCheckResponse_SERVICE_UNKNOWN

	if req.Service != "" {
		s := strings.Split(req.Service[1:], "/")
		_, ok := p.cluster.localNodes.Load(s[0] + "/" + s[1])
		if ok {
			status = pb.HealthCheckResponse_SERVING
		} else {
			status = pb.HealthCheckResponse_NOT_SERVING
		}
	}

	//fmt.Println("check", req.Service, status)

	return &pb.HealthCheckResponse{
		Status: status,
	}, nil
}

func (p *PeerRpc) Watch(req *pb.HealthCheckRequest, watcher pb.Health_WatchServer) error {
	//fmt.Println("watch", req.Service)
	return nil
}

func (p *PeerRpc) MountTime(ctx context.Context, req *pb.NodeName) (ack *pb.MountTimeAck, err error) {
	ack = &pb.MountTimeAck{}
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("PeerRpc.MountTime() panic: %v", errPanic)
		}
	}()

	node, err := p.cluster.Find(req.Str)
	if err != nil {
		err = errors.New("node not found: " + err.Error())
		return
	}

	if !node.IsLocal() {
		err = errors.New("not found node")
		return
	}

	ack.Unix = node.mTime
	err = nil
	return
}

func (p *PeerRpc) UnMount(ctx context.Context, req *pb.NodeName) (ack *pb.Empty, err error) {
	ack = &pb.Empty{}
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			printStack(errPanic)
			err = fmt.Errorf("PeerRpc.UnMount() panic: %v", errPanic)
		}
	}()

	node, err := p.cluster.Find(req.Str)
	if err != nil {
		err = errors.New("node not found: " + err.Error())
		return
	}

	if !node.IsLocal() {
		err = errors.New("not found node")
		return
	}

	// 取消挂载
	err = node.UnMount()
	return
}
