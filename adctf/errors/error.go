package errors

import "github.com/labstack/echo"

type ApplicationError struct {
	StatusCode int
	BaseError  error
}

type ApplicationErrorGenerator func(error) *ApplicationError

func (e ApplicationError) Error() string {
	return e.BaseError.Error()
}

func (e ApplicationError) HTTPError() *echo.HTTPError {
	return echo.NewHTTPError(e.StatusCode, e.Error())
}
