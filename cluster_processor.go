package snow

import (
	"errors"
	"github.com/jamestack/snow/pb"
	"sync"
)

// 本地挂载点处理器抽象
type ClusterProcessor struct {
	sync.Mutex
	localProcessor *LocalProcessor
	MasterAddr string
}

func (t *ClusterProcessor) Init(cluster *Cluster) error {
	pb.RegisterMasterRpcServer(cluster.GrpcServer(), &MasterRpc{
		cluster:   cluster,
		peerAddr:  "",
		masterKey: 0,
	})

	t.localProcessor = &LocalProcessor{}
	return t.localProcessor.Init(cluster)
}

// 挂载节点
func (t *ClusterProcessor) MountNode(serviceName string, nodeName string, address string) error {
	t.Lock()
	defer t.Unlock()

	newNode := &NodeInfo{
		NodeName: nodeName,
		Address:  address,
	}
	exService, ok := t.storage[serviceName]
	if !ok {
		t.storage[serviceName] = &ServiceInfo{
			ServiceName: serviceName,
			Nodes:       []*NodeInfo{newNode},
		}
		return nil
	}

	exService.Nodes = append(exService.Nodes, newNode)

	return nil
}

// 移除挂载节点
func (t *ClusterProcessor) UnMountNode(serviceName string, nodeName string) error {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return errors.New("service not found")
	}

	for i, node := range exService.Nodes {
		if node.NodeName == nodeName {
			exService.Nodes = append(exService.Nodes[:i], exService.Nodes[i+1:]...)
			return nil
		}
	}

	return errors.New("node not found")
}

// 查询单个节点
func (t *ClusterProcessor) Find(serviceName string, nodeName string) (*NodeInfo, error) {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return nil, errors.New("service not found")
	}

	for _, node := range exService.Nodes {
		if node.NodeName == nodeName {
			return node, nil
		}
	}

	return nil, errors.New("node not found")
}

// 查询所有节点
func (t *ClusterProcessor) FindAll(serviceName string) (*ServiceInfo, error) {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return nil, errors.New("service not found")
	}

	return exService, nil
}

// 查询所有节点
func (t *ClusterProcessor) FindAllService() ([]*ServiceInfo, error) {
	t.Lock()
	defer t.Unlock()

	list := make([]*ServiceInfo, len(t.storage))
	i := 0
	for _, item := range t.storage {
		list[i] = item
		i++
	}

	return list, nil
}
