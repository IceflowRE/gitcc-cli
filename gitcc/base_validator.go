package gitcc

type BaseValidator struct {
	Options map[string]string
}

func (v *BaseValidator) SetOptions(options map[string]string) error {
	v.Options = options
	return nil
}
