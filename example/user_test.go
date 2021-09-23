package example

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jamestack/snow"
)

type Root struct {
	*snow.Node
}

type User struct {
	*snow.Node
	name string
}

// 模拟全局消息总线
var GlobalProcessor = snow.NewGoPool(1024)

func (u *User) TestErr(req error) (string, error) {
	fmt.Println("call")
	return "hello", req
}

// 挂载回调
func (u *User) OnMount() {
	fmt.Println("mount")
}

//
func (u *User) OnUnMount() {
	fmt.Println("unmount")
}

func (u *User) Name() string {
	fmt.Println("u.Name() = ", u.name)
	return u.name
}

func (u *User) SetAge(age *int64) *int64 {
	fmt.Println("age =", *age)
	return age
}

func TestUserNode(t *testing.T) {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8001", "127.0.0.1:8001")
	_, _ = cluster.Serve()

	userManager, err := cluster.Mount("UserManager", &User{
		name: "james",
	})
	if err != nil {
		fmt.Println(err)
	}

	_ = userManager.CallAsync("Name").Then(func(name string) {
		fmt.Println("call done")
	})

	err = userManager.CallAsync("TestErr", errors.New("myErr")).Then(func(name string, err error) {
		fmt.Println(name, err)
	})

	fmt.Println("call err:", err)
}

func TestConsul(t *testing.T) {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8001", "127.0.0.1:8001")
	done, _ := cluster.Serve()
	defer func() {
		<-done
	}()

	userManager, err := cluster.Mount("UserManager", &User{
		name: "james",
	})

	fmt.Println(err)

	list, err := cluster.FindAll("UserManager")
	fmt.Println(err, list)

	err = userManager.CallAsync("Name").Then(func(name string) {
		fmt.Println("call done")
	})
	fmt.Println("call err:", err)

}
