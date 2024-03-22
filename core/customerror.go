package core

type CustomError struct {
	Code    int
	Message string
	Errors  any
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func (ce *CustomError) SetMessage(message string) *CustomError {
	ce.Message = message
	return ce
}

func (ce *CustomError) SetCode(code int) *CustomError {
	ce.Code = code
	return ce
}

func (ce *CustomError) SetErrors(errs any) *CustomError {
	ce.Errors = errs
	return ce
}
