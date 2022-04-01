package main

import (
	"fmt"
	"os/exec"
	"strconv"

	ps "github.com/shirou/gopsutil/process"
)

func InstallationManager() {
	//str := "tasklist /v /fo csv | findstr /i platform-agent-cor"

	// cmd := exec.Command("tasklist")
	// var stdoutBuf, stderrBuf bytes.Buffer
	// cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	// cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Printf("cmd.Run() failed with %s\n", err)
	// }
	// outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	// if strings.Contains(outStr, "platform-agent-core.exe") {
	// 	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	// }
	// err := kill.Output()

	processList, err := ps.Processes()
	if err != nil {
		fmt.Println("ps.Processes() Failed, are you using windows?")
		return
	}
	//fmt.Printf("%+v", processList)
	for _, v := range processList {
		//fmt.Println(v.Name())
		if name, err := v.Name(); err == nil && name == "platform-eventlog-plugin.exe" {
			fmt.Println(name, v.Pid)
			err = v.Kill()
			err = v.Terminate()
			kill := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(int(v.Pid)))
			err := kill.Run()
			if err != nil {
				fmt.Println("Error killing process", err)
			}
		}
	}
}

func main() {
	//var wg sync.WaitGrop

	InstallationManager()
}
