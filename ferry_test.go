package ferry_test

import (
	"sync/atomic"
	"sync"
	"testing"
	"time"

	. "github.com/dianelooney/ferry"
)

func TestItBlocks(t *testing.T) {
	var order []int
	b := Ferry{}
	go func() {
		time.Sleep(250 * time.Millisecond)
		order = append(order, 1)
		b.Done()
	}()
	b.Wait()
	order = append(order, 2)
	if order[0] != 1 || order[1] != 2 {
		t.Error("It did not block in the proper order")
	}
}

func TestItBlocksMany(t *testing.T) {
	const n = 1000
	done := 0
	mtx := sync.Mutex{}

	b := Ferry{}
	for i := 0; i < n; i++ {
		go func(i int) {
			b.Wait()

			mtx.Lock()
			defer mtx.Unlock()
			done++
		}(i)
	}
	time.Sleep(250 * time.Millisecond)
	if done != 0 {
		t.Errorf("It did not block all goroutines")
	}
	b.Done()
	time.Sleep(250 * time.Millisecond)

	mtx.Lock()
	defer mtx.Unlock()

	if done != n {
		t.Errorf("It did not unblock all goroutines")
	}
}

func TestItBlocksWhileUnblocking(t *testing.T) {
	const n = 1000
	fail := false
	b := Ferry{}
	for i := 0; i < n; i++ {
		go func() {
			for {
				b.Wait()
				b.Wait()
				fail = true
			}
		}()
	}

	done := false
	mtx := sync.Mutex{}
	go func() {
		b.Done()
		mtx.Lock()
		defer mtx.Unlock()
		done = true
	}()
	time.Sleep(250 * time.Millisecond)
	
	mtx.Lock()
	defer mtx.Unlock()

	if !done {
		t.Errorf("It did not return")
	}

	if fail {
		t.Errorf("It unlocked the same goroutine more than once")
	}
}

func TestMultipleUnlockCalls(t *testing.T) {
	const n = 1000
	b := Ferry{}
	var count int32
	for i := 0; i < n; i++ {
		go func() {
				b.Wait()
				b.Wait()
				atomic.AddInt32(&count, 1)
		}()
	}

	time.Sleep(250 * time.Millisecond)
	b.Done()
	time.Sleep(250 * time.Millisecond)
	b.Done()
	time.Sleep(250 * time.Millisecond)
	
	if n != atomic.LoadInt32(&count) {
		t.Errorf("It did not unblock all goroutines twice, expected '%v' got '%v'", n, count)
	}
}