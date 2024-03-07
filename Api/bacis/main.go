package main

import "fmt"

func main() {
	var a []int
	fmt.Println(a)
	aa := make([]int, 0)
	fmt.Println(aa)
	var bb [5]int
	fmt.Println(bb)
	var t int
	fmt.Println(t)
	var x = [5]int{}
	fmt.Println(x)

	arr := [5]int{1, 2, 3, 4, 5}
	s := arr[1:3]
	fmt.Println(len(s), cap(s), s)
	s = append(s, 80)
	s = append(s, 90)
	fmt.Println(len(s), cap(s), s)
	fmt.Println(arr)
	s = append(s, 100)
	fmt.Println(len(s), cap(s), s)
	fmt.Println(arr)
	s[0] = 500
	fmt.Println(len(s), cap(s), s)
	fmt.Println(arr)
}
