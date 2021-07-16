package snow

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func init() {
}

var randIndex = rand.New(rand.NewSource(time.Now().UnixNano()))

func checkPanic() {
	if err := recover(); err != nil {
		printStack(err)
	}
}

func printStack(err interface{}) {
	stack := make([]byte, 16384)
	n := runtime.Stack(stack, false)
	stack = stack[:n]
	fmt.Printf("[Snow] Found Panic Err: %v\n", err)
	fmt.Println(string(stack))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 随机字符串id生成器
func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randIndex.Intn(len(letterRunes))]
	}
	return string(b)
}
