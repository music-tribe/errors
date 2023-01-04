package handler

import (
	errs "errors"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors"
)

const correlationIDHeader = "X-Request-ID"

func NewCustomHTTPErrorHandler() func(error, echo.Context) {
	return func(err error, c echo.Context) {
		code := 500
		msg := err.Error()
		correlationID := c.Request().Header.Get(correlationIDHeader)

		ce := &errors.CloudError{}
		if errs.As(err, &ce) {
			ce.CorrelationID = correlationID
			if !isDevEnv() {
				ce.ErrorLocation = errors.ErrorLocation{}
			}
			_ = c.JSON(ce.StatusCode, ce)
			return
		}

		he := &echo.HTTPError{}
		if errs.As(err, &he) {
			if heMsg, ok := he.Message.(string); ok {
				msg = heMsg
			}

			heErr := errors.NewCloudError(he.Code, msg)
			if !isDevEnv() {
				heErr.ErrorLocation = errors.ErrorLocation{}
			}

			heErr.CorrelationID = correlationID

			_ = c.JSON(he.Code, heErr)
			return
		}

		outErr := errors.NewCloudError(code, msg)
		outErr.CorrelationID = correlationID
		if !isDevEnv() {
			outErr.ErrorLocation = errors.ErrorLocation{}
		}

		_ = c.JSON(code, outErr)
	}
}

func isDevEnv() bool {
	return os.Getenv("ENVIRONMENT") == "dev"
}
