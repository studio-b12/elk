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

	Error, ok := elk.As[elk.Error](err)
	if ok {
		message := Error.Message()
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

type DetailedError struct {
	elk.InnerError
	details any
}

func (t DetailedError) Details() any {
	return t.details
}

func ExampleJson() {
	strErr := errors.New("some error")
	mErr := elk.Wrap("some-error-code", strErr, "some message")

	json, _ := elk.Json(strErr, 0)
	fmt.Println(string(json))

	json, _ = elk.Json(strErr, 400)
	fmt.Println(string(json))

	json, _ = elk.Json(mErr, 0)
	fmt.Println(string(json))

	json, _ = elk.Json(mErr, 400)
	fmt.Println(string(json))

	dtErr := DetailedError{}
	dtErr.Inner = elk.NewError("some-error", "an error with details")
	dtErr.details = struct {
		Foo string
		Bar int
	}{
		Foo: "foo",
		Bar: 123,
	}

	json, _ = elk.Json(dtErr, 500)
	fmt.Println(string(json))

	dteErr := elk.Wrap("some-detailed-error-wrapped", dtErr, "some detailed error wrapped")
	json, _ = elk.Json(dteErr, 500)
	fmt.Println(string(json))

	// Output:
	// {
	//   "Code": "unexpected-error"
	// }
	// {
	//   "Code": "unexpected-error",
	//   "Status": 400
	// }
	// {
	//   "Code": "some-error-code",
	//   "Message": "some message"
	// }
	// {
	//   "Code": "some-error-code",
	//   "Message": "some message",
	//   "Status": 400
	// }
	// {
	//   "Code": "unexpected-error",
	//   "Status": 500,
	//   "Details": {
	//     "Foo": "foo",
	//     "Bar": 123
	//   }
	// }
	// {
	//   "Code": "some-detailed-error-wrapped",
	//   "Message": "some detailed error wrapped",
	//   "Status": 500,
	//   "Details": {
	//     "Foo": "foo",
	//     "Bar": 123
	//   }
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
