package contentpipeline

import "errors"

// Domain-specific errors for the ContentPipeline bounded context.
var (
	ErrTitleRequired      = errors.New("title is required")
	ErrProductURLRequired = errors.New("product URL is required")
	ErrJobNotFound        = errors.New("production job not found")
	ErrInvalidTransition  = errors.New("invalid job state transition")
)
