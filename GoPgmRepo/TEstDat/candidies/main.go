package main

import (
	"fmt"
	"sort"
)

func distributeCandies(candyType []int) int {
	n := len(candyType)
	if n%2 != 0 {
		return -1
	}
	sort.IntSlice(candyType).Sort()
	possibleCandidies := n / 2
	c := 0
	set := make(map[int]bool, 0)
	for i := 0; i < n; i++ {
		if !set[candyType[i]] {
			set[candyType[i]] = true
			c++
			if c == possibleCandidies {
				return possibleCandidies
			}
		}
	}
	return c
}

func distributeCandiesBruteForse(candyType []int) int {
	n := len(candyType)
	if n%2 != 0 {
		return -1
	}
	sort.IntSlice(candyType).Sort()
	possibleCandidies := n / 2
	uniqueValue := 0
	for i := 0; i < n; i++ {
		if i == n-1 {
			uniqueValue++
			break
		}
		if candyType[i] != candyType[i+1] {
			uniqueValue++

		}
	}
	newArray := make([]int, uniqueValue)
	j := 0
	for i := 0; i < n; i++ {
		if i == n-1 {
			newArray[j] = candyType[i]
			break
		}
		if candyType[i] != candyType[i+1] {

			newArray[j] = candyType[i]
			j++
		}
	}
	if uniqueValue <= possibleCandidies {
		return uniqueValue
	}

	return possibleCandidies

}

func main() {
	arr := []int{100000, 0, 100000, 0, 100000, 0, 100000, 0, 100000, 0, 100000, 0}
	fmt.Println(distributeCandies(arr))
}
