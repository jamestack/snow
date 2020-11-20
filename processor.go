package snow

import "time"

// 带缓存的任务队列
type Processor struct {
	cap int
	isStart bool
	ch chan func()
}

// 初始化一个带缓存的任务处理器
// @cap 任务处理器缓存容量。
// 也就是说cap=1时为同步任务阻塞Queue.Run()，cap>1则会变成异步任务直到缓存数被耗尽，建议不大于5
func NewProcessor(cap int) *Processor {
	processor :=  &Processor{
		cap: cap,
	}
	processor.start()
	return processor
}

func (p *Processor) start() {
	if p.isStart {
		return
	}
	p.isStart = true
	p.ch = make(chan func(), p.cap)
	go func() {
		for {
			fn,ok := <- p.ch
			if !ok {
				break
			}
			fn()
		}
	}()
}

// 使用完毕后必须执行销毁函数，不然会造成goroutine泄露
func (p *Processor) Destroy() {
	if !p.isStart {
		return
	}

	p.isStart = false
	close(p.ch)
}

// 实时任务
func (p *Processor) Run(fn func()){
	p.ch <- func() {
		// 防止panic导致影响后续任务处理
		defer checkPanic()
		fn()
	}
	return
}

// 延时任务
func (p *Processor) RunAfter(t time.Duration, fn func()) *time.Timer {
	return time.AfterFunc(t, func() {
		p.Run(fn)
	})
}
