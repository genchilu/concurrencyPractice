package main

import(
	"fmt"
)

func generater(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func(){
		defer close(intStream)

		for _, i:=range integers{
			select{
			case <- done:
				return
			case intStream <- i:
			}
		}
	}()

	return intStream;
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int{
	multipliedStream :=make(chan int)

	go func(){
		defer close(multipliedStream)

		for i:=range intStream {
			select{
			case <-done:
				return
			case multipliedStream<-i*multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int{
	addedStream := make(chan int)
	go func(){
		defer close(addedStream)
		for i:=range intStream{
			select{
			case <-done:
				return
			case addedStream<-i+additive:
			}
		}
	}()
	return addedStream
}

func main(){
	done:=make(chan interface{})
	defer close(done)

	intStream := generater(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v:=range pipeline{
		fmt.Println(v)
	}
}