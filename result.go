package result

import "errors"

type Result[T any] struct {
	data T
	err  error
}

// Ok returns a Result which contains the success value.
func Ok[T any](result T) Result[T] {
	return Result[T]{
		data: result,
		err:  nil,
	}
}

// Err returns a Result which contains the error value.
func Err[T any](msg string) Result[T] {
	var t T
	return Result[T]{
		data: t,
		err:  errors.New(msg),
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
		return Err[U](r.err.Error())
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
	return Err[T](f(r.err).Error())
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
