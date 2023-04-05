package whoops_test

import (
	"errors"
	"fmt"

	"github.com/studio-b12/whoops"
)

type StatusError struct {
	whoops.InnerError

	StatusCode int
}

func NewStatusError(inner error, status int) error {
	var s StatusError
	s.Inner = inner
	s.StatusCode = status
	return s
}

func (t StatusError) Error() string {
	return fmt.Sprintf("%s (%d)", t.Inner.Error(), t.StatusCode)
}

func Example_inner() {
	err := errors.New("not found")
	statusErr := NewStatusError(err, 404)
	fmt.Println(statusErr.Error())

	// Because whoops.InnerError implements the Error()
	// as well as the Unwrap() method, StatusError inherits
	// these methods unless they are overridden.
	inner := errors.Unwrap(statusErr)
	fmt.Println(inner == err)

	// Output:
	// not found (404)
	// true
}
