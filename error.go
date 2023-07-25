package elk

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

// Error contains a wrapped inner error,
// an optional message, optional details objects
// and a CallStack from where the error has been
// created.
type Error struct {
	InnerError

	code      ErrorCode
	message   string
	callStack *CallStack
}

var (
	_ HasMessage   = (*Error)(nil)
	_ HasCode      = (*Error)(nil)
	_ HasCallStack = (*Error)(nil)
)

// NewError creates a new Error with the given code and optional message.
func NewError(code ErrorCode, message ...string) Error {
	d := Wrap(code, errors.New(string(code)), message...)
	d.callStack.offset++
	return d
}

// NewError creates a new Error with the given code and message formatted
// according to the given format specification.
func NewErrorf(code ErrorCode, format string, a ...any) Error {
	e := NewError(code, fmt.Sprintf(format, a...))
	e.callStack.offset++
	return e
}

// Cast takes an arbitrary error and if it is not of type Error,
// it will be wrapped in a new Error which is then returned.
// If fallback is passed, it will be used as the ErrorCode of the new
// Error. Otherwise, CodeUnexpected is used.
//
// If err is of type Error, it is simply returned unchanged.
func Cast(err error, fallback ...ErrorCode) Error {
	d, ok := err.(Error)
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

// Wrap takes an ErrorCode, error and an optional message and creates a
// new wrapped Error containing the passed error.
func Wrap(code ErrorCode, err error, message ...string) Error {
	var d Error

	d.code = code
	d.Inner = err
	d.callStack = newCallStack(1, maxCallStackDepth)
	d.setMessage(message)

	return d
}

// Wrapf takes an ErrorCode, error and a message formatted according to the
// given format specification and creates a new wrapped Error containing the
// passed error.
func Wrapf(code ErrorCode, err error, format string, a ...any) Error {
	e := Wrap(code, err, fmt.Sprintf(format, a...))
	e.callStack.offset++
	return e
}

// WrapCopyCode wraps the error with an optional message keeping the error code
// of the wrapped error. If the wrapped error does not have a error code,
// CodeUnexpected is set insetad.
func WrapCopyCode(err error, message ...string) Error {
	e, ok := err.(Error)

	code := CodeUnexpected
	if ok {
		code = e.code
	}

	e = Wrap(code, err, message...)
	e.callStack.offset++

	return e
}

// WrapCopyCode wraps the error with a message formatted according to the given
// format specification keeping the error code of the wrapped error. If the
// wrapped error does not have a error code, CodeUnexpected is set insetad.
func WrapCopyCodef(err error, format string, a ...any) Error {
	e := WrapCopyCode(err, fmt.Sprintf(format, a...))
	e.callStack.offset++
	return e
}

// Error returns the error information as
// formatted string.
func (t Error) Error() string {
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
// printed. It shows the complete chain of errors wrapped in the Error
// with information about message, code, initiation origin and type of the error.
// With the precision parameter, you can define the depth of the unwrapping. The
// default value is 100, if not specified.
func (t Error) Format(s fmt.State, verb rune) {
	precision, hasPrecision := s.Precision()

	switch verb {
	case 'v':
		if s.Flag('+') {
			if !hasPrecision {
				precision = 1000
			}
			t.writeStack(s, precision)
		} else if s.Flag('#') {
			t.writeVerbose(s, precision)
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
func (t Error) Message() string {
	return t.message
}

// Code returns the inner ErrorCode of
// the error.
func (t Error) Code() ErrorCode {
	return t.code
}

// CallStack returns the errors CallStack
// starting from where the Error
// has been created.
func (t Error) CallStack() *CallStack {
	return t.callStack
}

func (t *Error) setMessage(message []string) {
	if len(message) > 0 {
		t.message = strings.Join(message, " ")
	}
}

func (t Error) writeTitle(w io.Writer, withError bool) {
	fmt.Fprintf(w, "<%s>", t.code)
	if t.message != "" {
		fmt.Fprintf(w, " %s", t.message)
	}
	if withError && t.Inner != nil {
		fmt.Fprintf(w, " (%s)", t.Inner)
	}
}

func (t Error) writeStack(w io.Writer, stack int) {
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

func (t Error) writeVerbose(w io.Writer, depth int) {
	if depth == 0 {
		depth = 1000
	}

	var err error = t
	i := 0

	for err != nil && i < depth {
		if d, ok := err.(Error); ok {
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
