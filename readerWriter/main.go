package main

import (
	"fmt"
	. "sync"
	. "time"
)

func read(id int, rwmutex *RWMutex, wg *WaitGroup, file *string) {
	fmt.Printf("id: %d in read\n", id)
	fmt.Printf("id: %d rwmutex status berofe get rlock: %v\n", id, rwmutex)
	rwmutex.RLock()
	fmt.Printf("id: %d rwmutex status after get rlock: %v\n", id, rwmutex)
	fmt.Printf("id: %d start read\n", id)
	Sleep(10 * Second)
	fmt.Printf("id: %d read: %s\n", id, *file)
	rwmutex.RUnlock()
	wg.Done()
}

func write(id int, rwmutex *RWMutex, wg *WaitGroup, file *string) {
	fmt.Printf("id: %d in write\n", id)
	fmt.Printf("id: %d rwmutex status berofe get lock: %v\n", id, rwmutex)
	rwmutex.Lock()
	fmt.Printf("id: %d rwmutex status after get lock: %v\n", id, rwmutex)
	fmt.Printf("id: %d start write\n", id)
	Sleep(3 * Second)
	fmt.Printf("id: %d writing...\n", id)
	*file += "1"
	rwmutex.Unlock()
	wg.Done()
}

func main() {
	var wg = &WaitGroup{}
	var rwmutex = &RWMutex{}
	file := "1"
	wg.Add(6)
	go read(1, rwmutex, wg, &file)
	Sleep(1 * Second)
	go read(2, rwmutex, wg, &file)
	Sleep(1 * Second)
	go write(3, rwmutex, wg, &file)
	Sleep(1 * Second)
	go read(4, rwmutex, wg, &file)
	go read(5, rwmutex, wg, &file)
	go read(6, rwmutex, wg, &file)
	wg.Wait()
}
