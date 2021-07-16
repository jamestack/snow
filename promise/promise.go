package promise

import "sync"

type Promise interface {
	Then(fn Fn) Promise
	Catch(fn ErrFn) <-chan bool
}

type Any interface{}
type Fn func(resolve Resolve, reject Reject, args ...Any) // func RPC_Login(a int, b int)
type ErrFn func(errs ...error)

type promise struct {
	lock     sync.Mutex
	then     []Fn
	catch    ErrFn
	resolves *[]Any
	rejects  *[]error
	done     chan bool
}

func (p *promise) Then(fn Fn) Promise {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.rejects != nil {
		return p
	}

	p.then = append(p.then, fn)

	if p.resolves != nil {
		p.run(*p.resolves...)
	}

	return p
}

func (p *promise) Catch(catch ErrFn) <-chan bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.catch != nil {
		return p.done
	}
	p.catch = catch
	p.done = make(chan bool, 1)

	if p.rejects != nil {
		catch(*p.rejects...)
		p.rejects = nil
		p.done <- true
	}

	if p.resolves != nil {
		p.done <- true
	}

	return p.done
}

func (p *promise) run(args ...Any) {
	var fn Fn
	if len(p.then) > 0 {
		fn = p.then[0]
		p.then = p.then[1:]
	} else {
		return
	}

	p.resolves = nil

	hasDone := false
	hasReject := false
	go fn(func(resolve ...Any) {
		p.lock.Lock()
		defer p.lock.Unlock()

		if hasReject {
			return
		}
		if hasDone {
			return
		}
		hasDone = true

		if len(p.then) > 0 {
			p.run(resolve...)
		} else {
			p.resolves = &resolve
			if p.catch != nil {
				p.done <- true
			}
		}
	}, func(err ...error) {
		p.lock.Lock()
		defer p.lock.Unlock()

		if hasReject {
			return
		}
		hasReject = true

		if p.catch != nil {
			p.catch(err...)
			p.done <- true
		} else {
			p.rejects = &err
		}
	}, args...)
}

// resolve() reject()
type Exit interface{}

type Resolve func(resolve ...Any)

type Reject func(err ...error)

func New(fn Fn, args ...Any) Promise {
	if fn == nil {
		panic("Promise Func is nil")
	}
	p := &promise{
		then: []Fn{fn},
	}
	p.run(args...)
	return p
}
