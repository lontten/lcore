package perfutil

import (
	"fmt"
	"time"
)

var lastTime = time.Now()

// Mark 标记当前时间并打印与上一次标记的时间间隔
func Mark(msg string) {
	now := time.Now()
	duration := now.Sub(lastTime)
	fmt.Printf("[PERF] %s: %v\n", msg, duration)
	lastTime = now
}

// Reset 重置计时器
func Reset() {
	lastTime = time.Now()
}
