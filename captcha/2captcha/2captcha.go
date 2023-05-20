package twocaptcha

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k773/utils"
	"github.com/k773/utils/fixedPoint"
	"time"
)

const (
	endpointIn  = "https://2captcha.com/in.php"
	endpointRes = "https://2captcha.com/res.php"
)

/*
	API
*/

type Client struct {
	Key string
	Ses *resty.Client
	// PollInterval defines frequency of requests for Wait() fn.
	// Default: 10s
	PollInterval time.Duration
	// MaxWaitDuration defines max amount of time Wait() is going to wait for the captcha to be solved.
	// 0 = duration is not limited.
	// Default: 0
	MaxWaitDuration time.Duration
}

func New(key string) *Client {
	return &Client{
		Key:             key,
		Ses:             resty.New(),
		PollInterval:    10 * time.Second,
		MaxWaitDuration: 0,
	}
}

func (c *CaptchaResponse) Report(ctx context.Context, good bool) (e error) {
	var action = utils.If(good, "reportgood", "reportbad")
	var req = &ActionRequest{Action: action, Id: c.Id}
	_, e = c.c.Execute(ctx, req, endpointRes)
	return
}

func (c *Client) GetBalance(ctx context.Context) (data fixedPoint.IntScaledP6, e error) {
	var req = &ActionRequest{Action: "getbalance"}
	r, e := c.Execute(ctx, req, endpointRes)
	if e == nil {
		data = fixedPoint.ParseIntScaledP6(r.Request)
	}
	return
}

func (c *Client) SolveCaptcha(ctx context.Context, request requestInterface) (response CaptchaResponse, e error) {
	response.c = c
	response.Response, e = c.Execute(ctx, request, endpointIn)
	if e == nil {
		response.Id = response.Request
		response.Response, e = c.Wait(ctx, response.Id)
	}
	return
}

func (c *Client) Wait(ctx context.Context, id string) (response Response, e error) {
	var startedAt = time.Now()
	var getTimeoutError = func() error {
		var e = &TimeoutError{
			TimeSpent:   time.Now().Sub(startedAt),
			TimeAllowed: c.MaxWaitDuration,
		}
		if c.MaxWaitDuration != 0 && e.TimeSpent > c.MaxWaitDuration {
			return e
		}
		return nil
	}
	for e == nil {
		if e = ctx.Err(); e != nil {
			continue
		}
		if e = getTimeoutError(); e != nil {
			continue
		}

		if e = utils.SleepWithContext(ctx, c.PollInterval); e == nil {
			var req = &ActionRequest{Id: id, Action: "get2"}
			if response, e = c.Execute(ctx, req, endpointRes); e == nil {
				if response.Request != "CAPCHA_NOT_READY" {
					break
				}
			}
		}
	}
	return
}

func (c *Client) Execute(ctx context.Context, request requestInterface, endpoint string) (response Response, e error) {
	request.fillInDefaults()
	request.setKey(c.Key)

	r, e := c.Ses.R().SetContext(ctx).SetFormData(structToMap(request)).Post(endpoint)
	if e == nil {
		if e = json.Unmarshal(r.Body(), &response); e == nil {
			e = response.GetError()
		}
	}
	return
}
