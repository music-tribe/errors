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

func TestNewStorageErrorBuilder(t *testing.T) {
	want := storageErrorBuilder{
		&StorageError{
			ErrorLocation: ErrorLocation{
				skip: 2,
			},
		},
	}
	if got := NewStorageErrorBuilder(); !reflect.DeepEqual(*got.err, *want.err) {
		t.Errorf("NewStorageErrorBuilder() = \n%v \nwant \n%v\n", *got.err, *want.err)
	}
}

func Test_storageErrorBuilder_Build(t *testing.T) {
	timeNow := time.Now().UTC()
	want := StorageError{
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

	setLine := func(se *StorageError) { se.ErrorLocation.Line = 37 }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }

	if got := NewStorageErrorBuilder().Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.Build() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_StatusCode(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 404
	status := http.StatusText(statusCode)
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }

	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := StorageError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    "NotFound",
		ErrorLocation: el,
	}

	if got := NewStorageErrorBuilder().StatusCode(statusCode).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v, \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_Message(t *testing.T) {
	timeNow := time.Now().UTC()
	msg := "not found"
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }

	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := StorageError{
		StatusCode:    500,
		Status:        http.StatusText(500),
		TimeStamp:     timeNow,
		Message:       msg,
		CustomCode:    InternalServerError,
		ErrorLocation: el,
	}

	if got := NewStorageErrorBuilder().Message(msg).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_Location(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	pkg := "testing"
	svc := "svc-presets"
	fnc := "storageErrorBuilder.ErrorLocation"
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Service: svc,
		Method:  funcCaller,
		Page:    builderTestPage,
		Line:    line,
		skip:    2,
	}

	want := StorageError{
		StatusCode:    500,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		ErrorLocation: el,
		CustomCode:    InternalServerError,
	}

	if got := NewStorageErrorBuilder().ErrorLocation(svc, pkg, fnc).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_CustomCode(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 403
	status := http.StatusText(statusCode)
	var ccode CustomCode = "FileIsInvalidType"
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := StorageError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    ccode,
		ErrorLocation: el,
	}

	if got := NewStorageErrorBuilder().StatusCode(statusCode).CustomCode(ccode).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_CorrelationID(t *testing.T) {
	timeNow := time.Now().UTC()
	statusCode := 502
	status := http.StatusText(statusCode)
	correlationID := uuid.New().String()
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := StorageError{
		StatusCode:    statusCode,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    "BadGateway",
		CorrelationID: correlationID,
		ErrorLocation: el,
	}

	if got := NewStorageErrorBuilder().StatusCode(statusCode).CorrelationID(correlationID).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_Tags(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	tags := []string{"blob", "invalid", "upload"}
	line := 37
	setLine := func(se *StorageError) { se.ErrorLocation.Line = line }
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   builderTestPage,
		Line:   line,
		skip:   2,
	}

	want := StorageError{
		StatusCode:    500,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    InternalServerError,
		Tags:          tags,
		ErrorLocation: el,
	}

	if got := NewStorageErrorBuilder().Tags(tags...).Build(timeNow, setLine, setPage); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = %v\n, \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_Options(t *testing.T) {
	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	svc := "someService"
	setPage := func(se *StorageError) { se.ErrorLocation.Page = builderTestPage }
	setLine := func(se *StorageError) { se.ErrorLocation.Line = 37 }
	options := []StorageErrorOption{
		func(se *StorageError) { se.ErrorLocation.Service = svc },
		setPage,
		setLine,
	}

	want := StorageError{
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

	if got := NewStorageErrorBuilder().Build(timeNow, options...); !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%v \nwant \n%v\n", *got, want)
	}
}

func Test_storageErrorBuilder_SkipCaller(t *testing.T) {
	pc, thisPage, _, _ := runtime.Caller(0)
	fnc := runtime.FuncForPC(pc)
	funcCaller := fnc.Name()

	timeNow := time.Now().UTC()
	status := http.StatusText(500)
	line := 291
	skip := 1
	setSkip := func(se *StorageError) { se.ErrorLocation.skip = skip }
	el := ErrorLocation{
		Method: funcCaller,
		Page:   thisPage,
		Line:   line,
		skip:   skip,
	}

	want := StorageError{
		StatusCode:    500,
		Status:        status,
		Message:       status,
		TimeStamp:     timeNow,
		CustomCode:    InternalServerError,
		ErrorLocation: el,
	}

	_, _, _line, _ := runtime.Caller(0)
	got := NewStorageErrorBuilder().Build(timeNow, setSkip)

	want.ErrorLocation.Line = _line + 1

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("storageErrorBuilder.StatusCode() = \n%+v\n \nwant \n%+v\n", *got, want)
	}
}
