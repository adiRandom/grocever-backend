package types

type Set[T comparable] interface {
	Add(element T)
	Remove(element T)
	Contains(element T) bool
	Size() int
	ToSlice() []T
	AddAll(elements []T)
}
