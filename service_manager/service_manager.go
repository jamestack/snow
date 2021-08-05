package service_manager

import (
	"errors"
	"net"
	"os"
	"runtime"
	"sort"

	"github.com/jamestack/snow"
)

type ServiceManager struct {
	*snow.Node
	listener      net.Listener
	Service       []ServiceInfo // 可启动的服务定义
	WebListenAddr string        // 是否启动web管理界面
}

// 参数信息
type ParamInfo struct {
	Name    string // 名称
	Remark  string // 备注
	Default string // 默认值
}

// 方法信息
type MethodInfo struct {
	Name   string      // 名称
	Remark string      // 备注
	Params []ParamInfo // 参数信息
}

// 服务信息
type ServiceInfo struct {
	Name    string             // 服务名
	Inode   func() snow.INode `json:"-"` // iNode对象
	Remark  string             // 备注
	Methods []MethodInfo       // 可执行的方法名
	Fields  []string           // 可查看的属性
}

type RuntimeInfo struct {
	GOOS         string
	GOARCH       string
	GoVersion    string
	NumCpu       int
	NumGoroutine int
	NumCgoCall   int64
}

type OsInfo struct {
	Hostname string // 主机名
	Ip       string // 本地ip
}

type ActiveService struct {
	Name      string
	MountTime int64
}

// 节点实例信息
type NodeInfo struct {
	Name        string
	Runtime     RuntimeInfo   // go运行时信息
	SnowVersion string        // Snow版本号
	Os          OsInfo        // 主机信息
	Services    []ServiceInfo // 已启动服务
	Active      []ActiveService
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("can not find the client ip address")
}

// 节点信息
func (s *ServiceManager) NodeInfo() NodeInfo {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "UnKnow"
	}
	ip, err := getClientIp()
	if err != nil {
		ip = "UnKnow"
	}
	res := NodeInfo{
		Name: s.Name(),
		Runtime: RuntimeInfo{
			GOOS:         runtime.GOOS,
			GOARCH:       runtime.GOARCH,
			GoVersion:    runtime.Version(),
			NumCpu:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCgoCall:   runtime.NumCgoCall(),
		},
		Os: OsInfo{
			Hostname: hostName,
			Ip:       ip,
		},
		SnowVersion: snow.Version,
		Services:    s.Service,
		Active:      []ActiveService{},
	}

	for _, item := range s.FindLocalAll() {
		if item.Name() == s.Name() {
			continue
		}

		res.Active = append(res.Active, ActiveService{
			Name:      item.Name(),
			MountTime: item.MountTime(),
		})
	}

	sort.Slice(res.Active, func(i, j int) bool {
		a, b := res.Active[i], res.Active[j]
		if a.MountTime == b.MountTime {
			return a.Name < b.Name
		} else {
			return a.MountTime < b.MountTime
		}
	})

	return res
}

// 服务名
func (s *ServiceManager) Services() []ServiceInfo {
	return s.Service
}
