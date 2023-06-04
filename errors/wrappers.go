package errors

import "errors"

// These variables are used to give access to existing functions in the standard
// library `errors` package
var (
	As  = errors.As
	Is  = errors.Is
	New = errors.New
)
