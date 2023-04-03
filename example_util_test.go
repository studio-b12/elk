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
	err = whoops.WrapMessage(originErr, "Some message")
	err = WrappedError{InnerError: whoops.InnerError{Inner: err}}

	err = whoops.UnwrapFull(err)
	fmt.Println(err == originErr)

	// Output: true
}

func ExampleAs() {
	type WrappedError struct {
		whoops.InnerError
	}

	err := errors.New("some error")
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

	innerError := errors.New("some error")
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

func ExampleMessage() {
	strErr := errors.New("some error")
	dErr := whoops.WrapMessage(strErr, "some message")

	fmt.Println(whoops.Message(strErr))
	fmt.Println(whoops.Message(dErr))

	// Output:
	// some error
	// some message
}

func ExampleJson() {
	strErr := errors.New("some error")
	dErr := whoops.WrapMessage(strErr, "some message", "some details")

	json, _ := whoops.Json(strErr)
	fmt.Println(json)

	json, _ = whoops.Json(dErr)
	fmt.Println(json)

	// Details are excluded by default, but you can
	// pass "true" for showDetails to include them
	// in the JSON output.
	json, _ = whoops.Json(dErr, true)
	fmt.Println(json)

	// Output:
	// {
	//   "error": "some error"
	// }
	// {
	//   "error": "some error",
	//   "message": "some message"
	// }
	// {
	//   "error": "some error",
	//   "message": "some message",
	//   "details": "some details"
	// }
}

func ExampleDetailsOfType() {
	type status struct {
		message string
		code    int
	}

	err := errors.New("some error")
	err = whoops.Wrap(err, "some string details")
	err = whoops.Wrap(err, status{"some message", 42})

	strDetails, _ := whoops.DetailsOfType[string](err)
	fmt.Printf("strDetails: %s\n", strDetails)

	statusDetails, _ := whoops.DetailsOfType[status](err)
	fmt.Printf("statusDetails: %+v\n", statusDetails)

	// Output:
	// strDetails: some string details
	// statusDetails: {message:some message code:42}
}
