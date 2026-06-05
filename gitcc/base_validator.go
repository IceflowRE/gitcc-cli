// Package gitcc provides the core functionality for gitcc.
package gitcc

// BaseValidator is a base implementation of the Validator interface.
// It provides a default implementation of the SetOptions method, which simply stores the options in a field.
// This can be embedded in custom validators to avoid having to implement SetOptions if no special handling is needed.
type BaseValidator struct {
	Options map[string]string
}

// SetOptions stores the provided options in the BaseValidator's Options field.
func (v *BaseValidator) SetOptions(options map[string]string) error {
	v.Options = options

	return nil
}
