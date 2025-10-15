package services

import "errors"

var (
	ErrorUserNotFound           = errors.New("user not found")
	ErrorActivityPeriodNotFound = errors.New("activity period not found")

	ErrorInternalError = errors.New("internal error")
)
