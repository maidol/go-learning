import "sort"

func arrayPairSum(nums []int) int {
	var r int
	sort.Ints(nums)
	for i := 0; i < len(nums)/2; i++ {
		r += nums[2*i]
	}
	return r
}
