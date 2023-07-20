package proxyDialer

import (
	"context"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
	"time"
)

type Dialer struct {
	d *net.Dialer

	ctx     context.Context
	timeout time.Duration
}

func (d *Dialer) Dial(network, addr string) (c net.Conn, err error) {
	return d.d.DialContext(d.ctx, network, addr)
}

func New(ctx context.Context, timeout time.Duration, proxyString string) (proxyDialer proxy.Dialer, e error) {
	var nextDialer = &Dialer{
		d: &net.Dialer{
			Timeout: timeout,
		},
		ctx:     ctx,
		timeout: timeout,
	}

	u, e := url.Parse(proxyString)
	if e == nil {
		proxyDialer, e = proxy.FromURL(u, nextDialer)
	}
	return
}
