package errors

import "fmt"

type ErrorCode int

const (
	ErrFile             ErrorCode = 1000 // file error
	ErrConfigValidation ErrorCode = 1001 // configuration file validation fail
	ErrConfigParsing    ErrorCode = 1002 // YAML parsing error
)

type Error struct {
	Code ErrorCode
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Msg)
}

// New creates an error with code and formatted message.
// Works like fmt.Errorf: New(code, "message") or New(code, "version %d", version)
func New(code ErrorCode, format string, a ...interface{}) *Error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, a...)}
}

// Wrap wraps a stdlib error into our application error (adds an error code)
func Wrap(code ErrorCode, err error) *Error {
	return &Error{Code: code, Msg: err.Error()}
}
