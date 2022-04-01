package main

import (
	"fmt"
	"strings"
)

const (
	ErrAgentMappingNotFound string = "ErrAgentMappingNotFound"
	Status                  string = "ErrAgentMappingNotFound : no record found"
)

func main() {
	err := fmt.Errorf("%s : no record found", ErrAgentMappingNotFound)
	if strings.EqualFold(err.Error(), Status) {
		fmt.Println("true")
	} else {
		fmt.Println("false")

	}

}
