package errors

import "fmt"

type ErrorCode int

type Data map[string]any

type Error struct {
	Data Data
}

func (e *Error) Error() string {
	return fmt.Sprintf("error %+v", e.Data)
}

func (e *Error) PrependPath(msg string, a ...any) {
	prev, _ := e.Data["path"].(string)
	e.Data["path"] = fmt.Sprintf(msg, a...) + "/" + prev
}

func New(key string, data Data) *Error {
	if data == nil {
		data = Data{}
	}

	data["key"] = key

	return &Error{Data: data}
}

// Wrap wraps a stdlib error into our application error (adds an error code)
func Wrap(key string, err error) *Error {
	// Don't wrap errors that are already our custom type
	if e, ok := err.(*Error); ok {
		return e
	}

	return &Error{
		Data: Data{
			"key":     key,
			"message": err.Error(),
		},
	}
}
