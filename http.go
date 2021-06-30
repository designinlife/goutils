package goutils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type WebClient struct {
	Debug          bool
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

func WithTimeout(timeout time.Duration) WebClientOption {
	return func(wc *WebClient) {
		wc.DefaultTimeout = timeout
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

func (wc *WebClient) doRequest(method, url string, queryParams map[string]string, formParams url2.Values, jsonData []byte, timeout time.Duration, filename string, showProgressBar bool) ([]byte, error) {
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

	if wc.Debug {
		logger.Debug(fmt.Sprintf("%s %s", method, url))
	}

	resp, err = client.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if filename != "" {
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			return nil, errors.New(fmt.Sprintf("A non-200 response status code was detected. (StatusCode: %d)", resp.StatusCode))
		}

		dirn := filepath.Dir(filename)

		if !IsDir(dirn) {
			os.MkdirAll(dirn, 0755)
		}

		out, err := os.Create(filename + ".tmp")
		if err != nil {
			return nil, err
		}

		// 此处不能使用 defer 方式关闭 out 资源，因为在 os.Rename 时资源句柄未释放造成重命名出错！
		// defer out.Close()

		var totalSize uint64

		totalSize = 0
		contentSize := resp.Header.Get("Content-Length")

		if contentSize != "" {
			intNum, _ := strconv.Atoi(contentSize)
			totalSize = uint64(intNum)
		}

		counter := &WriteCounter{ProgressBar: showProgressBar, TotalBytes: totalSize, OnlyShowPercentage: true}

		if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
			out.Close()
			return nil, err
		}

		out.Close()

		if showProgressBar {
			fmt.Print("\n")
		}

		if err = os.Rename(filename+".tmp", filename); err != nil {
			return nil, err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			return body, errors.New(fmt.Sprintf("检测到非 200 响应状态码。(StatusCode: %d)", resp.StatusCode))
		}

		return body, nil
	}

	return nil, nil
}

func (wc *WebClient) Get(url string) ([]byte, error) {
	return wc.doRequest("GET", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) GetWithQuery(url string, queryParams map[string]string) ([]byte, error) {
	return wc.doRequest("GET", url, queryParams, nil, nil, 0, "", false)
}

func (wc *WebClient) Head(url string) ([]byte, error) {
	return wc.doRequest("HEAD", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) Post(url string, formParams url2.Values) ([]byte, error) {
	return wc.doRequest("POST", url, nil, formParams, nil, 0, "", false)
}

func (wc *WebClient) PostWithJSON(url string, jsonData []byte) ([]byte, error) {
	return wc.doRequest("POST", url, nil, nil, jsonData, 0, "", false)
}

func (wc *WebClient) Put(url string) ([]byte, error) {
	return wc.doRequest("PUT", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) Delete(url string) ([]byte, error) {
	return wc.doRequest("DELETE", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) Options(url string) ([]byte, error) {
	return wc.doRequest("OPTIONS", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) Patch(url string) ([]byte, error) {
	return wc.doRequest("PATCH", url, nil, nil, nil, 0, "", false)
}

func (wc *WebClient) DownloadFile(url string, filename string, showProgress bool) error {
	_, err := wc.doRequest("GET", url, nil, nil, nil, 0, filename, showProgress)

	return err
}
