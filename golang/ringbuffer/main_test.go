package main

import (
	"os"
	"runtime/pprof"
	"sync"
	"testing"

	"./ringbuf"
	"./ringbufpadded"
)

var times = 1000000

func BenchmarkChan(b *testing.B) {
	f, _ := os.OpenFile("./profile/chancpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	msg := make(chan interface{}, times)
	wg.Add(times)
	b.ResetTimer()
	for i := 0; i < times; i++ {
		go func() {
			msg <- i
		}()
	}
	for i := 0; i < times; i++ {
		go func() {
			<-msg
			wg.Done()
		}()
	}
	wg.Wait()
	pprof.StopCPUProfile()
	f.Close()
}

func BenchmarkRingBuf(b *testing.B) {
	f, _ := os.OpenFile("./profile/rbcpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	rb := ringbuf.NewRingBuffer(uint64(times))
	wg.Add(times)
	b.ResetTimer()
	for i := 0; i < times; i++ {
		go func() {
			rb.Put(i)
		}()
	}
	for i := 0; i < times; i++ {
		go func() {
			rb.Get()
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
	wg.Add(times)
	b.ResetTimer()
	for i := 0; i < times; i++ {
		go func() {
			rb.Put(i)
		}()
	}
	for i := 0; i < times; i++ {
		go func() {
			rb.Get()
			wg.Done()
		}()
	}
	wg.Wait()
	pprof.StopCPUProfile()
	f.Close()
}
