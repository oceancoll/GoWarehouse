package main

import (
	"fmt"
)

//FAN-IN-OUT模式，扇入扇出模式
// 计算输入数的平方

func fan_producer(num int) <- chan int{
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i:=0;i<num;i++{
			ch <- i
		}
	}()
	return ch
}

func fan_cal(ch <-chan int) <-chan int{
	ch1 := make(chan int)
	go func() {
		defer close(ch1)
		for n:= range ch{
			ch1 <- n*n
		}
	}()
	return ch1
}

// 方法1：
func fan_merge(chs ...<-chan int) <-chan int{
	ch1 := make(chan int)
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

// 方法2：
//func fan_merge(chs ...<-chan int) <-chan int{
//	ch1 := make(chan int)
//	wg := sync.WaitGroup{}
//	collect := func(ch <-chan int) {
//		defer wg.Done()
//		for n:= range ch{
//			ch1 <- n
//		}
//	}
//	wg.Add(len(chs))
//	for _,n:=range chs{
//		go collect(n)
//	}
//
//	// 错误写法：close是同步操作，会先把ch1关闭掉，无法写入，陷入死锁
//	//	wg.Wait()
//	//	close(ch1)
//
//	go func() {
//		wg.Wait()
//		close(ch1)
//	}()
//	return ch1
//}

func main()  {
	ch := fan_producer(10)

	//FAN-OUT
	ch1 := fan_cal(ch)
	ch2 := fan_cal(ch)
	ch3 := fan_cal(ch)

	//FAN_IN
	for n:= range fan_merge(ch1, ch2, ch3){
		fmt.Println(n)
	}
}