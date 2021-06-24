package example

import (
	"fmt"
	"github.com/JamesWone/snow"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
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
	var num int32 = 0
	wg := sync.WaitGroup{}
	for i :=0; i<200000;i++ {
		wg.Add(1)
		p.AfterFunc(2*time.Second, func() {
			wg.Done()
			atomic.AddInt32(&num, 1)
		})
	}

	wg.Wait()
	fmt.Println("hello after 2s", num, time.Now())

	<-time.After(5*time.Second)
}
