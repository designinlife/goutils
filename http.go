package goutils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"time"
)

type WebClient struct {
	Proxy          string
	ViaProxy       bool
	Headers        []WebHeader
	DefaultTimeout time.Duration
}

type WebHeader struct {
	Name  string
	Value string
}

func NewWebClient() *WebClient {
	return &WebClient{}
}

func NewWebClientWithTimeout(timeout time.Duration) *WebClient {
	return &WebClient{
		DefaultTimeout: timeout,
	}
}

func (wc *WebClient) Do(method, url string, timeout time.Duration) ([]byte, error) {
	request, _ := http.NewRequest(method, url, nil)

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	if wc.Headers != nil {
		for i := 0; i < len(wc.Headers); i++ {
			request.Header.Set(wc.Headers[i].Name, wc.Headers[i].Value)
		}
	}

	if timeout <= 0 {
		timeout = wc.DefaultTimeout
	}
	if timeout <= 0 {
		timeout = time.Second * 60
	}

	var client *http.Client

	if wc.Proxy != "" && wc.ViaProxy {
		proxy, _ := url2.Parse(wc.Proxy)

		tr := &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		client = &http.Client{
			Transport: tr,
			Timeout:   timeout,
		}
	} else {
		client = &http.Client{
			Timeout: timeout,
		}
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (wc *WebClient) Get(url string) ([]byte, error) {
	return wc.Do("GET", url, 0)
}

func (wc *WebClient) Head(url string) ([]byte, error) {
	return wc.Do("HEAD", url, 0)
}

func (wc *WebClient) Post(url string) ([]byte, error) {
	return wc.Do("POST", url, 0)
}

func (wc *WebClient) Put(url string) ([]byte, error) {
	return wc.Do("PUT", url, 0)
}

func (wc *WebClient) Delete(url string) ([]byte, error) {
	return wc.Do("DELETE", url, 0)
}

func (wc *WebClient) Options(url string) ([]byte, error) {
	return wc.Do("OPTIONS", url, 0)
}

func (wc *WebClient) Patch(url string) ([]byte, error) {
	return wc.Do("PATCH", url, 0)
}
