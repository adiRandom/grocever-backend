package impl

type IdSet[T comparable, R any] struct {
	elements map[T]R
	id       func(R) T
}

func NewIdSet[T comparable, R any](id func(R) T) *IdSet[T, R] {
	return &IdSet[T, R]{elements: make(map[T]R), id: id}
}

func (s *IdSet[T, R]) Add(element R) {
	s.elements[s.id(element)] = element
}

func (s *IdSet[T, R]) Remove(element R) {
	delete(s.elements, s.id(element))
}

func (s *IdSet[T, R]) Contains(element R) bool {
	_, ok := s.elements[s.id(element)]
	return ok
}

func (s *IdSet[T, R]) Size() int {
	return len(s.elements)
}

func (s *IdSet[T, R]) ToSlice() []R {
	slice := make([]R, 0, len(s.elements))
	for _, element := range s.elements {
		slice = append(slice, element)
	}
	return slice
}

func (s *IdSet[T, R]) AddAll(elements []R) {
	for _, element := range elements {
		s.Add(element)
	}
}
