package lcore

import (
	"errors"
	"sync"
)

// Task 定义任务类型
type Task func()

// RejectPolicy 定义拒绝策略类型
type RejectPolicy func(task Task, pool *Pool)

// Pool 协程池结构体
type Pool struct {
	maxWorkers   int           // 最大工作协程数
	tasks        chan Task     // 任务队列
	rejectPolicy RejectPolicy  // 拒绝策略
	stop         chan struct{} // 停止信号
	wg           sync.WaitGroup
}

// NewPool 创建一个新的协程池
func NewPool(maxWorkers int, queueSize int, rejectPolicy RejectPolicy) *Pool {
	p := &Pool{
		maxWorkers:   maxWorkers,
		tasks:        make(chan Task, queueSize),
		rejectPolicy: rejectPolicy,
		stop:         make(chan struct{}),
	}
	p.wg.Add(maxWorkers)

	for i := 0; i < p.maxWorkers; i++ {
		go p.worker()
	}
	return p
}

// worker 工作协程
func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task()
		case <-p.stop:
			return
		}
	}
}

// Submit 提交任务
func (p *Pool) Submit(task Task) error {
	select {
	case p.tasks <- task:
		return nil
	default:
		if p.rejectPolicy != nil {
			p.rejectPolicy(task, p)
			return nil
		}
		return errors.New("task queue is full and no reject policy specified")
	}
}

// Shutdown 关闭协程池
func (p *Pool) Shutdown() {
	close(p.stop)  // 发送停止信号
	p.wg.Wait()    // 等待所有任务完成
	close(p.tasks) // 关闭任务队列
}

// 拒绝策略示例

// AbortPolicy 直接拒绝任务并抛出错误
func AbortPolicy(task Task, pool *Pool) {
}

// CallerRunsPolicy 由提交任务的 Goroutine 自己执行任务
func CallerRunsPolicy(task Task, pool *Pool) {
	task()
}

// DiscardPolicy 直接丢弃任务
func DiscardPolicy(task Task, pool *Pool) {
}

// DiscardOldestPolicy 丢弃队列中最老的任务，然后重新提交新任务
func DiscardOldestPolicy(task Task, pool *Pool) {
	for {
		select {
		case <-pool.tasks: // 丢弃最老的任务
		default:
		}

		// 尝试提交新任务
		select {
		case pool.tasks <- task:
			return // 提交成功，退出循环
		default:
			// tasks 已满，继续重试
		}
	}
}
