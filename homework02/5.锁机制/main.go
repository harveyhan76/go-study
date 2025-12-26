package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("锁机制示例")
	var (
		counter int
		mu      sync.Mutex // 互斥锁
		wg      sync.WaitGroup
	)

	// 启动10个协程
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				// 获取锁，保护对 counter 的访问
				mu.Lock()
				counter++
				mu.Unlock()

				// 注意：这里可以有其他不需要锁保护的操作
			}

			fmt.Printf("协程 %d 完成\n", id)
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 输出最终结果
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("期望值: %d\n", 10*1000)

	// 使用int64类型的原子计数器
	var counter2 int64
	var wg2 sync.WaitGroup

	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()

			// 每个goroutine递增1000次
			for j := 0; j < 1000; j++ {
				// 使用原子操作安全地递增计数器
				atomic.AddInt64(&counter2, 1)
			}

			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	// 等待所有goroutine完成
	wg2.Wait()

	// 使用原子操作读取最终值
	finalValue := atomic.LoadInt64(&counter2)
	fmt.Printf("Final counter value: %d\n", finalValue)
	fmt.Printf("Expected value: %d\n", 10*1000)
}
