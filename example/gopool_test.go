package example

import (
	"fmt"
	"github.com/JamesWone/snow"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	p := snow.GoPool{}
	p.SetMaxWorker(1)

	fmt.Println("start ", time.Now())
	p.Go(func() {
		fmt.Println("go 1", time.Now())
	})

	p.AfterFunc(3*time.Second, func() {
		fmt.Println("hello after 3s", time.Now())
	})

	fmt.Println("destroy start", time.Now())
	p.SetMaxWorker(0)
	fmt.Println("destroy end", time.Now())

	// panic

	p.Go(func() {
		fmt.Println("go 2", time.Now())
	})

	p.AfterFunc(3*time.Second, func() {
		fmt.Println("hello after 3s", time.Now())
	})
	p.AfterFunc(3*time.Second, func() {
		fmt.Println("hello after 3s", time.Now())
	})

	<-time.After(30 * time.Second)
}
