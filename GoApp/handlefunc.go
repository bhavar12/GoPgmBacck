package main

import "golang.org/x/sys/windows/registry"

import "fmt"

type AgentDetails struct {
	//Version plist/struct version
	ArbiterVersion string `json:"arbiterVersion" plist:"arbiterversion"`
	//EndpointID Plain endpoint ID of an agent
	EndpointID string `json:"endpointID" plist:"endpointid"`
}

const (
	cRegServiceImagePath = "SOFTWARE\\ITSPlatform"
)

func read() {
	//k, err := registry.OpenKey(registry.LOCAL_MACHINE, cRegServiceImagePath, registry.QUERY_VALUE|registry.READ|registry.WRITE|registry.ALL_ACCESS)
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, cRegServiceImagePath, registry.QUERY_VALUE|registry.READ)
	// if t != false {
	// 	fmt.Printf("false")
	// }
	if err != nil {
		fmt.Printf("Error during read", err)
	}
	defer k.Close()
	a, _, _ := k.GetStringValue("endpointid")
	b, _, _ := k.GetStringValue("arbiterversion")
	fmt.Printf("EndPointID%s", a)
	fmt.Printf("ArbiterVersion%s", b)
}
func write() {
	//k, err := registry.OpenKey(registry.LOCAL_MACHINE, cRegServiceImagePath, registry.QUERY_VALUE|registry.READ|registry.WRITE|registry.ALL_ACCESS)
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, cRegServiceImagePath, registry.QUERY_VALUE|registry.WRITE)
	if err != nil {
		fmt.Printf("Error during write", err)
	}
	defer k.Close()

	k.SetStringValue("arbiterversion", "1")
	k.SetStringValue("endpointid", "22222-5555")
}
func delete(){
	err := registry.DeleteKey(registry.LOCAL_MACHINE, cRegServiceImagePath)
	if err != nil {
		fmt.Printf("Error in Deleting key  ", err)
	}
}
func main() {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, cRegServiceImagePath, registry.QUERY_VALUE|registry.READ|registry.WRITE)
	if err != nil {
		fmt.Printf("Error in creating key  ", err)
	}
	fmt.Printf("Key Value%v", key)
	defer key.Close()
	key.SetStringValue("arbiterversion", "1")
	key.SetStringValue("endpointid", "22222-5555")
	read()
	delete()
}
