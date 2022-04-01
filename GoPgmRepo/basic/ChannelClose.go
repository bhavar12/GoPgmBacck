package main

import (
	"fmt"
)

func main() {

	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("Job Received", j)
			} else {
				fmt.Println("All Job Received")
				done <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("Sent Job", j)
	}
	close(jobs)
	fmt.Println("Sent All Jobs")
	<-done
}
