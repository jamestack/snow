package snow

import (
	"fmt"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	q := NewPriorityQueue()
	fmt.Println(q.Pop())
	q.Push(1,-1)
	q.Push(2,-5)
	q.Push(3, -2)

	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())

}
