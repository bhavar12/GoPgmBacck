package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/StackExchange/wmi"
)

type win32ComputerSystem struct {
	Manufacturer    string
	Model           string
	Name            string
	Domain          string
	DomainRole      int
	CurrentTimeZone int
}

const (
	q = "SELECT Manufacturer,Model,Name,CurrentTimeZone,Domain,DomainRole FROM Win32_ComputerSystem"
)

func main() {

	fmt.Println("Checking Virtual or Physical")
	bRet, err := checkSystemInfo()
	if nil == err {
		if bRet {
			fmt.Println("It is a virtual system")
		} else {
			bRet, err = checkModelInfo()
			if bRet {
				fmt.Println("It is a virtual system")
			} else {
				fmt.Println("It is a physical system")
			}
		}
	} else {
		bRet, err = checkModelInfo()
		if bRet {
			fmt.Println("It is a virtual system")
		} else {
			fmt.Println("It is a physical system")
		}
	}
}

func checkSystemInfo() (bRet bool, err error) {
	sysData, err := exec.Command("systeminfo").Output()

	if nil != err {
		return
	}
	strResult := string(sysData)
	strResult = strings.ToLower(strResult)
	if strings.Contains(strResult, "hypervisor has been detected") {
		return true, nil
	}
	return
}

func checkModelInfo() (bRet bool, err error) {
	var dst []win32ComputerSystem

	//wmiObj := wmi.GetWrapper()
	err = wmi.Query(q, &dst)
	if nil != err {
		return
	}

	for _, data := range dst {
		model := strings.ToLower(data.Model)
		if strings.Contains(model, "virtual") {
			return true, nil
		}
	}
	return
}
