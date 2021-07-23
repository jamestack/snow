package service_manager

import (
	"fmt"
	"testing"

	"github.com/jamestack/snow"
)

type User struct {
	*snow.Node
}

func TestServiceManager(t *testing.T) {
	cluster := snow.NewClusterWithLocal()
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		<-done
	}()

	_, _ = cluster.MountRandNode("SnowNodes", &ServiceManager{
		WebListenAddr: "127.0.0.1:8080",
		Service: []ServiceInfo{
			{
				Name: "Gate",
				Inode: func() interface{} {
					return &User{}
				},
				Remark:  "网关组件",
				Methods: nil,
				Fields:  nil,
			},
		},
	})

}
