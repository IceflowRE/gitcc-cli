package cli

import "github.com/IceflowRE/gitcc-cli/v3/gitcc"

// ExitError is a custom error type carrying an exit code. It does not carry any error.
type ExitError struct {
	Code int
}

func (*ExitError) Error() string {
	return ""
}

func getExitErrorFromStatus(status gitcc.Status) error {
	switch status {
	case gitcc.Invalid:
		return &ExitError{Code: 1}
	case gitcc.Warning:
		return &ExitError{Code: 2} //nolint:mnd
	case gitcc.Valid:
	default:
		return nil
	}

	return nil
}
