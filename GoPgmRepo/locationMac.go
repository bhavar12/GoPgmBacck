package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	cmd = "systemsetup -gettimezone"
)

func main() {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error in execution", err)
	} else {
		loc := strings.TrimSpace(string(out))
		locNew := strings.Split(loc, ":")
		fmt.Println("Time Zone location", locNew[1])
	}
}
