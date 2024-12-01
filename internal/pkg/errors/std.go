package errors

import stdErrors "errors"

// As is a convenience proxy method for errors.As.
func As[T error](err error, target *T) bool {
	return stdErrors.As(err, target)
}

// Is function is  a proxy method for errors.Is.
func Is(err, target error) bool {
	return stdErrors.Is(err, target)
}
