package main

import (
	"fmt"
	"sync"
	"time"
)

type Bank struct {
	sync.Mutex
	saving map[string] int
}

// 存钱
func (b *Bank) Deposit (name string, amount int) (bool, int) {
	b.Lock()
	defer b.Unlock()
	if amount<=0{
		return false, 0
	}
	if _,ok:=b.saving[name];!ok{
		b.saving[name] = 0
	}
	b.saving[name] += amount
	return true, b.saving[name]
}

// 取钱
func (b *Bank) Withdraw(name string, amount int) (bool, int) {
	b.Lock()
	defer b.Unlock()
	if amount<0 {
		return false, 0
	}
	if _, ok := b.saving[name]; !ok{
		return false, 0
	}
	
	if b.saving[name] < amount{
		return false, b.saving[name]
	}
	
	return true, b.saving[name]-amount

}

// 查询余额
func (b *Bank)Query(name string) (bool, int)  {
	b.Lock()
	defer b.Unlock()
	if _, ok:=b.saving[name];!ok{
		return false, 0
	}
	return true, b.saving[name]
}

func NewBank() *Bank {
	b := &Bank{
		saving:make(map[string]int),
	}
	return b
}
func main()  {
	b := NewBank()
	go b.Deposit("xiaoming", 200)
	go b.Withdraw("xiaoming", 1400)
	go b.Withdraw("xiaoming", 100)
	time.Sleep(time.Second)
	fmt.Println(b.Query("xiaoming"))
}

