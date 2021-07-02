package utils

import (
	"errors"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"strings"
	"time"
)

type SmsActivate struct {
	Ses   *gorequest.SuperAgent
	Token string
}

func NewSmsActivate(token string, proxy string, timeout time.Duration) SmsActivate {
	return SmsActivate{
		Ses:   gorequest.New().Proxy(proxy).Timeout(timeout),
		Token: token,
	}
}

func (s *SmsActivate) Balance() (b float64, e error) {
	r, resp, _ := s.Ses.Get("https://sms-activate.ru/stubs/handler_api.php?api_key=" + s.Token + "&action=getBalance").End()
	e = errors.New("nil response")
	if r != nil {
		e = errors.New(resp)
		if strings.HasPrefix(resp, "ACCESS_BALANCE:") {
			b, e = strconv.ParseFloat(strings.ReplaceAll(resp, "ACCESS_BALANCE:", ""), 64)
		}
	}
	return
}

type SmsActivateNumber struct {
	s *SmsActivate

	ID    string
	Phone string
}

// serviceCode, countryCode: https://sms-activate.ru/ru/api2/; russia=0, steam=mt
func (s *SmsActivate) GetNumber(serviceCode, countryCode string) (n SmsActivateNumber, e error) {
	n.s = s
retry:
	r, resp, _ := s.Ses.Get("https://sms-activate.ru/stubs/handler_api.php?api_key=" + s.Token + "&action=getNumber&service=" + serviceCode + "&country=" + countryCode).
		End()
	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()
		e = errors.New(resp)

		if resp == "NO_NUMBERS" {
			time.Sleep(time.Second)
			goto retry
		}
		if strings.HasPrefix(resp, "ACCESS_NUMBER:") {
			if spl := strings.Split(resp, ":"); len(spl) == 3 {
				e = nil
				n.ID, n.Phone = spl[1], spl[2]
			}
		}
	}

	return
}

type SmsActivateNumberStatus int

const (
	SmsActivateNumberStatusReady       = 1 // Sms were sent
	SmsActivateNumberStatusReceiveMore = 3 // Receive one more sms
	SmsActivateNumberStatusFinish      = 6 // Finish number activation
	SmsActivateNumberStatusCancel      = 8 // Cancel number activation
)

func (n *SmsActivateNumber) ChangeStatus(status SmsActivateNumberStatus) (s string, e error) {
	var r gorequest.Response
	r, s, _ = n.s.Ses.Get("https://sms-activate.ru/stubs/handler_api.php?api_key=" + n.s.Token + "&action=setStatus&status=" + strconv.Itoa(int(status)) + "&id=" + n.ID).End()
	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()
		e = nil
	}
	return
}

func (n *SmsActivateNumber) GetStatus() (status string, code string, e error) {
	r, resp, _ := n.s.Ses.Get("https://sms-activate.ru/stubs/handler_api.php?api_key=" + n.s.Token + "&action=getStatus&id=" + n.ID).End()
	e = errors.New("nil response")
	if r != nil {
		_ = r.Body.Close()
		e = nil

		spl := strings.Split(resp, ":")
		status = spl[0]
		if len(spl) == 2 {
			code = spl[1]
		}
	}

	return
}

func (n *SmsActivateNumber) GetSms(timeout time.Duration) (code string, e error) {
	var deadline = time.Now().Add(timeout)

	var status string
	for e == nil {
		if time.Now().After(deadline) {
			e = errors.New("timeout")
			continue
		}

		status, code, e = n.GetStatus()
		if e != nil {
			continue
		}

		if status == "STATUS_CANCEL" {
			e = errors.New("STATUS_CANCEL")
			break
		}
		if status == "STATUS_OK" {
			break
		}
		if status == "STATUS_WAIT_RESEND" {
			println("STATUS_WAIT_RESEND")
			//_, e = n.ChangeStatus(SmsActivateNumberStatusReceiveMore)
		}
		time.Sleep(2 * time.Second)
	}
	return
}
