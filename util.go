package whoops

import (
	"errors"
)

// As applies errors.As() on the given err
// using the given type T as target for the
// unwrapping.
//
// Refer to the documentation of errors.As()
// for more details:
// https://pkg.go.dev/errors#As
func As[T error](err error) (t T, ok bool) {
	ok = errors.As(err, &t)
	return t, ok
}

// IsOfType returns true when the given
// error is of the type of T.
//
// If not and the error can be unwrapped,
// the unwrapped error will be checked
// until it either matches the type T or
// can not be further unwrapped.
func IsOfType[T error](err error) bool {
	_, ok := err.(T)
	if ok {
		return true
	}

	err = errors.Unwrap(err)
	if err != nil {
		return IsOfType[T](err)
	}

	return false
}

// Format returns the formatted error message
// when err implements Formatted. Otherwise,
// the result of err.Error() is returned.
func Format(err error) string {
	if fErr, ok := err.(Formatted); ok {
		return fErr.Formatted()
	}
	return err.Error()
}
