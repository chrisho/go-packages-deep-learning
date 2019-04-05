package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type encoderFunc func(i int)
var encoderCache sync.Map

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	TestMap("string")(1)
	for i:=0; i<100; i++ {
		go TestMap("integer")(2)
	}
	time.Sleep(time.Second * 1)
	fmt.Println(t)
}

func TestMap(key string) encoderFunc {
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(key, encoderFunc(func(i int) {
		wg.Wait()
		f(i)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newTypeEncoder(key)
	wg.Done()
	encoderCache.Store(key, f)
	return f
}

var t int = 0
func newTypeEncoder(k string) encoderFunc {
	t++
	if k == "string" {
		return DoString
	}
	return DoInteger
}

func DoString(i int) {
	log.Println("String:", i)
}

func DoInteger(i int) {
	log.Println("Integer:", i)
}

