package simplehttp

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"time"
)

type ProtoMajor int

const (
	_ ProtoMajor = iota
	HTTP1
	HTTP2
)

var RoundTripFactory = map[ProtoMajor]func() http.RoundTripper{
	HTTP1: func() http.RoundTripper {
		return &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           DefaultDialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	},
	HTTP2: func() http.RoundTripper {
		return &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	},
}

func (t ProtoMajor) RoundTripper() http.RoundTripper {
	return RoundTripFactory[t]()
}
