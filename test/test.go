package test

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, test string, actual interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s: Must be equal, expected {%s} but got {%s}", test, expected, actual)
	}
}

func AssertNil(t *testing.T, test string, actual interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("%s: Must be nil, but is %s", test, actual)
	}
}

func AssertNotNil(t *testing.T, test string, actual interface{}) {
	t.Helper()
	if actual == nil {
		t.Errorf("%s: Must be not nil", test)
	}
}

func AssertError(t *testing.T, test string, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: error is nil", test)
	}
}

func AssertNoError(t *testing.T, test string, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: error isn't nil / %v", test, err)
	}
}

func AssertArrayNoError(t *testing.T, test string, err []error) {
	t.Helper()
	if len(err) > 0 {
		t.Errorf("%s: error isn't nil / %v", test, err)
	}
}

func AssertGreaterThanZero(t *testing.T, test string, value int) {
	t.Helper()
	if value <= 0 {
		t.Errorf("%s: error when value isn't greater than 0 / %v", test, value)
	}
}
