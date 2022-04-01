package main

import "fmt"

func evenodd(no int) {
	if no%2 == 0 {
		fmt.Println(no, " is even")
	} else {
		fmt.Println(no, " is odd")
	}
}

func checknum() {

	if num := 9; num < 0 {
		fmt.Println(num, " is negative")
	} else if num < 10 {
		fmt.Println(num, " is single digit")
	} else {
		fmt.Println(num, " is multiple digit")
	}
}

func main() {
	evenodd(1)
	checknum()
}
