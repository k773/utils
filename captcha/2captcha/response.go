package twocaptcha

import (
	"github.com/k773/utils/fixedPoint"
)

type CaptchaResponse struct {
	c *Client
	Response
	Id string `json:"id"`
}

func (c *CaptchaResponse) Result() string {
	return c.Request
}

type Response struct {
	// Status is only set to 1 if the request's succeeded.
	Status BoolInt `json:"status"`

	Price     fixedPoint.IntScaledP6 `json:"price"`
	ErrorText string                 `json:"error_text"`
	// Have no idea why it's called that way, but this field actually contains a response / error code.
	Request string `json:"request"`
}

func (r *Response) GetError() error {
	if r.Status {
		return nil
	}
	if r.Request == "CAPCHA_NOT_READY" {
		return nil
	}
	return &ErrorIncorrectResponseCode{Code: r.Request, Message: r.ErrorText}
}
