package main

import (
	"encoding/json"
	"fmt"
	"time"
)

//PhysicalMemory testing
type physicalMemory struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	SizeBytes    uint64 `json:"sizeBytes"`
	ECC          string `json:"ecc,omitempty"`
	Slot         string `json:"slot,omitempty"`
	Type         string `json:"type,omitempty"`
	Speed        string `json:"speed,omitempty"`
	Status       string `json:"status,omitempty"`
	PartNumber   string `json:"partNumber,omitempty"`
}
type assetBaseBoard struct {
	Product      string    `json:"product,omitempty" cql:"product"`
	Manufacturer string    `json:"manufacturer" cql:"manufacturer"`
	Model        string    `json:"model,omitempty" cql:"model"`
	SerialNumber string    `json:"serialNumber,omitempty" cql:"serial_number"`
	Version      string    `json:"version,omitempty" cql:"version"`
	Name         string    `json:"name,omitempty" cql:"name"`
	InstallDate  time.Time `json:"installDate,omitempty" cql:"install_date"`
	HardwareUUID string    `json:"hardwareUUID,omitempty" cql:"hardware_uuid"`
}
type assetCollection struct {
	Memory []physicalMemory `json:"physicalMemory"`
}

func main() {
	const prefix string = ",\"Action\":"
	fmt.Println(prefix)
	jsonData := `{
        "product": "00F6D3",
        "manufacturer": "Dell Inc.",
        "serialNumber": "/2QGF3M2/CN129637C801FC/",
        "version": "A00",
        "name": "Base Board",
        "installDate": "0001-01-01T00:00:00Z"
	}`
	var data assetBaseBoard
	//ptr := &data
	json.Unmarshal([]byte(jsonData), &data)
	fmt.Println(data)

}
