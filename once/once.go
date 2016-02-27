package main

import (
	"fmt"
	. "sync"
)

func testOnce() {
	fmt.Printf("print once\n")
}

func main() {
	var once Once
	onceBody := testOnce
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		tmp := <-done
		fmt.Printf("%v\n", tmp)
	}
}
