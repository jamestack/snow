package snow

import (
	"errors"
	"math"
	"sync"
)

// 无限缓存Channel
type Channel struct {
	lock sync.Mutex
	nw   *sRing
	nr   *sRing
	receiveCh chan struct{}
	close bool
	cap uint
	addCap uint
	maxCap uint
	chanLock sync.Mutex
}

func (q *Channel) getReceiveChan() chan struct{} {
	if q.receiveCh == nil {
		q.chanLock.Lock()
		if q.receiveCh == nil {
			q.receiveCh = make(chan struct{})
		}
		q.chanLock.Unlock()
	}
	return q.receiveCh
}

type sRing struct {
	value interface{}
	next  *sRing
}

var ErrEmpty = errors.New("[Channel: Empty]")
var ErrClosed = errors.New("[Channel: Closed]")

// 向Channel发送数据
func (q *Channel) Send(v interface{}) (ok bool) {
	if q.close {
		return false
	}

	if q.nw == nil {
		q.lock.Lock()
		if q.nw == nil {
			head := &sRing{}
			head.next = head
			q.nw = head
			q.nr = head
			q.cap = 1
			q.addCap = 0
			q.maxCap = math.MaxInt32
		}
		q.lock.Unlock()
	}

	q.lock.Lock()
	q.nw.value = v
	if q.nw.next == q.nr {
		if q.cap < q.maxCap {
			if q.cap < 1024 {
				q.addCap = q.cap
			}else {
				q.addCap = q.cap / 4
			}

			list := make([]sRing, q.addCap, q.addCap)
			for i:=uint(0);i<q.addCap-1;i++ {
				list[i].next = &list[i+1]
			}

			list[q.addCap-1].next = q.nw.next
			q.nw.next = &list[0]

			q.cap += q.addCap
		}else {
			q.lock.Unlock()
			// todo maxCap
			return false
		}
	}
	q.nw = q.nw.next

	q.lock.Unlock()

	select {
		case q.receiveCh <- struct{}{}:
	default:

	}

	return true
}

// 阻塞直到取到值或者该队列关闭
func (q *Channel) Receive() (v interface{}, ok bool) {
	value,err := q.Get()
	switch err {
	case nil:
		return value, true
	case ErrEmpty:
		<-q.getReceiveChan()
		return q.Receive()
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

	return true
}

// 从Chanel中取数据，不阻塞，如果队列为空或者已关闭则返回错误
func (q *Channel) Get() (v interface{}, err error) {
	q.lock.Lock()

	if q.nr == q.nw {
		if q.close {
			err = ErrClosed
		}else {
			err = ErrEmpty
		}
		q.lock.Unlock()
		return
	}

	v = q.nr.value
	q.nr.value = nil
	q.nr = q.nr.next

	q.lock.Unlock()

	// 完全关闭后，自动回收内存
	if q.close && q.nr == q.nw {
		go q.gc()
	}
	return
}

func (q *Channel) Peek() (v interface{}, err error) {
	q.lock.Lock()

	if q.nr == q.nw {
		if q.close {
			q.lock.Unlock()
			err = ErrClosed
		}else {
			q.lock.Unlock()
			err = ErrEmpty
		}
		return
	}

	v = q.nr.value
	q.lock.Unlock()
	return
}

// 内存回收，避免已close的queue占用大量内存无法被Runtime使用
func (q *Channel) gc() {
	//fmt.Println("cap", q.cap)
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

	if q.receiveCh != nil {
		close(q.receiveCh)
	}
}
