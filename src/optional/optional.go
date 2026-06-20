// Package optional provides a type-safe, generics-based Optional composite —
// the static-typing analogue of the Null-Object pattern.
//
// Unlike a bare *T or an interface{} that may be nil, an Optional[T] is always
// a real value. None and Some are both concrete; there is no nil to forget to
// check. The compiler tracks T, so OrElse / Map cannot be fed a value of the
// wrong type.
package optional

// Optional is either present (Some) carrying a T, or absent (None).
type Optional[T any] struct {
	value   T
	present bool
}

// Some constructs a present Optional carrying v.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, present: true}
}

// None constructs an absent Optional. It is a real value, never nil — the
// type-safe Null-Object.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// IsPresent reports whether a value is held.
func (o Optional[T]) IsPresent() bool { return o.present }

// Get returns the held value and true, or the zero value and false when absent.
func (o Optional[T]) Get() (T, bool) {
	if o.present {
		return o.value, true
	}
	var zero T
	return zero, false
}

// OrElse returns the held value when present, otherwise the supplied fallback.
func (o Optional[T]) OrElse(fallback T) T {
	if o.present {
		return o.value
	}
	return fallback
}

// Map transforms a present Optional[T] into an Optional[U]; an absent value is
// propagated as None[U]. It is a free function because Go methods cannot
// introduce new type parameters.
func Map[T, U any](o Optional[T], f func(T) U) Optional[U] {
	if v, ok := o.Get(); ok {
		return Some[U](f(v))
	}
	return None[U]()
}
