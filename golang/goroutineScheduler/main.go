package main

import (
	"fmt"
	"runtime"
	"sync"
)

func runForever(id int) {
	fmt.Printf("id: %d\n", id)
	for {
		runtime.Gosched()
	}
}
func main() {
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go runForever(i)
	}
	wg.Wait()
}
