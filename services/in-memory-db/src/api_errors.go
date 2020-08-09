package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorCode defines all possible errors codes of this service
type ErrorCode int

const (
	ErrorCodeInternal           ErrorCode = 0
	ErrorCodeRealmMissing                 = 1
	ErrorCodeKeyMissing                   = 2
	ErrorCodeEntityNotFound               = 3
	ErrorCodeInvalidRequestBody           = 4
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
