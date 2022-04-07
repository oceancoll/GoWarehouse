package demo

import (
	"sort"
)

type C interface {
	~int | ~int32 | ~int64 | ~float64
}
type warpSort[T C] struct {
	s   []T
	cmp func(a, b T) bool
}

func (s warpSort[T]) Len() int {
	return len(s.s)
}
func (s warpSort[T]) Less(a, b int) bool {
	return s.cmp(s.s[a], s.s[b])
}
func (s warpSort[T]) Swap(a, b int) {
	s.s[a], s.s[b] = s.s[b], s.s[a]
}
func Sort[T C](data []T, cmp func(a, b T) bool) {
	sort.Sort(warpSort[T]{data, cmp})
}
