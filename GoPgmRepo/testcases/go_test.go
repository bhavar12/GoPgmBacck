package main

import (
	"fmt"
	"testing"
)

func TestHelo(t *testing.T) {
	res := Hello("")
	if res != "Hello Dude!" {
		t.Error("Hello Fails")
	} else {
		t.Log("Hello Success")
	}
	fmt.Println("good", res)
}
