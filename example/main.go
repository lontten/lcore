package main

import (
	"fmt"
	"github.com/lontten/lcore"
	"time"
)

func main() {
	for range 40 {
		// 创建一个线程池，最大线程数为 3，任务队列大小为 5，使用 CallerRunsPolicy 拒绝策略
		pool := lcore.NewThreadPool(3, 4, lcore.CallerRunsPolicy)
		pool.Start()
		defer pool.Shutdown()

		// 提交 10 个任务
		for i := 0; i < 10; i++ {
			taskID := i
			err := pool.Submit(func() {
				fmt.Printf("Task %d executed by worker\n", taskID)
				time.Sleep(2 * time.Second)
			})
			if err != nil {
				fmt.Printf("Task %d rejected: %s\n", taskID, err)
			}
		}
	}

}
