package main

import "fmt"

func Hello(s string) string {
	if len(s) == 0 {
		return "Hello Dude!"
	} else {
		return fmt.Sprintf("Hello %v ", s)
	}
}
