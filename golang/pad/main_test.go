package main

import (
	"runtime"
	"sync"
	"testing"
)

func testAtomicIncrease(myatomic MyAtomic) {
	runtime.GOMAXPROCS(4)
	paraNum := 10000
	addTimes := 10000
	var wg sync.WaitGroup
	wg.Add(paraNum * 2)
	for i := 0; i < paraNum; i++ {
		go func() {
			for j := 0; j < addTimes; j++ {
				myatomic.IncreaseA()
			}
			wg.Done()
		}()
	}
	for i := 0; i < paraNum; i++ {
		go func() {
			for j := 0; j < addTimes; j++ {
				myatomic.IncreaseB()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkNoPad(b *testing.B) {
	myatomic := &NoPad{}
	b.ResetTimer()
	testAtomicIncrease(myatomic)
}

func BenchmarkPad(b *testing.B) {
	myatomic := &Pad{}
	b.ResetTimer()
	testAtomicIncrease(myatomic)
}
