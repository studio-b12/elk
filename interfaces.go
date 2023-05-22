package whoops

// HasFormat describes an error with additional
// information which can be accessed as a
// formatted string.
type HasFormat interface {
	error

	// Formatted returns the error details
	// as formatted string.
	Formatted() string
}

// HasMessage describes an error which has an
// additional message.
type HasMessage interface {
	error

	// Message returns the value for message.
	Message() string
}

// HasCode describes an error which has an
// ErrorCode.
type HasCode interface {
	error

	// Code returns the inner ErrorCode of
	// the error.
	Code() ErrorCode
}

type HasCallStack interface {
	error

	CallStack() CallStack
}
