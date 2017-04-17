package main

import (
	"fmt"
	. "sync"
)

func main() {
	var pool = Pool{New: func() interface{} { return "pool empty, this is new one" }}
	tmp := "hi, genchi"
	pool.Put(tmp)
	tmpout := pool.Get()
	pool.Put(tmpout)
	fmt.Println(tmpout)
	fmt.Println(pool.Get())
	fmt.Println(pool.Get())
}
