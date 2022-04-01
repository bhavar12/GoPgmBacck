package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/zcalusic/sysinfo"
)

var (
	cmd  = "timedatectl | grep 'Timezone'"
	cmd1 = "timedatectl | grep 'Time zone'"
)

func main() {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	fmt.Println("Time zone description using utility... ", si.Node.Timezone)
	if si.Node.Timezone == "" {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println("Error in execution", err)
			out, err := exec.Command("bash", "-c", cmd1).Output()
			if err != nil {
				fmt.Println("Error in execution usinf cmd1", err)
			} else {
				str := strings.TrimSpace(string(out))
				strArray := strings.Split(str, ":")
				strNew := strings.Split(strArray[1], "(")
				fmt.Println("Location using timedatectl command cmd1   ", strings.TrimSpace(strNew[0]))
			}
		} else {
			//str := "Time zone: America/New_York (EDT, -0400)"
			str := strings.TrimSpace(string(out))
			strArray := strings.Split(str, ":")
			strNew := strings.Split(strArray[1], "(")
			fmt.Println("Location using timedatectl command cmd  ", strings.TrimSpace(strNew[0]))
		}
	}
}
