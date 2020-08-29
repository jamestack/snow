package example

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"snow"
	"time"
)

func (u *User) Play(name string) string {
	// 模拟长连接
	//<-time.After(5*time.Second)

	return name
}

func TestBenchNodeA(t *testing.T) {
	err,done := snow.ServeMaster("127.0.0.1:8000", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		<-done
	}()

	_, err = snow.Mount("james", &User{name: "james"})
	fmt.Println(err)

	<-time.After(60*time.Minute)
}

func TestBenchNodeB(t *testing.T) {
	err,done := snow.ServePeer("127.0.0.1:8000", "127.0.0.1:8005", "127.0.0.1:8005")
	if err != nil {
		fmt.Println("serve err:", err)
		return
	}
	defer func() {
		return
		<-done
	}()

	_, err = snow.Mount("jacks", &User{name: "jacks"})
	if err != nil {
		fmt.Println("mount", err)
		return
	}

	var MAX_ROUND = 100000

	// 本地调用
	log.Println("local rpc start")
	for i:=0;i<MAX_ROUND;i++ {
		name := fmt.Sprintf("section-%d", i+1)
		_ = snow.Find("jacks").Call("Play", name, func(name string) {
			//fmt.Println("play done", name)
		})
	}
	log.Println("local rpc end")
	log.Println()

	// 远程调用
	log.Println("rpc start")
	for i:=0;i<MAX_ROUND;i++ {
		name := fmt.Sprintf("section-%d", i+1)
		err := snow.Find("james").Call("Play", name, func(name string) {
			//fmt.Println("play done", name)
		})
		if err != nil {
			fmt.Println("end", err)
			return
		}
	}
	log.Println("rpc end")
	log.Println()

	// 远程异步调用
	var wg sync.WaitGroup
	log.Println("async rpc start")
	for i:=0;i<MAX_ROUND;i++ {
		name := fmt.Sprintf("section-%d", i+1)
		wg.Add(1)
		go func() {
			err := snow.Find("james").Call("Play", name, func(name string) {
				//fmt.Println("play done", name)
			})
			wg.Done()
			if err != nil {
				fmt.Println("end", err)
			}
		}()
	}
	wg.Wait()
	log.Println("async rpc end")

	//<-time.After(30*time.Second)
}
