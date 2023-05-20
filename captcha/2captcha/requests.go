package twocaptcha

import (
	"encoding/json"
	"github.com/k773/utils"
	"github.com/k773/utils/fixedPoint"
	"strconv"
)

type requestInterface interface {
	setKey(key string)
	fillInDefaults()
}

/*
	Common request. Sent with every request to the api.
*/

type CommonRequest struct {
	Json   BoolInt `json:"json"`
	SoftId int     `json:"soft_id"`

	Key string `json:"key"`
}

func (c *CommonRequest) fillInDefaults() {
	c.Json = true
	c.SoftId = 0
}

func (c *CommonRequest) setKey(key string) {
	c.Key = key
}

/*
	Common captcha request. Sent with every
*/

type CommonCaptchaRequest struct {
	CommonRequest
	Method Method `json:"method"`
}

func (c *CommonCaptchaRequest) fillInDefaults() {
	c.CommonRequest.fillInDefaults()
}

type ProxyRequest struct {
	Cookies   string `json:"cookies"`
	UserAgent string `json:"userAgent"`
	// Proxy must have the following format: login:password@123.123.123.123:3128
	Proxy     string `json:"proxy"`
	ProxyType string `json:"proxy_type"`
}

func (r *ProxyRequest) fillInDefaults() {
}

func (r *ProxyRequest) SetProxy(proxy *utils.ProxyData) {
	r.Cookies = proxy.Cookies
	r.UserAgent = proxy.UserAgent
	if proxy.ProxyLogin != "" {
		r.Proxy = proxy.ProxyLogin + ":" + proxy.ProxyPassword + "@"
	}
	r.Proxy += proxy.ProxyAddress + ":" + strconv.Itoa(proxy.ProxyPort)
	r.ProxyType = proxy.ProxyType
}

/*
	Funcaptcha request
*/

type FuncaptchaRequest struct {
	CommonCaptchaRequest
	ProxyRequest

	PublicKey string `json:"publickey"`
	PageUrl   string `json:"pageurl"`
}

func (f *FuncaptchaRequest) fillInDefaults() {
	f.CommonCaptchaRequest.fillInDefaults()
	f.ProxyRequest.fillInDefaults()

	f.Method = MethodFuncaptcha
}

/*
	Recaptcha request
*/

type RecaptchaRequest struct {
	CommonCaptchaRequest
	ProxyRequest

	Enterprise BoolInt                `json:"enterprise"`
	GoogleKey  string                 `json:"googlekey"`
	PageUrl    string                 `json:"page_url"`
	Domain     string                 `json:"domain"` // Domain used to load the captcha
	Invisible  BoolInt                `json:"invisible"`
	DataS      string                 `json:"data-s"`
	MinScore   fixedPoint.IntScaledP6 `json:"min_score"`
}

func (r *RecaptchaRequest) fillInDefaults() {
	r.CommonCaptchaRequest.fillInDefaults()
	r.ProxyRequest.fillInDefaults()

	r.Method = MethodRecaptcha
}

/*
	Action request
*/

type ActionRequest struct {
	CommonRequest

	Id     string `json:"id"`
	Action string `json:"action"`
}

func (r *ActionRequest) fillInDefaults() {
	r.CommonRequest.fillInDefaults()
}

/*
	Method
*/

type Method string

const (
	MethodRecaptcha  Method = "userrecaptcha"
	MethodFuncaptcha Method = "funcaptcha"
)

/*
	Bool to int
*/

type BoolInt bool

func (b *BoolInt) String() string {
	return "1"
}

func (b *BoolInt) MarshalJSON() ([]byte, error) {
	if *b {
		return json.Marshal(1)
	}
	return json.Marshal(0)
}

func (b *BoolInt) UnmarshalJSON(data []byte) (e error) {
	var v int
	if e = json.Unmarshal(data, &v); e != nil {
		return e
	}
	*b = v == 1
	return e
}
