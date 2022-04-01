package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

var (
	wordSet       = make(map[string]bool)
	wordsCap      = make(map[string]string)
	wordsVowl     = make(map[string]string)
	index     int = 0
)

func duplicateZeros(arr []int) {
	l := len(arr)
	for i := 0; i < l-1; i++ {
		if arr[i] == 0 {
			for j := l - 1; j > i; j-- {
				arr[j] = arr[j-1]
			}
			arr[i+1] = 0
			i = i + 1
		}
	}
	fmt.Println("Array Element: ", arr)
}

func spellchecker(wordlist []string, queries []string) []string {
	res := make([]string, len(queries))
	for _, wordval := range wordlist {
		if !wordSet[wordval] {
			wordSet[wordval] = true
		}
		strLow := strings.ToLower(wordval)
		_, ok := wordsCap[strLow]
		if !ok {
			wordsCap[strLow] = wordval
		}
		worldLov := dovowel(strLow)
		_, ok = wordsVowl[worldLov]
		if !ok {
			wordsVowl[worldLov] = wordval
		}
	}

	for i, val := range queries {
		res[i] = solve(val)
	}
	return res
}

func solve(query string) string {

	if wordSet[query] {
		return query
	}
	queryLow := strings.ToLower(query)

	val, ok := wordsCap[queryLow]
	if ok {
		return val
	}
	queryLV := dovowel(queryLow)
	val1, ok := wordsVowl[queryLV]
	if ok {
		return val1
	}
	return ""
}

func dovowel(word string) string {

	for _, val := range word {
		if checkVowel(string(val)) {
			word = strings.Replace(word, string(val), "*", 1)
		}
	}
	return word
}

func checkVowel(c string) bool {
	if c == "a" || c == "e" || c == "i" || c == "o" || c == "u" {
		return true
	}
	return false
}

func duplicateZerosNewApproach(arr []int) {
	possbileDups := 0
	l := len(arr) - 1
	for left := 0; left <= l-possbileDups; left++ {
		if arr[left] == 0 {
			if left == l-possbileDups {
				arr[l] = 0
				l--
				break
			}
			possbileDups++
		}
	}
	last := l - possbileDups
	for i := last; i >= 0; i-- {
		if arr[i] == 0 {
			arr[i+possbileDups] = 0
			possbileDups--
			arr[i+possbileDups] = 0
		} else {
			arr[i+possbileDups] = arr[i]
		}
	}
	fmt.Println("Array Element: ", arr)
}

func advantageCount(A []int, B []int) []int {
	sort.IntSlice(A).Sort()
	for i := 0; i < len(A); i++ {
		for j := i; j < len(A); j++ {
			if B[i] < A[j] {
				if i != j {
					temp := A[j]
					A[j] = A[i]
					A[i] = temp
					break
				}
				if i == j {
					break
				}

			}
		}
	}
	return A
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func canFit(inner []int, outer []int) bool {

	if inner[0] < outer[0] && inner[1] < outer[1] {
		return true
	}
	return false
}

func maxEnvelopes(envelopes [][]int) int {

	if len(envelopes) == 0 {
		return 0
	}

	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] == envelopes[j][0] {
			return envelopes[i][1] < envelopes[j][1]
		}
		return envelopes[i][0] < envelopes[j][0]
	})

	dp := make([]int, len(envelopes))
	for i := 0; i < len(dp); i++ {
		dp[i] = 1
	}

	maxCount := 1
	for i := 1; i < len(envelopes); i++ {
		for j := 0; j < i; j++ {
			if canFit(envelopes[j], envelopes[i]) {
				dp[i] = max(dp[i], dp[j]+1)
				maxCount = max(maxCount, dp[i])
			}
		}
	}

	return maxCount
}

func longestValidParentheses(s string) int {
	maxLen := 0
	for i := 0; i < len(s); i++ {
		for j := i + 2; j <= len(s); j += 2 {
			if isValid(s[i:j]) {
				maxLen = int(math.Max(float64(maxLen), float64(j-i)))
			}
		}
	}
	return maxLen
}

func isValid(s string) bool {
	stack := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			stack = append(stack, "(")
		} else if len(stack) != 0 && stack[len(stack)-1] == "(" {
			stack = stack[0 : len(stack)-1]
		} else {
			return false
		}
	}
	if len(stack) == 0 {
		return true
	}
	return false
}

func longestValidParentheses1(s string) int {
	left := 0
	right := 0
	maxlength := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			left++
		} else {
			right++
		}
		if left == right {
			maxlength = int(math.Max(float64(maxlength), float64(2*right)))
		} else if right >= left {
			left = 0
			right = 0
		}
	}
	left = 0
	right = 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '(' {
			left++
		} else {
			right++
		}
		if left == right {
			maxlength = int(math.Max(float64(maxlength), float64(2*left)))
		} else if left >= right {
			left = 0
			right = 0
		}
	}
	return maxlength
}

