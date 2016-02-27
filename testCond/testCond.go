package main

import (
	"fmt"
	. "sync"
)

var ch = make(chan bool)

func testBroadcast(id int, cond *Cond, m *Mutex) {
	m.Lock()
	fmt.Printf("id: %d broadcast confd\n", id)
	cond.Broadcast()
	fmt.Printf("id: %d broadcast confd finish\n", id)
	m.Unlock()
}

func testSignal(id int, cond *Cond, m *Mutex, wg *WaitGroup) {
	m.Lock()
	fmt.Printf("id: %d signal confd\n", id)
	cond.Signal()
	fmt.Printf("id %d signal confd finish\n", id)
	m.Unlock()
	wg.Done()
}

func testWait(id int, cond *Cond, m *Mutex, wg *WaitGroup) {
	m.Lock()
	fmt.Printf("id: %d wait cond\n", id)
	ch <- true
	cond.Wait()
	fmt.Printf("id: %d wait cond finish\n", id)
	m.Unlock()
	wg.Done()
}

func main() {
	var wg WaitGroup
	var m Mutex
	cond := NewCond(&m)
	n := 3
	//test signal
	for i := 0; i < n; i++ {
		wg.Add(1)
		go testWait(i, cond, &m, &wg)
	}
	for i := 0; i < n; i++ {
		<-ch
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go testSignal(i, cond, &m, &wg)
	}
	wg.Wait()
	//test confd
	for i := 0; i < n; i++ {
		wg.Add(1)
		go testWait(i, cond, &m, &wg)
	}
	for i := 0; i < n; i++ {
		<-ch
	}
	testBroadcast(1, cond, &m)
}
