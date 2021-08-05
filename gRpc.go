package snow

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc/metadata"
	"reflect"
	"strings"
	"sync"
	"time"
)

// ======================== 主节点 ===========================
type MasterRpc struct {
	cluster  *Cluster
	mountMaster *ClusterMountProcessorMaster
	session sync.Map
}

type session struct {
	PeerAddr string
}

func (m *MasterRpc) getSession(ctx context.Context) *session {
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	ids := md.Get("session_id")
	if len(ids) == 0 {
		return nil
	}

	exSession,ok := m.session.Load(ids[0])
	if !ok {
		sess := &session{}
		m.session.Store(ids[0], sess)
		return sess
	}
	return exSession.(*session)
}

var emptyMsg = &pb.Empty{}

// 心跳
func (m *MasterRpc) Ping(ctx context.Context,req *pb.Empty) (ack *pb.Empty, err error) {
	ack = emptyMsg
	sess := m.getSession(ctx)
	if sess == nil {
		return ack, errors.New("session_id not found")
	}

	//root.iNode.(*rootNode).lastPing[m.peerAddr] = time.Now()

	return
}

// 注册
func (m *MasterRpc) Register(ctx context.Context, req *pb.RegisterReq) (ack *pb.RegisterAck, err error) {
	ack = &pb.RegisterAck{
		MasterKey: m.cluster.key,
	}

	sess := m.getSession(ctx)
	if sess == nil {
		err = errors.New("session_id not found")
		return
	}

	if sess.PeerAddr != "" {
		err = errors.New("already register")
		return
	}

	if req.MasterKey != "" && req.MasterKey != m.cluster.ClusterKey() {
		err = errors.New("register key is already expired")
		return
	}

	if req.SlaveKey == "" {
		err = errors.New("SlaveKey is empty")
		return
	}

	client, err := m.cluster.getRpcClient(req.PeerNode)
	if err != nil {
		err = errors.New("this rpc not allow dial")
		return
	}

	_, err = client.CheckKey(ctx, &pb.CheckKeyReq{
		Key: req.SlaveKey,
	})
	if err != nil {
		err = errors.New("CheckNum fail")
		return
	}

	sess.PeerAddr = req.PeerNode

	return
}

// 挂载节点
func (m *MasterRpc) Mount(ctx context.Context, req *pb.MountReq) (ack *pb.Empty, err error) {
	ack = emptyMsg
	sess := m.getSession(ctx)
	if sess == nil {
		err = errors.New("session_id not found")
		return
	}

	if sess.PeerAddr == "" {
		err = errors.New("not register")
		return
	}

	ns := strings.Split(req.Name, "/")
	if len(ns) != 2 {
		err = errors.New("name not valid")
		return
	}

	// 处理器处理
	err = m.mountMaster.MountNode(ns[0], ns[1], sess.PeerAddr, time.Now().Unix())

	return
}

// 移除节点
func (m *MasterRpc) UnMount(ctx context.Context, req *pb.MountReq) (ack *pb.Empty, err error) {
	ack = emptyMsg
	sess := m.getSession(ctx)
	if sess == nil {
		err = errors.New("session_id not found")
		return
	}

	if sess.PeerAddr == "" {
		err = errors.New("not register")
		return
	}

	ns := strings.Split(req.Name, "/")
	if len(ns) != 2 {
		err = errors.New("name not valid")
		return
	}

	// 处理器处理
	err = m.mountMaster.UnMountNode(ns[0], ns[1])

	return
}

// 同步挂载点
func (m *MasterRpc) Sync(req *pb.SyncReq, stream pb.MasterRpc_SyncServer) (err error) {
	done,err := m.mountMaster.syncLog.Sync(req.Id, stream)
	if err != nil {
		return err
	}

	<-done

	return nil
}

// ============================ 子节点 ===========================
type PeerRpc struct {
	cluster  *Cluster
	checkNum int64
}

func (p *PeerRpc) CheckKey(ctx context.Context, req *pb.CheckKeyReq) (ack *pb.Empty, err error) {
	ack = &pb.Empty{}
	if req.Key != p.cluster.key {
		err = errors.New("checkkey not valid")
		return
	}

	return
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
