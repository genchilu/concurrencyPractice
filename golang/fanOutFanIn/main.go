package main

import(
	"sync"
	"time"
	"runtime"
	"fmt"
	"math/big"
	"math/rand"
)

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{}{
	takeStream := make(chan interface{})
	go func(){
		defer close(takeStream)
		for i:=0;i<num;i++{
			select{
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func primeFinder(done <- chan interface{}, randIntStream <-chan int) <-chan int {
	primeStream := make(chan int)
	go func(){
		defer close(primeStream)
		for {
			select {
			case <-done:
				return
			case number:= <-randIntStream:
				if big.NewInt(int64(number)).ProbablyPrime(0) {
					primeStream<-number
				}
			}
		}
	}()

	return primeStream
}

func repeatFn(done <-chan interface{},fn func() interface{}) <-chan interface{}{
	valueStream := make(chan interface{})
	go func(){
		defer close(valueStream)
		for{
			select{
			case <-done:
				return
			case valueStream<-fn():
			}
		}
	}()

	return valueStream
}

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func(){
		defer close(intStream)

		for v:=range valueStream{
			select{
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()
	return intStream
}

func fanIn(done <-chan interface{}, channels ...<-chan int) <-chan interface{}{
	multiplexedStream := make(chan interface{})

	var wg sync.WaitGroup

	multiplex := func(c <-chan int){
		defer wg.Done()
		for i:=range c {
			select{
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplex(c)
	}

	go func(){
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main(){
	done:=make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func () interface{} {
		return rand.Intn(50000000)
	}

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()

	fmt.Printf("Spinning up %d prime finders.\n" , numFinders)

	finders := make([]<-chan int, numFinders)
	
	fmt.Println("Primes:")
	for i:=0;i<numFinders;i++{
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime:=range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime) 
	}
	
	fmt.Printf( "Search took: %v", time.Since(start)) 
}