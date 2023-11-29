package dto

import (
	"source/internal/errors"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Errors  []Error     `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	ID        string `json:"id,omitempty"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

type ResponseDelete struct {
	ID      int64  `json:"id"`
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
}

type Delete struct {
	ID int64
}

func Done(data interface{}, mess ...string) *Response {
	return OK(data, mess...)
}

func OK(data interface{}, mess ...string) *Response {
	var message string
	if len(mess) > 0 {
		message = mess[0]
	}
	return &Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

func Fail(errs ...error) *Response {
	var outErrs []Error
	for _, err := range errs {
		if err == nil {
			continue
		}
		outErr := Error{
			Message: err.Error(),
		}
		if errParse, ok := err.(errors.Error); ok {
			outErr = Error{
				ID:        errParse.ID(),
				Message:   errParse.Error(),
				ErrorCode: errParse.Code(),
			}
		}
		outErrs = append(outErrs, outErr)
	}
	return &Response{
		Status: false,
		Errors: outErrs,
	}
}

func FailWithString(messages ...string) *Response {
	var errs []Error
	for _, mess := range messages {
		errs = append(errs, Error{Message: mess})
	}
	return &Response{
		Status: false,
		Errors: errs,
	}
}

func MakeResponseSuccess(object interface{}, mess ...string) *Response {
	var message string
	if len(mess) > 0 {
		message = mess[0]
	}
	return &Response{
		Status:  true,
		Message: message,
		Data:    object,
	}
}

func MakeResponseError(errs ...error) *Response {
	var errors []Error
	for _, err := range errs {
		errors = append(errors, Error{Message: err.Error()})
	}
	return &Response{
		Status: false,
		Errors: errors,
	}
}

func MakeResponseErrorString(messages ...string) *Response {
	var errors []Error
	for _, message := range messages {
		errors = append(errors, Error{Message: message})
	}
	return &Response{
		Status: false,
		Errors: errors,
	}
}

func MakeResponseErrorWithID(err ...Error) *Response {
	return &Response{
		Status: false,
		Errors: err,
	}
}
