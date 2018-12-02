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

func defaultHttpRoundTripperFunc() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           DefaultDialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

var RoundTripFactory = map[ProtoMajor]func() http.RoundTripper{
	HTTP1: defaultHttpRoundTripperFunc,
	HTTP2: func() http.RoundTripper {
		tripper := defaultHttpRoundTripperFunc().(*http.Transport)
		_ = http2.ConfigureTransport(tripper)
		return tripper
	},
}

func (t ProtoMajor) RoundTripper() http.RoundTripper {
	return RoundTripFactory[t]()
}
