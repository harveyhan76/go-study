package main

import (
	"fmt"
	"sync"
)

func main() {
	//fmt.Println("Channel 示例")
	// 创建一个无缓冲的整数通道
	ch := make(chan int)

	// 使用WaitGroup等待所有协程完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者协程：生成1到10的整数并发送到通道
	go func() {
		defer wg.Done()
		defer close(ch) // 发送完成后关闭通道

		for i := 1; i <= 10; i++ {
			ch <- i
			fmt.Printf("生产者发送: %d\n", i)
		}
		fmt.Println("生产者完成发送")
	}()

	// 消费者协程：从通道接收整数并打印
	go func() {
		defer wg.Done()

		for num := range ch {
			fmt.Printf("消费者接收并打印: %d\n", num)
		}
		fmt.Println("消费者完成接收")
	}()

	// 等待两个协程都完成
	wg.Wait()
	fmt.Println("程序结束")

	// 题目二 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印
	// 创建一个缓冲通道，缓冲区大小为10
	ch2 := make(chan int, 10)

	// 使用WaitGroup等待所有协程完成
	var wg2 sync.WaitGroup

	// 生产者协程
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		defer close(ch2) // 生产完成后关闭通道

		for i := 1; i <= 100; i++ {
			ch2 <- i // 发送数据到通道
			fmt.Printf("生产者: 发送 %d\n", i)
		}
		fmt.Println("生产者: 已完成所有发送")
	}()

	// 消费者协程
	wg2.Add(1)
	go func() {
		defer wg2.Done()

		for num := range ch2 { // 从通道接收数据，直到通道关闭
			fmt.Printf("消费者: 接收 %d\n", num)
		}
		fmt.Println("消费者: 已完成所有接收")
	}()

	// 等待所有协程完成
	wg2.Wait()
	fmt.Println("程序结束")
}
