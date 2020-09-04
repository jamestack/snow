package mount_processor

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"strings"
	"sync"
)

// Consul挂载点处理器抽象
type ConsulProcessor struct {
	client *api.Client
	findCache sync.Map
	findAllCache sync.Map
}

func (t *ConsulProcessor) Init() error {
	t.client,_ = api.NewClient(api.DefaultConfig())

	//t.client.Agent().CheckRegister()

	return nil
}

// 挂载节点
func (t *ConsulProcessor) MountNode(serviceName string, nodeName string, address string) error {
	addr := strings.Split(address, ":")
	port,_ := strconv.Atoi(addr[1])
	if port <= 0 {
		return errors.New("address not validate")
	}

	return t.client.Agent().ServiceRegister( &api.AgentServiceRegistration{
		ID:    serviceName+"/"+nodeName,
		Name:  serviceName,
		Address: addr[0],
		Port:  port,
		Tags: []string{nodeName},
		Check: &api.AgentServiceCheck{
			Interval: "5s",
			Timeout: "1s",
			GRPC: address + "/:" +serviceName+"/"+nodeName,
			TLSSkipVerify: true,
		},
	})
}

// 移除挂载节点
func (t *ConsulProcessor) UnMountNode(serviceName string, nodeName string) error {

	return t.client.Agent().ServiceDeregister(serviceName+"/"+nodeName)
}

// 查询单个节点
func (t *ConsulProcessor) Find(serviceName string, nodeName string) (*Node, error) {
	list,_,err := t.client.Health().Service(serviceName, nodeName, true, nil)
	if err != nil {
		return nil, errors.New("node not found err: "+err.Error())
	}

	if len(list) == 0 {
		return nil, errors.New("node not found err: match=0")
	}

	node := list[0]
	return &Node{
		NodeName: strings.Split(node.Service.ID, "/")[1],
		Address:  fmt.Sprintf("%s:%d", node.Node.Address, node.Service.Port),
	}, nil
}

// 查询所有节点
func (t *ConsulProcessor) FindAll(serviceName string) (*Service, error) {
	list,_,err := t.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errors.New("node not found err: "+err.Error())
	}

	service := &Service{
		ServiceName: serviceName,
		Nodes:       make([]*Node, len(list)),
	}

	for i,node := range list {
		service.Nodes[i] = &Node{
			NodeName: strings.Split(node.Service.ID, "/")[1],
			Address:  fmt.Sprintf("%s:%d", node.Node.Address, node.Service.Port),
		}
	}

	return service, nil
}
