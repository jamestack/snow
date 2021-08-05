package snow

import (
	"context"
	"fmt"
	"github.com/jamestack/snow/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

// 本地挂载点处理器抽象
type ClusterMountProcessorSlave struct {
	localProcessor *LocalProcessor
	MasterAddr string
	rpcClient pb.MasterRpcClient
	masterKey string
	ctx context.Context
}

func (s *ClusterMountProcessorSlave) Init(cluster *Cluster) error {
	s.localProcessor = &LocalProcessor{}
	err := s.localProcessor.Init(cluster)
	if err != nil {
		return err
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "session_id", RandStr(24))
	conn, err := grpc.DialContext(ctx, s.MasterAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	masterRpc := pb.NewMasterRpcClient(conn)

	ack,err := masterRpc.Register(ctx, &pb.RegisterReq{
		PeerNode:  cluster.GetPeerAddr(),
		SlaveKey: cluster.ClusterKey(),
		MasterKey: s.masterKey,
	})

	if err != nil {
		return err
	}

	s.ctx = ctx
	s.rpcClient = masterRpc

	syncStream,err := s.rpcClient.Sync(ctx, &pb.SyncReq{
		Id: 0,
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
				return
			}

			if item.Id == -100 {
				done <- true
				continue
			}

			fmt.Println("SyncLog", item)

			ns := strings.Split(item.Name, "/")
			if item.IsAdd {
				_ = s.localProcessor.MountNode(ns[0], ns[1], item.PeerAddr, item.Time)
			}else {
				_ = s.localProcessor.UnMountNode(ns[0], ns[1])
			}

		}
	}()

	s.masterKey = ack.MasterKey

	<-done
	fmt.Println("First sync done")
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
