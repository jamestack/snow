package mount_processor

type Service struct {
	ServiceName string  // 服务名
	Nodes       []*Node // 节点列表
	CreateTime  int64   // 缓存时间
}

// 节点抽象
type Node struct {
	NodeName   string // 节点名
	Address    string // 节点地址
	CreateTime int64  // 缓存时间
}

// 服务挂载点抽象处理器抽象
type IMountProcessor interface {
	// 处理器初始化
	Init() error
	// 挂载某节点
	MountNode(serviceName string, nodeName string, address string) error
	// 删除某节点
	UnMountNode(serviceName string, nodeName string) error
	// 查找某个节点
	Find(serviceName string, nodeName string) (*Node, error)
	// 查询某服务的所有子节点
	FindAll(serviceName string) (*Service, error)
}
