package simplehttp

import (
	"github.com/ljun20160606/simplehttp/simplehttputil"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Method        string
	Url           *strings.Builder
	Header        http.Header
	Querys        [][2]string
	Forms         map[string][]string
	Body          io.Reader
	JsonData      interface{}
	Cookies       []*http.Cookie
	Charset       simplehttputil.Charset
	IsClearCookie bool
	Retry         int
	Client        Client
}

func NewRequest(client Client) *Request {
	return &Request{Header: http.Header{}, Client: client}
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

func (r *Request) SetUrl(url string) *Request {
	builder := strings.Builder{}
	builder.WriteString(url)
	r.Url = &builder
	return r
}

func (r *Request) ClearCookie() *Request {
	r.IsClearCookie = true
	return r
}

func (r *Request) Query(k, v string) *Request {
	if r.Querys == nil {
		r.Querys = [][2]string{}
	}
	r.Querys = append(r.Querys, [2]string{k, v})
	return r
}

func (r *Request) QueryArray(k string, vs []string) *Request {
	if r.Querys == nil {
		r.Querys = [][2]string{}
	}
	for _, v := range vs {
		r.Querys = append(r.Querys, [2]string{k, v})
	}
	return r
}

func (r *Request) SetQuerys(querys [][2]string) *Request {
	r.Querys = querys
	return r
}

func (r *Request) Form(k string, v string) *Request {
	if r.Forms == nil {
		r.Forms = make(map[string][]string)
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
	r.Forms[k] = append(r.Forms[k], v)
	return r
}

func (r *Request) FormForce(k string, v string) *Request {
	if r.Forms == nil {
		r.Forms = make(map[string][]string)
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
	r.Forms[k] = []string{v}
	return r
}

func (r *Request) FormArray(k string, v []string) *Request {
	if r.Forms == nil {
		r.Forms = make(map[string][]string)
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
	r.Forms[k] = v
	return r
}

func (r *Request) SetForms(forms map[string][]string) *Request {
	if r.Forms == nil {
		r.Head(ContentType, "application/x-www-form-urlencoded")
	}
	if forms == nil {
		r.Forms = make(map[string][]string)
	} else {
		r.Forms = forms
	}
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
	r.Retry = retry
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
	r.Charset = simplehttputil.GB18030
	return r
}

func (r *Request) UTF8() *Request {
	r.Charset = simplehttputil.UTF8
	return r
}

func (r *Request) Send() (resp *Response) {
	return r.Client.Send(r)
}
