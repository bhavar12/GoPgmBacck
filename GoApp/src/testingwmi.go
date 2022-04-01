package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/StackExchange/wmi"
)

func main() {
	var err error
	var win []struct {
		SerialNumber string
	}
	var arr []int
	fmt.Println(len(arr))
	err = wmi.Query("SELECT Name,NumberOfCores,MaxClockSpeed,Family,Manufacturer,ProcessorType,Level, NumberOfLogicalProcessors  FROM Win32_Processor", &win)
	if err == nil && len(win) > 0 {
		data := make([]string, len(win))
		for i := 0; i <= len(win)-1; i++ {
			data[i] = strings.Trim(win[i].SerialNumber, " ")
		}
		sort.Strings(data)
		fmt.Println(strings.Join(data, ","))
	}
}
