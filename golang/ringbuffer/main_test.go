package main

import (
	"runtime"
	"sync"
	"testing"

	"./ringbuf"
	"./ringbufpadded"
)

var times = 100000
var paral = 100
var maxProcess = 4

func BenchmarkChan(b *testing.B) {
	runtime.GOMAXPROCS(maxProcess)
	var wg sync.WaitGroup
	msg := make(chan interface{}, times/10)
	wg.Add(paral)
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
	runtime.GOMAXPROCS(maxProcess)
	var wg sync.WaitGroup
	rb := ringbuf.NewRingBuffer(uint64(times))
	wg.Add(paral)
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
}

func BenchmarkRingBufPadded(b *testing.B) {
	runtime.GOMAXPROCS(maxProcess)
	var wg sync.WaitGroup
	rb := ringbufpadded.NewRingBufferPadded(uint64(times))
	wg.Add(paral)
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
}
