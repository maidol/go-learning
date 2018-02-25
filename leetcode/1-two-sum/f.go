package main

import (
	"fmt"
)

func main() {
	fmt.Println(twoSum1([]int{4, 4, 4, 5, 15}, 9))
}

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target && i != j {
				return []int{i, j}
			}
		}
	}
	return nil
}

func twoSum1(nums []int, target int) []int {
	m := map[int]int{}
	for i := 0; i < len(nums); i++ {
		c := target - nums[i]
		o, ok := m[c]
		if ok {
			return []int{o, i}
		}
		_, ok = m[nums[i]]
		if !ok {
			m[nums[i]] = i
		}
	}
	return nil
}
