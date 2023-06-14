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

	detailedError, ok := whoops.As[whoops.Error](err)
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

	is := whoops.IsOfType[whoops.Error](innerError)
	fmt.Println("innerError:", is)

	is = whoops.IsOfType[whoops.Error](err)
	fmt.Println("err:", is)

	// Output:
	// innerError: false
	// err: true
}

func ExampleJson() {
	strErr := errors.New("some error")
	dErr := whoops.Wrap(whoops.ErrorCode("some-error-code"), strErr, "some message")

	json, _ := whoops.Json(strErr)
	fmt.Println(string(json))

	json, _ = whoops.Json(strErr, true)
	fmt.Println(string(json))

	json, _ = whoops.Json(dErr, true)
	fmt.Println(string(json))

	json, _ = whoops.Json(dErr)
	fmt.Println(string(json))

	// Output:
	// {
	//   "error": "internal error"
	// }
	// {
	//   "error": "some error"
	// }
	// {
	//   "error": "some error",
	//   "code": "some-error-code",
	//   "message": "some message"
	// }
	// {
	//   "error": "internal error",
	//   "code": "some-error-code",
	//   "message": "some message"
	// }
}

func Example_formatting() {
	err := errors.New("some normal error")
	fmt.Printf("%s\n", err)

	err = whoops.Wrap(whoops.CodeUnexpected, err, "Oh no!", "anyway")
	fmt.Printf("%s\n", err)

	// Print with callstack of depth 5
	fmt.Printf("%+5v\n", err)

	// Print detailed error stack
	fmt.Printf("%#v\n", err)
}
