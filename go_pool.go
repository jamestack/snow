package snow

import (
	"sync/atomic"
	"time"
)

type GoPool struct {
	workerNum int32  // 已启动的worker数
	activeNum int32  // 正在执行任务的worker数
	ch        *Channel      // 任务主线队列
	pq        *PriorityQueue
	peekTime  int64
	t         *time.Timer
}

//var ErrClosed = errors.New("Closed")

func NewGoPool(size uint32) *GoPool {
	p := &GoPool{
		pq: NewPriorityQueue(),
	}
	p.SetWorkerNum(size)
	return p
}

// 动态设置worker数,size必须大于0，否则会发回错误
func (p *GoPool) SetWorkerNum(size uint32) {
	switch {
	case int32(size) > p.workerNum:
		if p.ch == nil {
			p.SetChannel(&Channel{})
		}

		addNum := int32(size) - p.workerNum
		for i:=int32(0); i<addNum; i++ {
			w := worker{
				pool: p,
			}
			<-w.start()
		}

	case int32(size) < p.workerNum:
		decNum := p.workerNum - int32(size)
		for i:=int32(0);i<decNum;i++ {
			p.ch.Send(true)
		}
	}
}

// 关闭
func (p *GoPool) Close() {
	p.ch.Close()
	p.SetWorkerNum(0)
}

// 设置任务主线，达到多个worker组共享一个任务队列的目的
func (p *GoPool) SetChannel(ch *Channel) {
	p.ch = ch
}

// 线程池状态
// @workerNum 已启动的goroutine数
// @activeNum 正在工作的goroutine数
func (p *GoPool) Status() (workerNum int32, activeNum int32) {
	return p.workerNum, p.activeNum
}

// 线程池内执行异步任务
// @done 执行成功返回true
// @ok 成功插入待执行队列
func (p *GoPool) Go(fn func()) (done chan bool, err error) {
	done = make(chan bool, 2)
	ok := p.ch.Send(func() {
		defer func() {
			done <- true
		}()
		fn()
	})
	if !ok {
		return nil, ErrClosed
	}
	return done, nil
}

func (p *GoPool) AfterFunc(d time.Duration, fn func()) (done chan bool, err error) {
	tn := time.Now().Add(d).UnixNano()
	if p.t == nil {
		p.t = time.NewTimer(0)
		<-p.t.C
		p.peekTime = 1<<63 - 1
		_, err := p.Go(func() {
			for {
				_, ok := <-p.t.C
				if !ok {
					break
				}

				f,_,ok := p.pq.Pop()
				if !ok {
					p.t = nil
					p.peekTime = 0
					break
				}
				_,_ = p.Go(f.(func()))

				_,pt,ok := p.pq.Peek()
				if ok {
					p.t.Reset(time.Duration((-pt)-time.Now().UnixNano()))
					p.peekTime = -pt
				}else {
					p.peekTime = 1<<63 - 1
				}
			}
		})
		if err != nil {
			return nil, err
		}
	}

	p.pq.Push(fn, -tn)

	if tn < p.peekTime {
		p.t.Reset(time.Duration(tn-time.Now().UnixNano()))
	}

	return done, nil
}

type worker struct {
	pool       *GoPool
}

func (w *worker) start() <-chan bool {
	done := make(chan bool, 1)
	go func() {
		atomic.AddInt32(&w.pool.workerNum, 1)
		done <- true

		defer func() {
			atomic.AddInt32(&w.pool.workerNum, -1)
			w.pool = nil
		}()

		for {
			if w.pool.ch == nil {
				break
			}
			fn,ok := w.pool.ch.Receive()
			if !ok {
				break
			}

			if _,ok = fn.(bool);ok {
				break
			}

			func(){
				atomic.AddInt32(&w.pool.activeNum, 1)
				defer func() {
					atomic.AddInt32(&w.pool.activeNum, -1)
					checkPanic()
				}()

				fn.(func())()
			}()
		}
	}()

	return done
}
