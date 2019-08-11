package main

import (
	"fmt"
	"math/rand"
	"time"
)

func select_use_eat() chan string{
	out:= make(chan string)
	go func() {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		out <- "Mom call you eating"
		close(out)
	}()
	return out
}
func main()  {
	eatch := select_use_eat()
	sleep := time.NewTimer(time.Second*3)
	select {
	case s := <-eatch:
		fmt.Println(s)
	case <-sleep.C:
		fmt.Println("time to sleep")
	default:
		fmt.Println("beat doudou")
		time.Sleep(time.Second)
	}

}