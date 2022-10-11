package functional

import (
	types "dealScraper/lib/helpers"
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
	var mapped []R
	for _, el := range s {
		mapped = append(mapped, fn(el))
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
