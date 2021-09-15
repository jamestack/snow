package example

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/jamestack/snow"
)

func (u *User) Play(name string) string {
	// 模拟长连接
	//<-time.After(5*time.Second)

	return name
}

func init() {
	gob.Register(errors.New(""))
}

func TestBenchNodeA(t *testing.T) {
	cluster := snow.NewClusterMaster("127.0.0.1:8000", "127.0.0.1:8000")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = cluster.Mount("james", &User{name: "james"})
	fmt.Println(err)

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

	jacks, err := cluster.Mount("jacks", &User{name: "jacks"})
	if err != nil {
		fmt.Println("mount", err)
		return
	}

	var MAX_ROUND = 100000

	// 本地调用
	log.Println("local rpc start")
	for i := 0; i < MAX_ROUND; i++ {
		name := fmt.Sprintf("section-%d", i+1)
		err = jacks.Call("Play", name).Then(func(name string) {
			//fmt.Println("play done", name)
		})
		if err != nil {
			fmt.Println("end", err)
			return
		}
	}
	log.Println("local rpc end")
	log.Println()

	james, err := cluster.Find("james")
	if err != nil {
		panic(err)
	}

	err = james.Call("TestErr", errors.New("123")).Then(func(str string, err error) {
		fmt.Println(str, err, err == nil, reflect.TypeOf(err))
	})
	fmt.Println(err)

	pool := snow.NewGoPool(uint32(runtime.NumCPU()) * 4)
	// 远程异步调用
	var wg sync.WaitGroup
	log.Println("async rpc start")
	for i := 0; i < MAX_ROUND; i++ {
		name := fmt.Sprintf("section-%d", i+1)
		wg.Add(1)
		pool.Go(func() {
			err = james.Call("Play", name).Then(func(name string) {
				//fmt.Println("play done", name)
			})
			wg.Done()
			if err != nil {
				fmt.Println("end", err)
			}
		})
	}
	wg.Wait()
	log.Println("async rpc end")
	log.Println()

	// 远程调用
	log.Println("rpc start")
	for i := 0; i < MAX_ROUND; i++ {
		name := fmt.Sprintf("section-%d", i+1)
		err = james.Call("Play", name).Then(func(name string) {
			//fmt.Println("play done", name)
		})
		if err != nil {
			fmt.Println("end", err)
			return
		}
	}
	log.Println("rpc end")
	log.Println()

	//<-time.After(30*time.Second)
}
