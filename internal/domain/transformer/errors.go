package transformer

import "github.com/kolah/github-batch-updater/internal/pkg/errors"

func ProcessorNotFound() errors.SlugError {
	return errors.NotFoundError("processor not found", "transformer.processor_not_found")
}

func ProcessorError() errors.SlugError {
	return errors.InternalError("processor error", "processor.error")
}
