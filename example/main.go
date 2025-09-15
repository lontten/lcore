package main

import (
	"fmt"
	"time"

	"github.com/lontten/lcore"
)

func main() {
	line()
}
func line() {

	for range 22 {
		// 创建一个协程池，最大协程数为 3，任务队列大小为 5，使用 CallerRunsPolicy 拒绝策略
		pool := lcore.NewPool(10, 22, lcore.CallerRunsPolicy)
		defer pool.Shutdown()

		// 提交 10 个任务
		for i := 0; i < 100; i++ {
			taskID := i
			err := pool.Submit(func() {
				fmt.Printf("Task %d executed by worker\n", i)
				time.Sleep(200 * time.Millisecond)
			})
			if err != nil {
				fmt.Printf("Task %d rejected: %s\n", taskID, err)
			}
		}
		fmt.Println("All tasks completed")
	}

}
func lock() {
	list := make([]int, 0)
	kl := lcore.NewKeyLock(100)
	for i := range 100 {
		go func() {
			kl.Lock("num")
			list = append(list, i)
			kl.Unlock("num")
		}()
	}
	fmt.Println(list)
}
