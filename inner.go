package whoops

// InnerError wraps an inner error which's message
// is returned by calling Error() on it and
// which can be unwrapped using Unwrap().
//
// InnerError is mostly used as anonymous field by
// other errors to "inherit" the unwrap
// functionality of contained errors.
type InnerError struct {
	Inner error
}

func (t InnerError) Error() string {
	return t.Inner.Error()
}

func (t InnerError) Unwrap() error {
	return t.Inner
}
