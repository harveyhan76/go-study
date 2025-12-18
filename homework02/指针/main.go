package main

import "fmt"

// 定义一个函数，接收整数指针作为参数
func addTen(ptr *int) {
	// 通过指针修改原始变量的值
	*ptr += 10
	fmt.Printf("函数内部：指针指向的值增加10后为 %d\n", *ptr)
	fmt.Printf("函数内部：指针本身的值（内存地址）是 %p\n", ptr)
}

// 使用range遍历
func multiplyByTwoRange(slicePtr *[]int) {
	s := *slicePtr
	for i := range s {
		s[i] *= 2
	}
}

func main() {
	// 定义一个整数变量
	num := 5
	fmt.Printf("调用函数前：num = %d\n", num)
	fmt.Printf("调用函数前：num的内存地址是 %p\n", &num)
	fmt.Println("--- 调用addTen函数 ---")

	// 传递变量的地址（指针）给函数
	addTen(&num)

	fmt.Println("=== 使用切片指针 ===")
	// 创建一个整数切片
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始切片: %v, 长度: %d, 容量: %d\n",
		numbers, len(numbers), cap(numbers))
	fmt.Printf("切片的地址: %p\n", &numbers)

	multiplyByTwoRange(&numbers)
	fmt.Printf("修改后: %v\n", numbers)

}
