package whoops

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type ErrorCode string

const (
	CodeUnexpected = ErrorCode("unexpected-error")
)

const (
	maxCallStackDepth = 100
	indent            = "    "
)

// DetailedError contains a wrapped inner error,
// an optional message, optional details objects
// and a CallStack from where the error has been
// created.
type DetailedError struct {
	InnerError

	code      ErrorCode
	message   string
	callStack *CallStack
}

var (
	_ HasMessage   = (*DetailedError)(nil)
	_ HasCode      = (*DetailedError)(nil)
	_ HasCallStack = (*DetailedError)(nil)
)

// Detailed creates a new DetailedError with the
// given code and optional message.
func Detailed(code ErrorCode, message ...string) DetailedError {
	d := Wrap(code, errors.New(string(code)), message...)
	d.callStack.offset++
	return d
}

// Cast takes an arbitrary error and if it is not of type DetailedError,
// it will be wrapped in a new DetailedError which is then returned.
// If fallback is passed, it will be used as the ErrorCode of the new
// DetailedError. Otherwise, CodeUnexpected is used.
//
// If err is of type DetailedError, it is simply returned unchanged.
func Cast(err error, fallback ...ErrorCode) DetailedError {
	d, ok := err.(DetailedError)
	if !ok {
		fallbackCode := CodeUnexpected
		if len(fallback) > 0 {
			fallbackCode = fallback[0]
		}
		d = Wrap(fallbackCode, err)
		d.callStack.offset++
	}
	return d
}

// WrapMessage is the same as wrap including an
// additional message. The message will be shown
// in place of the wrapped errors result of the
// Error() method.
func Wrap(code ErrorCode, err error, message ...string) DetailedError {
	var d DetailedError

	d.code = code
	d.Inner = err
	d.callStack = newCallStack(1, maxCallStackDepth)

	if len(message) > 0 {
		d.message = strings.Join(message, " ")
	}

	return d
}

// Error returns the error information as
// formatted string.
func (t DetailedError) Error() string {
	return fmt.Sprintf("%s", t)
}

// Format implements custom formatting rules used with the formatting
// functionalities in the fmt package.
//
// %s, %q
//
// Prints the message of the error, if available. Otherwise, the
// %s format of the inner error is represented. If the inner error is nil
// and no message is set, the error code is printed.
//
// %v
//
// Prints a more detailed representation of the error. Without any flags,
// the error is printed in the format `<{errorCode}> {message} ({innerError})`.
//
// By passing the `+` flag, the inner error is represented in a seperate line.
// Also, by using the precision parameter, you can specify the depth of the
// represented callstack (i.E. `%+.5v` - prints a callstack of depth 5). Otherwise,
// no callstack will be printed.
//
// Bypassing the `#` flag, an even more verbose representation of the error is
// printed. It shows the complete chain of errors wrapped in the DetailedError
// with information about message, code, initiation origin and type of the error.
// With the precision parameter, you can define the depth of the unwrapping. The
// default value is 100, if not specified.
func (t DetailedError) Format(s fmt.State, verb rune) {
	width, _ := s.Precision()

	switch verb {
	case 'v':
		if s.Flag('+') {
			t.writeDetailed(s, width)
		} else if s.Flag('#') {
			t.writeStacked(s, width)
		} else {
			t.writeTitle(s, true)
		}
	case 's', 'q':
		if t.message != "" {
			fmt.Fprint(s, t.message)
		} else if t.Inner != nil {
			fmt.Fprintf(s, "%s", t.Inner)
		} else {
			fmt.Fprint(s, t.code)
		}
	}
}

// Message returns the errors message text,
// if specified.
func (t DetailedError) Message() string {
	return t.message
}

// Code returns the inner ErrorCode of
// the error.
func (t DetailedError) Code() ErrorCode {
	return t.code
}

// CallStack returns the errors CallStack
// starting from where the DetailedError
// has been created.
func (t DetailedError) CallStack() *CallStack {
	return t.callStack
}

func (t DetailedError) writeTitle(w io.Writer, withError bool) {
	fmt.Fprintf(w, "<%s>", t.code)
	if t.message != "" {
		fmt.Fprintf(w, " %s", t.message)
	}
	if withError && t.Inner != nil {
		fmt.Fprintf(w, " (%s)", t.Inner)
	}
}

func (t DetailedError) writeDetailed(w io.Writer, stack int) {
	t.writeTitle(w, false)

	fmt.Fprintln(w)

	if stack > 0 {
		fmt.Fprint(w, "stack:\n")

		// We only want to print the last callstack in the error
		// chain here, so we unwrap the error until we found the
		// last one which implements HasCallStack.
		var (
			e  error = t
			cs *CallStack
		)
		for e != nil {
			ecs, ok := e.(HasCallStack)
			if !ok {
				break
			}
			cs = ecs.CallStack()
			e = errors.Unwrap(e)
		}

		cs.WriteIndent(w, stack, "  ")
	}

	fmt.Fprintf(w, "inner error:\n  %s", t.Inner)
}

func (t DetailedError) writeStacked(w io.Writer, depth int) {
	if depth == 0 {
		depth = 100
	}

	var err error = t
	i := 0

	for err != nil && i < depth {
		if d, ok := err.(DetailedError); ok {
			d.writeTitle(w, false)

			fmt.Fprintln(w)

			if frame, ok := d.CallStack().At(0); ok {
				fmt.Fprintf(w, "originated:\n  %s\n", frame)
			}
		} else {
			fmt.Fprintf(w, "%+v\n", err)
		}

		fmt.Fprintf(w, "type:\n  %s\n", reflect.TypeOf(err))

		fmt.Fprintln(w, "----------")

		err = errors.Unwrap(err)
		i++
	}
}
