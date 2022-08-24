package result

import "errors"

type Result[T any] struct {
	data T
	err  error
}

// Ok returns a Result which contains the success value.
func Ok[T any](data T) Result[T] {
	return Result[T]{
		data: data,
		err:  nil,
	}
}

// Err returns a Result which contains the error value.
func Err[T any](err error) Result[T] {
	var t T
	return Result[T]{
		data: t,
		err:  err,
	}
}

// IsOk returns `true` if the result is `Ok`.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsOkAnd returns `true` if the result is `Ok` and the
// value inside of it matches a predicate.
func (r Result[T]) IsOkAnd(predicate func(T) bool) bool {
	if r.IsErr() {
		return false
	}
	return predicate(r.data)
}

// IsErr returns `true` if the result is `Err`.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// IsErrAnd returns `true` if the result is `Err` and the
// value inside of it matches a predicate.
func (r Result[T]) IsErrAnd(predicate func(error) bool) bool {
	if r.IsOk() {
		return false
	}
	return predicate(r.err)
}

// Map maps a Result[T] to Result[U] by applying a function
// to a contained `Ok` value, leaving an `Err` value untouched.
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	return Ok(f(r.data))
}

// MapOr returns the provided default (if `Err`), or applies
// a function to the contained value (if `Ok`).
func MapOr[T any, U any](r Result[T], fallback U, f func(T) U) U {
	if r.IsErr() {
		return fallback
	}
	return f(r.data)
}

// MapOr returns the called fallback function (if `Err`),
// or applies a function to the contained value (if `Ok`).
func MapOrElse[T any, U any](r Result[T], fallbackFn func(error) U, f func(T) U) U {
	if r.IsErr() {
		return fallbackFn(r.err)
	}
	return f(r.data)
}

// MapErr applies a function to the contained `Err` value,
// leaving an `Ok` value untouched.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsOk() {
		return r
	}
	return Err[T](f(r.err))
}

// Inspect calls the provided closure with the contained
// value (if `Ok`) and returns the unchanged Result.
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.data)
	}
	return r
}

// InspectErr calls the provided closure with the contained
// error (if `Err`) and returns the unchanged Result.
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	if r.IsErr() {
		f(r.err)
	}
	return r
}

// Expect returns the `Ok` value, or panics with
// the given message if `Err`.
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		panic(msg)
	}
	return r.data
}

// Unwrap returns the `Ok` value, or panics if `Err`.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic("Result is Err")
	}
	return r.data
}

// UnwrapOrDefault returns the `Ok` value, or the
// default value of type T if `Err`.
func (r Result[T]) UnwrapOrDefault() T {
	if r.IsErr() {
		var t T
		return t
	}
	return r.data
}

// ExpectErr returns the contained `Err` value, or
// panics if `Ok`.
func (r Result[T]) ExpectErr(msg string) error {
	if r.IsOk() {
		panic(msg)
	}
	return r.err
}

// UnwrapErr returns the contained `Err` or panics if `Ok`.
func (r Result[T]) UnwrapErr() error {
	if r.IsOk() {
		panic("Result is `Ok`")
	}
	return r.err
}

// And returns `other` if the first result is `Ok`, otherwise
// returns the `Err` value of the first result.
func And[T any, U any](r Result[T], other Result[U]) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	return other
}

// AndThen calls `f` if the result is `Ok`, otherwise returns
// the `Err` value of the given result. Also known as monadic bind.
func AndThen[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	return f(r.data)
}

// Or returns the result if it is `Ok`, otherwise returns `other`.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return other
}

// OrElse returns the result if it is `Ok`, otherwise
// returns the result of `f` applied with the `Err` value.
func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return f(r.err)
}

// UnwrapOr returns the contained `Ok` value or a provided default.
func (r Result[T]) UnwrapOr(fallback T) T {
	if r.IsErr() {
		return fallback
	}
	return r.data
}

// UnwrapOrElse returns the contained `Ok` value
// or computes it from a closure.
func (r Result[T]) UnwrapOrElse(fallbackFn func(error) T) T {
	if r.IsErr() {
		return fallbackFn(r.err)
	}
	return r.data
}

// Contains returns `true` if the result is an `Ok` value
// containing the given value.
func Contains[T comparable](r Result[T], x T) bool {
	if r.IsErr() {
		return false
	}
	return r.data == x
}

// ContainsErr return `true` if the result is an `Err` value
// containing the given error.
func (r Result[T]) ContainsErr(err error) bool {
	if r.IsOk() {
		return false
	}
	return errors.Is(r.err, err)
}

// Copy returns a value copy of the result.
func (r Result[T]) Copy() Result[T] {
	return r
}

// Flatten converts from a Result[Result[T]] to a [Result[T]]
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsErr() {
		return Err[T](r.err)
	}
	return r.data
}

// Equal tests deep equality of two results.
func Equal[T comparable](r, other Result[T]) bool {
	if r.IsOk() && other.IsOk() {
		return r.data == other.data
	}
	if r.IsErr() && other.IsErr() {
		return r.err.Error() == other.err.Error()
	}
	return false
}
