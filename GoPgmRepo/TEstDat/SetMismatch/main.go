package main

import (
	"fmt"
	"sort"
)

func findErrorNums(nums []int) []int {
	sort.IntSlice(nums).Sort()
	n := len(nums)
	duplicate := -1
	missing := 1
	for i := 1; i < n; i++ {
		if nums[i] == nums[i-1] {
			duplicate = nums[i]
		} else if nums[i] > (nums[i-1] + 1) {
			missing = (nums[i-1] + 1)
		}
	}
	calMiss := -1
	if nums[n-1] != n {
		calMiss = n
	} else {
		calMiss = missing
	}
	return []int{duplicate, calMiss}
}

func main() {
	arr := []int{3, 2, 2}
	fmt.Println(findErrorNums(arr))
}
