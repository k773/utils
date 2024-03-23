package utils

import (
	"errors"
	"strconv"
	"strings"
)

// ParseProxy parses both proxy formats:
//
//	proto://login:password@host:port,
//	host:port:login:password.
//
// In both cases login and password are optional.
// In the case second format was used, returned proto would be empty.
func ParseProxy(str string) (proto, login, password, host string, port int, e error) {
	var schemeSepIndex = strings.Index(str, "://")
	if schemeSepIndex >= 0 {
		proto = str[:schemeSepIndex]
		var spl = strings.Split(str[schemeSepIndex+3:], "@")
		if len(spl) == 1 {
			var credentials = strings.Split(spl[0], ":")
			host = credentials[0]
			port, e = strconv.Atoi(credentials[1])
		} else if len(spl) == 2 {
			var credentials = strings.Split(spl[0], ":")
			var address = strings.Split(spl[1], ":")

			login, password = credentials[0], credentials[1]
			host = address[0]
			port, e = strconv.Atoi(address[1])
		} else {
			e = errors.New("incorrect proxy format")
		}
	} else {
		var spl = strings.Split(str, ":")
		if len(spl) >= 2 {
			host = spl[0]
			port, e = strconv.Atoi(spl[1])
			if len(spl) == 4 {
				login, password = spl[2], spl[3]
			}
		}
		if len(spl) != 2 && len(spl) != 4 {
			e = errors.New("incorrect proxy format")
		}
	}
	return
}

func ParseProxyToProxyData(src string) (data *ProxyData, e error) {
	proto, login, password, host, port, e := ParseProxy(src)
	if e == nil {
		data = &ProxyData{
			ProxyType:     proto,
			ProxyAddress:  host,
			ProxyPort:     port,
			ProxyLogin:    login,
			ProxyPassword: password,
			UserAgent:     "",
			Cookies:       "",
		}
	}
	return
}
