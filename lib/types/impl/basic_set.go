package impl

type BasicSet[T comparable] struct {
	elements map[T]bool
}

func NewBasicSet[T comparable]() *BasicSet[T] {
	return &BasicSet[T]{elements: make(map[T]bool)}
}

func (s *BasicSet[T]) Add(element T) {
	s.elements[element] = true
}

func (s *BasicSet[T]) Remove(element T) {
	delete(s.elements, element)
}

func (s *BasicSet[T]) Contains(element T) bool {
	_, ok := s.elements[element]
	return ok
}

func (s *BasicSet[T]) Size() int {
	return len(s.elements)
}

func (s *BasicSet[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.elements))
	for element := range s.elements {
		slice = append(slice, element)
	}
	return slice
}

func (s *BasicSet[T]) AddAll(elements []T) {
	for _, element := range elements {
		s.Add(element)
	}
}
