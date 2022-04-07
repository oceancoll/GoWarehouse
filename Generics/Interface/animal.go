package Interface

// 接口支持自定义传参
type animal[T any] struct {
}

type Animal[T any] interface {
	GetAge(age int) int
	GetName(name T) T
}

// 实例化结构体时，需要指定可变参数类型。类似int/float这种。但是any也支持
func NewAnimalService() *animal[any] {
	return &animal[any]{}
}

func (s *animal[T]) GetAge(age int) int {
	return age
}

func (s *animal[T]) GetName(name T) T {
	return name
}
