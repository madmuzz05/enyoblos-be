package syserror

import "strings"

type NewError struct {
	Err     error
	Code    int
	Message string
}

type SysError interface {
	GetError() error
	GetStatusCode() int
	GetMessage() string
}

func CreateError(err error, code int, message string) SysError {
	message = strings.Trim(message, " ")
	sysErr := &NewError{
		Message: message,
		Code:    500,
		Err:     err,
	}

	if err != nil {
		sysErr.Err = err
	}

	if code != 0 {
		sysErr.Code = code
	}

	return sysErr
}

func (e *NewError) GetError() error {
	return e.Err
}

func (e *NewError) GetStatusCode() int {
	return e.Code
}

func (e *NewError) GetMessage() string {
	return e.Message
}
