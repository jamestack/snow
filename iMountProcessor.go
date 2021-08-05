package snow

type ServiceInfo struct {
	ServiceName string      // 服务名
	Nodes       []*NodeInfo // 节点列表
}

// 节点抽象
type NodeInfo struct {
	NodeName   string // 节点名
	Address    string // 节点地址
	CreateTime int64  // 节点创建时间
}

// 服务挂载点抽象处理器抽象
type IMountProcessor interface {
	// 处理器初始化
	Init(*Cluster) error
	// 挂载某节点
	MountNode(serviceName string, nodeName string, address string, createTime int64) error
	// 删除某节点
	UnMountNode(serviceName string, nodeName string) error
	// 查找某个节点
	Find(serviceName string, nodeName string) (*NodeInfo, error)
	// 查询某服务的所有子节点
	FindAll(serviceName string) (*ServiceInfo, error)
}
