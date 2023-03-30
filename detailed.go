package whoops

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

	message   string
	details   any
	callStack CallStack
}

var _ Formatted = (*DetailedError)(nil)

// Wrap takes an error and optionally some additional
// details and created a DetailedError from that. Also,
// the call stack from where this method has been called
// is embedded in the DetailedError.
func Wrap(err error, details ...any) DetailedError {
	dErr := WrapMessage(err, "", details...)
	dErr.callStack.offset++
	return dErr
}

// WrapMessage is the same as wrap including an
// additional message. The message will be shown
// in place of the wrapped errors result of the
// Error() method.
func WrapMessage(err error, message string, details ...any) DetailedError {
	var d DetailedError

	d.Inner = err
	d.message = message
	d.callStack = getCallFrames(1, maxCallStackDepth)

	if len(details) > 1 {
		d.details = details
	} else if len(details) == 1 {
		d.details = details[0]
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

	if t.details != nil {
		sb.WriteString("details:\n")
		if dSlice, ok := t.details.([]any); ok {
			for _, d := range dSlice {
				writeDetails(&sb, d, indent)
			}
		} else {
			writeDetails(&sb, t.details, indent)
		}
	}

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

// CallStack returns the errors CallStack
// starting from where the DetailedError
// has been created.
func (t DetailedError) CallStack() CallStack {
	return t.callStack
}

// Details returns the errros details,
// if specified.
func (t DetailedError) Details() any {
	return t.details
}

func (t DetailedError) writeTitle(w io.Writer) {
	if t.message != "" {
		fmt.Fprintf(w, "%s (%s)", t.message, t.Inner.Error())
	} else {
		fmt.Fprintf(w, "%s", t.Inner.Error())
	}
}

func writeDetails(w io.Writer, v any, prefix string) {
	var details string

	switch vt := v.(type) {
	case string:
		details = vt + "\n"
	case interface{ String() string }:
		details = vt.String() + "\n"
	case io.Reader:
		b, err := io.ReadAll(vt)
		if err == nil {
			details = string(b) + "\n"
		}
	default:
		details = spew.Sdump(v)
	}

	lines := strings.Split(details, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			break
		}
		lines[i] = prefix + line
	}

	w.Write([]byte(strings.Join(lines, "\n")))
}
