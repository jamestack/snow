package example

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jamestack/snow"
)

func (u *User) Play(name string) string {
	return snow.Response(u.Node, func() string {
		return name
	}, name)
}

func (u *User) Read(id int) (*snow.Stream, error) {
	return snow.StreamResponse(u.Node, func(stream *snow.Stream) {
		fmt.Println("accept", id, stream)
		// go func() {
		for {
			data, err := stream.ReadData()
			if err != nil {
				fmt.Println("stream.ReadData()", err)
				return
			}
			fmt.Println("user server receve", string(data))
			fmt.Println(stream.Write(data))
		}
		// }()
	}, id)
}

func init() {
	gob.Register(errors.New(""))
}

func TestBenchNodeA(t *testing.T) {
	cluster := snow.NewClusterMaster("127.0.0.1:8000", "127.0.0.1:8000")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println("err_str:[" + err.Error() + "]")
		return
	}

	james, err := snow.Mount(cluster, "james", &User{name: "james"})
	// james.RPC.node = james
	// _, err = cluster.Mount("james", &User{name: "james"})
	fmt.Println("Mount() james", err)

	stream, err := james.RPC.Read(100)
	if err != nil {
		fmt.Println("james.RPC.Read() err:", err)
	}

	fmt.Println("local stream start", stream)

	stream.Write([]byte("1"))
	stream.ReadData()
	stream.Write([]byte("2"))
	stream.ReadData()
	stream.Write([]byte("3"))
	stream.ReadData()
	stream.Close()

	fmt.Println("local stream end")

	s := <-done
	fmt.Println("Got signal:", s)
}

func TestBenchNodeB(t *testing.T) {
	cluster := snow.NewClusterSlave("127.0.0.1:8001", "127.0.0.1:8001", "127.0.0.1:8000")
	_, err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	jacks, err := snow.Mount(cluster, "jacks", &User{name: "jacks"})
	// jacks.RPC.node = jacks
	// jacks, err := cluster.Mount("jacks", &User{name: "jacks"})
	if err != nil {
		fmt.Println("mount", err)
		return
	}

	// var MAX_ROUND = 100000
	var MAX_ROUND = 1

	// 本地调用
	log.Println("local rpc start")
	for i := 0; i < MAX_ROUND; i++ {
		name := fmt.Sprintf("section-%d", i+1)
		res := jacks.RPC.Play(name)
		// _ = res
		fmt.Println("local play done", res)
	}
	log.Println("local rpc end")
	log.Println()

	// return

	james, err := snow.Find[User](cluster, "james")
	// james.RPC.node = james
	// james, err := cluster.Find("james")
	if err != nil {
		panic(err)
	}

	// err = james.Call("TestErr", errors.New("123")).Then(func(str string, err error) {
	// 	fmt.Println(str, err, err == nil, reflect.TypeOf(err))
	// })
	// fmt.Println(err)

	// 远程调用
	log.Println("remote rpc start")
	for i := 0; i < MAX_ROUND; i++ {
		name := fmt.Sprintf("section-%d", i+1)
		res := james.RPC.Play(name)
		// _ = res
		fmt.Println("remote rpc play done", res)
		if err != nil {
			fmt.Println("end", err)
			return
		}
	}
	log.Println("remote rpc end")
	log.Println()

	stream, err := james.RPC.Read(100)
	if err != nil {
		fmt.Println("remote james.RPC.Read() err:", err)
	}

	fmt.Println("remote stream start", stream)

	stream.Write([]byte("1"))
	fmt.Println(stream.ReadData())
	stream.Write([]byte("2"))
	fmt.Println(stream.ReadData())
	stream.Write([]byte("3"))
	fmt.Println(stream.ReadData())
	stream.Close()

	fmt.Println("remote stream end")

	<-time.After(30 * time.Second)
}
