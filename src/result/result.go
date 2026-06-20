// Package result provides a type-safe, generics-based Result composite.
//
// Unlike the dynamic go-composites/result (whose Payload() returns
// interface{}), this Result[T] carries a statically-typed payload. The
// compiler verifies that producers and consumers agree on T, catching the
// class of type mismatches the reflective variant can only discover at run
// time.
//
// Null-Object spirit: Ok and Err are real values, never nil. Err carries a
// non-nil error; Ok carries the (typed) payload. Methods are total — they
// never panic and never return a bare nil interface.
package result

import "errors"

// Result is the type-safe railway value: it is either a success carrying a
// payload of type T, or a failure carrying an error.
type Result[T any] interface {
	// HasError reports whether this Result is the failure case.
	HasError() bool
	// Payload returns the success value and true, or the zero value and
	// false when this Result is a failure.
	Payload() (T, bool)
	// Error returns the carried error, or nil when this Result is a success.
	Error() error
}

// result is the single concrete implementation. The ok flag distinguishes the
// two cases without resorting to a nil interface for either branch.
type result[T any] struct {
	value T
	err   error
	ok    bool
}

// Ok constructs a successful Result carrying v.
func Ok[T any](v T) Result[T] {
	return result[T]{value: v, ok: true}
}

// Err constructs a failed Result carrying e. A nil e is normalised to a
// non-nil sentinel so a failure is never silently indistinguishable from a
// success and Error() never returns nil on the failure branch.
func Err[T any](e error) Result[T] {
	if e == nil {
		e = errors.New("result: nil error")
	}
	return result[T]{err: e, ok: false}
}

func (r result[T]) HasError() bool { return !r.ok }

func (r result[T]) Payload() (T, bool) {
	if r.ok {
		return r.value, true
	}
	var zero T
	return zero, false
}

func (r result[T]) Error() error { return r.err }

// Map transforms the payload of a successful Result with f, producing a
// Result[U]. A failure is propagated unchanged (railway short-circuit).
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if v, ok := r.Payload(); ok {
		return Ok[U](f(v))
	}
	return Err[U](r.Error())
}

// FlatMap chains a fallible operation: f is invoked only on the success
// branch and may itself fail. Failures short-circuit. This is the monadic
// bind for Result.
func FlatMap[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if v, ok := r.Payload(); ok {
		return f(v)
	}
	return Err[U](r.Error())
}

// AndThen is a railway-oriented alias for FlatMap.
func AndThen[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	return FlatMap(r, f)
}
