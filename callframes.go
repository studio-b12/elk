package whoops

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// CallStack adds some additional utility
// to a slice of runtime.Frame.
type CallStack []runtime.Frame

// WriteIndent is an alias for write with the given
// indent string attached before each line of output.
func (t CallStack) WriteIndent(w io.Writer, max int, indent string) {
	if max > 0 && len(t) > max {
		t = t[:max]
	}

	maxLenFName := 0
	for _, frame := range t {
		if l := len(frame.Function); l > maxLenFName {
			maxLenFName = l
		}
	}
	for _, frame := range t {
		fmt.Fprintf(w, "%s%-"+strconv.Itoa(maxLenFName)+"s\t%s:%d\n",
			indent, frame.Function, frame.File, frame.Line)
	}
}

// Write formats the call stack into a table of called
// function and the file plus line number and writes
// the result into the writer w.
//
// max defines the number of stack frames which are
// printed starting from the original caller.
func (t CallStack) Write(w io.Writer, max int) {
	t.WriteIndent(w, max, "")
}

// String returns the formatted output of the callstack
// as string.
func (t CallStack) String() string {
	var b bytes.Buffer
	t.Write(&b, 0)
	return b.String()
}

func getCallFrames(skip, n int) CallStack {
	callerPtrs := make([]uintptr, n)
	nPtrs := runtime.Callers(skip+1, callerPtrs)
	frameCursor := runtime.CallersFrames(callerPtrs[:nPtrs])

	callFrames := make([]runtime.Frame, 0, nPtrs)
	for {
		frame, next := frameCursor.Next()
		if !next {
			break
		}
		callFrames = append(callFrames, frame)
	}

	return callFrames
}
