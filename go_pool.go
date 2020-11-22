package snow

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type GoPool struct {
	lock sync.Mutex
	workerNum int32  // 已启动的worker数
	activeNum int32  // 正在执行任务的worker数
	ch *Channel      // 任务主线队列
	afterFuncPool *GoPool  // AfterFunc异步任务队列
	afterDestroy chan bool
}

var ErrNotEnoughWorker = errors.New("[GoPool: NotEnoughWorker]")

// 动态设置worker数
func (p *GoPool) SetMaxWorker(size uint32) {
	defer checkPanic()

	switch {
	case size == 0:
		if p.ch != nil {
			p.destroy(true)
		}
	case int32(size) > p.workerNum:
		if p.ch == nil {
			p.ch = &Channel{}
		}

		addNum := int32(size) - p.workerNum
		for i:=int32(0);i<addNum;i++ {
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

// 销毁线程池
// 1. p.Go() 销毁后会全部执行
// 2. p.AfterFunc() 在pool销毁后未执行的将被全部丢弃
func (p *GoPool) destroy(waitAllJobDone ...bool) {
	if p.afterFuncPool != nil {
		p.afterFuncPool.destroy(false)
	}

	if p.afterDestroy != nil {
		p.afterDestroy <- true
	}

	p.ch.Close()
	if len(waitAllJobDone) == 0 || waitAllJobDone[0] == true {
		for {
			worker,active := p.Status()
			if worker == 0 && active == 0 {
				break
			}
			<-time.After(1)
		}
	}
	p.ch = nil
}

func (p *GoPool) Go(fn func()) <-chan bool {
	defer checkPanic()

	if p.workerNum == 0 {
		panic(ErrNotEnoughWorker)
		return nil
	}

	done := make(chan bool, 1)
	p.ch.Send(func() {
		defer func() {
			done <- true
		}()
		fn()
	})
	return done
}

func (p *GoPool) AfterFunc(d time.Duration, fn func()) <-chan bool {
	defer checkPanic()

	if p.workerNum == 0 {
		panic(ErrNotEnoughWorker)
		return nil
	}

	// 初始化
	if p.afterFuncPool == nil {
		p.afterFuncPool = &GoPool{}
		p.afterFuncPool.SetMaxWorker(1)
		p.afterFuncPool.afterDestroy = make(chan bool, 1)
	}

	endTime := time.Now().Add(d)

	return p.afterFuncPool.Go(func() {
		select {
		case <-p.afterFuncPool.afterDestroy:
		case <-time.After(endTime.Sub(time.Now())):
			p.Go(fn)
		}
	})
}

type worker struct {
	pool       *GoPool
}

func (w *worker) start() <-chan bool {
	done := make(chan bool, 1)
	go func() {
		atomic.AddInt32(&w.pool.workerNum, 1)
		done <- true

		defer w.stop()

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

func (w *worker) stop() {
	atomic.AddInt32(&w.pool.workerNum, -1)
	w.pool = nil
}
