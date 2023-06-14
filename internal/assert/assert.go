package assert

import (
	"reflect"
	"testing"
)

func fail(t *testing.T, expected any, actual any) {
	t.Helper()

	t.Errorf(
		"\nassertion failed:\n"+
			"\texpected: %+v\n"+
			"\tactual:   %+v",
		expected,
		actual)
}

func Equal[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()

	if expected != actual {
		fail(t, expected, actual)
	}
}

func Nil(t *testing.T, value any) {
	t.Helper()

	if !reflect.ValueOf(value).IsNil() {
		fail(t, nil, value)
	}
}

func True(t *testing.T, value bool) {
	t.Helper()

	Equal(t, true, value)
}

func False(t *testing.T, value bool) {
	t.Helper()

	Equal(t, false, value)
}
