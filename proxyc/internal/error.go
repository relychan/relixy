package internal

import "fmt"

// MetadataValidationError represents the error structure of the metadata validation error.
type MetadataValidationError struct {
	Path    string         `json:"path"`
	Method  string         `json:"method"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Error implements the error interface.
func (mve MetadataValidationError) Error() string {
	return fmt.Sprintf("%s %s: %s", mve.Method, mve.Path, mve.Message)
}
