package github

import "github.com/kolah/github-batch-updater/internal/pkg/errors"

func RepositoryNotFound() errors.SlugError {
	return errors.NotFoundError("repository not found", "github.repository.not-found")
}

func RepositoryAccessDenied() errors.SlugError {
	return errors.AuthorizationError("repository access denied", "github.repository.access-denied")
}

func RefAlreadyExists() errors.SlugError {
	return errors.ConflictError("ref already exists", "github.repository.ref-already-exists")
}

func UnknownError() errors.SlugError {
	return errors.InternalError("unknown error", "github.unknown")
}

func RefNotFound() errors.SlugError {
	return errors.NotFoundError("ref not found", "github.repository.ref-not-found")
}

func FileNotFound() errors.SlugError {
	return errors.NotFoundError("file not found", "github.repository.file-not-found")
}

func FileAccessDenied() errors.SlugError {
	return errors.AuthorizationError("file access denied", "github.repository.file-access-denied")
}

func UnableToDecodeFileContent() errors.SlugError {
	return errors.InternalError("unable to decode file content", "github.repository.unable-to-decode-file-content")
}

func FailedUpdatingFile() errors.SlugError {
	return errors.InternalError("failed updating file", "github.repository.failed-updating-file")
}
