package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
	//"github.com/araddon/dateparse"
)

func main() {

	//out, err := exec.Command("dumpe2fs /dev/sda1 | grep 'Filesystem created:'").Output()
	//out, err := exec.Command("dumpe2fs /dev/sda1").Output()
	//out, err := exec.Command("bash", "-c", "sudo sh ./test_hello.sh").Output()
	opt, err := exec.Command("bash", "-c", "fs=$(df / | tail -1 | cut -f1 -d' ') && tune2fs -l $fs | grep created").Output()

	if err != nil {
		fmt.Println("Error", err)
	} else {
		str := "Filesystem created:       Thu May 31 11:50:58 2018"
		str = strings.Trim(string(opt), "Filesystem created:       ")
		//str1 := strings.TrimSpace(str)
		//str1 := strings.Spli

		fmt.Println("Original Date", string(opt))
		//str := "Fri Sep 13 11:39:05 2019"
		const layout = "Mon Jan _2 15:04:05 2006"
		// t, err := time.ParseInLocation(layout, str1, time.UTC)
		// fmt.Println("Server: ", t.Format(time.RFC850))
		// // t, err := time.Parse("Fri Jan _2 15:04:05 2006", str)
		t, err := time.Parse(layout, str)

		// //t, err := dateparse.ParseLocal(str)
		// //t, err := dateparse.ParseIn(str, time.UTC)
		if err != nil {
			panic(err.Error())
		}
		//date, err := time.Parse("20060102", string(t))
		// t, err := time.Parse("2006-01-02T15:04:05Z", str)

		fmt.Println("Time", t)
		//fmt.Println("Checking Date parse utility", t)
	}
}
