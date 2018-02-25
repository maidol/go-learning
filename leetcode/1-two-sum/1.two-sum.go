func twoSum(nums []int, target int) []int {
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
