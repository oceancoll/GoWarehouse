package Struct

type C interface {
	~int32 | ~int64 | ~string | ~float64
}

type People[T C] struct {
	size  int
	value []T
}

func (s *People[T]) Push(v T) {
	s.value = append(s.value, v)
	s.size++
}

func (s *People[T]) Pop() T {
	val := s.value[s.size-1]
	if s.size > 0 {
		s.value = s.value[:s.size-1]
		s.size--
	}
	return val
}
