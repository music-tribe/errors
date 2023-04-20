package errors

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/music-tribe/uuid"
)

const (
	builderTestPage = "test/page/where/error/was/built.go"
	funcCaller      = "testing.tRunner"
	testPackage     = "testing"
)

func TestNewCloudErrorBuilder(t *testing.T) {
	want := cloudErrorBuilder{
		&CloudError{
			ErrorLocation: ErrorLocation{
				skip: 2,
			},
		},
	}
	if got := NewCloudErrorBuilder(); !reflect.DeepEqual(*got.err, *want.err) {
		t.Errorf("NewCloudErrorBuilder() = \n%v \nwant \n%v\n", *got.err, *want.err)
	}
}

func Test_cloudErrorBuilder_Build(t *testing.T) {
	timeNow := time.Now().UTC()
	want := CloudError{
		StatusCode: 500,
		Status:     "Internal Server Error",
		Message:    "Internal Server Error",
		Source:     "music-tribe",
		TimeStamp:  timeNow,
		CustomCode: InternalServerError,
		ErrorLocation: ErrorLocation{
			Method: funcCaller,
			Page:   builderTestPage,
			Line:   37,
			skip:   2,
		},
	}

	setLine := func(se *CloudError) { se.ErrorLocation.Line = 37 }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }

	if got := NewCloudErrorBuilder().Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.Build() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_StatusCode(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 404

	if got := NewCloudErrorBuilder().StatusCode(statusCode).Build(timeNow); got.StatusCode != statusCode {
		t.Errorf("cloudErrorBuilder.StatusCode() = %d, \nwant %d\n", got.StatusCode, statusCode)
	}
}

func Test_cloudErrorBuilder_Message(t *testing.T) {
	timeNow := time.Now().UTC()
	msg := "not found"

	if got := NewCloudErrorBuilder().Message(msg).Build(timeNow); got.Message != msg {
		t.Errorf("cloudErrorBuilder.StatusCode() = %s \nwant %s\n", got.Message, msg)
	}
}

func Test_cloudErrorBuilder_Location(t *testing.T) {
	timeNow := time.Now().UTC()
	pkg := "testing"
	svc := "svc-presets"
	fnc := "testing.tRunner"

	got := NewCloudErrorBuilder().ErrorLocation(svc, pkg, fnc).Build(timeNow)
	if got.ErrorLocation.Service != svc {
		t.Errorf("cloudErrorBuilder.StatusCode() = %s \nwant %s\n", got.ErrorLocation.Service, svc)
	}

	if got.ErrorLocation.Method != fnc {
		t.Errorf("cloudErrorBuilder.StatusCode() = %s \nwant %s\n", got.ErrorLocation.Method, fnc)
	}
}

func Test_cloudErrorBuilder_CustomCode(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 403
	var ccode CustomCode = "FileIsInvalidType"

	if got := NewCloudErrorBuilder().StatusCode(statusCode).CustomCode(ccode).Build(timeNow); got.CustomCode != ccode {
		t.Errorf("cloudErrorBuilder.StatusCode() = %s \nwant %s\n", got.CustomCode, ccode)
	}
}

func Test_cloudErrorBuilder_CorrelationID(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 502
	correlationID := uuid.New().String()

	got := NewCloudErrorBuilder().StatusCode(statusCode).CorrelationID(correlationID).Build(timeNow)
	if got.CorrelationID != correlationID {
		t.Errorf("cloudErrorBuilder.StatusCode() = %s \nwant %s\n", got.CorrelationID, correlationID)
	}
}

func Test_cloudErrorBuilder_Tags(t *testing.T) {
	timeNow := time.Now().UTC()
	tags := []string{"blob", "invalid", "upload"}

	if got := NewCloudErrorBuilder().Tags(tags...).Build(timeNow); !reflect.DeepEqual(got.Tags, tags) {
		t.Errorf("cloudErrorBuilder.StatusCode() = %v\n, \nwant \n%v\n", got.Tags, tags)
	}
}

