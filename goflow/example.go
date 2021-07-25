package main

import (
	"fmt"
	"goflow/flow"
	"time"
)

func main() {
	f1 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function1 started")
		time.Sleep(time.Millisecond * 1000)
		return 1, nil
	}

	f2 := func(r map[string]interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function2 started", r["f1"])
		return "some results", nil // errors.New("Some error")
	}

	f3 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function3 started", r["f1"])
		return nil, nil
	}

	f4 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function4 started", r)
		return nil, nil
	}

	res, err := flow.New().
		Add("f1", nil, f1). // 顶层节点，不依赖上游数据
		Add("f2", []string{"f1"}, f2). // f2的执行依赖f1的结果
		Add("f3", []string{"f1"}, f3). // f3的执行依赖f1的结果
		Add("f4", []string{"f2", "f3"}, f4). // f4的执行依赖f2,f3的结果
		Do()

	fmt.Println(res, err)
}

