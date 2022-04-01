package main

import "fmt"

type student struct {
	Name string
}

func (s *student) test() string {
	return s.Name + " test"
}

func test1(s student) string {
	return s.Name
}

func main() {
	var stud student
	stud.Name = "Chandrakant"
	newString := stud.test()
	fmt.Println(newString)
	finalString := test1(stud)
	fmt.Println(finalString)
}
