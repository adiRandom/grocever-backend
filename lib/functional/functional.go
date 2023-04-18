package functional

import (
	types "lib/helpers"
)

type ReduceCb[T, R any] func(acc R, current T) R

func Reduce[T any, R any](s []T, fn ReduceCb[T, R], initialValue R) R {
	partial := initialValue
	for _, el := range s {
		partial = fn(partial, el)
	}

	return partial
}

func Accumulate[T any](s []T, fn ReduceCb[T, T]) (T, error) {
	var partial T
	if len(s) == 0 {
		return partial, types.Error{"Empty slice", ""}
	}
	partial = (s)[0]
	for i, el := range s {
		if i == 0 {
			continue
		}
		partial = fn(partial, el)
	}

	return partial, nil
}

func Map[T any, R any](s []T, fn func(T) R) []R {
	mapped := make([]R, len(s))
	for index, el := range s {
		mapped[index] = fn(el)
	}

	return mapped
}

func IndexedMap[T any, R any](s []T, fn func(int, T) R) []R {
	mapped := make([]R, len(s))
	for i, el := range s {
		mapped[i] = fn(i, el)
	}

	return mapped
}

func Filter[T any](s []T, fn func(T) bool) []T {
	var filtered []T
	for _, el := range s {
		if fn(el) {
			filtered = append(filtered, el)
		}
	}

	return filtered
}

func IndexedFilter[T any](s []T, fn func(int, T) bool) []T {
	var filtered []T
	for i, el := range s {
		if fn(i, el) {
			filtered = append(filtered, el)
		}
	}

	return filtered
}

func Keys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func Exists[T any](s []T, fn func(T) bool) bool {
	for _, el := range s {
		if fn(el) {
			return true
		}
	}

	return false
}

func Find[T any](s []T, fn func(T) bool) *T {
	for _, el := range s {
		if fn(el) {
			return &el
		}
	}

	return nil
}

func GroupBy[T any, S comparable](s []T, fn func(T) S) map[S][]T {
	grouped := make(map[S][]T)
	for _, el := range s {
		key := fn(el)
		grouped[key] = append(grouped[key], el)
	}

	return grouped
}

func GroupByPointer[T any, S comparable](s []T, fn func(T) S) map[S][]*T {
	grouped := make(map[S][]*T)
	for _, el := range s {
		key := fn(el)
		grouped[key] = append(grouped[key], &el)
	}

	return grouped
}

func Reverse[T any](s []T) []T {
	reversed := make([]T, len(s))
	for i, el := range s {
		reversed[len(s)-1-i] = el
	}

	return reversed
}
