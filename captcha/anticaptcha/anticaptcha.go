package anticaptcha

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty"
	"github.com/k773/utils"
	"github.com/pkg/errors"
	"time"
)

type AntiCaptcha struct {
	logger *utils.Logger
	s      *resty.Client
	key    string
}

func New(s *resty.Client, key string, logger ...*utils.Logger) *AntiCaptcha {
	var ac = &AntiCaptcha{s: s, key: key}
	if len(logger) != 0 {
		ac.logger = logger[0]
	}
	return ac
}

const (
	antiCaptchaCreateTaskUrl    = "https://api.anti-captcha.com/createTask"
	antiCaptchaGetTaskResultUrl = "https://api.anti-captcha.com/getTaskResult"

	antiCaptchaTypeRecaptchaV2EnterpriseProxyless = "RecaptchaV2EnterpriseTaskProxyless"
	antiCaptchaTypeRecaptchaV2EnterpriseProxy     = "RecaptchaV2EnterpriseTask"

	antiCaptchaTypeRecaptchaV2Proxyless = "RecaptchaV2TaskProxyless"
	antiCaptchaTypeRecaptchaV2Proxy     = "RecaptchaV2Task"

	antiCaptchaTypeImageToText = "ImageToTextTask"
)

type antiCaptchaEnterprisePayload struct {
	S string `json:"s,omitempty"`
}

type antiCaptchaTaskRequest struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	WebsiteKey string `json:"websiteKey,omitempty"`

	EnterprisePayload *antiCaptchaEnterprisePayload `json:"enterprisePayload,omitempty"`

	Body      string `json:"body,omitempty"`
	Phrase    bool   `json:"phrase,omitempty"`
	Case      bool   `json:"case,omitempty"`
	Numeric   bool   `json:"numeric,omitempty"`
	Math      int    `json:"math,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`

	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Cookies       string `json:"cookies,omitempty"`

	ApiDomain string `json:"apiDomain,omitempty"`
}

type antiCaptchaNewTaskRequest struct {
	ClientKey    string                 `json:"clientKey"`
	Task         antiCaptchaTaskRequest `json:"task"`
	SoftID       int                    `json:"softId"`
	LanguagePool string                 `json:"languagePool"`
}

