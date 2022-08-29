package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/JustinKnueppel/go-result"
)

func TestIsOk(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[int]
		expected bool
	}{
		"success": {
			value:    result.Ok(1),
			expected: true,
		},
		"error": {
			value:    result.Err[int](errors.New("")),
			expected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsOk() != tc.expected {
				t.Fail()
			}
		})
	}
}

func TestIsOkAnd(t *testing.T) {
	tests := map[string]struct {
		value     result.Result[int]
		predicate func(int) bool
		expected  bool
	}{
		"success_true": {
			value:     result.Ok(1),
			predicate: func(i int) bool { return i == 1 },
			expected:  true,
		},
		"success_false": {
			value:     result.Ok(2),
			predicate: func(i int) bool { return i > 5 },
			expected:  false,
		},
		"error": {
			value:     result.Err[int](errors.New("")),
			predicate: func(i int) bool { return true },
			expected:  false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsOkAnd(tc.predicate) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestIsErr(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[int]
		expected bool
	}{
		"success": {
			value:    result.Ok(1),
			expected: false,
		},
		"error": {
			value:    result.Err[int](errors.New("")),
			expected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsErr() != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestIsErrAnd(t *testing.T) {
	tests := map[string]struct {
		value     result.Result[int]
		predicate func(error) bool
		expected  bool
	}{
		"success": {
			value:     result.Ok(1),
			predicate: func(err error) bool { return true },
			expected:  false,
		},
		"error_true": {
			value:     result.Err[int](errors.New("correct error")),
			predicate: func(err error) bool { return err.Error() == "correct error" },
			expected:  true,
		},
		"error_false": {
			value:     result.Err[int](errors.New("not this error")),
			predicate: func(err error) bool { return err.Error() == "" },
			expected:  false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsErrAnd(tc.predicate) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestMapSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		f        func(int) int
		expected result.Result[int]
	}{
		"success": {
			value:    result.Ok(1),
			f:        func(i int) int { return i * 2 },
			expected: result.Ok(2),
		},
		"success_closure": {
			value:    result.Ok(2),
			f:        func(i int) int { return i * closureConstant },
			expected: result.Ok(2 * closureConstant),
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			f:        func(i int) int { return i * 2 },
			expected: result.Err[int](errors.New("error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if !result.Equal(result.Map(tc.value, tc.f), tc.expected) {
				t.Fail()
			}
		})
	}
}

func TestMapDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		f        func(int) bool
		expected result.Result[bool]
	}{
		"success": {
			value:    result.Ok(1),
			f:        func(i int) bool { return i == 1 },
			expected: result.Ok(true),
		},
		"success_closure": {
			value:    result.Ok(2),
			f:        func(i int) bool { return i > closureConstant },
			expected: result.Ok(false),
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			f:        func(i int) bool { return true },
			expected: result.Err[bool](errors.New("error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if !result.Equal(result.Map(tc.value, tc.f), tc.expected) {
				t.Fail()
			}
		})
	}
}

func TestMapOrSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		fallback int
		f        func(int) int
		expected int
	}{
		"success": {
			value:    result.Ok(1),
			fallback: 10,
			f:        func(i int) int { return i * 2 },
			expected: 2,
		},
		"success_closure": {
			value:    result.Ok(2),
			fallback: 10,
			f:        func(i int) int { return i + closureConstant },
			expected: 7,
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			fallback: 10,
			f:        func(i int) int { return i * 2 },
			expected: 10,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.MapOr(tc.value, tc.fallback, tc.f) != tc.expected {
				t.Fail()
			}
		})
	}
}

func TestMapOrDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		fallback bool
		f        func(int) bool
		expected bool
	}{
		"success": {
			value:    result.Ok(1),
			fallback: true,
			f:        func(i int) bool { return i < 2 },
			expected: true,
		},
		"success_closure": {
			value:    result.Ok(2),
			fallback: true,
			f:        func(i int) bool { return i > closureConstant },
			expected: false,
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			fallback: true,
			f:        func(i int) bool { return i == 2 },
			expected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.MapOr(tc.value, tc.fallback, tc.f) != tc.expected {
				t.Fail()
			}
		})
	}
}

func TestMapOrElseSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		fallback func(error) int
		f        func(int) int
		expected int
	}{
		"success": {
			value:    result.Ok(1),
			fallback: func(err error) int { return 5 },
			f:        func(i int) int { return i * 2 },
			expected: 2,
		},
		"success_closure": {
			value:    result.Ok(2),
			fallback: func(err error) int { return closureConstant },
			f:        func(i int) int { return i + closureConstant },
			expected: 7,
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			fallback: func(err error) int { return len(err.Error()) },
			f:        func(i int) int { return i * 2 },
			expected: 5,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.MapOrElse(tc.value, tc.fallback, tc.f) != tc.expected {
				t.Fail()
			}
		})
	}
}

func TestMapOrElseDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		fallback func(error) bool
		f        func(int) bool
		expected bool
	}{
		"success": {
			value:    result.Ok(1),
			fallback: func(err error) bool { return true },
			f:        func(i int) bool { return i < 2 },
			expected: true,
		},
		"success_closure": {
			value:    result.Ok(2),
			fallback: func(err error) bool { return false },
			f:        func(i int) bool { return i > closureConstant },
			expected: false,
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			fallback: func(err error) bool { return err.Error() == "error" },
			f:        func(i int) bool { return i == 2 },
			expected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.MapOrElse(tc.value, tc.fallback, tc.f) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestMapErr(t *testing.T) {
	closureConstant := "test"
	tests := map[string]struct {
		value    result.Result[int]
		f        func(error) error
		expected result.Result[int]
	}{
		"error": {
			value:    result.Err[int](errors.New("error")),
			f:        func(err error) error { return fmt.Errorf("wrapped: %v", err) },
			expected: result.Err[int](errors.New("wrapped: error")),
		},
		"error_closure": {
			value:    result.Err[int](errors.New("error")),
			f:        func(err error) error { return fmt.Errorf("%s: %v", closureConstant, err) },
			expected: result.Err[int](errors.New(closureConstant + ": error")),
		},
		"success": {
			value:    result.Ok(1),
			f:        func(err error) error { return err },
			expected: result.Ok(1),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if !result.Equal(tc.value.MapErr(tc.f), tc.expected) {
				t.Fail()
			}
		})
	}
}
func TestInspect(t *testing.T) {
	initialClosureConstant := 5
	closureConstant := initialClosureConstant
	tests := map[string]struct {
		value           result.Result[int]
		f               func(int)
		closureExpected int
	}{
		"success": {
			value:           result.Ok(1),
			f:               func(i int) { closureConstant = closureConstant + i },
			closureExpected: closureConstant + 1,
		},
		"error": {
			value:           result.Err[int](errors.New("error")),
			f:               func(i int) {},
			closureExpected: closureConstant,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			closureConstant = initialClosureConstant
			copy := tc.value.Inspect(tc.f)
			if closureConstant != tc.closureExpected {
				t.Fail()
			}
			if !result.Equal(copy, tc.value) {
				t.Fail()
			}
		})
	}
}
func TestInspectErr(t *testing.T) {
	initialClosureConstant := "hello"
	closureConstant := initialClosureConstant
	tests := map[string]struct {
		value           result.Result[int]
		f               func(error)
		closureExpected string
	}{
		"success": {
			value:           result.Ok(1),
			f:               func(err error) {},
			closureExpected: closureConstant,
		},
		"error": {
			value:           result.Err[int](errors.New("error")),
			f:               func(err error) { closureConstant = closureConstant + err.Error() },
			closureExpected: closureConstant + "error",
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			closureConstant = initialClosureConstant
			copy := tc.value.InspectErr(tc.f)
			if closureConstant != tc.closureExpected {
				t.Fail()
			}
			if !result.Equal(copy, tc.value) {
				t.Fail()
			}
		})
	}
}
func TestExpect(t *testing.T) {
	tests := map[string]struct {
		value         result.Result[int]
		msg           string
		inner         int
		errorExpected bool
	}{
		"success": {
			value:         result.Ok(1),
			msg:           "",
			inner:         1,
			errorExpected: false,
		},
		"error": {
			value:         result.Err[int](errors.New("error")),
			msg:           "no value",
			inner:         0,
			errorExpected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			defer func() {
				panicMsg := recover()
				if tc.errorExpected && panicMsg != tc.msg {
					t.Fail()
				}
			}()
			val := tc.value.Expect(tc.msg)
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrap(t *testing.T) {
	tests := map[string]struct {
		value         result.Result[int]
		inner         int
		errorExpected bool
	}{
		"success": {
			value:         result.Ok(1),
			inner:         1,
			errorExpected: false,
		},
		"error": {
			value:         result.Err[int](errors.New("error")),
			inner:         0,
			errorExpected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			defer func() {
				panicMsg := recover()
				if (panicMsg != nil) != tc.errorExpected {
					t.Fail()
				}
			}()
			val := tc.value.Unwrap()
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOrDefault(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[int]
		expected int
	}{
		"success": {
			value:    result.Ok(1),
			expected: 1,
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			expected: 0,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val := tc.value.UnwrapOrDefault()
			if val != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestExpectErr(t *testing.T) {
	tests := map[string]struct {
		value         result.Result[int]
		msg           string
		inner         error
		panicExpected bool
	}{
		"success": {
			value:         result.Ok(1),
			msg:           "result had value",
			inner:         nil,
			panicExpected: true,
		},
		"error": {
			value:         result.Err[int](errors.New("error")),
			msg:           "",
			inner:         errors.New("error"),
			panicExpected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			defer func() {
				panicMsg := recover()
				if tc.panicExpected && panicMsg != tc.msg {
					t.Fail()
				}
			}()
			err := tc.value.ExpectErr(tc.msg)
			if err.Error() != tc.inner.Error() {
				t.Fail()
			}
		})
	}
}
func TestUnwrapErr(t *testing.T) {
	tests := map[string]struct {
		value         result.Result[int]
		inner         error
		panicExpected bool
	}{
		"success": {
			value:         result.Ok(1),
			inner:         nil,
			panicExpected: true,
		},
		"error": {
			value:         result.Err[int](errors.New("error")),
			inner:         errors.New("error"),
			panicExpected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			defer func() {
				panicMsg := recover()
				if (panicMsg != nil) != tc.panicExpected {
					t.Fail()
				}
			}()
			err := tc.value.UnwrapErr()
			if err.Error() != tc.inner.Error() {
				t.Fail()
			}
		})
	}
}
func TestAnd(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[int]
		other    result.Result[int]
		expected result.Result[int]
	}{
		"ok_ok": {
			value:    result.Ok(1),
			other:    result.Ok(2),
			expected: result.Ok(2),
		},
		"ok_err": {
			value:    result.Ok(1),
			other:    result.Err[int](errors.New("error")),
			expected: result.Err[int](errors.New("error")),
		},
		"err_ok": {
			value:    result.Err[int](errors.New("error")),
			other:    result.Ok(2),
			expected: result.Err[int](errors.New("error")),
		},
		"err_err": {
			value:    result.Err[int](errors.New("error")),
			other:    result.Err[int](errors.New("other error")),
			expected: result.Err[int](errors.New("error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			res := result.And(tc.value, tc.other)
			if !result.Equal(res, tc.expected) {
				t.Fail()
			}
		})
	}
}
func TestAndThenSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		f        func(int) result.Result[int]
		expected result.Result[int]
	}{
		"success": {
			value:    result.Ok(1),
			f:        func(i int) result.Result[int] { return result.Ok(i + 2) },
			expected: result.Ok(3),
		},
		"success_closure": {
			value:    result.Ok(2),
			f:        func(i int) result.Result[int] { return result.Ok(i * closureConstant) },
			expected: result.Ok(2 * closureConstant),
		},
		"success_return_err": {
			value:    result.Ok(1),
			f:        func(i int) result.Result[int] { return result.Err[int](errors.New("bad value")) },
			expected: result.Err[int](errors.New("bad value")),
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			f:        func(i int) result.Result[int] { return result.Ok(1) },
			expected: result.Err[int](errors.New("error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			res := result.AndThen(tc.value, tc.f)
			if !result.Equal(res, tc.expected) {
				t.Fail()
			}
		})
	}
}

func TestAndThenDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    result.Result[int]
		f        func(int) result.Result[bool]
		expected result.Result[bool]
	}{
		"success": {
			value:    result.Ok(1),
			f:        func(i int) result.Result[bool] { return result.Ok(i < 2) },
			expected: result.Ok(true),
		},
		"success_closure": {
			value:    result.Ok(2),
			f:        func(i int) result.Result[bool] { return result.Ok(i > closureConstant) },
			expected: result.Ok(false),
		},
		"success_return_err": {
			value:    result.Ok(1),
			f:        func(i int) result.Result[bool] { return result.Err[bool](errors.New("bad value")) },
			expected: result.Err[bool](errors.New("bad value")),
		},
		"error": {
			value:    result.Err[int](errors.New("error")),
			f:        func(i int) result.Result[bool] { return result.Ok(true) },
			expected: result.Err[bool](errors.New("error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			res := result.AndThen(tc.value, tc.f)
			if !result.Equal(res, tc.expected) {
				t.Fail()
			}
		})
	}
}
func TestOr(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[int]
		other    result.Result[int]
		expected result.Result[int]
	}{
		"ok_ok": {
			value:    result.Ok(1),
			other:    result.Ok(2),
			expected: result.Ok(1),
		},
		"ok_err": {
			value:    result.Ok(1),
			other:    result.Err[int](errors.New("error")),
			expected: result.Ok(1),
		},
		"err_ok": {
			value:    result.Err[int](errors.New("error")),
			other:    result.Ok(2),
			expected: result.Ok(2),
		},
		"err_err": {
			value:    result.Err[int](errors.New("error")),
			other:    result.Err[int](errors.New("other error")),
			expected: result.Err[int](errors.New("other error")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			res := tc.value.Or(tc.other)
			if !result.Equal(res, tc.expected) {
				t.Fail()
			}
		})
	}
}
func TestOrElse(t *testing.T) {
	tests := map[string]struct {
		value      result.Result[string]
		fallbackFn func(error) result.Result[string]
		expected   result.Result[string]
	}{
		"success": {
			value:      result.Ok("hello"),
			fallbackFn: func(err error) result.Result[string] { return result.Ok(err.Error() + "world") },
			expected:   result.Ok("hello"),
		},
		"error_success": {
			value:      result.Err[string](errors.New("failed")),
			fallbackFn: func(err error) result.Result[string] { return result.Ok(err.Error() + " the world") },
			expected:   result.Ok("failed the world"),
		},
		"error_error": {
			value:      result.Err[string](errors.New("failed")),
			fallbackFn: func(err error) result.Result[string] { return result.Err[string](fmt.Errorf("status: %v", err)) },
			expected:   result.Err[string](errors.New("status: failed")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			res := tc.value.OrElse(tc.fallbackFn)
			if !result.Equal(res, tc.expected) {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOr(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[string]
		fallback string
		expected string
	}{
		"success": {
			value:    result.Ok("hello"),
			fallback: "fallback",
			expected: "hello",
		},
		"error": {
			value:    result.Err[string](errors.New("failed")),
			fallback: "fallback",
			expected: "fallback",
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.UnwrapOr(tc.fallback) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOrElse(t *testing.T) {
	tests := map[string]struct {
		value      result.Result[string]
		fallbackFn func(error) string
		expected   string
	}{
		"success": {
			value:      result.Ok("hello"),
			fallbackFn: func(err error) string { return err.Error() + " world" },
			expected:   "hello",
		},
		"error": {
			value:      result.Err[string](errors.New("failed")),
			fallbackFn: func(err error) string { return err.Error() + " world" },
			expected:   "failed world",
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.UnwrapOrElse(tc.fallbackFn) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestContains(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[string]
		target   string
		expected bool
	}{
		"success_true": {
			value:    result.Ok("hello"),
			target:   "hello",
			expected: true,
		},
		"success_false": {
			value:    result.Ok("hello"),
			target:   "world",
			expected: false,
		},
		"error": {
			value:    result.Err[string](errors.New("failed")),
			target:   "failed",
			expected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.Contains(tc.value, tc.target) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestContainsErr(t *testing.T) {
	err1 := errors.New("failed")
	err2 := errors.New("second failure")
	tests := map[string]struct {
		value    result.Result[string]
		target   error
		expected bool
	}{
		"success": {
			value:    result.Ok("hello"),
			target:   errors.New("false"),
			expected: false,
		},
		"error_true": {
			value:    result.Err[string](err1),
			target:   err1,
			expected: true,
		},
		"error_true_suberror": {
			value:    result.Err[string](fmt.Errorf("attempted parsing: %w", err1)),
			target:   err1,
			expected: true,
		},
		"error_false": {
			value:    result.Err[string](err1),
			target:   err2,
			expected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.ContainsErr(tc.target) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestCopy(t *testing.T) {
	tests := map[string]struct {
		value result.Result[string]
	}{
		"success": {
			value: result.Ok("hello"),
		},
		"error": {
			value: result.Err[string](errors.New("failed")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			copy := tc.value.Copy()
			if !result.Equal(tc.value, copy) {
				t.Fail()
			}
			if &copy == &tc.value {
				t.Fail()
			}
		})
	}
}
func TestFlatten(t *testing.T) {
	tests := map[string]struct {
		value    result.Result[result.Result[string]]
		expected result.Result[string]
	}{
		"ok_ok": {
			value:    result.Ok(result.Ok("hello")),
			expected: result.Ok("hello"),
		},
		"ok_error": {
			value:    result.Ok(result.Err[string](errors.New("failed"))),
			expected: result.Err[string](errors.New("failed")),
		},
		"error": {
			value:    result.Err[result.Result[string]](errors.New("failed")),
			expected: result.Err[string](errors.New("failed")),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if !result.Equal(result.Flatten(tc.value), tc.expected) {
				t.Fail()
			}
		})
	}
}

func TestEqual(t *testing.T) {
	err := errors.New("failed")
	tests := map[string]struct {
		value    result.Result[string]
		other    result.Result[string]
		expected bool
	}{
		"ok_ok_true": {
			value:    result.Ok("hello"),
			other:    result.Ok("hello"),
			expected: true,
		},
		"ok_ok_false": {
			value:    result.Ok("hello"),
			other:    result.Ok("world"),
			expected: false,
		},
		"ok_error": {
			value:    result.Ok("hello"),
			other:    result.Err[string](err),
			expected: false,
		},
		"error_error_true": {
			value:    result.Err[string](err),
			other:    result.Err[string](err),
			expected: true,
		},
		"error_error_false": {
			value:    result.Err[string](err),
			other:    result.Err[string](errors.New("not this")),
			expected: false,
		},
		"error_error_false_suberror": {
			value:    result.Err[string](err),
			other:    result.Err[string](fmt.Errorf("suberror: %w", err)),
			expected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if result.Equal(tc.value, tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}
