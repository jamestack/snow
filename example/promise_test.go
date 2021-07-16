package example

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jamestack/snow/promise"
)

type Person struct {
	name string
}

// 设置名字
// @name string 名称
func (p *Person) SetName(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) {
	name := args[0].(string)
	fmt.Println("Start SetName...2s")
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("SetName:", args)
		p.name = name
		resolve("ok")
	})
}

func (p *Person) Say(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) {
	fmt.Println("Say:", args)
	reject(errors.New("this is a err"))
	resolve("hello, my name is " + p.name)
}

func (p *Person) Eat(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) {
	fmt.Println("Eat:", args)
	resolve("i am eating.")
}

func TestPromise(t *testing.T) {
	p := Person{}
	<-promise.New(p.SetName, "james").
		Then(p.Say).
		Then(p.Eat).Then(func(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) {
		fmt.Println("qwwe")
		resolve(123)
	}).Catch(func(errs ...error) {
		fmt.Println("catch:", errs)
	})

	//promise.New(p.Eat).Catch(nil)
	//<-time.After(3*time.Second)
}

func TestPromise2(t *testing.T) {
	p := Person{}

	ps := promise.New(p.SetName, "james")
	ps.Then(p.Say)
	<-time.After(time.Second)
	ps.Then(func(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) {
		fmt.Println("hello")
		resolve()
	})
	<-time.After(time.Second)
	<-ps.Catch(func(errs ...error) {
		fmt.Println("catch:", errs)
	})
}
