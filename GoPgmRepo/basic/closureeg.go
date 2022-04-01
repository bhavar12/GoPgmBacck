package main

import "fmt"

//This function intSeq returns another function, which we define anonymously in the body of nextInt.
//The returned function closes over the variable i to form a closure.

func nextInt() func() int {

	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	nextint := nextInt()
	fmt.Println(nextint())
	fmt.Println(nextint())
	fmt.Println(nextint())

	newInt := nextInt()
	fmt.Println("new int", newInt())

}
