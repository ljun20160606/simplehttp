package simplehttp

import (
	"bytes"
	"encoding/json"
	"github.com/LFZJun/cookiejar"
	"github.com/LFZJun/simplehttp/simplehttputil"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"io/ioutil"
	"net/http"
	"strings"
)

type Driver interface {
	Send(req *Request) *Response
}

var DefaultDriver = &HttpDriver{Client: DefaultClient}

type HttpDriver struct {
	Client      *http.Client
	StoreCookie StoreCookie
}

func (h *HttpDriver) send(realReq *http.Request) (resp *Response) {
	response, err := h.Client.Do(realReq)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	resp = new(Response)
	if err != nil {
		resp.err = err
		return
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp.err = err
		return
	}
	if h.StoreCookie != nil {
		h.StoreCookie(h.Client.Jar)
	}
	var enc encoding.Encoding
	contentType := response.Header.Get(ContentType)
	if len(contentType) > 0 {
		subMatch := ContentTypeMatchCharset.FindStringSubmatch(contentType)
		var name string
		if len(subMatch) == 2 {
			if Verbose {
				logger.Println("find html Encode ", subMatch[1])
			}
			name = subMatch[1]
		}
		if name != "" {
			enc, _ = htmlindex.Get(name)
		}
	}
	return &Response{code: response.StatusCode, body: data, header: response.Header, url: response.Request.URL, encoding: enc}
}

func (h *HttpDriver) Send(r *Request) (resp *Response) {
	var err error
	if r.querys != nil {
		r.url.WriteByte('?')
		r.url.Write(simplehttputil.BuildQueryEncoded(r.querys, r.charset))
	}
	switch {
	case r.body != nil:
	case r.forms != nil:
		r.body = bytes.NewReader(simplehttputil.BuildFormEncoded(r.forms, r.charset))
	case r.jsonData != nil:
		body, err := json.Marshal(r.jsonData)
		if err != nil {
			resp.err = err
			return
		}
		r.body = bytes.NewReader(body)
	}
	resp = new(Response)
	realReq, err := http.NewRequest(r.method, r.url.String(), r.body)
	if err != nil {
		resp.err = err
		return
	}
	realReq.Header = r.header
	if r.clearCookies || h.Client.Jar == nil {
		h.Client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	}
	if r.cookies != nil {
		h.Client.Jar.SetCookies(realReq.URL, r.cookies)
	}
	switch r.retry {
	case 0:
		resp = h.send(realReq)
	default:
		for times := -1; times < r.retry; times++ {
			resp = h.send(realReq)
			if resp.err == nil || !strings.Contains(resp.err.Error(), "request canceled") {
				break
			}
		}
	}
	return resp
}
