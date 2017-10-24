package simplehttputil

import (
	"fmt"
	"net/http"
	"net/url"
)

// no request.url valid check
func SimpleHttpProxy(proxy string) func(*http.Request) (*url.URL, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil ||
		(proxyURL.Scheme != "http" &&
			proxyURL.Scheme != "https" &&
			proxyURL.Scheme != "socks5") {
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return func(request *http.Request) (*url.URL, error) {
				return proxyURL, nil
			}
		}
	}
	if err != nil {
		return func(request *http.Request) (*url.URL, error) {
			return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
		}
	}
	return func(request *http.Request) (*url.URL, error) {
		return proxyURL, nil
	}
}
