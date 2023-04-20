package errors

import (
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
	status := http.StatusText(statusCode)
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }

	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := CloudError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    "NotFound",
		ErrorLocation: el,
	}

	if got := NewCloudErrorBuilder().StatusCode(statusCode).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v, \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_Message(t *testing.T) {
	timeNow := time.Now().UTC()
	msg := "not found"
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }

	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := CloudError{
		StatusCode:    500,
		Status:        http.StatusText(500),
		TimeStamp:     timeNow,
		Message:       msg,
		CustomCode:    InternalServerError,
		ErrorLocation: el,
	}

	if got := NewCloudErrorBuilder().Message(msg).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_Location(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	pkg := "testing"
	svc := "svc-presets"
	fnc := "cloudErrorBuilder.ErrorLocation"
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Service: svc,
		Method:  funcCaller,
		Page:    builderTestPage,
		Line:    line,
		skip:    2,
	}

	want := CloudError{
		StatusCode:    500,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		ErrorLocation: el,
		CustomCode:    InternalServerError,
	}

	if got := NewCloudErrorBuilder().ErrorLocation(svc, pkg, fnc).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_CustomCode(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 403
	status := http.StatusText(statusCode)
	var ccode CustomCode = "FileIsInvalidType"
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := CloudError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    ccode,
		ErrorLocation: el,
	}

	if got := NewCloudErrorBuilder().StatusCode(statusCode).CustomCode(ccode).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_CorrelationID(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 502
	status := http.StatusText(statusCode)
	correlationID := uuid.New().String()
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := CloudError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    "BadGateway",
		CorrelationID: correlationID,
		ErrorLocation: el,
	}

	if got := NewCloudErrorBuilder().StatusCode(statusCode).CorrelationID(correlationID).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_cloudErrorBuilder_Tags(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	tags := []string{"blob", "invalid", "upload"}
	line := 37
	setLine := func(se *CloudError) { se.ErrorLocation.Line = line }
	setPage := func(se *CloudError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := CloudError{
		StatusCode:    500,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    InternalServerError,
		Tags:          tags,
		ErrorLocation: el,
	}

	if got := NewCloudErrorBuilder().Tags(tags...).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("cloudErrorBuilder.StatusCode() = %v\n, \nwant \n%v\n", *got, want)
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
