# errors
Standardized error reporting for music tribe golang projects

## Installation
To use this package in your go program, open your terminal and run the command... 
```
go get github.com/music-tribe/errors
```

## In use
To init a new storage error...
```golang
import (
  "github.com/music-tribe/errors"
  "github.com/music-tribe/uuid"

  "some/local/path/database"
)

func (svc *service)someMethod(id) error {
  if err := svc.db.Get(id); err != nil {
    if err == database.NotFoundError {
      return errors.NewCloudError(404, "add your own error message here")
    }
    return errors.NewCloudError(500, err.Error())
  }
}
```

## Functional Options
We have the ability to use functional options when initializing an error. These options passed to the `NewCloudError` method via the `CloudErrorOption` type...
```golang
type CloudErrorOption func(*CloudError)
```
These options offer the chance to alter any of the fields within the `CloudError` object (or any of it's child objects)...
```golang
type CloudError struct {
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
```
You can create and pass as many options as is necessary to customise your `CloudError`. The `NewCloudError` method is actually a variadic function, allowing you to pass an arbitrary number of `CloudErrorOptions` into the 3rd parameter...
```golang
func setErrorLocationService(svc string) CloudErrorOption {
  return func(ce *CloudError) {
    ce.ErrorLocation.Service = svc
  }
}

func doSomethingAuthy() {
  ...

  if err != nil {
    return errors.NewCloudError(403, "some auth issue", setErrorLocationService("myService"))
  }
  ...
}
```

## Custom Echo Error Handler
By using the custom error handler in this package, any errors from requests made via our echo router will be returned in the `CloudError` JSON format.
If the `ENVIRONMENT` env var is set to `dev`, you will recieve a detailed error location object as well. (page, line, method)
```golang

import (
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors/handlers"
)

func main() {
	e := echo.New()
	
	// init our custom error handler here - this will return all errors from this router - 
	// error, echo.HTTPError and CloudError - as CloudError JSON objects
	e.HTTPErrorHandler = handler.NewCustomHTTPErrorHandler()
	...
}

```

## Contributing
Contribution to this package will only be permitted for Music Tribe employees.

