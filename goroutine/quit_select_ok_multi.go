package main

import (
	"fmt"
	"time"
)

// 多通道
//协程的退出方式，select_ok,使用,ok来检测通道的关闭.
//如果某个通道关闭了，不再处理该通道，而是继续处理其他case，退出是等待所有的可读通道关闭。
// 我们需要使用select的一个特征：select不会在nil的通道上进行等待。这种情况，把只读通道设置为nil即可解决。

func quit_select_ok_multi_producer(num int) <- chan int {
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

func quit_select_ok_multi_cal(ch <- chan int, ch1 <- chan int) <-chan struct{}{
	ch2 := make(chan struct{})
	t := time.NewTicker(time.Millisecond * 500)
	processedCnt := 0

	go func() {
		defer func() {
			fmt.Println("worker exit")
			ch2 <- struct{}{}
			close(ch2)
		}()
		for {
			select {
			case n, ok:= <-ch:
				if !ok{
					ch = nil
				}
				processedCnt += 1
				fmt.Printf("Process %d\n", n)
			case k, ok:= <-ch1:
				if !ok{
					ch1 = nil
				}
				processedCnt += 1
				fmt.Printf("Process1 %d\n", k)
			case <-t.C:
				fmt.Printf("Working, processedCnt = %d\n", processedCnt)
			}
			if ch == nil && ch1 == nil{
				fmt.Println("all done")
				return
			}
		}
	}()
	return ch2
}
func main()  {
	ch := quit_select_ok_multi_producer(5)
	ch1 := quit_select_ok_multi_producer(7)
	ch2 := quit_select_ok_multi_cal(ch, ch1)
	<-ch2
	fmt.Println("main exit")
}
