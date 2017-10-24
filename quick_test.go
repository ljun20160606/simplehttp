package simplehttp

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	fmt.Println(Get("https://github.com/search").
		Query("utf8", "✓").
		Query("q", "httpclient").
		Send().
		String())
}
