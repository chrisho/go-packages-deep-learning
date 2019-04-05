package main

import (
	"fmt"
	"github.com/chrisho/go-packages-deep-learning/encode"
)

func main() {
	var a encode.A
	a.Print1()
	a.Print2()

	fmt.Printf("%p", a)
}

