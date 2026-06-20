// Package slice provides a type-safe, generics-based Slice composite — a typed
// Array.
//
// Unlike the dynamic go-composites/array (a []interface{} requiring a cast on
// every read), Slice[T] is a []T. Map/Filter/Reduce/Find/Any/All are all
// statically checked: there is no interface{} in the public API and no runtime
// type assertion.
//
// Find returns an optional.Optional[T] rather than a (T, bool) pair so a
// "not found" result is a real Null-Object value, in keeping with the
// composition-oriented spirit.
package slice

import "github.com/go-composites/typed/src/optional"

// Slice is a typed, immutable-by-convention sequence of T.
type Slice[T any] []T

// Of builds a Slice from the given elements.
func Of[T any](items ...T) Slice[T] { return Slice[T](items) }

// Len reports the number of elements.
func (s Slice[T]) Len() int { return len(s) }

// Map applies f to every element, yielding a Slice[U]. Free function because
// it introduces a second type parameter.
func Map[T, U any](s Slice[T], f func(T) U) Slice[U] {
	out := make(Slice[U], len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

// Filter returns the elements for which keep reports true.
func (s Slice[T]) Filter(keep func(T) bool) Slice[T] {
	out := make(Slice[T], 0, len(s))
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reduce folds the slice left-to-right from an initial accumulator. Free
// function because the accumulator type A is independent of T.
func Reduce[T, A any](s Slice[T], initial A, combine func(A, T) A) A {
	acc := initial
	for _, v := range s {
		acc = combine(acc, v)
	}
	return acc
}

// Find returns the first element satisfying pred, wrapped in an Optional —
// optional.None when nothing matches.
func (s Slice[T]) Find(pred func(T) bool) optional.Optional[T] {
	for _, v := range s {
		if pred(v) {
			return optional.Some(v)
		}
	}
	return optional.None[T]()
}

// Any reports whether at least one element satisfies pred.
func (s Slice[T]) Any(pred func(T) bool) bool {
	for _, v := range s {
		if pred(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies pred (vacuously true for empty).
func (s Slice[T]) All(pred func(T) bool) bool {
	for _, v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}
