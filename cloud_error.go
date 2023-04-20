package errors

import (
	"encoding/json"
	"time"
)

const (
	InternalServerError CustomCode = "InternalServerError"
	NotFound            CustomCode = "NotFound"
)

type CloudError struct {
	StatusCode    int           `json:"status_code"`
	Status        string        `json:"status"`
	Message       string        `json:"message"`
	Source        string        `json:"source"`
	TimeStamp     time.Time     `json:"timestamp"`
	CustomCode    CustomCode    `json:"custom_code"`
	ErrorLocation ErrorLocation `json:"location,omitempty"`
	CorrelationID string        `json:"correlation_id"`
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

func (se *CloudError) Error() string {
	byt, _ := json.MarshalIndent(se, "", "  ")

	return string(byt)
}

type CloudErrorOption func(*CloudError)

func NewCloudError(statusCode int, message string, options ...CloudErrorOption) *CloudError {
	se := NewCloudErrorBuilder().
		StatusCode(statusCode).
		Message(message).
		Build(time.Now().UTC(), options...)

	return se
}
