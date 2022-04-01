package main

import (
	"fmt"
	"time"
)

func basicswtich() {
	i := 2
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println(" Default")
	}
}

func switchwithmultipleexpression() {
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")

	default:
		fmt.Println("It's Week")
	}
}

func comparetypesinsteadofvalue(i interface{}) {

	switch t := i.(type) {
	case bool:
		fmt.Println("I am bool")
	case int:
		fmt.Println("I am int")
	default:
		fmt.Println("Don't know type %T", t)
	}

}

func main() {
	basicswtich()
	switchwithmultipleexpression()
	comparetypesinsteadofvalue(true)
	comparetypesinsteadofvalue(1)
	comparetypesinsteadofvalue("hey")
}
