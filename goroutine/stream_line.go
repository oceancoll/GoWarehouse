package main

import "fmt"

// 流水线模式
// 计算输入数的平方

// 切片
func stream_line_producer(num int) <- chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i:=0; i<num; i++{
			ch <- i
		}
	}()
	return ch
}

// 计算平方数
func stream_line_cal(ch <-chan int) <-chan int{
	ch1 := make(chan int)
	go func() {
		defer close(ch1)
		for n:= range ch{
			ch1 <- n*n
		}
	}()
	return ch1
}
func main()  {
	ch := stream_line_producer(5)
	ch1 := stream_line_cal(ch)
	for n:= range ch1{
		fmt.Println(n)
	}
}
