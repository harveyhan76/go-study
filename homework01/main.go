package main

import "fmt"

func main() {
	// 测试只出现一次的数字
	nums1 := []int{2, 2, 1}
	fmt.Printf("数组 %v 中只出现一次的数字是: %d\n", nums1, singleNumber(nums1))
	fmt.Printf("数组 %v 中只出现一次的数字是: %d\n", nums1, singleNumberXOR(nums1))

	// 测试回文数
	fmt.Printf("数字 %d 是否是回文数: %v\n", 121, siPalidrome(121))
	fmt.Printf("数字 %d 是否是回文数: %v\n", -121, siPalidrome(-121))
	fmt.Printf("数字 %d 是否是回文数: %v\n", 10, siPalidrome(10))

	// 测试有效的括号字符串
	testCases := []string{
		"()",     // true
		"()[]{}", // true
		"(]",     // false
		"([)]",   // false
		"{[]}",   // true
		"((()))", // true
		"([{}])", // true
		"",       // true
		"((())",  // false
	}
	for _, s := range testCases {
		fmt.Printf("isValid(\"%s\") = %v\n", s, isValid(s))
	}

	// 测试查找字符串数组中的最长公共前缀
	prefixStrs := [][]string{
		{"flower", "flow", "flight"},
		{"dog", "racecar", "car"},
		{"", "abc"},
		{"abc"},
		{"abc", "abc", "abc"},
		{"prefix", "preface", "preview"},
	}

	fmt.Println("横向扫描:")
	for _, strs := range prefixStrs {
		fmt.Printf("%v -> \"%s\"\n", strs, longestCommonPrefix(strs))
	}

	// 测试给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
	testCases2 := [][]int{
		{1, 2, 3},
		{4, 3, 2, 1},
		{9},
		{9, 9, 9},
		{0},
	}

	for _, digits := range testCases2 {
		result := plusOne(digits)
		fmt.Printf("%v -> %v\n", digits, result)
	}

	// 测试原地删除有序数组重复元素
	nums2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	len2 := removeDuplicates(nums2)
	fmt.Println(len2, nums2[:len2]) // 输出: 5 [0 1 2 3 4]

	// 测试合并区间
	testCases3 := [][][]int{
		{{1, 4}, {0, 2}, {3, 5}},
		{},                       // 空数组
		{{1, 4}},                 // 只有一个区间
		{{5, 8}, {1, 3}, {2, 4}}, // 需要排序的
		{{4, 5}, {1, 4}, {0, 1}}, // 边界情况
	}

	fmt.Println("测试合并区间功能：")
	fmt.Println("=================")

	for i, intervals := range testCases3 {
		fmt.Printf("测试合并区间 %d:\n", i+1)
		fmt.Printf("输入: %v\n", intervals)
		result := merge(intervals)
		fmt.Printf("输出: %v\n\n", result)
	}

	// 测试两数之和
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Println(result) // 输出 [0 1]
}

// singleNumber 使用哈希表统计每个数字的出现次数，返回只出现一次的数字
func singleNumber(nums []int) int {
	numsMap := make(map[int]int)

	for _, num := range nums {
		numsMap[num]++
	}
	for num, count := range numsMap {
		if count == 1 {
			return num
		}
	}
	return 0
}

// singleNumberXOR 使用异或运算查找只出现一次的数字
func singleNumberXOR(nums []int) int {
	result := 0
	// 遍历数组，对每个数字进行异或运算
	// 相同的数字异或后会抵消，最终结果是只出现一次的数字
	for _, num := range nums {
		result ^= num
	}
	return result
}

func siPalidrome(x int) bool {
	// 负数不是回文数，个位为0的非0数也不是回文数
	if x < 0 && (x%10 == 0 && x != 0) {
		return false
	}
	reversed := 0
	// 只反转数字的一半
	for x > reversed {
		// 取出x的个位数并添加到reversed末尾
		reversed = reversed*10 + x%10
		// 去掉x的个位数
		x /= 10
	}

	// 对于偶数位数字：x == reversed
	// 对于奇数位数字：x == reversed / 10（去掉中间位）
	return x == reversed
}

