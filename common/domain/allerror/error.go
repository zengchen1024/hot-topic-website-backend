/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package allerror provides a set of error codes and error types used in the application.
package allerror

import (
	"errors"
	"strings"
)

// New creates a new error with the specified code and message.
func New(code string, msg string, err error) errorImpl {
	v := errorImpl{
		code:     code,
		innerErr: err,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

func IsError(err error, code string) bool {
	v := errorImpl{}
	if errors.As(err, &v) && v.code == code {
		return true
	}

	return false
}

// errorImpl
type errorImpl struct {
	code     string
	msg      string
	innerErr error // error info for diagnostic
}

// Error returns the error message.
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode returns the error code.
func (e errorImpl) ErrorCode() string {
	return e.code
}

// InnerErr returns the inner error.
func (e errorImpl) InnerError() error {
	return e.innerErr
}
