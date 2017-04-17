package main

import (
	"fmt"
	"sync"
)

var buf = make(chan int, 5)
var wg = sync.WaitGroup{}

func producer(v int) {
	buf <- v
	fmt.Printf("producer %d product %d\n", v, v)
	wg.Done()
}

func cosumer(v int) {
	fmt.Printf("cosumer %d cosume %d\n", v, <-buf)
	wg.Done()
}

func main() {
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go cosumer(i)
	}
	for i := 0; i < 10; i++ {
		go producer(i)
	}
	wg.Wait()
}
