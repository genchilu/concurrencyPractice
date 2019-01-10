package main

import(
	"fmt"
	"time"
)

func doWork(done <-chan interface{}, strings <-chan string) <-chan interface{} {
	terminated := make(chan interface{})
	
	go func() {
		defer fmt.Println("Dowork finished.")
		defer close(terminated)

		for{
			select {
			case s := <-strings:
				fmt.Println(s)
			case <-done:
				return
			}
		}
	}()

	return terminated
}

func main(){
	done:=make(chan interface{})
	terminated:=doWork(done, nil)

	go func(){
		time.Sleep(2*time.Second)
		fmt.Println("Cancel task")
		close(done)
	}()

	<-terminated
	fmt.Println("Finish main")
}