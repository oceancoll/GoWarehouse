package main

import (
	"fmt"
	"time"
)

// 协程池
// 协程池最简要（核心）的逻辑是所有协程从任务读取任务，处理后把结果存放到结果队列。
func goroutine_pool_produce(num int) <-chan int{
	ch := make(chan int, 200)
	go func() {
		defer close(ch)
		for i:=0;i<num;i++{
			ch <- i
		}
	}()
	return ch
}

func goroutine_pool_workpool(poolnum int, jobch <- chan int, retch chan <- string) {
	for i:=0;i<poolnum;i++{
		go goroutine_pool_work(i, jobch, retch)
	}
}

func goroutine_pool_work(jobnum int, jobch <- chan int, retch chan <- string)  {
	cnt := 0
	for job := range jobch{
		cnt++
		ret := fmt.Sprintf("worker %d processed job: %d, it's the %dth processed by me.", jobnum, job, cnt)
		retch <- ret
	}
}
func main()  {
	ch := goroutine_pool_produce(1000)
	retch := make(chan string, 1000)
	goroutine_pool_workpool(5, ch, retch)
	time.Sleep(time.Second)
	close(retch)
	for n:= range retch{
		fmt.Println(n)
	}
}
