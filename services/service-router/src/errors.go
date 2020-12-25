/*
api_errors.go
Defines all error types and messages of the the RESTful API.

###################################################################################

MIT License

Copyright (c) 2020 Bruno Hautzenberger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

// ErrorCode defines all possible errors codes of this service
const (
	ErrorCodeInternal             = 0
	ErrorCodeInvalidRequestBody   = 1
	ErrorCodeUnknownTargetService = 2
)

// RaiseError logs and returns a given error via on the current http request
func RaiseError(ctx echo.Context, message string, statusCode int, code int) error {
	errorMessage := Error{
		Message: &message,
		Status:  &statusCode,
		Code:    &code,
	}

	log.Printf("Error: %v. HTTP Status: %v Code: %v", errorMessage.Message, errorMessage.Code, errorMessage.Code)

	return ctx.JSON(statusCode, errorMessage)
}
