package main

import "fmt"

type Sensor struct {
	//Name indicate the name of the sensor
	Name string `json:"name,omitempty"`
	//HealthStatus property indicate the status of the sensor eg. Green
	HealthStatus string `json:"healthStatus,omitempty"`
	//BaseUnit property indicate the base unit of the sensor eg. RPM , Watts
	BaseUnit string `json:"baseUnit,omitempty"`
	//Type property indicate the type of the sensor eg. Fan,Power, Voltage
	Type string `json:"type,omitempty"`
	//Reading property indicate the current Reading of the sensor
	Reading string `json:"reading,omitempty"`
}
type Virtualization struct {
	//Sensors property indicate VMWare Hardware Sensor info
	Sensors []Sensor `json:"sensors,omitempty"`
}

func getValidIps(arr []string) []string {
	newArray := make([]string, 0)
	for _, val := range arr {
		if val != "0.0.0.0" {
			newArray = append(newArray, val)
		}
	}
	return newArray

}

func main() {
	arr := []string{"0.0.0.0", "0.0.0.0", "0.0.0.0"}
	fmt.Println(getValidIps(arr))
	//	var arr []*Sensor
	//arr = append(arr, s)
	// wInt, sParseErr1 := strconv.ParseInt("", 10, 64)
	// if sParseErr1 != nil {
	// 	fmt.Println("Error", sParseErr1)
	// }
	// fmt.Println(wInt)

}
