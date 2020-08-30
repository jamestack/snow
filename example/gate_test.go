package example

import (
	"fmt"
	"snow"
	"testing"
)

// 启动网关节点
func TestGate(t *testing.T)  {
	err,done := snow.ServeMaster("127.0.0.1:8000", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		<-done
	}()

	// 挂载3个网关节点
	_, err = snow.Mount("Gate_1", &Gate{
		ListenAddr: "127.0.0.1:81",
		TargetNode: "/Game_1",
	})
	_, err = snow.Mount("Gate_2", &Gate{
		ListenAddr: "127.0.0.1:82",
		TargetNode: "/Game_1",
	})
	_, err = snow.Mount("Gate_3", &Gate{
		ListenAddr: "127.0.0.1:83",
		TargetNode: "/Game_1",
	})
	fmt.Println("mount err:", err)

}

// 启动游戏节点
func TestGame(t *testing.T)  {
	err,done := snow.ServePeer("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		<-done
	}()

	game := &Game{}

	// 挂载游戏节点
	_, err = snow.Mount("Game_1", game)
	fmt.Println("mount err:", err)

}
