package errors

import "fmt"

type ErrorCode int

const (
	ErrFile             ErrorCode = 1000 // file error
	ErrConfigValidation ErrorCode = 1001 // configuration file validation fail
	ErrConfigParsing    ErrorCode = 1002 // YAML parsing error
)

type Data map[string]any

type Error struct {
	Data Data
}

func (e *Error) Error() string {
	return fmt.Sprintf("error %+v", e.Data)
}

func New(code ErrorCode, message string, data Data) *Error {
	if data == nil {
		data = Data{}
	}

	data["code"] = code
	data["msg"] = message

	return &Error{Data: data}
}

// Wrap wraps a stdlib error into our application error (adds an error code)
func Wrap(code ErrorCode, err error) *Error {
	return New(code, err.Error(), nil)
}
