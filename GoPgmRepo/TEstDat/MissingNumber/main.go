package main

import (
	"fmt"
	"sort"
)

func missingNumber(nums []int) int {
	missing := -1
	n := len(nums)
	sort.IntSlice(nums).Sort()
	if n == 1 {
		if nums[0] == 0 {
			return 1
		}
		return 0

	}
	for i := 0; i < n; i++ {
		if nums[i] != i {
			return i
		}

	}
	return missing
}

func main() {
	arr := []int{0, 2}
	fmt.Println(missingNumber(arr))
}
