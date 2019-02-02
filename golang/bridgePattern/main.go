package main

import "fmt"

func bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range stream {
				select {
				case valStream <- val:
				case <-done:
				}
			}

		}
	}()
	return valStream
}

func genVals() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)

		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
