package main

import "fmt"

func main() {
	// slice of 1 and 0 s.  sort the slice in Dec order. complexity should be o(n).
	// no need to use any sorting algo. don't use inbuild function
	s := []int{1, 0, 0, 1, 0, 1, 0, 0, 0, 1}

	i := 0
	j := len(s) - 1
	for i < j {
		if s[i] == 1 {
			i++
			continue
		}
		if s[i] == 0 {
			if s[j] == 0 {
				j--
				continue
			} else {
				s[i] = s[j]
				i++
				s[j] = 0
				j--
			}
		}
	}
	fmt.Println(s)
}
