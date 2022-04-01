package main

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

type win32Processor struct {
	Name                      string
	NumberOfCores             int
	MaxClockSpeed             int
	Family                    int
	Manufacturer              string
	DataWidth                 int
	ProcessorType             int
	Level                     int
	NumberOfLogicalProcessors int
}
type win32ProcessorXP2003 struct {
	Name          string
	MaxClockSpeed int
	Family        int
	Manufacturer  string
	DataWidth     int
	Level         int
}

type win32TimeZone struct {
	Caption     string
	Description string
	Bias        int32
}

const (
	//q             = "SELECT Name,NumberOfCores,MaxClockSpeed,Family,Manufacturer,ProcessorType,Level, NumberOfLogicalProcessors, DataWidth  FROM Win32_Processor"
	q1            = "SELECT Name,NumberOfCores,MaxClockSpeed,Family,Manufacturer,DataWidth,Level FROM Win32_Processor"
	q3            = "SELECT Name,NumberOfCores,MaxClockSpeed,Family,Manufacturer,DataWidth,Level FROM Win32_Processor"
	timeZoneQuery = "SELECT Bias, Caption, Description FROM Win32_TimeZone"
)

func main() {

	var dst []win32TimeZone
	err := wmi.QueryNamespace(timeZoneQuery, &dst, "root\\cimv2")
	if nil != err {
		fmt.Println("error", err)
	}
	fmt.Println(dst)
}
