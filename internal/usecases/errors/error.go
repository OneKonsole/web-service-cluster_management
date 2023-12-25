package errors

type ErrorNotFound struct {
	Err error
}

func (e *ErrorNotFound) Error() string {
	return e.Err.Error()
}

func NewErrorNotFound(err error) *ErrorNotFound {
	return &ErrorNotFound{Err: err}
}

type ErrorBusinessException struct {
	Err error
}

func (e *ErrorBusinessException) Error() string {
	return e.Err.Error()
}

func NewErrorBusinessException(err error) *ErrorBusinessException {
	return &ErrorBusinessException{Err: err}
}
