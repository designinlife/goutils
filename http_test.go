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

func TestWebClient_DownloadFile(t *testing.T) {
	client := NewHTTPClient()
	err := client.DownloadFile("https://www.php.net/distributions/php-8.0.7.tar.gz", "D:/tmp/php-8.0.7.tar.gz", true)

	if err != nil {
		t.Errorf("Download errors. (%v)", err)
	}
}
