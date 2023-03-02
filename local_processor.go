package snow

import (
	"errors"
	"sync"
)

// 本地挂载点处理器抽象
type LocalProcessor struct {
	sync.Mutex
	storage map[string]*ServiceInfo //
}

func (t *LocalProcessor) Init(cluster *Cluster) error {
	t.Lock()
	defer t.Unlock()
	t.storage = make(map[string]*ServiceInfo)
	return nil
}

// 挂载节点
func (t *LocalProcessor) MountNode(serviceName string, nodeName string, address string, createTime int64) error {
	t.Lock()
	defer t.Unlock()

	newNode := &NodeInfo{
		NodeName:   nodeName,
		Address:    address,
		CreateTime: createTime,
	}
	exService, ok := t.storage[serviceName]
	if !ok {
		t.storage[serviceName] = &ServiceInfo{
			ServiceName: serviceName,
			Nodes:       []*NodeInfo{newNode},
		}
		return nil
	}

	for _, item := range exService.Nodes {
		if item.NodeName == nodeName {
			return errors.New("already mounted")
		}
	}

	exService.Nodes = append(exService.Nodes, newNode)

	return nil
}

// 移除挂载节点
func (t *LocalProcessor) UnMountNode(serviceName string, nodeName string) error {
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
func (t *LocalProcessor) Find(serviceName string, nodeName string) (*NodeInfo, error) {
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
func (t *LocalProcessor) FindAll(serviceName string) (*ServiceInfo, error) {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return nil, errors.New("service not found")
	}

	return exService, nil
}

// 查询所有节点
func (t *LocalProcessor) FindAllService() ([]*ServiceInfo, error) {
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
