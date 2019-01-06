package main

import (
	"time"
	"fmt"
)

func or(channels ... <-chan interface{}) <-chan interface{} {
	switch len( channels) {
		case 0: 
			return nil
		case 1: 
			return channels[0]
	} 

	orDone := make(chan interface{})
	go func() { 
		defer close(orDone) 
		select { 
			case <-channels[0]: 
			case <-channels[1]: 
			case <-or(append(channels[2:], orDone)...):
		}
	}()
	return orDone
}

func sig (after time.Duration) <-chan interface{}{
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main(){
	start := time.Now()
	
	<-or(
		sig(2* time.Second),
		sig(1* time.Second),
		sig(3* time.Second),
	)

	fmt.Printf("done after %v", time.Since(start))
}

