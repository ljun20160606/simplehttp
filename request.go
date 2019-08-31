package simplehttp

import (
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Method   string
	Url      *strings.Builder
	Header   http.Header
	Querys   [][2]string
	Forms    [][2]string
	Body     io.Reader
	JsonData interface{}
	Cookies  []*http.Cookie
	Charset  Charset
	Client   Client
	Config   RequestConfig
}

type RequestConfig struct {
	IsClearCookie bool
	Retry         int
}

func NewRequest(client Client) *Request {
	return NewPureRequest().SetClient(client)
}

func NewPureRequest() *Request {
	return &Request{Header: http.Header{}, Charset: UTF8}
}

func (r *Request) SetClient(client Client) *Request {
	r.Client = client
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) Get() *Request {
	r.Method = http.MethodGet
	return r
}

func (r *Request) Post() *Request {
	r.Method = http.MethodPost
	return r
}

func (r *Request) Patch() *Request {
	r.Method = http.MethodPatch
	return r
}

func (r *Request) Connect() *Request {
	r.Method = http.MethodConnect
	return r
}

func (r *Request) Delete() *Request {
	r.Method = http.MethodDelete
	return r
}

func (r *Request) MethodHead() *Request {
	r.Method = http.MethodHead
	return r
}

func (r *Request) Options() *Request {
	r.Method = http.MethodOptions
	return r
}

func (r *Request) SetCookies(c []*http.Cookie) *Request {
	r.Cookies = c
	return r
}

func (r *Request) SetStringCookies(cookie string) *Request {
	return r.SetCookies(ReadCookies([]string{cookie}, ""))
}

func (r *Request) AddCookie(ck *http.Cookie) *Request {
	if r.Cookies == nil {
		r.Cookies = []*http.Cookie{ck}
		return r
	}
	r.Cookies = append(r.Cookies, ck)
	return r
}

func (r *Request) SetUrl(rawurl string) *Request {
	builder := strings.Builder{}
	builder.WriteString(rawurl)
	r.Url = &builder
	return r
}

func (r *Request) ClearCookie() *Request {
	r.Config.IsClearCookie = true
	return r
}

func (r *Request) Query(k, v string) *Request {
	if r.Querys == nil {
		r.Querys = [][2]string{}
	}
	r.Querys = append(r.Querys, [2]string{k, v})
	return r
}

func (r *Request) SetQuerys(querys [][2]string) *Request {
	r.Querys = querys
	return r
}

func (r *Request) Form(k string, v string) *Request {
	r.formInit()
	r.Forms = append(r.Forms, [2]string{k, v})
	return r
}

func (r *Request) formInit() {
	if r.Forms == nil {
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
}

func (r *Request) SetForms(forms [][2]string) *Request {
	if forms == nil {
		return r
	}
	if r.Forms == nil {
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
	r.Forms = forms
	return r
}

func (r *Request) SetJSON(data interface{}) *Request {
	r.JsonData = data
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	r.Body = body
	return r
}

func (r *Request) SetRetry(retry int) *Request {
	r.Config.Retry = retry
	return r
}

func (r *Request) RefererInHeader(refer string) *Request {
	return r.Head(Referer, refer)
}

func (r *Request) OriginInHeader(origin string) *Request {
	return r.Head(Origin, origin)
}

func (r *Request) Head(k, v string) *Request {
	r.Header.Set(k, v)
	return r
}

func (r *Request) GB18030() *Request {
	r.Charset = GB18030
	return r
}

func (r *Request) UTF8() *Request {
	r.Charset = UTF8
	return r
}

func (r *Request) Send() (resp *Response) {
	return r.Client.Send(r)
}
