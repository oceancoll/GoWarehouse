package main

import (
	"fmt"
	"time"
)

// 主动传递退出信号
// 协程的退出方式，select通过检测主动传入的退出信号，判断是否退出
// 使用一个专门的通道，发送退出的信号，可以解决这类问题。
// 以第2个场景为例，协程入参包含一个停止通道stopCh，当stopCh被关闭，case <-stopCh会执行，直接返回即可。
// 当我启动了100个worker时，只要main()执行关闭stopCh，
// 每一个worker都会都到信号，进而关闭。如果main()向stopCh发送100个数据，这种就低效了。

func quit_chan_signal_worke(stopCh <-chan struct{})  {
	go func() {
		defer fmt.Println("worker exit")

		t := time.NewTicker(time.Millisecond * 500)

		// Using stop channel explicit exit
		for {
			select {
			case <-stopCh:
				fmt.Println("Recv stop signal")
				return
			case <-t.C:
				fmt.Println("Working .")
			}
		}
	}()
	return
}
func main()  {
	stopch := make(chan struct{})
	quit_chan_signal_worke(stopch)

	time.Sleep(time.Second * 2)
	close(stopch)

	// wait some print
	time.Sleep(time.Second)
	fmt.Println("main exit")
}