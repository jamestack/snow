package promise

type Promise interface {
	Then(fn Fn) Promise
	Catch(fn ErrFn) <-chan bool
}

type Any interface {}
type Fn func(resolve Resolve, reject Reject, args ...Any) Exit // func RPC_Login(a int, b int)
type ErrFn func(errs ...error)

type promise struct {
	fn       Fn
	catch    ErrFn
	next     *promise
	resolves   *[]Any
	done chan bool
}

func (p *promise) Then(fn Fn) Promise {
	if fn == nil {
		panic("Promise Func is nil")
	}
	newPromise := &promise{
		fn:    fn,
		catch: nil,
		next:  nil,
	}
	p.next = newPromise

	if p.resolves != nil {
		newPromise.run(*p.resolves...)
		p.resolves = nil
	}
	return newPromise
}

func (p *promise) Catch(fn ErrFn) <-chan bool{
	newPromise := &promise{
		fn:    nil,
		catch: fn,
		next:  nil,
		done: make(chan bool, 1),
	}
	p.next = newPromise
	return newPromise.done
}

func (p *promise) resolve(resolve ...Any) Exit {
	if p.fn == nil {
		return struct {}{}
	}

	if p.next != nil {
		if p.next.fn != nil {
			p.next.run(resolve...)
		}
		if p.next.done != nil {
			p.next.done <- true
		}
	} else {
		p.resolves = &resolve
	}

	p.fn = nil
	p.catch = nil
	return struct {}{}
}

func (p *promise) reject(err ...error) Exit {
	p.fn = nil
	p.catch = nil
	for n := p.next; n != nil; {
		next := n.next
		catch := n.catch
		n.fn = nil
		n.catch = nil
		n.next = nil

		if catch != nil {
			catch(err...)
			n.done <- true
		}
		n = next
	}
	p.next = nil
	return struct {}{}
}

func (p *promise) run(args ...Any) {
	go func() {
		if p.fn != nil {
			p.fn(p.resolve, p.reject, args...)
		}
	}()

}

type Exit interface {}

type Resolve func(resolve ...Any) Exit

type Reject func(err ...error) Exit

func New(fn Fn, args ...Any) Promise {
	if fn == nil {
		panic("Promise Func is nil")
	}
	p := &promise{
		fn: fn,
	}
	p.run(args...)
	return p
}
