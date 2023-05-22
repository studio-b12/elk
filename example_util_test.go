package whoops_test

import (
	"errors"
	"fmt"

	"github.com/studio-b12/whoops"
)

func ExampleUnwrapFull() {
	type WrappedError struct {
		whoops.InnerError
	}

	var err error
	originErr := errors.New("some error")
	err = whoops.Wrap(whoops.CodeUnexpected, originErr, "Some message")
	err = WrappedError{InnerError: whoops.InnerError{Inner: err}}

	err = whoops.UnwrapFull(err)
	fmt.Println(err == originErr)

	// Output: true
}

func ExampleAs() {
	const ErrUnexpected = whoops.ErrorCode("unexpected-error")

	type WrappedError struct {
		whoops.InnerError
	}

	err := errors.New("some error")
	err = whoops.Wrap(whoops.CodeUnexpected, err, "Some message")
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

	innerError := errors.New("some error")
	var err error = whoops.Wrap(whoops.CodeUnexpected, innerError, "Some message")
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

	err = whoops.Wrap(whoops.CodeUnexpected, err, "Oh no!", "anyway")
	msg = whoops.Format(err)
	fmt.Println(msg)
}

func ExampleMessage() {
	strErr := errors.New("some error")
	dErr := whoops.Wrap(whoops.CodeUnexpected, strErr, "some message")

	fmt.Println(whoops.Message(strErr))
	fmt.Println(whoops.Message(dErr))

	// Output:
	// some error
	// some message
}

func ExampleJson() {
	strErr := errors.New("some error")
	dErr := whoops.Wrap(whoops.ErrorCode("some-error-code"), strErr, "some message")

	json, _ := whoops.Json(strErr, true)
	fmt.Println(json)

	json, _ = whoops.Json(dErr, true)
	fmt.Println(json)

	// Output:
	// {
	//   "error": "some error"
	// }
	// {
	//   "error": "some error",
	//   "code": "some-error-code",
	//   "message": "some message"
	// }
}
