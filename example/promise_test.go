package example

import (
	"fmt"
	"github.com/JamesWone/snow/promise"
	"testing"
)

type Person struct {
	name string
}

// 设置名字
// @name string 名称
func (p *Person) SetName(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) promise.Exit {
	name := args[0].(string)
	fmt.Println("SetName:", args)
	p.name = name
	return resolve("edew")
}

func (p *Person) Say(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) promise.Exit {
	fmt.Println("Say:", args)
	//reject(errors.New("err"))
	return resolve("hello, my name is " + p.name)
}

func (p *Person) Eat(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) promise.Exit {
	fmt.Println("Eat:", args)
	return resolve("i am eating.")
}

func TestPromise(t *testing.T) {
	p := Person{}
	<-promise.New(p.SetName, "james").
		Then(p.Say).
		Then(p.Eat).Then(func(resolve promise.Resolve, reject promise.Reject, args ...promise.Any) promise.Exit {
		fmt.Println("qwwe")
		return resolve(123)
	}).Catch(func(errs ...error) {
		fmt.Println("catch:", errs)
	})

	<-promise.New(p.Eat).Catch(nil)
	//<-time.After(3*time.Second)
}
