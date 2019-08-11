package main

import (
	"fmt"
	"time"
)

// 协程退出的方式，fo-range,range能够感知channel的关闭，当channel被发送数据的协程关闭时，range就会结束，接着退出for循环。
// 它在并发中的使用场景是：当协程只从1个channel读取数据，然后进行处理，处理后协程退出。下面这个示例程序，当in通道被关闭时，协程可自动退出。

func quit_range_producer(num int) <- chan int{
	ch := make(chan int)
	go func() {
		defer func() {
			close(ch)
			fmt.Println("producer exit")
		}()
		for i:=0;i<num;i++{
			fmt.Printf("send %d\n", i)
			ch <- i
			time.Sleep(time.Millisecond) //为了顺序输出
		}
	}()
	return ch
}

func quit_range_cal(ch <- chan int) <-chan struct{}{
	ch1 := make(chan struct{})
	go func() {
		defer func() {
			fmt.Println("worker exit")
			ch1 <- struct{}{}
			close(ch1)
		}()
		// for-range模式
		for n:= range ch{
			fmt.Printf("Process %d\n", n*n)
		}
	}()
	return ch1
}
func main()  {
	ch := quit_range_producer(5)
	ch1 := quit_range_cal(ch)
	<-ch1
	fmt.Println("main exit")
}