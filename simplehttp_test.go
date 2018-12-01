package simplehttp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	Verbose = true
}

func TestGet(t *testing.T) {
	body, err := Get().
		SetUrl("https://github.com/search").
		Query("q", "ljun20160606").
		Query("type", "Users").
		Query("utf-8", "✓").
		Send().
		Body()
	assert.NoError(t, err)
	assert.True(t, bytes.Contains(body, []byte("LJun")))
}

func TestGetHttp2(t *testing.T) {
	client := NewClient(&Builder{
		ProtoMajor: HTTP2,
	})
	request := NewRequest(client)
	get := request.SetUrl("https://v.qq.com").Get().Send()
	assert.Equal(t, get.ProtoMajor, int(HTTP2))
}
