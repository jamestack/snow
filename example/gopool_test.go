package example

import (
	"fmt"
	"github.com/JamesWone/snow"
	"sync"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	p0 := snow.NewGoPool(32)
	p := snow.NewGoPool(64)

	fmt.Println("start ", time.Now())
	p.Go(func() {
		fmt.Println("go 1", time.Now())
	})

	p.AfterFunc(3*time.Second, func() {
		fmt.Println("hello after 3s", time.Now())
	})

	fmt.Println("destroy start", time.Now())
	//p.SetWorkerNum(0)
	fmt.Println("destroy end", time.Now())

	// panic

	p.Go(func() {
		fmt.Println("go 2", time.Now())
	})

	p.AfterFunc(5*time.Second, func() {
		fmt.Println("hello after 5s", time.Now())
	})
	num := 0
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	for i :=0; i<200000;i++ {
		// 异步安全性测试
		wg.Add(1)
		p0.Go(func() {
			p.AfterFunc(2*time.Second, func() {
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				num ++
				//fmt.Println("hello after 2s", time.Now())
			})
		})
	}


	wg.Wait()
	fmt.Println(num, time.Now())

	<-time.After(5*time.Second)
}
