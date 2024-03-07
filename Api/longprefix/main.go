package main

import (
	"fmt"
	"strings"
)

func main() {
	s := []string{"flower", "flie", "flow"}
	fmt.Println(longestCommonPrefix(s))
}

func longestCommonPrefix(s []string) string {
	pref := s[0]
	for i := 1; i < len(s); i++ {
		for !strings.HasPrefix(s[i], pref) {
			pref = pref[:len(pref)-1]
		}
		//fmt.Println(s[i])
	}
	return pref
}
