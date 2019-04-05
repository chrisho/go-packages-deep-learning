package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

type Student struct {
	Name string `json:"string"`
	Number int `json:"number"`
	Sex bool `json:"sex"` // true:男 false:女
	Age int `json:"age"`
}

func testA(i int) {
	fmt.Println("A", i)
}
func testB(i int) {
	fmt.Println("B", i)
}
func testC(i int) {
	fmt.Println("C", i)
}
func testD(i int) {
	fmt.Println("D", i)
}
func testE(i int) {
	fmt.Println("E", i)
}
func testF(i int) {
	fmt.Println("F", i)
}
func testG(i int) {
	fmt.Println("G", i)
}
func testH(i int) {
	fmt.Println("H", i)
}
func testI(i int) {
	fmt.Println("I", i)
}
func testJ(i int) {
	fmt.Println("J", i)
}
func testK(i int) {
	fmt.Println("K", i)
}

type test func(i int)


var smap sync.Map

func loop(i int) test {
	//if fi, ok := sync.Map{}

	index := i % 10
	key := fmt.Sprintf("test %d", index)
	var (
		wg sync.WaitGroup
		t test
	)
	wg.Add(1)
	fi, loaded := smap.LoadOrStore(key, test(func(i int) {
		wg.Wait()
		fmt.Println(uintptr(unsafe.Pointer(&wg)), i)
		t(i)
	}))
	if loaded {
		fmt.Println("loaded", i)
		return fi.(test)
	}
	fmt.Println("!loaded", i)

	if index == 0 {
		t = testA
	} else if index == 1 {
		t = testB
	} else if index == 2 {
		t = testC
	} else if index == 3 {
		t = testD
	} else if index == 4 {
		t = testF
	} else if index == 5 {
		t = testG
	} else if index == 6 {
		t = testH
	} else if index == 7 {
		t = testI
	} else if index == 8 {
		t = testJ
	} else if index == 9 {
		t = testK
	}
	wg.Done()
	smap.Store(key, t)
	return t
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i:=0; i<20; i++ {
		go func(i int){
			loop(i)(i)
		}(i)
	}
	time.Sleep(time.Second * 3)
	// 编译json
	//s, err := json.Marshal(stu) // func Marshal(v interface{}) ([]byte, error)
	//if err != nil {
	//	panic(err)
	//}
	//// 编译json
	//_, err = json.Marshal(stu) // func Marshal(v interface{}) ([]byte, error)
	//if err != nil {
	//	panic(err)
	//}
	//v := reflect.ValueOf(stu)
	//log.Println(v, v.Type(), math.Abs(math.Inf(0)))
	//log.Println("字节:", s)
	//log.Println("字符串:", string(s))

	//// 解析json
	//var stu1 Student // 映射到Student
	//if err := json.Unmarshal(s, &stu1); err != nil {
	//	panic(err)
	//}
	//
	//var stu2 interface{}
	//if err := json.Unmarshal(s, &stu2); err != nil {
	//	panic(err)
	//}
	//
	//log.Println("Student:", stu1)
	//log.Println("interface{}:", stu2) // 输出为Map[string]interface{}
	//
	//log.Println(reflect.TypeOf(stu2))
}
