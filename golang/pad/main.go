package main

import (
	"fmt"
	"sync/atomic"
)

type MyAtomic interface {
	IncreaseAllEles()
}
type NoPad struct {
	a uint64
	b uint64
	c uint64
}

func (myatomic *NoPad) IncreaseAllEles() {
	atomic.AddUint64(&myatomic.a, 1)
	atomic.AddUint64(&myatomic.b, 1)
	atomic.AddUint64(&myatomic.c, 1)
}

type Pad struct {
	a   uint64
	_p1 [8]uint64
	b   uint64
	_p2 [8]uint64
	c   uint64
	_p3 [8]uint64
}

func (myatomic *Pad) IncreaseAllEles() {
	atomic.AddUint64(&myatomic.a, 1)
	atomic.AddUint64(&myatomic.b, 1)
	atomic.AddUint64(&myatomic.c, 1)
}
func main() {
	myatomic := &Pad{}
	myatomic.IncreaseAllEles()
	fmt.Printf("%d\n", myatomic.a)
}
