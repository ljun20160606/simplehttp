package simplehttp

import (
	"bytes"
	"testing"
)

func init() {
	Verbose = true
}

func TestGet(t *testing.T) {
	if !bytes.Contains(
		Get("https://github.com/search").
			Query("q", "ljun20160606").
			Query("type", "Users").
			Query("utf-8", "✓").
			Send().
			Bytes(), []byte("LJun")) {
		t.Fail()
	}
}
