package function

// T 称为类型参数
// 紧跟在T后面的类型值称为类型限制
func Sum[T int | int32 | int64](nums []T) T {
	var res T
	for _, num := range nums {
		res += num
	}
	return res
}

type C interface {
	// 通过接口定义支持的类型，～T的作用是：包含底层类型T的所有类型集合
	~int | ~int32 | ~int64
}

func Sum1[T C](nums []T) T {
	var res T
	for _, num := range nums {
		res += num
	}
	return res
}

// comparable 是Go预声明的，表示任何能做 == 和 != 操作的类型
// type comparable interface{ comparable }
func Sum2[K comparable, V int32 | int64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func Map[T1, T2 interface{}](s []T1, f func(T1) T2) []T2 {
	r := make([]T2, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// type any = interface{}，any等同于interface{}，是个万能类型
func Reduce[T1, T2 any](s []T1, init T2, f func(T2, T1) T2) T2 {
	r := init
	for _, v := range s {
		r = f(r, v)
	}
	return r
}
