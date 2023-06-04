package errors

// publicError defines a new type that wraps the standard `error` type with an
// custom public error message
type publicError struct {
	err error
	msg string
}

// Public returns the custom-type `publicError`
func Public(err error, msg string) error {
	return publicError{err, msg}
}

// Error returns the string representation of the custom-type `publicError`.
// This is the implementation of the standard-library `Error()` interface
func (pe publicError) Error() string {
	return pe.err.Error()
}

// Public returns the public error message associated with a `publicError`
func (pe publicError) Public() string {
	return pe.msg
}

// Unwrap unwraps the underlying `error` from the custom-type `publicError`
func (pe publicError) Unwrap() error {
	return pe.err
}
