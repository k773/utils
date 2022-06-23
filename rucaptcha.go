package utils

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/go-resty/resty/v2"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type RuCaptcha struct {
//	Ses    *resty.Client
//	ApiKey string
//}
//
//type captchaResponseStruct struct {
//	Status  int    `json:"status"`
//	Request string `json:"request"`
//}
//
//type RuCaptchaResponse struct {
//	RuCaptchaInstance    *RuCaptcha
//	ResolvedSuccessfully bool
//	CaptchaResponse      string
//	CaptchaID            string
//}
//
//const (
//	RuCaptchaErrorCaptchaNotReady   = "CAPCHA_NOT_READY" // not typo
//	RuCaptchaErrorCaptchaUnsolvable = "ERROR_CAPTCHA_UNSOLVABLE"
//	RuCaptchaErrorWrongCaptchaID    = "ERROR_WRONG_CAPTCHA_ID"
//)
//
//func (a *RuCaptcha) getCaptchaStatus(id string) (res captchaResponseStruct, e error) {
//	r, e := a.Ses.R().G(fmt.Sprintf("https://2captcha.com/res.php?key=%s&action=get&taskinfo=0&json=1&id=%s", a.ApiKey, id)).EndBytes()
//	if r != nil {
//		_ = r.Body.Close()
//		e = json.Unmarshal(resp, &res)
//	} else {
//		e = errors.New("nil response")
//	}
//	return res, e
//}
//
//func (a *RuCaptcha) waitForResponse(capResp captchaResponseStruct) (res RuCaptchaResponse, e error) {
//	//fmt.Println(capResp)
//	if capResp.Status != 1 {
//		e = errors.New("wrong rucaptcha status: " + strconv.Itoa(capResp.Status))
//	} else {
//		res.CaptchaID = capResp.Request
//
//		capResp, e = a.getCaptchaStatus(res.CaptchaID)
//		for e == nil && capResp.Request == RuCaptchaErrorCaptchaNotReady {
//			capResp, e = a.getCaptchaStatus(res.CaptchaID)
//			time.Sleep(2000 * time.Millisecond)
//		}
//		res.CaptchaResponse = capResp.Request
//		res.ResolvedSuccessfully = !strings.HasPrefix(res.CaptchaResponse, "ERROR")
//	}
//	return
//}
//
//func (a *RuCaptcha) SolveImageCaptcha(file []byte) (res RuCaptchaResponse, e error) {
//	res.RuCaptchaInstance = a
//
//begin:
//	r, resp, _ := a.Ses.Post("http://2captcha.com/in.php?json=1").
//		Type("multipart").
//		Send(`{"key": "`+a.ApiKey+`"}`).
//		SendFile(file, "captcha", "file", true).
//		EndBytes()
//	if r != nil {
//		_ = r.Body.Close()
//
//		var capResp captchaResponseStruct
//		_ = json.Unmarshal(resp, &capResp)
//		res, e = a.waitForResponse(capResp)
//		if res.CaptchaResponse == RuCaptchaErrorCaptchaUnsolvable {
//			goto begin
//		}
//	}
//	return res, e
//}
//
//// Returns: capResponse, capID, error
//func (a *RuCaptcha) SolveRecaptchaEnterpriseV2(url, key, dataS string) (res RuCaptchaResponse, e error) {
//	res.RuCaptchaInstance = a
//
//begin:
//	r, response, _ := a.Ses.Get(fmt.Sprintf("http://2captcha.com/in.php?key=%v&method=userrecaptcha&googlekey=%v&pageurl=%v&data-s=%v&json=1&enterprise=1", a.ApiKey, key, url, dataS)).EndBytes()
//	if r != nil {
//		_ = r.Body.Close()
//
//		var capResp captchaResponseStruct
//		_ = json.Unmarshal(response, &capResp)
//		res, e = a.waitForResponse(capResp)
//		if res.CaptchaResponse == RuCaptchaErrorCaptchaUnsolvable {
//			goto begin
//		}
//	}
//
//	return
//}
//
//// Not working: repair before using
//func (a *RuCaptcha) SolveRecaptchaV3(url, action, key string) (res RuCaptchaResponse, e error) {
//	res.RuCaptchaInstance = a
//
//begin:
//	r, response, _ := a.Ses.Get(fmt.Sprintf("https://2captcha.com/in.php?key=%s&method=userrecaptcha&version=v3&action=%s&min_score=0.9&googlekey=%s&pageurl=%s&json=1", a.ApiKey, action, key, url)).EndBytes()
//	if r != nil {
//		_ = r.Body.Close()
//	}
//
//	var capResp captchaResponseStruct
//	_ = json.Unmarshal(response, &capResp)
//	if capResp.Status != 1 {
//		e = errors.New("wrong rucaptcha status")
//	} else {
//		res.CaptchaID = capResp.Request
//
//		var capResp, e = a.getCaptchaStatus(res.CaptchaID)
//		for ; e == nil && capResp.Request == RuCaptchaErrorCaptchaNotReady; capResp, e = a.getCaptchaStatus(capResp.Request) {
//			time.Sleep(2000 * time.Millisecond)
//		}
//		if capResp.Request == RuCaptchaErrorCaptchaUnsolvable {
//			goto begin
//		}
//		res.CaptchaResponse = capResp.Request
//	}
//
//	return res, e
//}
//
//// Deprecated
//// Not working
//func (a *RuCaptcha) SolveRecaptchaV2(url, key string) (string, string) { //returns cap-response, cap-id
//startCaptcha:
//	gurl := "https://2captcha.com/in.php?key=%s&method=userrecaptcha&version=v2&" +
//		"googlekey=%s&pageurl=%s&json=1"
//	gurl = fmt.Sprintf(gurl, a.ApiKey, key, url)
//	_, response, _ := a.Ses.Get(gurl).EndBytes()
//	var captchaResponse1 captchaResponseStruct
//	_ = json.Unmarshal(response, &captchaResponse1)
//	if captchaResponse1.Status != 1 {
//		goto startCaptcha
//	}
//
//	var capchaResponse2 captchaResponseStruct
//waitForCaptcha:
//	switch capchaResponse2.Request {
//	case "CAPCHA_NOT_READY", "":
//		_, capchaResponse2B, _ := a.Ses.Get(fmt.Sprintf("https://2captcha.com/res.php?key=%s&action=get&taskinfo=0&json=1&id=%s", a.ApiKey, captchaResponse1.Request)).EndBytes()
//		_ = json.Unmarshal(capchaResponse2B, &capchaResponse2)
//		time.Sleep(2000 * time.Millisecond)
//		goto waitForCaptcha
//	case "ERROR_CAPTCHA_UNSOLVABLE":
//		goto startCaptcha
//	}
//	return capchaResponse2.Request, captchaResponse1.Request
//}
//
//func (ru *RuCaptchaResponse) CapReport(good bool) {
//	if ru.RuCaptchaInstance == nil {
//		panic("1")
//	}
//
//	var action string
//	if good {
//		action = "reportgood"
//	} else {
//		action = "reportbad"
//	}
//	r, aga, _ := ru.RuCaptchaInstance.Ses.Get(fmt.Sprintf("https://2captcha.com/res.php?key=%s&action=%s&id=%s", ru.RuCaptchaInstance.ApiKey, action, ru.CaptchaID)).End()
//	if r != nil {
//		_ = r.Body.Close()
//	}
//	fmt.Println(aga)
//}
//
////F
