package errors

import (
	"errors"
)

type Error interface {
	error
	Code() string
	ID() string
	Origin() error
	Line() int
	File() string
}

type appError struct {
	id      string
	message string
	code    string
	file    string
	line    int
	err     error
}

func (e *appError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *appError) Unwrap() error {
	return e.err
}

func (e *appError) ID() string {
	return e.id
}

func (e *appError) Code() string {
	return e.code
}

func (e *appError) File() string {
	return e.file
}

func (e *appError) Line() int {
	return e.line
}

func (e *appError) Origin() error {
	return e.err
}

func New(message string, args ...string) error {
	var id, code, file string
	var line int
	//_, file, line, _ = runtime.Caller(1)
	if len(args) == 1 {
		code = args[0]
	}
	if len(args) == 2 {
		code = args[0]
		id = args[1]
	}
	return &appError{
		id:      id,
		message: message,
		err:     errors.New(message),
		line:    line,
		file:    file,
		code:    code,
	}
}

func NewWithID(message string, id string, args ...string) error {
	var code string
	if len(args) == 1 {
		code = args[0]
	}
	return &appError{
		id:      id,
		message: message,
		code:    code,
		err:     errors.New(message),
	}
}

func Wrap(err error, args ...string) error {
	var (
		totalArgs         = len(args)
		message, id, code string
	)

	if totalArgs >= 3 {
		message = args[0]
		id = args[1]
		code = args[2]
	} else if totalArgs == 2 {
		message = args[0]
		id = args[1]
	} else if totalArgs == 1 {
		message = args[0]
	}

	if err == nil {
		err = errors.New(message)
	}
	return &appError{
		message: message,
		id:      id,
		code:    code,
		err:     err,
	}
}

func Origin(err error) error {
	if errOrigin, ok := err.(Error); ok {
		return errOrigin.Origin()
	}
	return nil
}
