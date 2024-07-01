package Datatypes

type Set[T comparable] struct {
	list map[T]struct{} //empty structs occupy 0 memory
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set[T]) Add(v T) {
	s.list[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.list, v)
}

func (s *Set[T]) Clear() {
	s.list = make(map[T]struct{})
}

func (s *Set[T]) Size() int {
	return len(s.list)
}

func NewSet[T comparable]() *Set[T] {
	s := &Set[T]{}
	s.list = make(map[T]struct{})
	return s
}

func (s Set[T]) Get() []T {
	elements := make([]T, 0, len(s.list))
	for key := range s.list {
		elements = append(elements, key)
	}
	return elements
}
