package main

import (
	"fmt"
)

//FAN-IN-OUT模式，扇入扇出模式，优化版本，fan_in是程序的瓶颈(多数据写入)，增加缓冲区
// 计算输入数的平方

func fan_optimize_producer(num int) <- chan int{
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i:=0;i<num;i++{
			ch <- i
		}
	}()
	return ch
}

func fan_optimize_cal(ch <-chan int) <-chan int{
	ch1 := make(chan int)
	go func() {
		defer close(ch1)
		for n:= range ch{
			ch1 <- n*n
		}
	}()
	return ch1
}


func fan_optimize_merge(chs ...<-chan int) <-chan int{
	ch1 := make(chan int, 10)
	collect := func(ch <-chan int) {
		for n:= range ch{
			ch1 <- n
		}
	}
	go func() {
		defer close(ch1)
		for _, n:= range chs{
			collect(n)
		}
	}()
	return ch1
}


func main()  {
	ch := fan_optimize_producer(10)

	//FAN-OUT
	ch1 := fan_optimize_cal(ch)
	ch2 := fan_optimize_cal(ch)
	ch3 := fan_optimize_cal(ch)

	//FAN_IN
	for n:= range fan_optimize_merge(ch1, ch2, ch3){
		fmt.Println(n)
	}
}