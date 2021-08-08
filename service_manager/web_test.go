package service_manager

import (
	"fmt"
	"testing"

	"github.com/jamestack/snow"
)

type User struct {
	*snow.Node
}

func TestServiceManagerA(t *testing.T) {
	cluster := snow.NewClusterMaster("127.0.0.1:8000", "127.0.0.1:8000")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = cluster.MountRandNode("SnowNodes", &ServiceManager{
		WebListenAddr: "127.0.0.1:8080",
		Service: []ServiceInfo{
			{
				Name: "Gate",
				Inode: func() snow.INode {
					return &User{}
				},
				Remark:  "网关组件",
				Methods: nil,
				Fields:  nil,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	<-done
}

func TestServiceManagerB(t *testing.T) {
	cluster := snow.NewClusterSlave("127.0.0.1:8001", "127.0.0.1:8001", "127.0.0.1:8000")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = cluster.MountRandNode("SnowNodes", &ServiceManager{
		WebListenAddr: "",
		Service: []ServiceInfo{
			{
				Name: "Game_1",
				Inode: func() snow.INode {
					return &User{}
				},
				Remark:  "游戏服务1",
				Methods: nil,
				Fields:  nil,
			},
			{
				Name: "Game_2",
				Inode: func() snow.INode {
					return &User{}
				},
				Remark:  "游戏服务2",
				Methods: nil,
				Fields:  nil,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	<-done
}
