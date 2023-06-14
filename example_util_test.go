package elk_test

import (
	"errors"
	"fmt"

	"github.com/studio-b12/elk"
)

func ExampleUnwrapFull() {
	type WrappedError struct {
		elk.InnerError
	}

	var err error
	originErr := errors.New("some error")
	err = elk.Wrap(elk.CodeUnexpected, originErr, "Some message")
	err = WrappedError{InnerError: elk.InnerError{Inner: err}}

	err = elk.UnwrapFull(err)
	fmt.Println(err == originErr)

	// Output: true
}

func ExampleAs() {
	const ErrUnexpected = elk.ErrorCode("unexpected-error")

	type WrappedError struct {
		elk.InnerError
	}

	err := errors.New("some error")
	err = elk.Wrap(elk.CodeUnexpected, err, "Some message")
	err = WrappedError{InnerError: elk.InnerError{Inner: err}}

	detailedError, ok := elk.As[elk.Error](err)
	if ok {
		message := detailedError.Message()
		fmt.Println(message)
	}

	// Output: Some message
}

func ExampleIsOfType() {
	type WrappedError struct {
		elk.InnerError
	}

	innerError := errors.New("some error")
	var err error = elk.Wrap(elk.CodeUnexpected, innerError, "Some message")
	err = WrappedError{InnerError: elk.InnerError{Inner: err}}

	is := elk.IsOfType[elk.Error](innerError)
	fmt.Println("innerError:", is)

	is = elk.IsOfType[elk.Error](err)
	fmt.Println("err:", is)

	// Output:
	// innerError: false
	// err: true
}

func ExampleJson() {
	strErr := errors.New("some error")
	dErr := elk.Wrap(elk.ErrorCode("some-error-code"), strErr, "some message")

	json, _ := elk.Json(strErr)
	fmt.Println(string(json))

	json, _ = elk.Json(strErr, true)
	fmt.Println(string(json))

	json, _ = elk.Json(dErr, true)
	fmt.Println(string(json))

	json, _ = elk.Json(dErr)
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

	err = elk.Wrap(elk.CodeUnexpected, err, "Oh no!", "anyway")
	fmt.Printf("%s\n", err)

	// Print with callstack of depth 5
	fmt.Printf("%+5v\n", err)

	// Print detailed error stack
	fmt.Printf("%#v\n", err)
}
