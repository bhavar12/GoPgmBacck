package main

import (
	"fmt"
	"time"
)

// Write a program to have keepalive between 2 service every 700ms
// The service checking the keep alive terminates after 16s and terminate all service.
func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan bool)
	ter := make(chan bool)
	go serviceA(ch1, ch2, ter)
	<-ter
	fmt.Println("All services are closed...")
}

func serviceA(ch1 chan struct{}, ch2, ter chan bool) {
	i := 0
	go serviceB(ch1)
	go func(i int) {
		for {
			ch1 <- struct{}{}
			time.Sleep(time.Millisecond * 700)
			fmt.Printf("calling....%d", i)
			i++
		}
	}(i)
	time.Sleep(time.Second * 16)
	ter <- true
}

func serviceB(ch1 chan struct{}) {
	for {
		select {
		case <-ch1:
			fmt.Println("service alive")
		}
	}
}
