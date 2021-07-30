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

// 自定义消息处理器，模拟全局消息总线
//func (u *User) OnCall(name string, call func([]reflect.Value) []reflect.Value, args []reflect.Value) []reflect.Value {
//	var res []reflect.Value
//	ch := make(chan bool)
//
//	// 丢到全局消息总线中去执行
//	GlobalProcessor.Go(func() {
//		fmt.Println("Call", name, "Start")
//		res = call(args)
//		fmt.Println("Call", name, "End")
//		ch <- true
//	})
//
//	<-ch
//	return res
//}

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

	_ = userManager.Call("Name", func(name string) {
		fmt.Println("call done")
	})

	err = userManager.Call("TestErr", errors.New("myErr"), func(name string, err error) {
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

	err = userManager.Call("Name", func(name string) {
		fmt.Println("call done")
	})
	fmt.Println("call err:", err)

}

//
//// Master
//func TestServeMaster(t *testing.T) {
//	err,done := snow.ServeMaster("127.0.0.1:8001", "127.0.0.1:8001")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	fmt.Println(err)
//}
//
//// Peer
//func TestServePeer(t *testing.T) {
//	err,done := snow.ServePeer("127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8002")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	_, err = snow.Mount("james", &User{name: "james"})
//
//	<-time.After(3000*time.Second)
//	fmt.Println(err)
//}
//
//func TestServePeer2(t *testing.T) {
//	err,done := snow.ServePeer("127.0.0.1:8000", "127.0.0.1:8003", "127.0.0.1:8003")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	_, err = snow.Mount("smith", &User{})
//
//	<-time.After(5*time.Second)
//	_ = snow.UnMount("smith")
//
//	snow.Find("")
//	<-time.After(3000*time.Second)
//	fmt.Println(err)
//
//}
//
//// 路径搜索算法
//func TestFindPath(t *testing.T) {
//	err,done := snow.ServePeer("127.0.0.1:8001", "127.0.0.1:8004", "127.0.0.1:8004")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	_, _ = snow.Mount("lucy", &User{name: "lucy"})
//
//	lucy := snow.Find("lucy")
//	fmt.Println(lucy)
//	err = lucy.Call("Name", func(name string) {
//		fmt.Println("name =", name)
//	})
//
//	lucy.Mount("jack", &User{name: "jack"})
//	lucy.Find("jack").Call("Name", func(name string) {
//		fmt.Println("jack =", name)
//	})
//
//	err = snow.Find("lucy/jack").Call("Name", func(name string) {
//		fmt.Println("lucy/jack =", name)
//	})
//	fmt.Println(err)
//
//	err = snow.Find("/lucy/jack").Call("Name", func(name string) {
//		fmt.Println("/lucy/jack =", name)
//	})
//	fmt.Println(err)
//
//	//snow.Find("james").Call("Name", func(name string) {
//	//	fmt.Println("james =", name)
//	//})
//
//	//snow.Find("/james/stack").Call("Name", func(name string) {
//	//	fmt.Println("/james/stack =", name)
//	//})
//
//	//snow.Find("james").Find("stack").Call("Name", func(name string) {
//	//	fmt.Println("/james/stack =", name)
//	//})
//
//	//notfound := snow.Find("/dsa/sadf")
//	//fmt.Println("notfound =", notfound)
//
//}
//
//// 测试远程call
//func TestPeerCall(t *testing.T) {
//	err,done := snow.ServePeer("127.0.0.1:8001", "127.0.0.1:8004", "127.0.0.1:8004")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	err = snow.Find("james").Call("Name", func(name string) {
//		fmt.Println("call =", name)
//	})
//	fmt.Println("call err =", err)
//
//	err = snow.Find("james").Call("SetAge", 17, func(age int64) {
//		fmt.Println("call SetAge =", age)
//	})
//	fmt.Println("call err =", err)
//
//
//}
//
//func (u *User) StreamTest(name string, stream *snow.Stream) {
//	fmt.Println("user.StreamTest", name)
//	for {
//		data,err := stream.ReadData()
//		if err != nil {
//			fmt.Println("u read err:", err)
//			break
//		}
//		fmt.Println("u read", string(data), err)
//		fmt.Println(stream.Write(data))
//		fmt.Println(stream.Write(data))
//		fmt.Println("close err:", stream.Close())
//		fmt.Println("close err:", stream.Close())
//	}
//}
//
////
//func TestStream(t *testing.T) {
//	err,done := snow.ServeMaster("", "")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func() {
//		<-done
//	}()
//
//	node,err := snow.Mount("test", &User{})
//	fmt.Println("mount err:", err)
//
//	var stream *snow.Stream
//	fmt.Println("stream start")
//	stream,err = node.Stream("StreamTest", "james")
//	fmt.Println("stream err:", err)
//
//	fmt.Println("stream write start")
//	_,err = stream.Write([]byte("hello"))
//	fmt.Println("stream write err:", err)
//	for {
//		msg,err := stream.ReadData()
//		if err != nil {
//			fmt.Println("t read err:", err)
//			_, err = stream.Write([]byte("123"))
//			fmt.Println("t write err:", err)
//			break
//		}
//		fmt.Println("t read", string(msg), err)
//	}
//
//}
