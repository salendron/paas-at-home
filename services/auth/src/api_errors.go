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
	"encoding/json"
	"log"
	"net/http"
)

// ErrorCode defines all possible errors codes of this service
type ErrorCode int

// ErrorCodes
const (
	ErrorCodeInternal                ErrorCode = 0
	ErrorCodeLoginFailed                       = 1
	ErrorCodeRefreshFailed                     = 2
	ErrorCodeIDIsMissing                       = 3
	ErrorCodeInvalidRequestBody                = 4
	ErrorCodeUnexpectedSigningMethod           = 5
	ErrorCodeInvalidToken                      = 6
	ErrorCodeTokenExpired                      = 7
)

// ErrorMessage holds all information of a certain error
type ErrorMessage struct {
	Message    string    `json:"message"`
	StatusCode int       `json:"status"`
	Code       ErrorCode `json:"code"`
}

// RaiseError logs and returns a given error via on the current http request
func RaiseError(w http.ResponseWriter, message string, statusCode int, code ErrorCode) {
	errorMessage := ErrorMessage{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}

	log.Printf("Error: %v. HTTP Status: %v Code: %v", errorMessage.Message, errorMessage.StatusCode, errorMessage.Code)

	response := ErrorMessageType{
		Error: errorMessage,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
