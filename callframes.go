package whoops

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// CallStack contains the list of called
// runtime.Frames in the call chain with
// an offset from which frames are
// reported.
type CallStack struct {
	frames []runtime.Frame
	offset int
}

// Frames returns the offset slice of called
// runtime.Frame's in the recorded call stack.
func (t CallStack) Frames() []runtime.Frame {
	if len(t.frames) < t.offset {
		return nil
	}
	return t.frames[t.offset:]
}

// WriteIndent is an alias for write with the given
// indent string attached before each line of output.
func (t CallStack) WriteIndent(w io.Writer, max int, indent string) {
	frames := t.Frames()

	if max > 0 && len(frames) > max {
		frames = frames[:max]
	}

	maxLenFName := 0
	for _, frame := range frames {
		if l := len(frame.Function); l > maxLenFName {
			maxLenFName = l
		}
	}
	for _, frame := range frames {
		format := "%s%-" + strconv.Itoa(maxLenFName) + "s\t%s:%d\n"
		fmt.Fprintf(w, format, indent, frame.Function, frame.File, frame.Line)
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

func getCallFrames(offset, n int) CallStack {
	callerPtrs := make([]uintptr, n)
	nPtrs := runtime.Callers(2, callerPtrs)
	frameCursor := runtime.CallersFrames(callerPtrs[:nPtrs])

	callFrames := make([]runtime.Frame, 0, nPtrs)
	for {
		frame, more := frameCursor.Next()
		callFrames = append(callFrames, frame)
		if !more {
			break
		}
	}

	return CallStack{
		frames: callFrames,
		offset: offset,
	}
}
