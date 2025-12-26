package main

import (
	"fmt"
	"sync"
	"time"
)

// 任务类型：一个无参数无返回值的函数
type Task func()

// 任务调度器
func scheduleTasks(tasks []Task) {
	var wg sync.WaitGroup
	results := make(chan string, len(tasks))

	// 为每个任务启动一个协程
	for i, task := range tasks {
		wg.Add(1)
		go func(taskID int, taskFunc Task) {
			defer wg.Done()

			// 记录开始时间
			start := time.Now()

			// 执行任务
			taskFunc()

			// 记录结束时间
			duration := time.Since(start)

			// 发送结果
			results <- fmt.Sprintf("任务%d 执行时间: %v", taskID+1, duration)
		}(i, task)
	}

	// 等待所有协程完成
	wg.Wait()
	close(results)

	// 打印所有结果
	fmt.Println("\n任务执行统计：")
	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	// 使用WaitGroup等待两个协程完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 协程1：打印奇数
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Printf("奇数: %d\n", i)
			}
		}
	}()

	// 协程2：打印偶数
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("偶数: %d\n", i)
			}
		}
	}()

	// 等待两个协程执行完毕
	wg.Wait()
	fmt.Println("所有数字打印完成！")

	//time.Sleep(time.Second)

	// 定义一些测试任务
	tasks := []Task{
		func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("任务1: 处理数据完成")
		},
		func() {
			time.Sleep(200 * time.Millisecond)
			fmt.Println("任务2: 发送网络请求完成")
		},
		func() {
			time.Sleep(150 * time.Millisecond)
			fmt.Println("任务3: 计算完成")
		},
		func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Println("任务4: 保存文件完成")
		},
	}

	fmt.Println("开始执行任务...")
	start := time.Now()

	// 调度执行所有任务
	scheduleTasks(tasks)

	// 计算总耗时
	totalTime := time.Since(start)
	fmt.Printf("\n所有任务完成，总耗时: %v\n", totalTime)

}
