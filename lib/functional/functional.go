package functional

type MappedSlice[T, R any] []T
type ReduceCb[T, R any] func(acc R, current T) R

type Error string

func (e Error) Error() string {
	return string(e)
}

func (s *MappedSlice[T, R]) Reduce(fn ReduceCb[T, R], initialValue R) R {
	partial := initialValue
	for _, el := range *s {
		partial = fn(partial, el)
	}

	return partial
}

func (s *MappedSlice[T, T]) Accumulate(fn ReduceCb[T, T]) (T, error) {
	var partial T
	if len(*s) == 0 {
		return partial, Error("Empty slice")
	}
	partial = (*s)[0]
	for i, el := range *s {
		if i == 0 {
			continue
		}
		partial = fn(partial, el)
	}

	return partial, nil
}
