package goutils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type WebClient struct {
	Proxy          string
	Headers        []WebHeader
	DefaultTimeout time.Duration
}

type WebHeader struct {
	Name  string
	Value string
}

type WebClientOption func(*WebClient)

func NewHTTPClient() *WebClient {
	return &WebClient{}
}

func NewHTTPClientWithTimeout(timeout time.Duration) *WebClient {
	return &WebClient{
		DefaultTimeout: timeout,
	}
}

func NewHTTPClientWithOptions(opts ...WebClientOption) *WebClient {
	wc := &WebClient{}

	for _, opt := range opts {
		opt(wc)
	}

	return wc
}

func WithProxy(proxyUrl string) WebClientOption {
	return func(wc *WebClient) {
		wc.Proxy = proxyUrl
	}
}

func WithHeaders(headers map[string]string) WebClientOption {
	return func(wc *WebClient) {
		for k, v := range headers {
			wc.Headers = append(wc.Headers, WebHeader{
				Name:  k,
				Value: v,
			})
		}
	}
}

func (wc *WebClient) Do(method, url string, queryParams map[string]string, formParams url2.Values, jsonData []byte, timeout time.Duration) ([]byte, error) {
	if queryParams != nil {
		query := make([]string, 0)

		for k, v := range queryParams {
			query = append(query, fmt.Sprintf("%s=%s", k, url2.QueryEscape(v)))
		}

		if strings.Contains(url, "?") {
			url = url + "&" + strings.Join(query, "&")
		} else {
			url = url + "?" + strings.Join(query, "&")
		}
	}

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

	if wc.Proxy != "" {
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

	var resp *http.Response
	var err error

	switch method {
	case "POST":
		if formParams != nil {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			request.PostForm = formParams
		}
		break
	}

	if jsonData != nil {
		request.Header.Set("Content-Type", "application/json")
		request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
	}

	resp, err = client.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (wc *WebClient) Get(url string) ([]byte, error) {
	return wc.Do("GET", url, nil, nil, nil, 0)
}

func (wc *WebClient) GetWithQuery(url string, queryParams map[string]string) ([]byte, error) {
	return wc.Do("GET", url, queryParams, nil, nil, 0)
}

func (wc *WebClient) Head(url string) ([]byte, error) {
	return wc.Do("HEAD", url, nil, nil, nil, 0)
}

func (wc *WebClient) Post(url string, formParams url2.Values) ([]byte, error) {
	return wc.Do("POST", url, nil, formParams, nil, 0)
}

func (wc *WebClient) PostWithJSON(url string, jsonData []byte) ([]byte, error) {
	return wc.Do("POST", url, nil, nil, jsonData, 0)
}

func (wc *WebClient) Put(url string) ([]byte, error) {
	return wc.Do("PUT", url, nil, nil, nil, 0)
}

func (wc *WebClient) Delete(url string) ([]byte, error) {
	return wc.Do("DELETE", url, nil, nil, nil, 0)
}

func (wc *WebClient) Options(url string) ([]byte, error) {
	return wc.Do("OPTIONS", url, nil, nil, nil, 0)
}

func (wc *WebClient) Patch(url string) ([]byte, error) {
	return wc.Do("PATCH", url, nil, nil, nil, 0)
}
