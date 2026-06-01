package cli

import "github.com/IceflowRE/gitcc/v3/standalone/gitcc"

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
		return &ExitError{Code: 2}
	case gitcc.Valid:
	default:
		return nil
	}

	return nil
}
