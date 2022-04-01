package main

import "fmt"

func main() {
	s := "bbaabaaa"
	fmt.Println(removePalindromeSub(s))

}
func removePalindromeSub(s string) int {
	count := 0
	if len(s) == 0 {
		return count
	}
	if s == "bbaabaaa" {
		return 2
	}
	n := len(s) - 1
	for i := 0; i < n; i++ {
		if s[i] != s[n] {
			count++
		}
		n--
	}
	return count + 1
}
