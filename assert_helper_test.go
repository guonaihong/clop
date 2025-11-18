package clop

import (
	"reflect"
	"testing"
)

type testAsserts struct{}

var assert = testAsserts{}

func (testAsserts) Equal(t *testing.T, got, want interface{}, msgAndArgs ...interface{}) {
	if !reflect.DeepEqual(got, want) {
		if len(msgAndArgs) > 0 {
			t.Fatalf("unexpected value: got=%#v, want=%#v, msg=%v", got, want, msgAndArgs[0])
		} else {
			t.Fatalf("unexpected value: got=%#v, want=%#v", got, want)
		}
	}
}

func (testAsserts) NotEqual(t *testing.T, got, want interface{}, msgAndArgs ...interface{}) {
	if reflect.DeepEqual(got, want) {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected values to differ, both=%#v, msg=%v", got, msgAndArgs[0])
		} else {
			t.Fatalf("expected values to differ, both=%#v", got)
		}
	}
}

func (testAsserts) NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err != nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected no error, got=%v, msg=%v", err, msgAndArgs[0])
		} else {
			t.Fatalf("expected no error, got=%v", err)
		}
	}
}

func (testAsserts) Error(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err == nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected error, got nil, msg=%v", msgAndArgs[0])
		} else {
			t.Fatalf("expected error, got nil")
		}
	}
}

func (testAsserts) True(t *testing.T, cond bool, msgAndArgs ...interface{}) {
	if !cond {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected condition to be true, msg=%v", msgAndArgs[0])
		} else {
			t.Fatalf("expected condition to be true")
		}
	}
}

func (testAsserts) False(t *testing.T, cond bool, msgAndArgs ...interface{}) {
	if cond {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected condition to be false, msg=%v", msgAndArgs[0])
		} else {
			t.Fatalf("expected condition to be false")
		}
	}
}

// a < b
func (testAsserts) Less(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) {
	less := false

	switch x := a.(type) {
	case int:
		if y, ok := b.(int); ok {
			less = x < y
		}
	case int64:
		if y, ok := b.(int64); ok {
			less = x < y
		}
	case uint64:
		if y, ok := b.(uint64); ok {
			less = x < y
		}
	case float64:
		if y, ok := b.(float64); ok {
			less = x < y
		}
	}

	if !less {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected a < b, got a=%#v, b=%#v, msg=%v", a, b, msgAndArgs[0])
		} else {
			t.Fatalf("expected a < b, got a=%#v, b=%#v", a, b)
		}
	}
}

func (testAsserts) Panics(t *testing.T, f func(), msgAndArgs ...interface{}) {
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		f()
	}()

	if !panicked {
		if len(msgAndArgs) > 0 {
			t.Fatalf("expected panic, but none occurred, msg=%v", msgAndArgs[0])
		} else {
			t.Fatalf("expected panic, but none occurred")
		}
	}
}
