package impl

type Queue[T any] struct {
	elements []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: make([]T, 0)}
}

func (q *Queue[T]) Push(element T) {
	q.elements = append(q.elements, element)
}

func (q *Queue[T]) Pop() T {
	element := q.elements[0]
	q.elements = q.elements[1:]
	return element
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue[T]) Top() T {
	return q.elements[0]
}
