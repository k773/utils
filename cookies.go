package utils

import (
	"github.com/go-resty/resty/v2"
	"net/url"
)

func GetCookie(ses *resty.Client, urlStr, cookieName string) string {
	var cookie, _, _ = GetCookieMore(ses, urlStr, cookieName)
	return cookie
}

func GetCookieMore(ses *resty.Client, urlStr, cookieName string) (cookie string, has bool, e error) {
	u, e := url.Parse(urlStr)
	if e == nil {
		for _, cookie := range ses.GetClient().Jar.Cookies(u) {
			if cookie.Name == cookieName {
				return cookie.Value, true, nil
			}
		}
	}
	return
}
