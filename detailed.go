package whoops

import (
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

// Wrap takes an error and optionally some additional
// details
func Wrap(err error, details ...any) error {
	return WrapMessage(err, "", details...)
}

func WrapMessage(err error, message string, details ...any) error {
	var d DetailedError

	d.Inner = err
	d.message = message
	d.callStack = getCallFrames(3, maxCallStackDepth)

	if len(details) > 1 {
		d.details = details
	} else if len(details) == 1 {
		d.details = details[0]
	}

	return d
}

func (t DetailedError) Error() string {
	var sb strings.Builder

	if t.message != "" {
		fmt.Fprintf(&sb, "%s (%s)", t.message, t.Inner.Error())
	} else {
		sb.WriteString(t.Inner.Error())
	}

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

	if len(t.callStack) > 0 {
		sb.WriteString("caused at:\n")
		t.callStack.WriteIndent(&sb, 5, indent)
	}

	return sb.String()
}

func (t DetailedError) Message() string {
	return t.message
}

func (t DetailedError) CallStack() CallStack {
	return t.callStack
}

func (t DetailedError) Details() any {
	return t.details
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
