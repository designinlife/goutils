package goutils

import (
	"fmt"
	"strings"
	"testing"
)

func TestWebClient_Get(t *testing.T) {
	client := NewHTTPClient()
	content, err := client.Get("https://icanhazip.com/")

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(strings.TrimSpace(string(content)))
}
