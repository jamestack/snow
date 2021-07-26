package snow

import (
	"errors"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
)

// Consul挂载点处理器抽象
type ConsulProcessor struct {
	Client       *api.Client
	findCache    sync.Map
	findAllCache sync.Map
}

func (t *ConsulProcessor) Init(cluster *Cluster) (err error) {
	if t.Client == nil {
		t.Client, err = api.NewClient(api.DefaultConfig())
	}

	return
}

// 挂载节点
func (t *ConsulProcessor) MountNode(serviceName string, nodeName string, address string) error {
	return t.Client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      serviceName + "/" + nodeName,
		Name:    serviceName,
		Address: address,
		Tags:    []string{nodeName},
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "1s",
			GRPC:                           address + "/:" + serviceName + "/" + nodeName,
			TLSSkipVerify:                  true,
			DeregisterCriticalServiceAfter: "1m",
		},
	})
}

// 移除挂载节点
func (t *ConsulProcessor) UnMountNode(serviceName string, nodeName string) error {

	return t.Client.Agent().ServiceDeregister(serviceName + "/" + nodeName)
}

// 查询单个节点
func (t *ConsulProcessor) Find(serviceName string, nodeName string) (*NodeInfo, error) {
	list, _, err := t.Client.Health().Service(serviceName, nodeName, true, nil)
	if err != nil {
		return nil, errors.New("node not found err: " + err.Error())
	}

	if len(list) == 0 {
		return nil, errors.New("node not found err: match=0")
	}

	node := list[0]
	return &NodeInfo{
		NodeName: strings.Split(node.Service.ID, "/")[1],
		Address:  node.Service.Address,
	}, nil
}

// 查询所有节点
func (t *ConsulProcessor) FindAll(serviceName string) (*ServiceInfo, error) {
	list, _, err := t.Client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errors.New("node not found err: " + err.Error())
	}

	service := &ServiceInfo{
		ServiceName: serviceName,
		Nodes:       make([]*NodeInfo, len(list)),
	}

	for i, node := range list {
		service.Nodes[i] = &NodeInfo{
			NodeName: strings.Split(node.Service.ID, "/")[1],
			Address:  node.Service.Address,
		}
	}

	return service, nil
}
