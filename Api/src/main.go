package main

import "fmt"

func main1() {
	a := "abc"
	result := gen(a)
	fmt.Println("Combinations are : ", result)
}
func gen(a string) []string {
	result := []string{}
	genhelp([]byte(a), 0, &result)
	return result
}
func genhelp(s []byte, start int, result *[]string) {
	if start == len(s) {
		*result = append(*result, string(s))
		fmt.Println("iteration:= ", string(s))
		return
	}
	for i := start; i < len(s); i++ {
		s[start], s[i] = s[i], s[start]
		genhelp(s, start+1, result)
		s[start], s[i] = s[i], s[start]
	}
}
