package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"testing/quick"
	"time"
)

func TestNewCloudError(t *testing.T) {
	_, testPage, _, _ := runtime.Caller(0)

	When("the status code is not provided", t,
		Then("it should return a storage error that reports an internal server error",
			func(t *testing.T) {
				sc := 0
				wantSc := 500
				wantStatus := http.StatusText(wantSc)
				msg := "unknown error"
				timeNow := time.Now().UTC()

				pc, _, _, _ := runtime.Caller(0)
				rtFunc := runtime.FuncForPC(pc)
				line := 37
				setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
				el := ErrorLocation{
					Method: rtFunc.Name(),
					Page:   testPage,
					Line:   line,
					skip:   2,
				}

				setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
				want := CloudError{
					StatusCode:    wantSc,
					Status:        wantStatus,
					Message:       msg,
					TimeStamp:     timeNow,
					CustomCode:    InternalServerError,
					ErrorLocation: el,
				}

				if got := NewCloudError(sc, msg, setTimeOpt, setLine); !reflect.DeepEqual(*got, want) {
					t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
					return
				}
			},
		),
	)

	When("the status code is provided", t,
		Then("it should return a storage error with the correct status code", func(t *testing.T) {
			sc := 404
			wantSc := 404
			wantStatus := http.StatusText(wantSc)
			msg := "sorry, we couldn't find the blob you requested"
			timeNow := time.Now().UTC()

			pc, _, _, _ := runtime.Caller(0)
			rtFunc := runtime.FuncForPC(pc)
			line := 37
			setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
			el := ErrorLocation{
				Method: rtFunc.Name(),
				Page:   testPage,
				Line:   line,
				skip:   2,
			}

			setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
			want := CloudError{
				StatusCode:    wantSc,
				Status:        wantStatus,
				Message:       msg,
				TimeStamp:     timeNow,
				CustomCode:    NotFound,
				ErrorLocation: el,
			}

			if got := NewCloudError(sc, msg, setTimeOpt, setLine); !reflect.DeepEqual(*got, want) {
				t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
			}
		}),
	)

	When("the error message is not provided", t,
		Then("it should return a storage error with a message that defaults to the status", func(t *testing.T) {
			inputMsg := ""

			wantSc := 403
			wantStatus := http.StatusText(wantSc)
			timeNow := time.Now().UTC()

			pc, _, _, _ := runtime.Caller(0)
			rtFunc := runtime.FuncForPC(pc)
			line := 37
			setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
			el := ErrorLocation{
				Method: rtFunc.Name(),
				Page:   testPage,
				Line:   line,
				skip:   2,
			}

			setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
			want := CloudError{
				StatusCode:    wantSc,
				Status:        wantStatus,
				Message:       wantStatus,
				TimeStamp:     timeNow,
				CustomCode:    "Forbidden",
				ErrorLocation: el,
			}

			if got := NewCloudError(wantSc, inputMsg, setTimeOpt, setLine); !reflect.DeepEqual(*got, want) {
				t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
			}
		}),
	)

	When("the error message is provided", t,
		Then("it should return a storage error with that provided message", func(t *testing.T) {
			inputMsg := "hello"

			wantSc := 403
			wantStatus := http.StatusText(wantSc)
			timeNow := time.Now().UTC()

			pc, _, _, _ := runtime.Caller(0)
			rtFunc := runtime.FuncForPC(pc)
			line := 37
			setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
			el := ErrorLocation{
				Method: rtFunc.Name(),
				Page:   testPage,
				Line:   line,
				skip:   2,
			}

			setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
			want := CloudError{
				StatusCode:    wantSc,
				Status:        wantStatus,
				Message:       inputMsg,
				TimeStamp:     timeNow,
				CustomCode:    "Forbidden",
				ErrorLocation: el,
			}

			if got := NewCloudError(wantSc, inputMsg, setTimeOpt, setLine); !reflect.DeepEqual(*got, want) {
				t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
			}
		}),
	)

	When("a storage error option is passed", t,
		Then("it should return a storage error that reflects the change set by that option", func(t *testing.T) {
			svc := "presets"
			fnc := "NewUploadStorage"
			errLoc := ErrorLocation{
				Service: svc,
				Method:  fnc,
				skip:    2,
			}

			setErrorLocation := func(se *CloudError) {
				se.ErrorLocation = errLoc
			}

			wantSc := 403
			wantStatus := http.StatusText(wantSc)
			timeNow := time.Now().UTC()

			setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
			want := CloudError{
				StatusCode:    wantSc,
				Status:        wantStatus,
				Message:       wantStatus,
				TimeStamp:     timeNow,
				CustomCode:    "Forbidden",
				ErrorLocation: errLoc,
			}

			if got := NewCloudError(wantSc, "", setTimeOpt, setErrorLocation); !reflect.DeepEqual(*got, want) {
				t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
			}
		}),
	)

	When("the skip level is changed to 3", t,
		Then("the error location should reflect that change", func(t *testing.T) {
			skip := 3

			errLoc := ErrorLocation{
				Page: testPage,
				skip: skip,
			}

			wantSc := 403
			wantStatus := http.StatusText(wantSc)
			timeNow := time.Now().UTC()

			setSkip := func(se *CloudError) { se.ErrorLocation.skip = skip }
			setTimeOpt := func(se *CloudError) { se.TimeStamp = timeNow }
			want := CloudError{
				StatusCode:    wantSc,
				Status:        wantStatus,
				Message:       wantStatus,
				TimeStamp:     timeNow,
				CustomCode:    "Forbidden",
				ErrorLocation: errLoc,
			}

			pc, _, _line, _ := runtime.Caller(1)
			got := NewCloudError(wantSc, "", setTimeOpt, setSkip)

			want.ErrorLocation.Line = _line
			want.ErrorLocation.Method = runtime.FuncForPC(pc).Name()

			if !reflect.DeepEqual(*got, want) {
				t.Errorf("NewCloudError() = \n%v\n but want \n%v\n", *got, want)
			}
		}),
	)

}

func When(description string, t *testing.T, then ...func(t *testing.T)) {
	// return func(t *testing.T) {
	t.Run(fmt.Sprintf("When %s\n", description), func(t *testing.T) {
		for _, test := range then {
			test(t)
		}
	})
}

func Then(description string, do func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("\t\tThen %s\n", description)
		do(t)
	}
}

func And(description string, do func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("\t\tThen %s\n", description)
		do(t)
	}
}

func TestQuick(t *testing.T) {
	var se *CloudError
	f := func(statusCode int, message string) bool {
		se = NewCloudError(statusCode, message)

		return !se.TimeStamp.IsZero() &&
			(se.StatusCode == statusCode || se.StatusCode == 500) &&
			se.Status != "" &&
			(se.Message == message || se.Message == http.StatusText(500)) &&
			se.CustomCode != "" &&
			!reflect.DeepEqual(se.ErrorLocation, ErrorLocation{})
	}

	err := quick.Check(f, &quick.Config{MaxCount: 300})
	if err != nil {
		t.Errorf("something stupid happened: %v st err = \n%+v\n", err, se)
	}
}
