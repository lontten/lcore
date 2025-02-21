package lcore

import (
	"errors"
	"sync"
	"time"
)

// Task 定义任务类型
type Task func()

// RejectPolicy 定义拒绝策略类型
type RejectPolicy func(task Task, pool *ThreadPool)

// ThreadPool 线程池结构体
type ThreadPool struct {
	maxWorkers   int           // 最大工作线程数
	taskQueue    chan Task     // 任务队列
	rejectPolicy RejectPolicy  // 拒绝策略
	stop         chan struct{} // 停止信号
	wg           sync.WaitGroup

	mainRun bool // 主线程是否运行任务中
}

// NewThreadPool 创建一个新的线程池
func NewThreadPool(maxWorkers int, queueSize int, rejectPolicy RejectPolicy) *ThreadPool {
	return &ThreadPool{
		maxWorkers:   maxWorkers,
		taskQueue:    make(chan Task, queueSize),
		rejectPolicy: rejectPolicy,
		stop:         make(chan struct{}),
	}
}

// Start 启动线程池
func (pool *ThreadPool) Start() {
	for i := 0; i < pool.maxWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}
}

// worker 工作线程
func (pool *ThreadPool) worker() {
	defer pool.wg.Done()
	for {
		select {
		case task, ok := <-pool.taskQueue:
			if !ok {
				return
			}
			task()
		case <-pool.stop:
			return
		}
	}
}

// Submit 提交任务
func (pool *ThreadPool) Submit(task Task) error {
	select {
	case pool.taskQueue <- task:
		return nil
	default:
		if pool.rejectPolicy != nil {
			pool.rejectPolicy(task, pool)
			return nil
		}
		return errors.New("task queue is full and no reject policy specified")
	}
}

// Shutdown 关闭线程池
func (pool *ThreadPool) Shutdown() {
	for len(pool.taskQueue) > 0 || pool.mainRun {
		time.Sleep(1 * time.Second)
	}
	close(pool.stop)      // 发送停止信号
	pool.wg.Wait()        // 等待所有任务完成
	close(pool.taskQueue) // 关闭任务队列
}

// 拒绝策略示例

// AbortPolicy 直接拒绝任务并抛出错误
func AbortPolicy(task Task, pool *ThreadPool) {
}

// CallerRunsPolicy 由提交任务的 Goroutine 自己执行任务
func CallerRunsPolicy(task Task, pool *ThreadPool) {
	pool.mainRun = true
	task()
	pool.mainRun = false
}

// DiscardPolicy 直接丢弃任务
func DiscardPolicy(task Task, pool *ThreadPool) {
}

// DiscardOldestPolicy 丢弃队列中最老的任务，然后重新提交新任务
func DiscardOldestPolicy(task Task, pool *ThreadPool) {
	for {
		select {
		case <-pool.taskQueue: // 丢弃最老的任务
		default:
		}

		// 尝试提交新任务
		select {
		case pool.taskQueue <- task:
			return // 提交成功，退出循环
		default:
			// taskQueue 已满，继续重试
		}
	}
}
