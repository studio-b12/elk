package whoops

// Formatted describes an error with additional
// information which can be accessed as a
// formatted string.
type Formatted interface {
	error

	// Formatted returns the error details
	// as formatted string.
	Formatted() string
}
