package twocaptcha

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/k773/utils"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type TwoCaptcha struct {
	s   *resty.Client
	key string
}

const endPointIn = "https://2captcha.com/in.php"
const endPointRes = "https://2captcha.com/res.php"

const DomainRecaptchaNet = "recaptcha.net"
const DomainGoogleCom = "google.com"

func New(s *resty.Client, key string) *TwoCaptcha {
	return &TwoCaptcha{
		s:   s,
		key: key,
	}
}

type CaptchaResult struct {
	cap    *TwoCaptcha
	id     string
	result string
}

func (cr *CaptchaResult) Result() string {
	return cr.result
}

func (cr *CaptchaResult) ReportGood() {
	_, _ = cr.cap.s.R().SetQueryParams(map[string]string{"key": cr.cap.key, "action": "reportgood", "id": cr.id}).Get(endPointRes)
}

func (cr *CaptchaResult) ReportBad() {
	_, _ = cr.cap.s.R().SetQueryParams(map[string]string{"key": cr.cap.key, "action": "reportbad", "id": cr.id}).Get(endPointRes)
}

type epInRes struct {
	Status    int    `json:"status"`
	Request   string `json:"request"`
	ErrorText string `json:"error_text"`
}

func (c *TwoCaptcha) GetBalance() (balance float64, e error) {
	res, e := c.s.R().SetQueryParams(map[string]string{"key": c.key, "action": "getbalance", "json": "1"}).Get(endPointRes)
	if e == nil {
		var resp epInRes
		if e = json.Unmarshal(res.Body(), &resp); e == nil {
			if resp.Status != 1 {
				e = errors.New(fmt.Sprintf("%v: %v", resp.Request, resp.ErrorText))
			} else {
				balance, e = strconv.ParseFloat(resp.Request, 64)
			}
		}
	}
	return balance, errors.Wrap(e, "GetBalance")
}

func (c *TwoCaptcha) SolveRecaptchaEnterpriseV2(siteKey, pageUrl, dataS, captchaDomain string, proxy *utils.ProxyData) (cap *CaptchaResult, e error) {
	cap = &CaptchaResult{cap: c}
	var m = map[string]string{"key": c.key, "method": "userrecaptcha", "enterprise": "1", "googlekey": siteKey, "pageurl": pageUrl, "min_score": "0.9", "domain": captchaDomain, "json": "1"}
	if len(dataS) != 0 {
		m["stoken"] = dataS
		m["data-s"] = dataS
	}
	if proxy != nil {
		m["proxy"] = proxy.StringNoType()
		m["proxytype"] = strings.ToUpper(proxy.ProxyType)
		if proxy.UserAgent != "" {
			m["userAgent"] = proxy.UserAgent
		}
	}
	res, e := c.s.R().SetFormData(m).Post(endPointIn)
	if e == nil {
		var resp epInRes
		if e = json.Unmarshal(res.Body(), &resp); e == nil {
			if resp.Status != 1 {
				e = errors.New(fmt.Sprintf("%v: %v", resp.Request, resp.ErrorText))
			} else {
				cap.id = resp.Request
				e = c.WaitForResult(20*time.Second, cap)
			}
		}
	}

	return cap, errors.Wrap(e, "SolveRecaptchaEnterprise")
}

func (c *TwoCaptcha) WaitForResult(timeout time.Duration, cap *CaptchaResult) (e error) {
	for e == nil && cap.result == "" {
		time.Sleep(timeout)
		var res *resty.Response
		res, e = c.s.R().SetFormData(map[string]string{"key": c.key, "action": "get", "id": cap.id, "json": "1"}).Post(endPointRes)
		var resp epInRes
		if e = json.Unmarshal(res.Body(), &resp); e == nil {
			if resp.Status != 1 {
				if resp.Request != "CAPCHA_NOT_READY" {
					e = errors.New(fmt.Sprintf("%v: %v", resp.Request, resp.ErrorText))
				}
			} else {
				cap.result = resp.Request
			}
		}
	}
	return
}
