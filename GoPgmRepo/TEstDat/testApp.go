package main

import "fmt"

func main() {
	a := []int{1, 2}
	b := []int{11, 22}
	a = append(a, b) // a == [1 2 11 22]
	for _, val := range a {
		fmt.Println(val)
	}
}
