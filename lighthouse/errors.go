package lighthouse

import (
	"errors"
	"fmt"
)

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Body    []byte `json:"-"`
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("lighthouse: status=%d code=%s msg=%s", e.Status, e.Code, e.Message)
}

func AsError(err error) (*Error, bool) {
	var le *Error
	ok := errors.As(err, &le)
	return le, ok
}
