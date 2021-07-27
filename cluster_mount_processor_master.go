package snow

import (
	"github.com/jamestack/snow/pb"
)

// 本地挂载点处理器抽象
type ClusterMountProcessorMaster struct {
	localProcessor *LocalProcessor
	syncLog *SyncLog
}

func (m *ClusterMountProcessorMaster) Init(cluster *Cluster) error {
	m.syncLog = NewSyncLog()

	pb.RegisterMasterRpcServer(cluster.GrpcServer(), &MasterRpc{
		cluster:     cluster,
		mountMaster: m,
	})

	m.localProcessor = &LocalProcessor{}
	return m.localProcessor.Init(cluster)
}

// 挂载节点
func (m *ClusterMountProcessorMaster) MountNode(serviceName string, nodeName string, address string) (err error) {
	err = m.localProcessor.MountNode(serviceName, nodeName, address)
	if err != nil {
		return
	}
	m.syncLog.AddLog(&pb.MountLogItem{
		Id:                   0,
		IsAdd:                true,
		Name:                 serviceName + "/" + nodeName,
		PeerAddr:             address,
	})
	return
}

// 移除挂载节点
func (m *ClusterMountProcessorMaster) UnMountNode(serviceName string, nodeName string) (err error) {
	err = m.localProcessor.UnMountNode(serviceName, nodeName)
	if err != nil {
		return
	}
	m.syncLog.AddLog(&pb.MountLogItem{
		Id:                   0,
		IsAdd:                false,
		Name:                 serviceName + "/" + nodeName,
		PeerAddr:             "",
	})
	return
}

// 查询单个节点
func (m *ClusterMountProcessorMaster) Find(serviceName string, nodeName string) (*NodeInfo, error) {
	return m.localProcessor.Find(serviceName, nodeName)
}

// 查询所有节点
func (m *ClusterMountProcessorMaster) FindAll(serviceName string) (*ServiceInfo, error) {
	return m.localProcessor.FindAll(serviceName)
}

// 查询所有节点
func (m *ClusterMountProcessorMaster) FindAllService() ([]*ServiceInfo, error) {
	return m.localProcessor.FindAllService()
}
