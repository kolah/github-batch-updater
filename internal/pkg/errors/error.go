package errors

import (
	stdErrors "errors"
)

type SlugError struct {
	message      string
	kind         ErrorKind
	slug         string
	wrappedError error
}

func NewSlugError(message string, slug string, kind ErrorKind) SlugError {
	return SlugError{
		message: message,
		kind:    kind,
		slug:    slug,
	}
}

func NotFoundError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorKindNotFound())
}

func InternalError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorKindInternal())
}

func AuthorizationError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorKindAuthorization())
}

func ConflictError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorKindConflict())
}

func (e SlugError) Error() string {
	if e.wrappedError != nil {
		return e.message + ": " + e.wrappedError.Error()
	}

	return e.message
}

func (e SlugError) Slug() string {
	return e.slug
}

func (e SlugError) Kind() ErrorKind {
	return e.kind
}

func (e SlugError) Is(err error) bool {
	if err == nil {
		return false
	}

	var slugErr SlugError
	return stdErrors.As(err, &slugErr) &&
		slugErr.message == e.message &&
		slugErr.Slug() == e.Slug() &&
		slugErr.Kind() == e.Kind()
}

func (e SlugError) WrapError(err error) SlugError {
	return SlugError{
		message:      e.message,
		kind:         e.kind,
		slug:         e.slug,
		wrappedError: err,
	}
}
