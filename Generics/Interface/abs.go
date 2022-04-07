package Interface

type NumberAbs[T C] interface {
	Abs() T
}

type C interface {
	~int | ~int64
}

type DDD[T C] struct {
	Number T
}

func (d *DDD[T]) Abs() T {
	return d.Number + d.Number
}

func AbsDifference[T C](a, b T) T {
	c := a - b
	d := &DDD[T]{
		Number: c,
	}
	return d.Abs()
}
