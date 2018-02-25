package main

import (
	"fmt"
)

func main() {
	matrix := [][]int{[]int{1, 2, 3, 4}, []int{5, 1, 2, 3}, []int{9, 5, 1, 2}}
	fmt.Println(tmatrix(matrix))
}

func tmatrix(matrix [][]int) bool {
	rn := len(matrix)
	cn := len(matrix[0])

	for ci := 0; ci < cn; ci++ {
		ele := matrix[rn-1][ci]
		r := rn - 1
		c := ci
		if !check(r, c, ele, matrix) {
			return false
		}
	}

	for ri := 0; ri < rn; ri++ {
		ele := matrix[ri][cn-1]
		r := ri
		c := cn - 1
		if !check(r, c, ele, matrix) {
			return false
		}
	}

	return true
}

func check(r, c int, ele int, matrix [][]int) bool {
	if r == 0 || c == 0 {
		return true
	}
	if matrix[r-1][c-1] != ele {
		return false
	}
	return check(r-1, c-1, matrix[r-1][c-1], matrix)
}
