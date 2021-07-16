package example

import (
	"fmt"
	"testing"
	"time"

	"github.com/jamestack/snow"
)

func TestChannel(t *testing.T) {
	ch := &snow.Channel{}
	//fmt.Println(ch.Close())

	fmt.Println(ch.Get())
	fmt.Println(ch.Send(1))
	fmt.Println(ch.Send(2))
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())
	fmt.Println(ch.Send(3))
	fmt.Println(ch.Send(4))
	fmt.Println(ch.Send(5))
	fmt.Println(ch.Send(6))
	fmt.Println(ch.Close())
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())
	fmt.Println(ch.Get())

}

func TestChannelBench(t *testing.T) {
	ch := snow.Channel{}

	fmt.Println("send start", time.Now())
	for i := 0; i < 100; i++ {
		ch.Send(i)
	}
	fmt.Println("send end", time.Now())

	fmt.Println("get start", time.Now())
	for i := 0; i < 100; i++ {
		fmt.Println(ch.Get())
	}
	fmt.Println("get end", time.Now())

	ch.Close()

	fmt.Println("close end")

	for {
		_, ok := ch.Receive()
		if !ok {
			fmt.Println("closed")
			break
		}
	}

	fmt.Println("receive end")

	//<-time.After(30*time.Second)
}

func benchSysChanel(num int64) (res int64) {
	ch := make(chan int64, 1024)
	done := make(chan bool)
	go func() {
		for range ch {
			res += 1
		}
		done <- true
	}()
	for i := int64(1); i <= num; i++ {
		ch <- i
	}
	close(ch)

	<-done
	return
}

func benchSnowChannel(num int64) (res int64) {
	ch := snow.NewChannel()
	done := make(chan bool)

	go func() {
		for {
			_, ok := ch.Receive()
			if !ok {
				done <- true
				return
			}

			res += 1
		}
	}()

	for i := int64(1); i <= num; i++ {
		ch.Send(i)
	}
	ch.Close()

	<-done
	return
}

func benchTwoChannel(num int64) {
	fmt.Println("benchSysChanel start")
	start := time.Now()
	fmt.Println(benchSysChanel(num))
	fmt.Println("benchSysChanel end", time.Since(start))

	fmt.Println("benchSnowChannel start")
	start = time.Now()
	fmt.Println(benchSnowChannel(num))
	fmt.Println("benchSnowChannel end", time.Since(start))
}

func TestBenchChanelChannel(t *testing.T) {
	fmt.Println("-----------bench 1000 1k ---------")
	benchTwoChannel(1000)
	fmt.Println("-----------bench 100000 10w---------")
	benchTwoChannel(100000)
	fmt.Println("-----------bench 10000000 1000w ---------")
	benchTwoChannel(10000000)
	fmt.Println("-----------bench 100000000 1e ---------")
	benchTwoChannel(100000000)
	fmt.Println("-----------bench 1000000000 10e ---------")
	benchTwoChannel(1000000000)
}