func main() {
	//arr := [][]int{{5, 4}, {6, 4}, {6, 7}, {2, 3}}
	//arr1 := []int{1, 10, 4, 11}
	arr1 := "(()"
	fmt.Println(longestValidParentheses(arr1))
	//arr1 := []int{0, 0, 0, 0, 0, 0, 0, 0}
	//arr1 := []int{0, 1, 2, 2, 3, 4, 4, 5}
	//duplicateZeros(arr2)
	// len := removeDuplicates(arr1)
	// for i := 0; i < len; i++ {
	// 	fmt.Printf("%d\n", arr1[i])
	// }
	// arr1 := []int{-2, 0, 10, -19, 4, 6, -8}
	// fmt.Println(checkIfExist(arr1))

	//arr1 := []int{1, 2, 5}
	//amount := 4
	//words := []string{"Yellow"}
	//queries := []string{"yellow"}
	//fmt.Println(spellchecker(words, queries))
}

func removeElement(arr []int, val int) int {
	j := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] != val {
			arr[j] = arr[i]
			j++
		}
	}
	return j
}

func removeElementNewApproach(arr []int, val int) int {
	i := 0
	n := len(arr)
	for i < n {
		if arr[i] == val {
			arr[i] = arr[n-1]
			n--
		} else {
			i++
		}
	}
	return n
}

func removeDuplicates(nums []int) int {
	n := len(nums)
	t := 0
	for j := 0; j < n; j++ {
		if j == n-1 {
			nums[t] = nums[j]
			t++
			break
		}
		if nums[j] != nums[j+1] {
			nums[t] = nums[j]
			t++
		}
	}
	return t
	// newArr := make([]int, possibleLenght)
	// k := 0
	// for i := 0; i < n; i++ {
	// 	if i == 0 {
	// 		newArr[k] = nums[i]
	// 		k++
	// 	}
	// 	if k == possibleLenght {
	// 		break
	// 	}
	// 	if newArr[k-1] != nums[i] {
	// 		newArr[k] = nums[i]
	// 		k++
	// 	}
	// }
	// for i := 0; i < possibleLenght; i++ {
	// 	nums[i] = newArr[i]
	// }
	//return t
}

func removeDuplicatesBruteForce(nums []int) int {
	n := len(nums)
	possibleLenght := 0
	for j := 0; j < n; j++ {
		if j == n-1 {
			possibleLenght++
			break
		}
		if nums[j] != nums[j+1] {
			possibleLenght++
		}
	}
	newArr := make([]int, possibleLenght)
	k := 0
	for i := 0; i < n; i++ {
		if i == 0 {
			newArr[k] = nums[i]
			k++
		}
		if k == possibleLenght {
			break
		}
		if newArr[k-1] != nums[i] {
			newArr[k] = nums[i]
			k++
		}
	}
	for i := 0; i < possibleLenght; i++ {
		nums[i] = newArr[i]
	}
	return possibleLenght
}

func checkIfExist(arr []int) bool {
	if len(arr) == 0 || arr == nil {
		return false
	}
	length := len(arr)
	i := 0
	for i < length {
		for j := 0; j < length; j++ {
			if i != j {
				if arr[i] == arr[j]*2 {
					return true
				}
			}
		}
		i++
	}
	return false
}

func validMountailn(arr []int) bool {
	if len(arr) < 3 || arr == nil {
		return false
	}
	length := len(arr)
	i := 0
	//walk up
	for i+1 < length && arr[i] < arr[i+1] {
		i++
	}

	if i == 0 || i == length-1 {
		return false
	}

	//walk down
	for i+1 < length && arr[i] > arr[i+1] {
		i++
	}
	return i == length-1
}

func replaceElements(arr []int) []int {

	if len(arr) == 0 || arr == nil {
		return nil
	}
	length := len(arr)
	i := 0
	ele := 0
	for i < length {
		for j := i + 1; j < length; j++ {
			if i == length-1 {
				ele = -1
				break
			}
			if j == length-1 {
				break
			}
			if arr[j] > arr[j+1] {
				ele = arr[j]
			}
		}
		arr[i] = ele
		i++
	}

	return arr
}
func coinChange(coins []int, amount int) int {
	n := len(coins)
	newArray := make([]int, amount+1)

	for i := 0; i <= amount; i++ {
		newArray[i] = amount + 1
	}
	newArray[0] = 0
	for i := 1; i <= amount; i++ {
		for j := 0; j < n; j++ {
			if coins[j] <= i {
				newArray[i] = int(math.Min(float64(newArray[i]), float64(newArray[i-coins[j]]+1)))
			}
		}
	}
	if newArray[amount] > amount {
		return -1
	}
	return newArray[amount]
}

func intToRoman(num int) string {
	var val = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	var rom = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	ans := ""
	for i := 0; num > 0; i++ {
		for num >= val[i] {
			ans = ans + rom[i]
			num = num - val[i]
		}
	}
	return ans
}

func letterCombinations(digits string) []string {
	var val = []int{2, 3, 4, 5, 6, 7, 8, 9}
	var rom = []string{"abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"}
	var output = make([]string, 0)
	n := len(digits)
	if n == 0 {
		return []string{}
	}
	for i := 0; i < n; i++ {
		if digits[i] == '2' {

		}
	}
	return output
}
