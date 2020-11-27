package snow

import (
	"errors"
	"sync"
)

// 无限缓存Channel
type Channel struct {
	lock sync.Mutex
	nw   *sRing
	nr   *sRing
	receiveCh chan struct{}
	close bool
	cap int
	maxCap int
}

type sRing struct {
	value interface{}
	next  *sRing
}

var ErrEmpty = errors.New("[Channel: Empty]")
var ErrClosed = errors.New("[Channel: Closed]")

// todo 动态设置最大容量
//func (q *Channel) SetMaxSize(size uint) {
//
//}

// 向Channel发送数据
func (q *Channel) Send(v interface{}) (ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.close {
		return false
	}

	if q.nw == nil {
		head := &sRing{}
		head.next = head
		q.nw = head
		q.nr = head
		q.cap = 1
		q.receiveCh = make(chan struct{}, 1)
	}

	q.nw.value = v
	next := q.nw.next
	if next == q.nr {
		list := make([]sRing, q.cap, q.cap)
		for i:=0;i<q.cap-1;i++ {
			list[i].next = &list[i+1]
		}

		list[q.cap-1].next = next
		q.nw.next = &list[0]

		next = &list[0]

		q.cap += q.cap
	}
	q.nw = next

	select {
		case q.receiveCh <- struct{}{}:
	default:

	}

	return true
}

// 阻塞直到取到值或者该队列关闭
// @sleep 抢占式调度时的频率控制，达到任务或者线程优先级控制的目的
func (q *Channel) Receive() (v interface{}, ok bool) {
	value,err := q.Get()
	switch err {
	case nil:
		return value, true
	case ErrEmpty:
		_,ok := <-q.receiveCh
		if ok {
			return q.Receive()
		}
	case ErrClosed:
		return nil, false
	}
	return nil, false
}

// 关闭队列并回收队列
func (q *Channel) Close() (ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.close {
		return false
	}

	// 更新关闭状态
	q.close = true
	close(q.receiveCh)

	// 回收内存
	q.gc()

	return true
}

// 从Chanel中取数据，不阻塞，如果队列为空或者已关闭则返回错误
func (q *Channel) Get() (v interface{}, err error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.nr == q.nw {
		if q.close {
			err = ErrClosed
		}else {
			err = ErrEmpty
		}
		return
	}

	v = q.nr.value
	q.nr.value = nil
	q.nr = q.nr.next

	// 完全关闭后，自动回收内存
	if q.close && q.nr == q.nw {
		q.gc()
	}
	return
}

func (q *Channel) Peek() (v interface{}, err error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.nr == q.nw {
		if q.close {
			err = ErrClosed
		}else {
			err = ErrEmpty
		}
		return
	}

	v = q.nr.value
	return
}

// 内存回收，避免已close的queue占用大量内存无法被Runtime使用
func (q *Channel) gc() {
	if q.nw != nil && q.nr != nil {
		var next *sRing
		for n:=q.nw; n.next != q.nw && n.next != q.nr;n = n.next {
			next = n.next
			q.cap -= 1
		}
		if next != nil {
			q.nw.next = next
		}
	}
}
