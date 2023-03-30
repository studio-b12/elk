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

// HasDetails describes an error which has
// additional details of any type.
type HasDetails interface {
	error

	// Details returns the value for details.
	Details() any
}

type HasCallStack interface {
	error

	CallStack() CallStack
}
