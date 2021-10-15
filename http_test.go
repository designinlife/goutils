package goutils

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"net/http/cookiejar"
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

func TestHttpClient_GetCookieJar(t *testing.T) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	client := NewHttpClient()
	resp1, err := client.Get("https://live.kuaishou.com/profile/3xr4nqdfsbgxyy6", &HttpRequest{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Referer":      "https://live.kuaishou.com",
		},
		CookieJar: jar,
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

type tProgressBar struct {
	FileName    string
	LoadedBytes uint64
	TotalBytes  uint64
}

func (bar *tProgressBar) SetTotalBytes(totalBytes uint64) {
	bar.TotalBytes = totalBytes
}

func (bar *tProgressBar) Write(p []byte) (int, error) {
	n := len(p)
	bar.LoadedBytes += uint64(n)

	// fmt.Println(bar.LoadedBytes)
	// fmt.Println(bar.TotalBytes)

	fmt.Printf("\r%s (%.2f%%)", bar.FileName, float64(bar.LoadedBytes)*100.00/float64(bar.TotalBytes))

	return n, nil
}

func (bar *tProgressBar) Close() error {
	fmt.Println("Completed. [bar]")

	return nil
}

func TestHttpClient_DownloadFile(t *testing.T) {
	bar := &tProgressBar{
		FileName: "composer-2.1.3.phar",
	}

	client := NewHttpClient(HttpClientOptionWithProgressBar(bar))
	resp, err := client.Get("https://www.python.org/ftp/python/3.9.6/python-3.9.6-amd64.exe", &HttpRequest{
		ToFile:      filepath.Join(os.TempDir(), "python-3.9.6-amd64.exe"),
		ProgressBar: true,
		Proxy:       "http://127.0.0.1:3128",
		Timeout:     time.Second * 300,
	})
	// resp, err := client.Get("https://getcomposer.org/download/2.1.3/composer.phar", &HttpRequest{
	// 	ToFile:      filepath.Join(os.TempDir(), "composer-2.1.3.phar"),
	// 	ProgressBar: true,
	// 	Proxy:       "http://127.0.0.1:3128",
	// 	Timeout:     time.Second * 300,
	// })

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}
