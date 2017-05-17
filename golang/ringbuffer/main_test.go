package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"testing"

	"./ringbuf"
	"./ringbufpadded"
)

var times = 100000
var paral = 100
var sumTimes int64 = 10000000

func BenchmarkLock(b *testing.B) {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var sum int64
	wg.Add(int(sumTimes))
	b.ResetTimer()
	for i := int64(1); i <= sumTimes; i++ {
		go func(i int64) {
			mutex.Lock()
			sum += i
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkCAS(b *testing.B) {
	var sum int64
	var wg sync.WaitGroup
	wg.Add(int(sumTimes))
	for i := int64(1); i <= sumTimes; i++ {
		go func(i int64) {
			for {
				if atomic.CompareAndSwapInt64(&sum, sum, sum+i) {
					break
				}
				if i == 10000 {
					runtime.Gosched() // free up the cpu before the next iteration
					i = 0
				} else {
					i++
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func BenchmarkChan(b *testing.B) {
	var wg sync.WaitGroup
	msg := make(chan interface{}, times/10)
	wg.Add(paral)
	b.ResetTimer()
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				msg <- i
			}
		}()
	}
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				<-msg
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRingBuf(b *testing.B) {
	f, _ := os.OpenFile("./profile/rbcpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	rb := ringbuf.NewRingBuffer(uint64(times))
	wg.Add(paral)
	b.ResetTimer()
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				rb.Put(i)
			}
		}()
	}
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				rb.Get()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pprof.StopCPUProfile()
	f.Close()
}

func BenchmarkRingBufPadded(b *testing.B) {
	f, _ := os.OpenFile("./profile/rbpadcpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	rb := ringbufpadded.NewRingBufferPadded(uint64(times))
	wg.Add(paral)
	b.ResetTimer()
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				rb.Put(i)
			}
		}()
	}
	for i := 0; i < paral; i++ {
		go func() {
			for j := 0; j < times; j++ {
				rb.Get()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pprof.StopCPUProfile()
	f.Close()
}
