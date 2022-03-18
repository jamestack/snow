package example

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type Number interface {
	int | int32 | float32 | float64 | string
}

func max[T Number](a T, b T) bool {
	return a > b
}

func any_test(a any) {
	fmt.Println(a)
}

type IEat interface {
	eat()
}

type myInt int

func (this myInt) eat() {
	fmt.Println("eat:", this)
}

type Man interface {
	IEat
	~int | ~float32
}

func gen_test[T Man](s T) {
	s.eat()
}

func generecs_test[T any](in T) T {
	return in
}

type Persons[T any] struct {
	list []T
}

func (p *Persons[T]) res() []T {
	return p.list
}

type LoginRequest interface {
	Req()
}

func Call[Res any](method string, req any) (Res, error) {
	var a Res
	return a, nil
}

type LoginServer struct {
}

func (s *LoginServer) Login(username string) error {
	return errors.New("123")
}

func (s *LoginServer) Logout(username string) error {
	return errors.New("456")
}

func TestGenerics(t *testing.T) {
	println(max("3", "3"))
	println(max(3, 2))
	any_test(123)

	// 泛型类型集测试
	gen_test[myInt](123)

	// 模板函数
	fmt.Println(generecs_test("1234"))

	// 模板结构体
	james := Persons[int]{
		list: []int{1, 2, 3},
	}
	fmt.Println(james.res())

	res, err := Call[int]("Login", 123)
	fmt.Println(res, err)

	login_server := &LoginServer{}

	res_2, err := CallFunc(login_server.Login, "123")
	fmt.Println(res_2, err)

	// fmt.Println(CallFunc(login_server.Logout, "1"))
}

// 普通处理器
type Handler[Req any, Res any] func(Req) Res

// 流处理器
type StreamHandler[Req any, Res any] struct {
}

// 异步处理器
type AsyncHandler[Req any, Res any] struct {
}

func CallFunc[Req any, Res any](handler Handler[Req, Res], args Req) (Res, error) {
	name := getFuncName(handler)
	fmt.Println(name)

	return handler(args), nil
}

func getFuncName(fn any) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	s_l, s_r := strings.LastIndexByte(name, '*'), strings.LastIndexByte(name, '.')

	server_name := name[s_l+1 : s_r-1]
	method_name := name[s_r+1 : len(name)-3]

	return server_name + "/" + method_name
}

type Promise[Req any, Res any] struct {
	handler Handler[Req, Res]
	res     Res
}

func Call_2[Req any, Res any](handler Handler[Req, Res]) *Promise[Req, Res] {
	return &Promise[Req, Res]{
		handler: handler,
	}
}

// 异步回调
func (p *Promise[Req, Res]) Then(func(res Res, err error)) {

}

// 异步调用，没有返回值
func (p *Promise[Req, Res]) Async() (err error) {

	return nil
}

// 同步调用，等待返回
// func (p *Promise[Req, Res]) Await() (res Res, err error) {
// 	return nil, nil
// }