type antiCaptchaNewTaskResponse struct {
	ErrorID          int    `json:"errorId"`
	TaskID           int    `json:"taskId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
}

type antiCaptchaGetTaskResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

type CaptchaResult struct {
	cap      *AntiCaptcha
	TaskType string
	id       int

	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
		Text               string `json:"text"`
		URL                string `json:"url"`
	} `json:"solution"`
	Cost       string `json:"cost"`
	IP         string `json:"ip"`
	CreateTime int    `json:"createTime"`
	EndTime    int    `json:"endTime"`
	SolveCount int    `json:"solveCount"`
}

type antiCaptchaReportResult struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
}

func (a *AntiCaptcha) waitForResponse(ctx context.Context, acType, sitekey, siteUrl string, newTaskResponseB []byte) (antiCaptchaResponse CaptchaResult, e error) {
	antiCaptchaResponse.cap = a
	antiCaptchaResponse.TaskType = acType
	antiCaptchaResponse.Status = "processing"
	var newTaskResponse antiCaptchaNewTaskResponse

	if e = json.Unmarshal(newTaskResponseB, &newTaskResponse); e == nil {
		if newTaskResponse.ErrorID != 0 {
			e = errors.New(newTaskResponse.ErrorCode + ": " + newTaskResponse.ErrorDescription)
			_ = utils.SleepWithContext(ctx, 20*time.Second)
		} else {
			antiCaptchaResponse.id = newTaskResponse.TaskID

			for e == nil && antiCaptchaResponse.Status == "processing" {
				if e = utils.SleepWithContext(ctx, 20*time.Second); e == nil {
					var resp *resty.Response
					resp, e = a.s.R().SetContext(ctx).
						SetBody(antiCaptchaGetTaskResultRequest{
							ClientKey: a.key,
							TaskID:    newTaskResponse.TaskID,
						}).
						Post(antiCaptchaGetTaskResultUrl)
					if e == nil {
						if e = json.Unmarshal(resp.Body(), &antiCaptchaResponse); e == nil {
							if antiCaptchaResponse.ErrorID != 0 {
								e = errors.New(antiCaptchaResponse.ErrorCode + ": " + antiCaptchaResponse.ErrorDescription)
							}
						}
					}
				}
			}
		}
	}

	if e == nil {
		if a.logger != nil {
			a.logger.Log(acType, "info", "sitekey:", sitekey, ", site url:", siteUrl, "; response:", string(utils.Marshal(antiCaptchaResponse)))
		}
	} else {
		if a.logger != nil {
			a.logger.Log(acType, "error", "sitekey:", sitekey, ", site url:", siteUrl, "; error:", e)
		}
	}
	return
}

func (a *AntiCaptcha) SolveRecaptchaV2(ctx context.Context, websiteKey, websiteUrl string, proxyData *utils.ProxyData) (antiCaptchaResponse CaptchaResult, e error) {
	var taskType = antiCaptchaTypeRecaptchaV2Proxy
	if proxyData == nil {
		taskType = antiCaptchaTypeRecaptchaV2Proxyless
		proxyData = &utils.ProxyData{}
	}

	resp, e := a.s.R().SetContext(ctx).
		SetBody(antiCaptchaNewTaskRequest{
			ClientKey: a.key,
			Task: antiCaptchaTaskRequest{
				Type:          taskType,
				WebsiteURL:    websiteUrl,
				WebsiteKey:    websiteKey,
				ProxyType:     proxyData.ProxyType,
				ProxyAddress:  proxyData.ProxyAddress,
				ProxyPort:     proxyData.ProxyPort,
				ProxyLogin:    proxyData.ProxyLogin,
				ProxyPassword: proxyData.ProxyPassword,
				UserAgent:     proxyData.UserAgent,
				Cookies:       proxyData.Cookies,
			},
			SoftID:       994,
			LanguagePool: "en",
		}).Post(antiCaptchaCreateTaskUrl)

	if e == nil {
		antiCaptchaResponse, e = a.waitForResponse(ctx, taskType, websiteKey, websiteUrl, resp.Body())
	}
	return antiCaptchaResponse, errors.Wrap(e, "SolveRecaptchaV2")
}

func (a *AntiCaptcha) SolveRecaptchaEnterpriseV2(ctx context.Context, websiteKey, websiteUrl, s string, proxyData *utils.ProxyData) (antiCaptchaResponse CaptchaResult, e error) {
	var taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxy
	if proxyData == nil {
		taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxyless
		proxyData = &utils.ProxyData{}
	}
	var epl *antiCaptchaEnterprisePayload
	if s != "" {
		epl = &antiCaptchaEnterprisePayload{
			S: s,
		}
	}

	resp, e := a.s.R().SetContext(ctx).
		SetBody(antiCaptchaNewTaskRequest{
			ClientKey: a.key,
			Task: antiCaptchaTaskRequest{
				Type:              taskType,
				WebsiteURL:        websiteUrl,
				WebsiteKey:        websiteKey,
				EnterprisePayload: epl,
				ProxyType:         proxyData.ProxyType,
				ProxyAddress:      proxyData.ProxyAddress,
				ProxyPort:         proxyData.ProxyPort,
				ProxyLogin:        proxyData.ProxyLogin,
				ProxyPassword:     proxyData.ProxyPassword,
				UserAgent:         proxyData.UserAgent,
				Cookies:           proxyData.Cookies,
			},
			SoftID:       994,
			LanguagePool: "en",
		}).Post(antiCaptchaCreateTaskUrl)

	if e == nil {
		antiCaptchaResponse, e = a.waitForResponse(ctx, taskType, websiteKey+"/"+s, websiteUrl, resp.Body())
	}
	return antiCaptchaResponse, errors.Wrap(e, "SolveRecaptchaEnterpriseV2")
}

func (a *AntiCaptcha) SolveRecaptchaEnterpriseV2Domain(ctx context.Context, websiteKey, websiteUrl, s, domain string, proxyData *utils.ProxyData) (antiCaptchaResponse CaptchaResult, e error) {
	var taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxy
	if proxyData == nil {
		taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxyless
		proxyData = &utils.ProxyData{}
	}
	var epl *antiCaptchaEnterprisePayload
	if s != "" {
		epl = &antiCaptchaEnterprisePayload{
			S: s,
		}
	}

	resp, e := a.s.R().SetContext(ctx).
		SetBody(antiCaptchaNewTaskRequest{
			ClientKey: a.key,
			Task: antiCaptchaTaskRequest{
				Type:              taskType,
				WebsiteURL:        websiteUrl,
				WebsiteKey:        websiteKey,
				EnterprisePayload: epl,
				ProxyType:         proxyData.ProxyType,
				ProxyAddress:      proxyData.ProxyAddress,
				ProxyPort:         proxyData.ProxyPort,
				ProxyLogin:        proxyData.ProxyLogin,
				ProxyPassword:     proxyData.ProxyPassword,
				UserAgent:         proxyData.UserAgent,
				Cookies:           proxyData.Cookies,
				ApiDomain:         domain,
			},
			SoftID:       994,
			LanguagePool: "en",
		}).Post(antiCaptchaCreateTaskUrl)

	if e == nil {
		antiCaptchaResponse, e = a.waitForResponse(ctx, taskType, websiteKey+"/"+s, websiteUrl, resp.Body())
	}
	return antiCaptchaResponse, errors.Wrap(e, "SolveRecaptchaEnterpriseV2")
}

func (a *AntiCaptcha) SolveImageCaptcha(ctx context.Context, img []byte) (antiCaptchaResponse CaptchaResult, e error) {
	resp, e := a.s.R().SetContext(ctx).
		SetBody(antiCaptchaNewTaskRequest{
			ClientKey: a.key,
			Task: antiCaptchaTaskRequest{
				Type:      antiCaptchaTypeImageToText,
				Body:      base64.StdEncoding.EncodeToString(img),
				Phrase:    false,
				Case:      false,
				Numeric:   false,
				Math:      0,
				MinLength: 0,
				MaxLength: 0,
			},
			SoftID: 994,
		}).Post(antiCaptchaCreateTaskUrl)

	if e == nil {
		antiCaptchaResponse, e = a.waitForResponse(ctx, antiCaptchaTypeImageToText, "none(image)", "none(image)", resp.Body())
	}
	return antiCaptchaResponse, errors.Wrap(e, "SolveImageCaptcha")
}

func (cr *CaptchaResult) Report(ctx context.Context, good bool) error {
	if cr.cap == nil {
		return errors.New("nil captcha instance")
	}

	var url string
	switch cr.TaskType {
	case antiCaptchaTypeRecaptchaV2EnterpriseProxy, antiCaptchaTypeRecaptchaV2EnterpriseProxyless:
		if good {
			url = "https://api.anti-captcha.com/reportCorrectRecaptcha"
		} else {
			url = "https://api.anti-captcha.com/reportIncorrectRecaptcha"
		}
	case antiCaptchaTypeImageToText:
		if !good {
			url = "https://api.anti-captcha.com/reportIncorrectImageCaptcha"
		}
	}

	if url == "" {
		return errors.New("method is not supported")
	}
	resp, e := cr.cap.s.R().SetContext(ctx).
		SetBody(antiCaptchaGetTaskResultRequest{
			ClientKey: cr.cap.key,
			TaskID:    cr.id,
		}).Post(url)
	if e == nil {
		var reportResult antiCaptchaReportResult
		if e = json.Unmarshal(resp.Body(), &reportResult); e == nil {
			if reportResult.ErrorID != 0 {
				e = errors.New(reportResult.ErrorCode + ": " + reportResult.ErrorDescription)
			}
		}
	}
	return errors.Wrap(e, "AntiCaptchaResponse.Report")
}

func (cr *CaptchaResult) Result() string {
	switch cr.TaskType {
	case antiCaptchaTypeRecaptchaV2EnterpriseProxy, antiCaptchaTypeRecaptchaV2EnterpriseProxyless,
		antiCaptchaTypeRecaptchaV2Proxy, antiCaptchaTypeRecaptchaV2Proxyless:
		return cr.Solution.GRecaptchaResponse
	case antiCaptchaTypeImageToText:
		return cr.Solution.Text
	default:
		return cr.Solution.Text
	}
}

func (cr *CaptchaResult) ReportGood(ctx context.Context) {
	_ = cr.Report(ctx, true)
}

func (cr *CaptchaResult) ReportBad(ctx context.Context) {
	_ = cr.Report(ctx, false)
}

//F
