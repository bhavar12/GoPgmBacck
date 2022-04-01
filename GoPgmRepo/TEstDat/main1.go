package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	pathSep   = "/"
	logDir    = "log"
	pluginDir = "plugin"
	path      = "../plugin/vmware/platform-vmware-plugin"
)

type Dog struct {
	Breed string
	// The first comma below is to separate the name tag from the omitempty tag
	WeightKg int  `json:",omitempty"`
	Status   bool `json:",omitempty"`
}

func main() {

	// _, pluginName := filepath.Split(path)
	// fmt.Println("Plugin Name  ", pluginName)
	// pluginType := pluginName[strings.Index(pluginName, "-")+1 : strings.LastIndex(pluginName, "-")]
	// // retDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// // if err != nil {
	// // 	fmt.Println("error", err)
	// // }
	// //pathToAppend := filepath.Dir(retDir)
	// fmt.Println("Plugin type   ", pluginType)
	// logFileName := fmt.Sprintf("%s_agent_plugin.log", pluginType)
	// configFileName := fmt.Sprintf("%s_agent_plugin_cfg.json", pluginType)
	// //fmt.Println("Path to Append", pathToAppend)
	// fmt.Println("log file name  ", logFileName)
	// fmt.Println("config file name  ", configFileName)
	// a := parseName(path)
	// fmt.Printf("array data %+v", a)
	d := Dog{
		Breed: "dalmation",
	}
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
func parseName(name string) []string {
	return strings.Split(strings.Trim(name, pathSep), pathSep)
}
