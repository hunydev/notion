package api

import (
	"encoding/json"
	"fmt"
	"io"
)

type Error struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s", err.Status, err.Code, err.Message)
}

func (err *Error) String() string {
	return err.Error()
}

func ReadError(r io.Reader) error {
	err := &Error{}

	d := json.NewDecoder(r)
	if decodeError := d.Decode(err); decodeError != nil {
		return decodeError
	}

	return err
}

func ParseError(buf []byte) error {
	err := &Error{}

	if jsonError := json.Unmarshal(buf, err); jsonError != nil {
		return jsonError
	}

	return err
}
