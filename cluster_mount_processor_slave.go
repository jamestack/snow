package snow

import (
	"context"
	"errors"
	"fmt"
	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

// 本地挂载点处理器抽象
type ClusterMountProcessorSlave struct {
	localProcessor *LocalProcessor
	MasterAddr string
	rpcClient pb.MasterRpcClient
	masterKey string
	ctx context.Context
}

var ErrGRpcDialFail = errors.New("GRpc Connect Fail")
var ErrGRpcRegisterFail = errors.New("GRpc Register Fail")

func (s *ClusterMountProcessorSlave) Init(cluster *Cluster) error {
	if s.localProcessor == nil {
		s.localProcessor = &LocalProcessor{}
		err := s.localProcessor.Init(cluster)
		if err != nil {
			return err
		}
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "session_id", RandStr(24))
	conn, err := grpc.DialContext(ctx, s.MasterAddr, grpc.WithInsecure())
	if err != nil {
		return ErrGRpcDialFail
	}

	masterRpc := pb.NewMasterRpcClient(conn)

	ack,err := masterRpc.Register(ctx, &pb.RegisterReq{
		PeerNode:  cluster.GetPeerAddr(),
		SlaveKey: cluster.ClusterKey(),
		MasterKey: "",  // 暂时不验证MasterKey，每次都全量拉取
	})

	if err != nil {
		return ErrGRpcRegisterFail
	}

	s.ctx = ctx
	s.rpcClient = masterRpc

	syncStream,err := s.rpcClient.Sync(ctx)
	if err != nil {
		return err
	}

	err = syncStream.Send(&pb.SyncReq{
		Id: 0,
		Addr: cluster.GetPeerAddr(),
	})
	if err != nil {
		return err
	}

	done := make(chan bool)
	go func() {
		for {
			item,err := syncStream.Recv()
			if err != nil {
				fmt.Println("syncStream.Recv()", err)
				_ = conn.Close()
				return
			}

			if item.Id == -100 {
				done <- true
				continue
			}

			ns := strings.Split(item.Name, "/")
			if item.IsAdd {
				fmt.Println("[Snow] Sync Mount Node", item.Name , "Success")
				_ = s.localProcessor.MountNode(ns[0], ns[1], item.PeerAddr, item.Time)
			}else {
				fmt.Println("[Snow] Sync UnMount Node", item.Name , "Success")
				_ = s.localProcessor.UnMountNode(ns[0], ns[1])
			}

		}
	}()

	if s.masterKey != ack.MasterKey {
		//fmt.Println("MasterKey not match, close this cluster now..")
		// 自动关闭集群
	}

	s.masterKey = ack.MasterKey

	<-done
	fmt.Println("[Snow] First sync done")

	go func() {
		// 检测链接状态
		conn.WaitForStateChange(context.Background(), connectivity.Ready)
		_ = conn.Close()
		fmt.Println("Master GRpc Conn Disconnect", conn.GetState())
		// 删除本地节点
		s.localProcessor.Init(cluster)
		for i := 1;;i++ {
			err := s.Init(cluster)

			if err != nil {
				if err == ErrGRpcDialFail || err == ErrGRpcRegisterFail {
					fmt.Println("Reconnect", s.MasterAddr, i, "s Fail, Continue, err:", err)
					<-time.After(1*time.Second)
					continue
				}
				fmt.Println("Reconnect", s.MasterAddr, i, "s", "Fail err:", err)
			}else {
				fmt.Println("Reconnect", s.MasterAddr, i, "s", "Success")
				// 同步所有本地节点
				for _,node := range cluster.FindLocalAll() {
					strs := strings.Split(node.Name(), "/")
					_ = s.MountNode(strs[0], strs[1], node.GetPeerAddr(), node.mTime)
				}
			}

			break
		}

	}()
	return nil
}

// 挂载节点
func (s *ClusterMountProcessorSlave) MountNode(serviceName string, nodeName string, address string, createTime int64) (err error) {
	_,err = s.rpcClient.Mount(s.ctx, &pb.MountReq{
		Name: serviceName+"/"+nodeName,
	})
	return
}

// 移除挂载节点
func (s *ClusterMountProcessorSlave) UnMountNode(serviceName string, nodeName string) (err error) {
	_,err = s.rpcClient.UnMount(s.ctx, &pb.MountReq{
		Name: serviceName+"/"+nodeName,
	})
	return
}

// 查询单个节点
func (s *ClusterMountProcessorSlave) Find(serviceName string, nodeName string) (*NodeInfo, error) {
	return s.localProcessor.Find(serviceName, nodeName)
}

// 查询所有节点
func (s *ClusterMountProcessorSlave) FindAll(serviceName string) (*ServiceInfo, error) {
	return s.localProcessor.FindAll(serviceName)
}

// 查询所有节点
func (s *ClusterMountProcessorSlave) FindAllService() ([]*ServiceInfo, error) {
	return s.localProcessor.FindAllService()
}
