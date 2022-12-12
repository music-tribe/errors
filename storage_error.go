package errors

import (
	"encoding/json"
	"time"
)

const (
	InternalServerError CustomCode = "InternalServerError"
	NotFound            CustomCode = "NotFound"
)

type StorageError struct {
	StatusCode    int           `json:"statusCode"`
	Status        string        `json:"status"`
	Message       string        `json:"message"`
	TimeStamp     time.Time     `json:"timeStamp"`
	CustomCode    CustomCode    `json:"customCode"`
	ErrorLocation ErrorLocation `json:"location,omitempty"`
	CorrelationID string        `json:"correlationID"`
	Tags          []string      `json:"tags,omitempty"`
}

type ErrorLocation struct {
	Service string `json:"service,omitempty"`
	Method  string `json:"method,omitempty"`
	Page    string `json:"page,omitempty"`
	Line    int    `json:"line,omitempty"`
	skip    int    `json:"-"`
}

type CustomCode string

func (se *StorageError) Error() string {
	byt, _ := json.MarshalIndent(se, "", "  ")

	return string(byt)
}

type StorageErrorOption func(*StorageError)

func NewStorageError(statusCode int, message string, options ...StorageErrorOption) *StorageError {
	se := NewStorageErrorBuilder().
		StatusCode(statusCode).
		Message(message).
		Build(time.Now().UTC(), options...)

	return se
}
