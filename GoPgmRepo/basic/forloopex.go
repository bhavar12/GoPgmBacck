package main

import "fmt"

func forloopwithoutcondition() {
	fmt.Println("In forloopwithoutcondition ")
	for {
		fmt.Println("loop")
		break
	}
}

func forloopwithcondition() {
	fmt.Println("In forloopwithcondition ")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
}

func forlikewhile() {
	fmt.Println("In forlikewhile")
	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}
}

func forloopwithconandbreak() {
	fmt.Println("In forloopwithcontandbreak")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println(i)

		if i == 7 {
			break
		}
	}
}
func main() {

	forloopwithcondition()
	forloopwithoutcondition()
	forlikewhile()
	forloopwithconandbreak()
}
