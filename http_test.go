package goutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHttpClient_Get(t *testing.T) {
	client := NewHttpClient()
	resp1, err := client.Get("https://postman-echo.com/get?foo1=bar1&foo2=bar2", &HttpRequest{
		Proxy: "http://127.0.0.1:3128",
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp1)

	resp2, err := client.Get("https://postman-echo.com/get?foo1=bar1&foo2=bar2", &HttpRequest{
		Proxy: "http://127.0.0.1:3128",
		Query: map[string]interface{}{"foo3": "bar3", "foo4": "bar4", "custom": "1"},
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp2)
}

func TestHttpClient_Post(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Post("https://postman-echo.com/post", &HttpRequest{
		FormParams: map[string]interface{}{"foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
		Proxy:      "http://127.0.0.1:3128",
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}

func TestHttpClient_Put(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Put("https://postman-echo.com/put", &HttpRequest{
		JSON:  map[string]interface{}{"foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
		Proxy: "http://127.0.0.1:3128",
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}

func TestHttpClient_DownloadFile(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Get("https://www.python.org/ftp/python/3.9.6/python-3.9.6-amd64.exe", &HttpRequest{
		ToFile:      filepath.Join(os.TempDir(), "/python-3.9.6-amd64.exe"),
		ProgressBar: true,
		Proxy:       "http://127.0.0.1:3128",
		Timeout:     time.Second * 300,
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}
