package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(arrayPairSum([]int{4, 6, 2, 7, 34, 1}))
}

func arrayPairSum(nums []int) int {
	var r int
	sort.Ints(nums)
	for i := 0; i < len(nums)/2; i++ {
		r += nums[2*i]
	}
	return r
}
