package encode

import "fmt"

type A struct {
	i int
}

func (s *A) Print1() {
	s.i += 1
	fmt.Println("*A", s.i)
}

func (s A) Print2() {
	s.i += 1
	fmt.Println("A", s.i)
}


