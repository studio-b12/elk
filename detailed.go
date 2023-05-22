package whoops

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	callStack CallStack
}

var (
	_ HasFormat    = (*DetailedError)(nil)
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

// WrapMessage is the same as wrap including an
// additional message. The message will be shown
// in place of the wrapped errors result of the
// Error() method.
func Wrap(code ErrorCode, err error, message ...string) DetailedError {
	var d DetailedError
	d.code = code

	if dErr, ok := As[HasCallStack](err); ok {
		d.callStack = dErr.CallStack()
	} else {
		d.callStack = getCallFrames(1, maxCallStackDepth)
	}

	d.Inner = err

	if len(message) > 0 {
		d.message = strings.Join(message, " ")
	}

	return d
}

// Error returns the error information as
// formatted string.
func (t DetailedError) Error() string {
	var b bytes.Buffer
	t.writeTitle(&b)
	return b.String()
}

// Formatted returns the errors detailed
// information as formatted string.
func (t DetailedError) Formatted() string {
	var sb strings.Builder

	t.writeTitle(&sb)

	sb.WriteByte('\n')

	fmt.Fprintf(&sb, "error code: %s\n", t.code)

	if len(t.callStack.Frames()) > 0 {
		sb.WriteString("caused at:\n")
		t.callStack.WriteIndent(&sb, 5, indent)
	}

	return sb.String()
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
func (t DetailedError) CallStack() CallStack {
	return t.callStack
}

func (t DetailedError) writeTitle(w io.Writer) {
	if t.message != "" {
		fmt.Fprintf(w, "%s (%s)", t.message, t.Inner.Error())
	} else {
		fmt.Fprintf(w, "%s", t.Inner.Error())
	}
}
