package main

import (
	"errors"
	"fmt"
)

func sum(a int, b int) int {
	return a + b
}

func sumplus(a, b, c int) int {
	return a + b + c
}

// returning two values
func swap(a, b string) (string, string) {
	return b, a
}

func witherror() (string, error) {
	return "With Error", errors.New("Just checking how to handle error")
}

func variadic(nums ...int) {
	fmt.Print(nums, " addition of variadic ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {
	res := sum(2, 4)
	fmt.Println("2+4:- ", res)

	res1 := sumplus(3, 4, 5)
	fmt.Println("3+4+5:- ", res1)

	str1, str2 := swap("hello", "world")
	fmt.Println(str1)
	fmt.Println(str2)

	// ommiting first string
	_, strsecond := swap("Hello", "world")
	fmt.Println(strsecond)

	value, err := witherror()
	if err != nil {
		fmt.Println("Error occured ", err)
	}
	fmt.Println(value)

	// calling Variadic functions
	variadic(1, 2)
	variadic(1, 2, 3)
	num := []int{1, 2, 3, 4, 5}
	//If you already have multiple args in a slice,
	//apply them to a variadic function using func(slice...) like this.
	variadic(num...)

}
