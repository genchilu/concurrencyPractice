package main

import "fmt"

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for _, v := range values {
			//fmt.Println(v)
			select {
			case <-done:
				return
			case valueStream <- v:
			}
		}
	}()

	return valueStream
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
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

func tee(done <-chan interface{}, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)
		fmt.Println("hio")
		for val := range in {
			fmt.Println(val)
			var localOut1, localOut2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case localOut1 <- val:
					localOut1 = nil
				case localOut2 <- val:
					localOut2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	done := make(chan interface{})
	defer close(done)

	repeatStream := repeat(done, 1, 2, 3, 4, 5)
	takeStream := take(done, repeatStream, 4)

	out1, out2 := tee(done, takeStream)

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}

}
