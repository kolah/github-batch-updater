package errors

type ErrorKind struct {
	kind string
}

func ErrorKindNotFound() ErrorKind {
	return ErrorKind{"not-found"}
}

func ErrorKindInternal() ErrorKind {
	return ErrorKind{"internal"}
}

func ErrorKindAuthorization() ErrorKind {
	return ErrorKind{"authorization"}
}
