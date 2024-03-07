package main

import (
	"fmt"
	"math"
)

func maxprofit(prices []int, start, end int) int {

	if end <= start {
		return 0
	}
	profit := 0.0

	// the stock at witch we bought
	for i := start; i < end; i++ {

		// the stock at witch we need sell

		for j := i + 1; j <= end; j++ {

			if prices[j] > prices[1] {
				//update the current profit
				curr_profit := prices[j] - prices[i] + maxprofit(prices, start, i-1) + maxprofit(prices, j+1, end)

				// update the profit

				profit = math.Max(float64(profit), float64(curr_profit))
			}

		}
	}
	return int(profit)
}
func main() {
	prices := []int{100, 180, 260, 310, 40, 535, 695}
	fmt.Println(maxprofit(prices, 0, len(prices)-1))
}
