package main

import (
	"fmt"
	"github.com/ljun20160606/simplehttp"
	"time"
)

type C map[string][]byte

func (c *C) Set(key string, value []byte, exp time.Duration) {
	(*c)[key] = value
}

func (c *C) Get(key string) ([]byte, bool) {
	val, ok := (*c)[key]
	return val, ok
}

func main() {
	simplehttp.Verbose = true
	// http get https://github.com/search
	fmt.Println(
		simplehttp.
			Get().
			SetUrl("https://github.com/search").
			Query("q", "simplehttp").
			Query("utf8", "✓").
			Send().
			String())

	var c C = make(map[string][]byte)
	user := "user1"
	h := simplehttp.NewClient(&simplehttp.Builder{Cache: &c, SessionID: user})

	// login github
	req := simplehttp.NewRequest(h)
	dom, err := req.Get().SetUrl("http://github.com/login").Send().Document()
	if err != nil {
		panic(err)
	}
	authenticityToken, _ := dom.Find("input[name=authenticity_token]").Attr("value")
	fmt.Println(
		req.Post().
			SetUrl("https://github.com/session").
			Form("commit", "Sign in").
			Form("utf8", "✓").
			Form("authenticity_token", authenticityToken).
			Form("login", "").
			Form("password", "").
			Send().
			String())
}
