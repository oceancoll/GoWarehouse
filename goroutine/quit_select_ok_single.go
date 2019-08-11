package main

import (
	"fmt"
	"time"
)

// 单通道
//协程的退出方式，select_ok,使用,ok来检测通道的关闭.
//如果某个通道关闭后，需要退出协程，直接return即可。示例代码中，该协程需要从in通道读数据，还需要定时打印已经处理的数量，
// 有2件事要做，所有不能使用for-range，需要使用for-select，当in关闭时，ok=false，我们直接返回。

func quit_select_ok_single_producer(num int) <- chan int {
	ch := make(chan int)
	go func() {
		defer func() {
			close(ch)
			fmt.Println("producer exit")
		}()
		for i:=0; i<num; i++{
			fmt.Printf("send %d\n", i)
			ch <- i
			time.Sleep(time.Second)
		}
	}()
	return ch
}

func quit_select_ok_single_cal(ch <- chan int) <-chan struct{}{
	ch1 := make(chan struct{})
	t := time.NewTicker(time.Millisecond * 500)
	processedCnt := 0

	go func() {
		defer func() {
			fmt.Println("worker exit")
			ch1 <- struct{}{}
			close(ch1)
		}()
		for {
			select {
			case n, ok:= <-ch:
				if !ok{
					return
				}
				processedCnt += 1
				fmt.Printf("Process %d\n", n)
			case <-t.C:
				fmt.Printf("Working, processedCnt = %d\n", processedCnt)
			}
		}
	}()
	return ch1
}
func main()  {
	ch := quit_select_ok_single_producer(5)
	ch1 := quit_select_ok_single_cal(ch)
	<-ch1
	fmt.Println("main exit")
}
