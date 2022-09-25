/*
	Interfaces for the captcha providers
*/

package types

import (
	"context"
	"github.com/k773/utils"
)

type CaptchaSolverInstance interface {
	SolveRecaptchaEnterpriseV2(ctx context.Context, websiteKey, websiteUrl, s string, proxyData *utils.ProxyData) (antiCaptchaResponse CaptchaResult, e error)
	SolveRecaptchaEnterpriseV2Domain(ctx context.Context, websiteKey, websiteUrl, s, domain string, proxyData *utils.ProxyData) (antiCaptchaResponse CaptchaResult, e error)
}

type CaptchaResult interface {
	Result() string
	Report(ctx context.Context, good bool) error
	ReportGood(ctx context.Context)
	ReportBad(ctx context.Context)
	IsZeroBalance() bool
	Cost() float64
}
