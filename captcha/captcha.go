package captcha

import "github.com/k773/utils"

type Instance interface {
	SolveRecaptchaEnterprise(siteKey, dataS string, proxy *utils.ProxyData)
}
