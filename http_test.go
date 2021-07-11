package goutils

import (
	"fmt"
	"testing"
)

func TestHttpClient_Get(t *testing.T) {
	client := NewHttpClient()
	resp1, err := client.Get("https://postman-echo.com/get?foo1=bar1&foo2=bar2", nil)

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp1)

	resp2, err := client.Get("https://postman-echo.com/get?foo1=bar1&foo2=bar2", &HttpRequest{
		Proxy: "http://127.0.0.1:3128",
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp2)
}

func TestHttpClient_Post(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Post("https://postman-echo.com/post", &HttpRequest{
		JSON: map[string]interface{}{"foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}

func TestHttpClient_Put(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Put("https://postman-echo.com/put", &HttpRequest{
		JSON: map[string]interface{}{"foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
	})

	if err != nil {
		t.Errorf("Request errors. (%v)", err)
	}

	fmt.Println(resp)
}

// func TestWebClient_DownloadFile(t *testing.T) {
// 	// client := NewHTTPClient()
// 	client := NewHTTPClientWithOptions(HTTPOptionWithProxy("http://127.0.0.1:3128"), HTTPOptionWithTimeout(300*time.Second))
// 	err := client.DownloadFile("https://www.python.org/ftp/python/3.9.6/python-3.9.6-amd64.exe", filepath.Join(os.TempDir(), "/python-3.9.6-amd64.exe"), true)
//
// 	if err != nil {
// 		t.Errorf("Download errors. (%v)", err)
// 	}
// }
