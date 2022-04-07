package main

import (
	"Generics/Interface"
	"Generics/Struct"
	"Generics/demo"
	"Generics/function"
	"Generics/service"
	"fmt"
)

func init() {
	service.SetUp()
}

func main() {
	// func
	// []int32
	int32Slice := []int32{1, 2, 3}
	fmt.Println(function.Sum(int32Slice))
	fmt.Println(function.Sum1(int32Slice))
	fmt.Println(function.Sum1[int32](int32Slice)) // 指定func使用的类型
	// []int64
	int64Slice := []int64{1, 2, 3}
	fmt.Println(function.Sum(int64Slice))
	fmt.Println(function.Sum1(int64Slice))
	fmt.Println(function.Sum1[int64](int64Slice)) // 指定func使用的类型
	m := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	fmt.Println(function.Sum2(m))

	// Map
	// double
	fmt.Println(function.Map[int, int]([]int{1, 2, 3}, func(t1 int) int {
		return t1 * t1
	}))

	// Reduce
	// add
	fmt.Println(function.Reduce[int, int]([]int{1, 2, 3, 4}, 5, func(t1, t2 int) int {
		return t1 + t2
	}))

	// struct
	p1 := Struct.People[int64]{}
	p1.Push(1)
	p1.Push(2)
	p1.Push(3)
	fmt.Println(p1.Pop())
	fmt.Println(p1.Pop())
	fmt.Println(p1.Pop())

	p2 := Struct.People[float64]{}
	p2.Push(1.1)
	p2.Push(2.2)
	p2.Push(3.3)
	fmt.Println(p2.Pop())
	fmt.Println(p2.Pop())
	fmt.Println(p2.Pop())

	// interface
	animal := service.AnimalService.GetAge(19)
	fmt.Println(animal)
	name := service.AnimalService.GetName("abc")
	fmt.Println(name)
	// sort
	sortData := []float64{6.6, 1.1, 4.4}
	demo.Sort(sortData, func(a, b float64) bool {
		if a < b {
			return true
		}
		return false
	})
	fmt.Println(sortData)

	a := Interface.DDD[int]{5}
	fmt.Println(a.Abs())
	b := Interface.AbsDifference[int64](int64(1), int64(3))
	fmt.Println(b)
}
