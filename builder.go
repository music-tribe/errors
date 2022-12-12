package errors

import (
	"net/http"
	"runtime"
	"strings"
	"time"
)

type storageErrorBuilder struct {
	err *StorageError
}

func NewStorageErrorBuilder() *storageErrorBuilder {
	return &storageErrorBuilder{
		&StorageError{
			ErrorLocation: ErrorLocation{
				skip: 2,
			},
		},
	}
}

func (s *storageErrorBuilder) StatusCode(statusCode int) *storageErrorBuilder {
	if statusCode < 100 || statusCode > 599 {
		statusCode = 500
	}
	s.err.StatusCode = statusCode
	s.err.Status = http.StatusText(statusCode)
	return s
}

func (s *storageErrorBuilder) Message(errMsg string) *storageErrorBuilder {
	s.err.Message = errMsg
	return s
}

func (s *storageErrorBuilder) ErrorLocation(svc, pkg, fnc string) *storageErrorBuilder {
	s.err.ErrorLocation.Service = svc
	s.err.ErrorLocation.Method = fnc
	return s
}

func (s *storageErrorBuilder) CustomCode(code CustomCode) *storageErrorBuilder {
	s.err.CustomCode = code
	return s
}

func (s *storageErrorBuilder) CorrelationID(id string) *storageErrorBuilder {
	s.err.CorrelationID = id
	return s
}

func (s *storageErrorBuilder) Tags(tags ...string) *storageErrorBuilder {
	s.err.Tags = append(s.err.Tags, tags...)
	return s
}

// SkipCaller allows you to skip levels of the trace when trying to determine in which
// method the errors was called.
func (s *storageErrorBuilder) SkipCaller(skip int) *storageErrorBuilder {
	s.err.ErrorLocation.skip = skip
	return s
}

func (s *storageErrorBuilder) Build(t time.Time, options ...StorageErrorOption) *StorageError {
	if s.err.StatusCode == 0 {
		s.err.StatusCode = 500
		s.err.Status = http.StatusText(500)
	}
	if s.err.Message == "" {
		s.err.Message = s.err.Status
	}
	if s.err.CustomCode == "" {
		s.err.CustomCode = CustomCode(strings.ReplaceAll(s.err.Status, " ", ""))
	}

	pc, page, line, _ := runtime.Caller(s.err.ErrorLocation.skip)
	funcDetails := runtime.FuncForPC(pc)

	name := funcDetails.Name()
	s.err.ErrorLocation.Page = page
	s.err.ErrorLocation.Line = line
	s.err.ErrorLocation.Method = name

	s.err.TimeStamp = t

	for _, option := range options {
		option(s.err)
	}

	// if the call trace skip level has not been changed, return the error
	if s.err.ErrorLocation.skip == 2 {
		return s.err
	}

	pc, s.err.ErrorLocation.Page, s.err.ErrorLocation.Line, _ = runtime.Caller(s.err.ErrorLocation.skip)
	funcDetails = runtime.FuncForPC(pc)

	s.err.ErrorLocation.Method = funcDetails.Name()

	return s.err
}
