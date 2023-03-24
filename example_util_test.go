package whoops_test

import (
	"errors"
	"fmt"

	"github.com/studio-b12/whoops"
)

func ExampleAs() {
	type WrappedError struct {
		whoops.InnerError
	}

	err := errors.New("Some error")
	err = whoops.WrapMessage(err, "Some message")
	err = WrappedError{InnerError: whoops.InnerError{Inner: err}}

	detailedError, ok := whoops.As[whoops.DetailedError](err)
	if ok {
		message := detailedError.Message()
		fmt.Println(message)
	}

	// Output: Some message
}

func ExampleIsOfType() {
	type WrappedError struct {
		whoops.InnerError
	}

	innerError := errors.New("Some error")
	var err error = whoops.WrapMessage(innerError, "Some message")
	err = WrappedError{InnerError: whoops.InnerError{Inner: err}}

	is := whoops.IsOfType[whoops.DetailedError](innerError)
	fmt.Println("innerError:", is)

	is = whoops.IsOfType[whoops.DetailedError](err)
	fmt.Println("err:", is)

	// Output:
	// innerError: false
	// err: true
}

func ExampleFormat() {
	err := errors.New("some normal error")
	msg := whoops.Format(err)
	fmt.Println(msg)

	err = whoops.WrapMessage(err, "Oh no!", "anyway")
	msg = whoops.Format(err)
	fmt.Println(msg)
}
