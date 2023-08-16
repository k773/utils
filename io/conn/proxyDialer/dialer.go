package proxyDialer

import (
	"errors"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
)

func New(proxyString string) (proxyDialer proxy.ContextDialer, e error) {
	u, e := url.Parse(proxyString)
	if e == nil {
		var proxyDialerWithoutContext proxy.Dialer
		proxyDialerWithoutContext, e = proxy.FromURL(u, &net.Dialer{})
		if e == nil {
			var ok bool
			if proxyDialer, ok = proxyDialerWithoutContext.(proxy.ContextDialer); !ok {
				e = errors.New("dialer is not assignable to proxy.DialerContext")
			}
		}
	}
	return
}
