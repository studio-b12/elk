package whoops

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// CallFrame is a type alias for runtime.Frame with additional formatting
// functionailty used in downstream functions.
type CallFrame runtime.Frame

func (t CallFrame) String() string {
	return fmt.Sprintf("%s", t)
}

func (t CallFrame) Format(s fmt.State, verb rune) {
	width, hasWidth := s.Width()

	switch verb {
	case 's':
		format := "%s %s:%d"
		if hasWidth {
			format = "%-" + strconv.Itoa(width) + "s\t%s:%d"
		}
		fmt.Fprintf(s, format, t.Function, t.File, t.Line)
	case 'v':
		fmt.Fprintf(s, "%v", runtime.Frame(t))
	}
}

// CallStack contains the list of called
// runtime.Frames in the call chain with
// an offset from which frames are
// reported.
type CallStack struct {
	frames []CallFrame
	offset int
}

// Frames returns the offset slice of called
// runtime.Frame's in the recorded call stack.
func (t CallStack) Frames() []CallFrame {
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
		fmt.Fprintf(w, "%s%"+strconv.Itoa(maxLenFName)+"s\n", indent, frame)
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

// At returns the formatted call frame at the given position n
// if existent.
func (t CallStack) At(n int) (s string, ok bool) {
	n += t.offset

	if n >= len(t.frames) || n < 0 {
		return "", false
	}

	frame := t.frames[n]
	return frame.String(), true
}

func getCallFrames(offset, n int) CallStack {
	callerPtrs := make([]uintptr, n)
	nPtrs := runtime.Callers(2, callerPtrs)
	frameCursor := runtime.CallersFrames(callerPtrs[:nPtrs])

	callFrames := make([]CallFrame, 0, nPtrs)
	for {
		frame, more := frameCursor.Next()
		callFrames = append(callFrames, CallFrame(frame))
		if !more {
			break
		}
	}

	return CallStack{
		frames: callFrames,
		offset: offset,
	}
}
