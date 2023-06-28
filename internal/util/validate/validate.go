package validate

type Validatable interface {
	Validate() error
}

func Each(v ...Validatable) error {
	for _, validatable := range v {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ErrorAsValidatable(err error) Validatable {
	return errorValidatable{err: err}
}

type errorValidatable struct {
	err error
}

func (e errorValidatable) Validate() error {
	return e.err
}
