package ferry_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/dianelooney/ferry"
)

func TestValue_ItBlocks(t *testing.T) {
	b := NewValue()
	go func() {
		b.Wait()
		t.Error("It did not block")
	}()
	time.Sleep(10 * time.Millisecond)
}

func TestValue_ItBlocksMany(t *testing.T) {
	const n = 1000
	b := NewValue()
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()
			wg.Done()
		}()
	}
	time.Sleep(10 * time.Millisecond)
	b.Done(nil)
	wg.Wait()
}

func TestValue_ItUnblocksGoroutinesOnce(t *testing.T) {
	const n = 1000
	b := NewValue()
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()
			b.Wait()
			t.Errorf("It unlocked the same goroutine more than once")
		}()
	}
	b.Done(nil)
	time.Sleep(10 * time.Millisecond)
}

func TestValue_ItSendsTheSameValue(t *testing.T) {
    const n = 1000
    f := NewValue()
    ch := make(chan time.Time)
    for i := 0; i < n; i++ {
        go func() {
            ch <- f.Wait().(time.Time)
        }()
    }
		time.Sleep(10 * time.Millisecond)
		d := time.Now()
		f.Done(d)
    for i := 0; i < n; i++ {
        if w := <- ch; d != w {
            t.Errorf("Mismatched return times: '%v' from Done and '%v' from Wait", d, w)
        }
    }
}

func BenchmarkValue_Done(b *testing.B) {
	f := NewValue()
	for i := 0; i < b.N; i++ {
		f.Done(nil)
	}
}

func benchmarkValue_Done_NWaiters(b *testing.B, n int) {
	arr := make([]Value, b.N)
	for i := 0; i < b.N; i++ {
		f := NewValue()
		arr[i] = f
		for j := 0; j < n; j++ {
			go f.Wait()
		}
	}
	for i := 0; i < b.N; i++ {
		arr[i].Done(nil)
	}
}
func BenchmarkValue_Done_0Waiters(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 0)
}
func BenchmarkValue_Done_1Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 1)
}
func BenchmarkValue_Done_2Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 2)
}
func BenchmarkValue_Done_5Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 5)
}
func BenchmarkValue_Done_10Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 10)
}
func BenchmarkValue_Done_20Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 20)
}
func BenchmarkValue_Done_50Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 50)
}
func BenchmarkValue_Done_100Waiter(b *testing.B) {
	benchmarkValue_Done_NWaiters(b, 100)
}
