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

}
func TestExpectErr(t *testing.T) {

}
func TestUnwrapErr(t *testing.T) {

}
func TestAnd(t *testing.T) {

}
func TestAndThen(t *testing.T) {

}
func TestOr(t *testing.T) {

}
func TestOrElse(t *testing.T) {

}
func TestUnwrapOr(t *testing.T) {

}
func TestUnwrapOrElse(t *testing.T) {

}
func TestContains(t *testing.T) {

}
func TestContainsErr(t *testing.T) {

}
func TestCopy(t *testing.T) {

}
func TestFlatten(t *testing.T) {

}
