package main

import (
	"fmt"
)

func main() {
	// wg := &sync.WaitGroup{}
	// evenChan := make(chan int)
	// oddChan := make(chan int)
	// go evenChan1(evenChan)
	// go oddChan1(oddChan)
	// wg.Add(1)
	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	for i := 1; i <= 5; i++ {
	// 		if i%2 == 0 {
	// 			evenChan <- i
	// 		} else {
	// 			oddChan <- i
	// 		}
	// 	}
	// }(wg)

	// wg.Wait()
	s := []int{1, 2, 3, 4, 5}
	//testSlice(s)

	for i, val := range s {
		s[i] = 500
		fmt.Println(val)
	}

	var a []int
	b := []int{}
	fmt.Println("hello...  ", s)
	fmt.Println(a)
	fmt.Println(b)

	m := make(map[int]struct{})
	m[12] = struct{}{}
	m[0] = struct{}{}
	m[1] = struct{}{}
	m[2] = struct{}{}

	for i, val := range m {
		fmt.Println(i, val)
	}
}

func evenChan1(i chan int) {
	for val := range i {
		fmt.Println(val)
	}

}

func oddChan1(i chan int) {
	for val := range i {
		fmt.Println(val)
	}
}

func testSlice(s []int) {
	s[0] = 100
}
