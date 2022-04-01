package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5}
	b := a[2:4]
	d := a[2:4]
	d = append(d, 100)
	b = append(b, 101)
	fmt.Println(b) //[3,4]
	fmt.Println(d)
	fmt.Println(a)
	fmt.Println(b)   //[3,4,3]
	fmt.Println(a)   //[1,2,3,4,3]
	b = append(b, 4) //[3,4,3,4] // here capacity became 6, so going forward we can add 2 more ele
	fmt.Println(a)   //[1,2,3,4,3]
	b = append(b, 8)
	c := b[:]
	c = append(c, 9) // [3,4,3,4,8,9]
	c[2] = 50
	fmt.Println(b) //[50,4,3,4,8]
	fmt.Println(c)
	fmt.Println(b)
	fmt.Println(a)
	fmt.Println(len(b))
	fmt.Println(cap(b))
	fmt.Println(len(c))
	fmt.Println(cap(c))
}
