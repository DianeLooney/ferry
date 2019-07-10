package ferry

import "sync"

func NewValue() *Value {
	return &Value{
		ch: make(chan interface{}),
	}
}

// Time is just a Ferry that passes through a value listener
//
// It is less performant than a Ferry if no value is required to be sent
type Value struct {
    ch  chan interface{}
    i   int
    mtx sync.Mutex
}

func (f *Value) Wait() interface{} {
    f.mtx.Lock()
    f.i++
    ch := f.ch
    f.mtx.Unlock()

    v, ok := <-ch
    if !ok {
        panic("should not occur")
    }
    return v
}

func (f *Value) Done(v interface{}) {
    f.mtx.Lock()
    listenerCount := f.i
    f.i = 0

    ch := f.ch
    f.ch = make(chan interface{})
    f.mtx.Unlock()

    for sent := 0; sent < listenerCount; sent++ {
        ch <- v
		}
		close(ch)
}
