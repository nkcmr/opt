package opt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Join allows for two Options to be used to create a new value if they are both
// present. If either is not present, then a None[R] will be returned.
func Join[A, B, R any](a Option[A], b Option[B], joinfn func(A, B) R) Option[R] {
	if a.Some() && b.Some() {
		return Some(joinfn(a.Unwrap(), b.Unwrap()))
	}
	return None[R]()
}

// Map allows a function to be run on the present value of an option if it is
// actually present and then optionally return something else from that value.
func Map[I, O any](in Option[I], mapfn func(I) Option[O]) Option[O] {
	if in.Some() {
		return mapfn(in.Unwrap())
	}
	return None[O]()
}

// Coalesce will take 0 or more Option and will return the first one that is
// Some value.
func Coalesce[T any](os ...Option[T]) Option[T] {
	for i := range os {
		if os[i].Some() {
			return os[i]
		}
	}
	return None[T]()
}

// Equal will compare the value in two options and check if their equal. If both
// are none, that is interpretted as "equal."
func Equal[T comparable](a, b Option[T]) bool {
	if a.None() && b.None() {
		return true
	}
	if a.None() || b.None() {
		return false
	}
	return a.Unwrap() == b.Unwrap()
}

// FromPointer will take in a pointer to a value and dereference it if it is
// not nil and return a Some[T](), if it is nil it will return None[T]().
func FromPointer[T any](v *T) Option[T] {
	if v == nil {
		return Option[T]{}
	}
	return Some(*v)
}

// FromMaybe will take a tuple of a value and a bool representing the value's
// status and convert it to an Option. If true, a Some[T] is returned, otherwise
// None[T] is returned.
func FromMaybe[T any](v T, ok bool) Option[T] {
	if ok {
		return Some(v)
	}
	return None[T]()
}

// None will return an Option[T] that has no value
func None[T any]() Option[T] {
	return Option[T]{}
}

// Some will return an Option[T] that contains the given value
func Some[T any](v T) Option[T] {
	return Option[T]{ok: true, v: v}
}

// Option represents an optional value. Every Option has either has something or
// has none.
// If None() reports true, then calls to Unwrap() will panic. To prevent this,
// Option[T] should be treated like a pointer and Unwrap() like a dereference.
//
// UnwrapOr() is a handy function to reduce the verbosity of checking an
// Option[T]'s state. UnwrapOr() will return the contained value if present, OR
// it will return the value provided.
//
// Option[T] is immutable once created.
//
// The zero-value of Option[T] is safe and will just report None() => true
//
// Inspired by Rust's Option<T>: https://doc.rust-lang.org/std/option/index.html
type Option[T any] struct {
	ok bool
	v  T
}

// Some reports whether there is a value contained or not. A returned
// value of `true` means there is a value and it can be retrieved with
// Unwrap() without it panicking.
// A returned value of `false` means there is no value and any call to
// Unwrap() will cause a panic.
func (o Option[T]) Some() bool {
	return o.ok
}

// None is just the opposite of Some(). True means Unwrap() panicks.
// False means Unwrap() will not panic.
func (o Option[T]) None() bool {
	return !o.ok
}

// Unwrap retrieves the underlying value if there is one. Unwrap WILL PANIC
// if there is no value.
func (o Option[T]) Unwrap() T {
	if o.ok {
		return o.v
	}
	panic(fmt.Sprintf("%T.Unwrap: no value to unwrap", o))
}

// UnwrapOr is a safer version of Unwrap() that will return the provided
// fallback value if the Option[T] does not contain a value.
func (o Option[T]) UnwrapOr(v T) T {
	if !o.ok {
		return v
	}
	return o.v
}

// MaybeUnwrap allows the underlying value to be retrieved in a more idiomatic
// way by returning a tuple of the possible underlying value and a bool that
// will be `true` if the value was present. The returned value will be the
// zeroed value of T if it is not present.
func (o Option[T]) MaybeUnwrap() (T, bool) {
	if o.ok {
		return o.v, true
	}
	var zv T
	return zv, false
}

// UnwrapOrZero allows the either the underlying value or a zeroed value of T
// to be returned. If the value is present it will be returned, otherwise a zero
// value will be returned.
func (o Option[T]) UnwrapOrZero() T {
	if o.ok {
		return o.v
	}
	var zv T
	return zv
}

// MarshalJSON implements json.Marshaler
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.ok {
		return json.Marshal(o.v)
	}
	return json.RawMessage("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var v T
	o.v = v
	if bytes.Equal(data, []byte("null")) {
		o.ok = false
		return nil
	}
	o.ok = true
	return json.Unmarshal(data, &o.v)
}
