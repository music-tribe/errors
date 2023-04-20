package errors

import (
	"net/http"
	"runtime"
	"strings"
	"time"
)

type cloudErrorBuilder struct {
	err *CloudError
}

func NewCloudErrorBuilder() *cloudErrorBuilder {
	return &cloudErrorBuilder{
		&CloudError{
			ErrorLocation: ErrorLocation{
				skip: 2,
			},
		},
	}
}

func (s *cloudErrorBuilder) StatusCode(statusCode int) *cloudErrorBuilder {
	if statusCode < 100 || statusCode > 599 {
		statusCode = 500
	}
	s.err.StatusCode = statusCode
	s.err.Status = http.StatusText(statusCode)
	return s
}

func (s *cloudErrorBuilder) Message(errMsg string) *cloudErrorBuilder {
	s.err.Message = errMsg
	return s
}

func (s *cloudErrorBuilder) ErrorLocation(svc, pkg, fnc string) *cloudErrorBuilder {
	s.err.ErrorLocation.Service = svc
	s.err.ErrorLocation.Method = fnc
	return s
}

func (s *cloudErrorBuilder) CustomCode(code CustomCode) *cloudErrorBuilder {
	s.err.CustomCode = code
	return s
}

func (s *cloudErrorBuilder) CorrelationID(id string) *cloudErrorBuilder {
	s.err.CorrelationID = id
	return s
}

func (s *cloudErrorBuilder) Source(name string) *cloudErrorBuilder {
	s.err.Source = name
	return s
}

func (s *cloudErrorBuilder) Tags(tags ...string) *cloudErrorBuilder {
	s.err.Tags = append(s.err.Tags, tags...)
	return s
}

// SkipCaller allows you to skip levels of the trace when trying to determine in which
// method the errors was called.
func (s *cloudErrorBuilder) SkipCaller(skip int) *cloudErrorBuilder {
	s.err.ErrorLocation.skip = skip
	return s
}

func (s *cloudErrorBuilder) Build(t time.Time, options ...CloudErrorOption) *CloudError {
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
	if s.err.Source == "" {
		s.err.Source = "music-tribe"
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
