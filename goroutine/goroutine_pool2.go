package main

import (
	"fmt"
	"sync"
	"time"
)

// Pool 协程池
type Pool struct {
	TaskChannel chan func() // fuc类型任务队列
}

// NewPool 创建一个协程池
func NewPool(n int) *Pool {
	// 初始化 Pool.TaskChannel
	p := &Pool{
		TaskChannel: make(chan func()),
	}

	// 创建指定数量 worker 从任务队列取出任务执行
	for i := 0; i < n; i++ {
		go func() {
			for task := range p.TaskChannel {
				task() // 取出的即位 func 类型，直接加括号即运行
			}
		}()
	}
	return p
}

// Submit 提交任务
func (p *Pool) Submit(f func()) {
	p.TaskChannel <- f
}

func main() {
	p := NewPool(2)
	var wg sync.WaitGroup
	wg.Add(3)
	task1 := func() {
		fmt.Println("eat cost 3 seconds")
		time.Sleep(3 * time.Second)
		wg.Done()
	}
	task2 := func() {
		defer wg.Done()
		fmt.Println("wash feet cost 3 seconds")
		time.Sleep(3 * time.Second)
	}
	task3 := func() {
		fmt.Println("watch tv cost 3 seconds")
		time.Sleep(3 * time.Second)
		wg.Done()
	}
	p.Submit(task1)
	p.Submit(task2)
	p.Submit(task3)
	// 等待所有任务执行完成
	wg.Wait()
}
