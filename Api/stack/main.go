package main

import "fmt"

type Stack struct {
	data []int
}

func (s *Stack) push(n int) {
	s.data = append(s.data, n)
	fmt.Println("Elemented added to stack", n)
}

func (s *Stack) pop() int {
	if len(s.data) > 0 {
		n := s.data[len(s.data)-1]
		s.data = s.data[0 : len(s.data)-1]
		return n
	} else {
		fmt.Println("Stack is empty")
	}
	return -1
}

func main() {
	s := &Stack{}
	s.push(1)
	s.push(2)
	s.pop()
	s.push(3)
	s.push(4)
	s.pop()
	fmt.Println("final element int stack..", s.data)
}
