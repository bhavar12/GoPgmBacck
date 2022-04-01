package main

import (
	"fmt"
)

func main() {

	stud := make(map[int]string)

	stud[1] = "sachin"

	stud[2] = "mahesh"
	stud[3] = "ravi"
	fmt.Println("student map:- ", stud)
	fmt.Println("len:- ", len(stud))

	v1 := stud[1]

	fmt.Println(v1)

	delete(stud, 1)
	_, prs := stud[1]

	fmt.Println("present:- ", prs)

	newMap := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("new Map:-", newMap)

	for key, value := range newMap {
		fmt.Printf("%s -> %d\n", key, value)
	}

}
