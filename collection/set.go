package collections

type void struct{}

var member void

type Set[T comparable] map[T]void

func NewSet[T comparable]() Set[T] {
	return make(map[T]void)
}

func (s Set[T]) Add(elem T) {
	s[elem] = member
}

func (s Set[T]) Delete(elem T) {
	delete(s, elem)
}

func (s Set[T]) Slice() []T {
	r := make([]T, 0, len(s))
	for k := range s {
		r = append(r, k)
	}
	return r
}
