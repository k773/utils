package utils

import (
	"errors"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"strings"
	"time"
)

type CheapSMS struct {
	Ses   *gorequest.SuperAgent
	Token string
}

func NewCheapSms(token string, proxy string, timeout time.Duration) CheapSMS {
	return CheapSMS{
		Ses:   gorequest.New().Proxy(proxy).Timeout(timeout),
		Token: token,
	}
}

func (c *CheapSMS) Balance() (b float64, e error) {
	r, resp, _ := c.Ses.Get("https://cheapsms.pro/handler/index?api_key=" + c.Token + "&action=getBalance").End()
	e = errors.New("nil response")
	if r != nil {
		e = errors.New(resp)
		if strings.HasPrefix(resp, "ACCESS_BALANCE:") {
			b, e = strconv.ParseFloat(strings.ReplaceAll(resp, "ACCESS_BALANCE:", ""), 64)
		}
	}
	return
}

type CheapSMSNumber struct {
	ID string
}

func (c *CheapSMS) GetNumber(serviceCode, ref string) {

}