// 测试有效的括号字符串
func isValid(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	stack := []byte{}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	for i := 0; i < len(str); i++ {
		char := str[i]
		// 如果是右括号
		if value, exists := pairs[char]; exists {
			// 如果栈为空，或者栈顶元素不匹配当前右括号对应的左括号
			if len(stack) == 0 || stack[len(stack)-1] != value {
				return false
			}

			// 匹配成功，弹出栈顶元素
			stack = stack[:len(stack)-1]
		} else {
			// 如果是左括号，压入栈中
			stack = append(stack, char)
		}

	}
	return len(stack) == 0
}

// 查找字符串数组中的最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		prefix = commonPrefix(prefix, strs[i])
		if prefix == "" {
			return ""
		}
	}
	return prefix
}

// 对比俩字符串最大公共前缀
func commonPrefix(str1, str2 string) string {
	length := min(len(str1), len(str2))
	index := 0
	for index < length && str1[index] == str2[index] {
		index++
	}
	return str1[:index]
}

// 取最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func plusOne(digits []int) []int {
	// 创建一个新数组，复制原始数组的值
	newDigits := append([]int{}, digits...)
	n := len(newDigits)
	for i := n - 1; i >= 0; i-- {
		if newDigits[i] < 9 {
			newDigits[i]++
			return newDigits
		}
		newDigits[i] = 0
	}
	// 如果所有数字都是9，则需要增加一位
	return append([]int{1}, newDigits...)
}

// 原地删除有序数组重复元素
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// 慢指针，指向当前不重复元素的最后一个位置
	i := 0

	// 快指针，遍历整个数组
	for j := 1; j < len(nums); j++ {
		// 当发现不同的元素时
		if nums[j] != nums[i] {
			// 将慢指针后移
			i++
			// 将新元素复制到慢指针的位置
			nums[i] = nums[j]
		}
	}

	// 返回新长度（i是索引，需要+1）
	return i + 1
}

// 测试合并区间
// 插入排序实现
func insertionSort(intervals [][]int) {
	n := len(intervals)
	for i := 1; i < n; i++ {
		key := intervals[i]
		j := i - 1

		// 将 intervals[i] 插入到已排序的 intervals[0..i-1] 中
		for j >= 0 && (intervals[j][0] > key[0] ||
			(intervals[j][0] == key[0] && intervals[j][1] > key[1])) {
			intervals[j+1] = intervals[j]
			j--
		}
		intervals[j+1] = key
	}
}

// 合并区间的主函数
func merge(intervals [][]int) [][]int {
	// 空数组或只有一个区间，直接返回
	if len(intervals) <= 1 {
		return intervals
	}

	// 1. 先排序
	insertionSort(intervals)

	// 2. 合并区间
	result := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		// 当前区间
		curr := intervals[i]
		// 结果中的最后一个区间
		last := result[len(result)-1]

		// 判断是否有重叠
		if curr[0] <= last[1] {
			// 有重叠，合并（取结束位置的最大值）
			if curr[1] > last[1] {
				last[1] = curr[1]
			}
			result[len(result)-1] = last
		} else {
			// 没有重叠，直接加入结果
			result = append(result, curr)
		}
	}

	return result
}

// 两数之和
func twoSum(nums []int, target int) []int {
	// 哈希表，key 为数组元素值，value 为该元素的索引
	indexMap := make(map[int]int)

	for i, num := range nums {
		// 计算当前数字需要的互补数
		complement := target - num

		// 检查互补数是否已经在哈希表中
		if idx, ok := indexMap[complement]; ok {
			// 找到结果，返回两个索引
			return []int{idx, i}
		}

		// 没找到则将当前数字和索引存入哈希表
		indexMap[num] = i
	}

	// 未找到符合条件的两个数，返回空切片
	return []int{}
}
