package main

import (
	"fmt"
)

const (
	max_len = 1000
)

// Problem: Inplace URL Encoding of String

// Input: Good Morning India : byte[1024]
// Output: Good%32Morning%32India
// Condition: Inplace replacing, O(n), Expandable buffer.
func main() {
	str := []byte("Good Morning India     ")
	n := replace(str)
	for i := 0; i < n; i++ {
		fmt.Printf("%c", str[i])
	}

}

func replace(str []byte) int {
	var count_spaces, i int
	for i = 0; i < len(str); i++ {
		if str[i] == ' ' {
			count_spaces++
		}
	}
	fmt.Println("no of spaces", count_spaces)

	for str[i-1] == ' ' {
		count_spaces--
		i--
	}
	// finding new lenght
	new_length := i + count_spaces*2 + 1

	// new lenght must be smalled
	if new_length > max_len {
		return -1
	}

	// start the filling chat at the end
	index := new_length - 1
	str[index] = '0'
	index--

	// fill the rest of string from end

	for j := i - 1; j >= 0; j-- {
		// insert %20 in a place
		if str[j] == ' ' {
			str[index] = '0'
			str[index-1] = '2'
			str[index-2] = '%'
			index = index - 3
		} else {
			str[index] = str[j]
			index = index - 1
		}
	}
	fmt.Println(string(str))
	return new_length
}
