/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"errors"
	"net/http"
)

const (
	errorSystemError      = "system_error"
	errorBadRequestBody   = "bad_request_body"
	errorValidationFailed = "validation_failed"
	errorBadRequestParam  = "bad_request_param"
)

type errorCode interface {
	ErrorCode() string
}

type errorNotFound interface {
	errorCode

	NotFound()
}

type errorNoPermission interface {
	errorCode

	NoPermission()
}

func httpError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	sc := http.StatusInternalServerError
	code := errorSystemError
	var t errorCode

	if ok := errors.As(err, &t); ok {
		code = t.ErrorCode()

		var n errorNotFound
		var p errorNoPermission

		if ok := errors.As(err, &n); ok {
			sc = http.StatusNotFound

		} else if ok := errors.As(err, &p); ok {
			sc = http.StatusForbidden

		} else {
			switch code {
			default:
				sc = http.StatusBadRequest
			}
		}
	}

	return sc, code
}