func Test_cloudErrorBuilder_Options(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	svc := "someService"
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }
	setLine := func(se *CloudError) { se.ErrorLocation.Line = 37 }
	options := []CloudErrorOption{
		func(se *CloudError) { se.ErrorLocation.Service = svc },
		setPage,
		setLine,
	}

	want := CloudError{
		StatusCode: 500,
		Status:     status,
		Message:    status,
		TimeStamp:  timeNow,
		CustomCode: InternalServerError,
		Source:     "music-tribe",
		ErrorLocation: ErrorLocation{
			Method:  funcCaller,
			Service: svc,
			Page:    builderTestPage,
			Line:    37,
			skip:    2,
		},
	}

	if got := NewCloudErrorBuilder().Build(timeNow, options...); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_SkipCaller(t *testing.T) {
	pc, thisPage, _, _ := runtime.Caller(0)
	fnc := runtime.FuncForPC(pc)
	funcCaller := fnc.Name()

	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	line := 291
	skip := 1
	setSkip := func(se *CloudError) { se.ErrorLocation.skip = skip }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   thisPage,
		Line:   line,
		skip:   skip,
	}

	want := CloudError{
		StatusCode:    500,
		Status:        status,
		Source:        "music-tribe",
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    InternalServerError,
		ErrorLocation: el,
	}

	_, _, _line, _ := runtime.Caller(0)
	got := NewCloudErrorBuilder().Build(timeNow, setSkip)

	want.ErrorLocation.Line = _line + 1

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%+v\n \nwant \n%+v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_Source(t *testing.T) {
	timeNow := time.Now().UTC()
	t.Run("when we add a source to error it should show up in the CloudError object", func(t *testing.T) {
		srcname := "azure"
		got := NewCloudErrorBuilder().Source(srcname).Build(timeNow)
		if got.Source != srcname {
			t.Errorf("expected source to be %s but got %s", srcname, got.Source)
		}
	})

	t.Run("when no source is provided, it defaults to music-tribe", func(t *testing.T) {
		want := "music-tribe"
		got := NewCloudErrorBuilder().Source("").Build(timeNow)
		if got.Source != want {
			t.Errorf("expected source to be %s but got %s", want, got.Source)
		}
	})
}

func Test_cloudErrorBuilder_Error(t *testing.T) {
	timeNow := time.Now().UTC()
	t.Run("when we pass an error to the Error builder, we should store this error as the Internal error", func(t *testing.T) {
		want := errors.New("simple error")
		got := NewCloudErrorBuilder().Error(want).Build(timeNow)
		if got.InternalError != want {
			t.Errorf("expected InternalError to be %v but got %v", want, got.InternalError)
		}
	})

	t.Run("when we pass an error to the Error builder, we should store the message within the CloudError message field", func(t *testing.T) {
		want := errors.New("simple error")
		got := NewCloudErrorBuilder().Error(want).Build(timeNow)
		if got.Message != want.Error() {
			t.Errorf("expected Message to be %s but got %s", want, got.Message)
		}
	})

	t.Run("when we pass a string to the Error builder, we should store this string within message field", func(t *testing.T) {
		want := "simple error"
		got := NewCloudErrorBuilder().Error(want).Build(timeNow)
		if got.Message != want {
			t.Errorf("expected Message to be %v but got %v", want, got.Message)
		}
	})

	t.Run("when we pass a string to the Error builder, we expect the internal error to be nil", func(t *testing.T) {
		want := "simple error"
		got := NewCloudErrorBuilder().Error(want).Build(timeNow)
		if got.InternalError != nil {
			t.Errorf("expected InternalError to be nil but got %v", got.InternalError)
		}
	})

	t.Run("when we pass an unknown object to the Error builder, we expect to get the default error code msg back", func(t *testing.T) {
		want := "Internal Server Error"
		got := NewCloudErrorBuilder().Error(struct{ name string }{"hello"}).Build(timeNow)
		if got.Message != want {
			t.Errorf("expected InternalError to be nil but got %v", got.Message)
		}
	})
}
