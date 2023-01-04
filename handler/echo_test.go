package handler

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors"
)

const (
	testCorrelationID string = "5f1aa5d0-bdb1-4cd7-a807-6d673f49f871"
)

func TestCustomHTTPErrorHandler(t *testing.T) {
	type args struct {
		err error
		env string
	}
	tests := []struct {
		name              string
		args              args
		wantStatusCode    int
		wantErrMsg        string
		wantLocation      bool
		wantCorrelationID string
	}{
		{
			name: "When it's a standard golang error",
			args: args{
				fmt.Errorf("this is a standard error"),
				"dev",
			},
			wantStatusCode:    500,
			wantErrMsg:        "this is a standard error",
			wantLocation:      true,
			wantCorrelationID: testCorrelationID,
		},
		{
			name: "When it's a standard golang error on production env",
			args: args{
				fmt.Errorf("this is a standard prod error"),
				"production",
			},
			wantStatusCode:    500,
			wantErrMsg:        "this is a standard prod error",
			wantLocation:      false,
			wantCorrelationID: testCorrelationID,
		},
		{
			name: "When it's a standard golang error",
			args: args{
				echo.NewHTTPError(405, "Method Not Allowed"),
				"dev",
			},
			wantStatusCode:    405,
			wantErrMsg:        "Method Not Allowed",
			wantLocation:      true,
			wantCorrelationID: testCorrelationID,
		},
		{
			name: "When it's a standard golang error on production env",
			args: args{
				echo.NewHTTPError(405, "Method Not Allowed"),
				"production",
			},
			wantStatusCode:    405,
			wantErrMsg:        "Method Not Allowed",
			wantLocation:      false,
			wantCorrelationID: testCorrelationID,
		},
		{
			name: "When the error is an MT cloud error",
			args: args{
				errors.NewCloudError(403, "do one!"),
				"dev",
			},
			wantStatusCode:    403,
			wantErrMsg:        "do one!",
			wantLocation:      true,
			wantCorrelationID: testCorrelationID,
		},
		{
			name: "When the error is an MT cloud error on staging",
			args: args{
				errors.NewCloudError(403, "do one!"),
				"staging",
			},
			wantStatusCode:    403,
			wantErrMsg:        "do one!",
			wantLocation:      false,
			wantCorrelationID: testCorrelationID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENVIRONMENT", tt.args.env)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(correlationIDHeader, testCorrelationID)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			h := NewCustomHTTPErrorHandler()
			h(tt.args.err, ctx)

			fmt.Printf("%s", rec.Body.Bytes())

			ce := new(errors.CloudError)
			if err := json.Unmarshal(rec.Body.Bytes(), ce); err != nil {
				t.Fatal(err)
			}

			if ce.StatusCode != tt.wantStatusCode {
				t.Errorf("want status code to be %d but got %d\n", tt.wantStatusCode, ce.StatusCode)
				return
			}

			if ce.Message != tt.wantErrMsg {
				t.Errorf("want errors msg to be %s but got %s\n", tt.wantErrMsg, ce.Message)
				return
			}

			if ce.CorrelationID != tt.wantCorrelationID {
				t.Errorf("want correlation id to be %s but got %s\n", tt.wantCorrelationID, ce.CorrelationID)
				return
			}

			if !reflect.DeepEqual(ce.ErrorLocation, errors.ErrorLocation{}) != tt.wantLocation {
				t.Errorf("want location to be %v but got %+v\n", tt.wantLocation, ce.ErrorLocation)
				return
			}
		})
	}
}
