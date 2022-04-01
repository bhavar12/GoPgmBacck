package main

import "fmt"

func another(f func(string) string) {
	result := f("David")
	fmt.Println(result)
}

func main() {

	anon := func(str string) string {
		return str
	}
	another(anon)
}
