/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ResponseData is a struct that holds the response data for an API request.
type ResponseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponseData(data interface{}) ResponseData {
	return ResponseData{
		Data: data,
	}
}

// nolint:golint,unused
func newResponseCodeError(code string, err error) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  msg,
	}
}

// SendBadRequestBody sends a bad request body error response.
func SendBadRequestBody(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		resp := newResponseCodeMsg(errorBadRequestBody, err.Error())
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var fields []string
			for _, fieldError := range ve {
				fields = append(fields, fieldError.StructField())
			}
			resp = newResponseCodeMsg(errorValidationFailed,
				fmt.Sprintf("validate:%s failed", strings.Join(fields, ", ")))
		}
		ctx.JSON(http.StatusBadRequest, resp)
	}
}

// SendRespOfPost sends a successful POST response with data if provided.
func SendRespOfPost(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(data))
	}
}

// SendError sends an error response based on the given error.
func SendError(ctx *gin.Context, err error) {
	sc, code := httpError(err)

	//_ = ctx.AbortWithError(sc, allerror.InnerErr(err))

	ctx.JSON(sc, newResponseCodeMsg(code, err.Error()))
}
