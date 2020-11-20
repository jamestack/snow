package mount_processor

import (
	"errors"
	"sync"
)

// 本地挂载点处理器抽象
type LocalProcessor struct {
	sync.Mutex
	storage map[string]*Service //
}

func (t *LocalProcessor) Init() error {
	t.storage = make(map[string]*Service)
	return nil
}

// 挂载节点
func (t *LocalProcessor) MountNode(serviceName string, nodeName string, address string) error {
	t.Lock()
	defer t.Unlock()

	newNode := &Node{
		NodeName: nodeName,
		Address:  address,
	}
	exService,ok := t.storage[serviceName]
	if !ok {
		t.storage[serviceName] = &Service{
			ServiceName: serviceName,
			Nodes:       []*Node{newNode},
		}
		return nil
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

	for i,node := range exService.Nodes {
		if node.NodeName == nodeName {
			exService.Nodes = append(exService.Nodes[:i], exService.Nodes[i+1:]...)
			return nil
		}
	}

	return errors.New("node not found")
}

// 查询单个节点
func (t *LocalProcessor) Find(serviceName string, nodeName string) (*Node, error) {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return nil, errors.New("service not found")
	}

	for _,node := range exService.Nodes {
		if node.NodeName == nodeName {
			return node, nil
		}
	}

	return nil, errors.New("node not found")
}

// 查询所有节点
func (t *LocalProcessor) FindAll(serviceName string) (*Service, error) {
	t.Lock()
	defer t.Unlock()

	exService := t.storage[serviceName]
	if exService == nil {
		return nil, errors.New("service not found")
	}

	return exService, nil
}

// 查询所有节点
func (t *LocalProcessor) FindAllService() ([]*Service, error) {
	t.Lock()
	defer t.Unlock()

	list := make([]*Service, len(t.storage))
	i := 0
	for _,item := range t.storage {
		list[i] = item
		i++
	}

	return list, nil
}
