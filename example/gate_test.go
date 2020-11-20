package example

import (
	"fmt"
	"snow"
	"testing"
)

// 启动网关节点
func TestGate(t *testing.T)  {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8000","127.0.0.1:8000")
	done,err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		<-done
	}()

	// 挂载3个网关节点
	_, err = cluster.Mount("Gate/1", &Gate{
		ListenAddr: "127.0.0.1:81",
		TargetNode: "Game/1",
	})
	_, err = cluster.Mount("Gate/2", &Gate{
		ListenAddr: "127.0.0.1:82",
		TargetNode: "Game/1",
	})
	_, err = cluster.Mount("Gate/3", &Gate{
		ListenAddr: "127.0.0.1:83",
		TargetNode: "Game/1",
	})
	fmt.Println("mount err:", err)

}

// 启动游戏节点
func TestGame(t *testing.T)  {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8001","127.0.0.1:8001")
	done,err := cluster.Serve()
	fmt.Println(err)
	defer func() {
		<-done
	}()

	game := &Game{}

	// 挂载游戏节点
	_, err = cluster.Mount("Game/1", game)
	fmt.Println("mount err:", err)

	fmt.Println(cluster.FindAll("Game"))

}
