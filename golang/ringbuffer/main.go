package main

import (
	"fmt"

	"./ringbuf"
)

func main() {
	rb := ringbuf.NewRingBuffer(4)
	err := rb.Put(1)
	fmt.Print(err)
}
