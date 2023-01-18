package pixiv_api_go

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound means the api return 404
	ErrNotFound = errors.New("NotFound")
)

type ErrorJsonUnmarshal struct {
	err    error
	rawStr string
}

func NewJsonUnmarshalErr(date []byte, err error) error {
	return &ErrorJsonUnmarshal{err: err, rawStr: string(date)}
}

func (j *ErrorJsonUnmarshal) Error() string {
	return fmt.Sprintf("failed to unmarshal json, err: %s, raw: %s", j.err, j.rawStr)
}
