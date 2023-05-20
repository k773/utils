package twocaptcha

import (
	"fmt"
	"time"
)

type ErrorIncorrectResponseCode struct {
	Code    string
	Message string
}

func (e *ErrorIncorrectResponseCode) Error() string {
	if e.Code == "" {
		return "the request has failed, but an empty error code was received"
	}
	return e.Code + ": " + e.Message
}

type TimeoutError struct {
	TimeSpent   time.Duration
	TimeAllowed time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout: spent %v (max allowed: %v)", e.TimeSpent, e.TimeAllowed)
}
