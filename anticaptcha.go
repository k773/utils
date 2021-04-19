package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"time"
)

type AntiCaptcha struct {
	Ses    *gorequest.SuperAgent
	ApiKey string
}

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

type AntiCaptchaResponse struct {
	antiCaptchaInstance *AntiCaptcha
	TaskType            string
	TaskID              int

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

const (
	antiCaptchaCreateTaskUrl = "https://api.anti-captcha.com/createTask"

	antiCaptchaTypeRecaptchaV2EnterpriseProxyless = "RecaptchaV2EnterpriseTaskProxyless"
	antiCaptchaTypeRecaptchaV2EnterpriseProxy     = "RecaptchaV2EnterpriseTask"

	antiCaptchaTypeRecaptchaV2Proxyless = "RecaptchaV2TaskProxyless"
	antiCaptchaTypeRecaptchaV2Proxy     = "RecaptchaV2Task"

	antiCaptchaTypeImageToText = "ImageToTextTask"
)

func (a *AntiCaptcha) waitForResponse(newTaskResponseB []byte) (antiCaptchaResponse AntiCaptchaResponse, e error) {
	antiCaptchaResponse.antiCaptchaInstance = a
	var newTaskResponse antiCaptchaNewTaskResponse

	if e = json.Unmarshal(newTaskResponseB, &newTaskResponse); e == nil {
		if newTaskResponse.ErrorID != 0 {
			e = errors.New(newTaskResponse.ErrorCode + ": " + newTaskResponse.ErrorDescription)
		} else {
			antiCaptchaResponse.TaskID = newTaskResponse.TaskID
		retry:
			r, resp, _ := a.Ses.Get("https://api.anti-captcha.com/getTaskResult").
				Send(antiCaptchaGetTaskResultRequest{
					ClientKey: a.ApiKey,
					TaskID:    newTaskResponse.TaskID,
				}).EndBytes()

			e = errors.New("nil response")
			if r != nil {
				_ = r.Body.Close()

				if e = json.Unmarshal(resp, &antiCaptchaResponse); e == nil {
					if antiCaptchaResponse.ErrorID != 0 {
						e = errors.New(antiCaptchaResponse.ErrorCode + ": " + antiCaptchaResponse.ErrorDescription)
					} else if antiCaptchaResponse.Status == "processing" {
						time.Sleep(2 * time.Second)
						goto retry
					}
				}
			}
		}
	}
	return
}

func (a *AntiCaptcha) SolveRecaptchaV2(websiteUrl, websiteKey string, proxyData *ProxyData) (antiCaptchaResponse AntiCaptchaResponse, e error) {
	var taskType = antiCaptchaTypeRecaptchaV2Proxy
	if proxyData == nil {
		taskType = antiCaptchaTypeRecaptchaV2Proxyless
		proxyData = &ProxyData{}
	}

	r, resp, _ := a.Ses.Post(antiCaptchaCreateTaskUrl).
		Send(antiCaptchaNewTaskRequest{
			ClientKey: a.ApiKey,
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
			SoftID:       0,
			LanguagePool: "en",
		}).EndBytes()

	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()

		antiCaptchaResponse, e = a.waitForResponse(resp)
		antiCaptchaResponse.TaskType = taskType
	}
	return
}

func (a *AntiCaptcha) SolveRecaptchaEnterpriseV2(websiteUrl, websiteKey, s string, proxyData *ProxyData) (antiCaptchaResponse AntiCaptchaResponse, e error) {
	var taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxy
	if proxyData == nil {
		taskType = antiCaptchaTypeRecaptchaV2EnterpriseProxyless
		proxyData = &ProxyData{}
	}

	r, resp, _ := a.Ses.Post(antiCaptchaCreateTaskUrl).
		Send(antiCaptchaNewTaskRequest{
			ClientKey: a.ApiKey,
			Task: antiCaptchaTaskRequest{
				Type:       taskType,
				WebsiteURL: websiteUrl,
				WebsiteKey: websiteKey,
				EnterprisePayload: &antiCaptchaEnterprisePayload{
					S: s,
				},
				ProxyType:     proxyData.ProxyType,
				ProxyAddress:  proxyData.ProxyAddress,
				ProxyPort:     proxyData.ProxyPort,
				ProxyLogin:    proxyData.ProxyLogin,
				ProxyPassword: proxyData.ProxyPassword,
				UserAgent:     proxyData.UserAgent,
				Cookies:       proxyData.Cookies,
			},
			SoftID:       0,
			LanguagePool: "en",
		}).EndBytes()

	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()

		antiCaptchaResponse, e = a.waitForResponse(resp)
		antiCaptchaResponse.TaskType = taskType
	}
	return
}

func (a *AntiCaptcha) SolveImageCaptcha(img []byte) (antiCaptchaResponse AntiCaptchaResponse, e error) {
	r, resp, _ := a.Ses.Post(antiCaptchaCreateTaskUrl).
		Send(antiCaptchaNewTaskRequest{
			ClientKey: a.ApiKey,
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
		}).EndBytes()

	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()

		antiCaptchaResponse, e = a.waitForResponse(resp)
		antiCaptchaResponse.TaskType = antiCaptchaTypeImageToText
	}
	return
}

func (r *AntiCaptchaResponse) Report(good bool) error {
	if r.antiCaptchaInstance == nil {
		return errors.New("nil captcha instance")
	}

	var url string
	switch r.TaskType {
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
	re, resp, _ := r.antiCaptchaInstance.Ses.Post(url).
		Send(antiCaptchaGetTaskResultRequest{
			ClientKey: r.antiCaptchaInstance.ApiKey,
			TaskID:    r.TaskID,
		}).EndBytes()
	e := errors.New("nil response")
	if re != nil {
		_ = re.Body.Close()

		var reportResult antiCaptchaReportResult
		if e = json.Unmarshal(resp, &reportResult); e == nil {
			if reportResult.ErrorID != 0 {
				e = errors.New(reportResult.ErrorCode + ": " + reportResult.ErrorDescription)
			}
		}
	}
	return e
}

//F
