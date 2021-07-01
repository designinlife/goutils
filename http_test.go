package goutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	// client := NewHTTPClient()
	client := NewHTTPClientWithOptions(HTTPOptionWithProxy("http://127.0.0.1:3128"), HTTPOptionWithTimeout(300*time.Second))
	err := client.DownloadFile("https://www.python.org/ftp/python/3.9.6/python-3.9.6-amd64.exe", filepath.Join(os.TempDir(), "/python-3.9.6-amd64.exe"), true)

	if err != nil {
		t.Errorf("Download errors. (%v)", err)
	}
}
