package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}
func test() {
	test1 := false
	go func(msg string) {
		test1 = true
		defer func() {
			test1 = false
			fmt.Println("Completed go routine")
		}()
		fmt.Println(msg)
	}("going")
	time.Sleep(time.Second)
	fmt.Println("done")
}

func main() {
	var err error
	var win []struct {
		SerialNumber string
	}
	err = wmi.Query("select SerialNumber from Win32_operatingsystem where SerialNumber is not null", &win)
	if err == nil && len(win) > 0 {
		data := make([]string, len(win))
		for i := 0; i <= len(win)-1; i++ {
			data[i] = strings.Trim(win[i].SerialNumber, " ")
		}
		sort.Strings(data)
		fmt.Println(strings.Join(data, ","))
	}
}
