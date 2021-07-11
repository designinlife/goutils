package goutils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HttpClient struct {
}

type HttpRequest struct {
	// 超时 (默认值: 60s)
	Timeout time.Duration
	// 查询参数
	Query interface{}
	// 请求头信息
	Headers map[string]interface{}
	// Cookie 信息
	Cookies interface{}
	// POST 表单参数
	FormParams map[string]interface{}
	// JSON 数据参数
	JSON interface{}
	// XML 数据参数
	XML interface{}
	// 代理服务器地址
	Proxy string
	// 是否服务器响应非 200 状态时，返回 error？
	AllowNon200Response bool
	// 下载到本地文件
	ToFile string
	// 是否显示进度条？
	ProgressBar bool
}

type HttpResponse struct {
	StatusCode    int
	RequestURI    string
	Header        http.Header
	ContentType   string
	ContentLength int64
	Body          []byte
}

func (h HttpResponse) String() string {
	return fmt.Sprintf("%d %s\n%s (%d)\n%s", h.StatusCode, h.RequestURI, h.ContentType, h.ContentLength, string(h.Body))
}

func (h HttpResponse) ToString() string {
	return string(h.Body)
}

func (h HttpResponse) ToFile(filename string) error {
	return ioutil.WriteFile(filename, h.Body, 0644)
}

func (h HttpResponse) ToJson(v interface{}) error {
	return json.Unmarshal(h.Body, v)
}

func (h HttpResponse) ToXml(v interface{}) error {
	return xml.Unmarshal(h.Body, v)
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) Request(method, uri string, r *HttpRequest) (*HttpResponse, error) {
	// 创建 HTTP 客户端实例
	req, _ := http.NewRequest(method, uri, nil)

	// 设置 Headers
	if r != nil && r.Headers != nil {
		for k, v := range r.Headers {
			if vv, ok := v.(string); ok {
				req.Header.Set(k, vv)
				continue
			}
			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					req.Header.Add(k, vvv)
				}
			}
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}

	// 设置 Cookies
	if r != nil {
		switch r.Cookies.(type) {
		case string:
			cookies := r.Cookies.(string)
			req.Header.Add("Cookie", cookies)
		case map[string]string:
			cookies := r.Cookies.(map[string]string)
			for k, v := range cookies {
				req.AddCookie(&http.Cookie{
					Name:  k,
					Value: v,
				})
			}
		case []*http.Cookie:
			cookies := r.Cookies.([]*http.Cookie)
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
		}
	}

	// 设置 Query 查询参数
	if r != nil {
		switch r.Query.(type) {
		case string:
			str := r.Query.(string)
			req.URL.RawQuery = str
		case map[string]interface{}:
			q := req.URL.Query()
			for k, v := range r.Query.(map[string]interface{}) {
				if vv, ok := v.(string); ok {
					q.Set(k, vv)
					continue
				}
				if vv, ok := v.([]string); ok {
					for _, vvv := range vv {
						q.Add(k, vvv)
					}
				}
			}
			req.URL.RawQuery = q.Encode()
		}
	}

	// 设置 Form 表单参数
	if r != nil {
		if r.FormParams != nil {
			if _, ok := r.Headers["Content-Type"]; !ok {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			values := url.Values{}
			for k, v := range r.FormParams {
				if vv, ok := v.(string); ok {
					values.Set(k, vv)
				}
				if vv, ok := v.([]string); ok {
					for _, vvv := range vv {
						values.Add(k, vvv)
					}
				}
			}
			req.Body = ioutil.NopCloser(strings.NewReader(values.Encode()))
		}
	}

	// 设置 JSON 请求
	if r != nil {
		if r.JSON != nil {
			if _, ok := r.Headers["Content-Type"]; !ok {
				req.Header.Set("Content-Type", "application/json")
			}

			b, err := json.Marshal(r.JSON)
			if err == nil {
				req.Body = ioutil.NopCloser(bytes.NewReader(b))
			}
		}
	}

	// 设置 XML 请求
	if r != nil {
		if r.XML != nil {
			if _, ok := r.Headers["Content-Type"]; !ok {
				req.Header.Set("Content-Type", "application/xml")
			}

			switch r.XML.(type) {
			case map[string]string:
				b, err := map2XML(r.XML.(map[string]string))
				if err == nil {
					req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
				}
			default:
				b, err := xml.Marshal(r.JSON)
				if err == nil {
					req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
				}
			}
		}
	}

	// 创建客户端并发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if r != nil && r.Proxy != "" {
		proxy, err := url.Parse(r.Proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxy)
		}
	}

	timeout := time.Second * 60

	if r != nil && r.Timeout > 0 {
		timeout = r.Timeout
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// 检查非 200 响应状态
	allowNon200 := false

	if r != nil {
		allowNon200 = r.AllowNon200Response
	}

	if !allowNon200 && !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, errors.New(fmt.Sprintf("A non-200 response status code was detected. (StatusCode: %d)", resp.StatusCode))
	}

	var content []byte

	if r != nil && r.ToFile != "" {
		filename := r.ToFile

		dirn := filepath.Dir(filename)

		if !IsDir(dirn) {
			if err := os.MkdirAll(dirn, 0755); err != nil {
				return nil, err
			}
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

		showProgressBar := false

		if isatty.IsTerminal(os.Stdout.Fd()) && r.ProgressBar {
			showProgressBar = true
		}

		if showProgressBar {
			counter := &progressBarCounter{ProgressBar: true, TotalBytes: totalSize, SimpleBarStyle: true}

			if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
				out.Close()
				return nil, err
			}
		} else {
			if _, err = io.Copy(out, resp.Body); err != nil {
				out.Close()
				return nil, err
			}
		}

		out.Close()

		if err = os.Rename(filename+".tmp", filename); err != nil {
			return nil, err
		}
	} else {
		content, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}
	}

	ret := &HttpResponse{
		RequestURI:    resp.Request.RequestURI,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		Body:          content,
	}

	return ret, nil
}

func (h *HttpClient) Get(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("GET", uri, r)
}

func (h *HttpClient) Post(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("POST", uri, r)
}

func (h *HttpClient) Put(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("PUT", uri, r)
}

func (h *HttpClient) Delete(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("DELETE", uri, r)
}

func (h *HttpClient) Options(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("OPTIONS", uri, r)
}

func (h *HttpClient) Head(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("HEAD", uri, r)
}

func (h *HttpClient) Patch(uri string, r *HttpRequest) (*HttpResponse, error) {
	return h.Request("PATCH", uri, r)
}

func map2XML(m map[string]string, opts ...interface{}) ([]byte, error) {
	rootTag := "xml"
	if len(opts) > 0 {
		// the first opts is the root tag of xml struct
		if v, ok := opts[0].(string); ok {
			rootTag = v
		}
	}

	d := mdata{xml.Name{Local: rootTag}, m}

	bt, err := xml.Marshal(d)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

type mdata struct {
	XMLName xml.Name
	data    map[string]string
}

// MarshalXML xml encode
func (m mdata) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m.data) == 0 {
		return nil
	}

	start.Name.Local = m.XMLName.Local

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m.data {
		if strings.HasPrefix(v, "cdata:") {
			v = strings.Replace(v, "cdata:", "", 1)
			xs := struct {
				XMLName xml.Name
				Value   interface{} `xml:",cdata"`
			}{xml.Name{Local: k}, v}
			e.Encode(xs)
		} else {
			xs := struct {
				XMLName xml.Name
				Value   interface{} `xml:",chardata"`
			}{xml.Name{Local: k}, v}
			e.Encode(xs)
		}
	}

	return e.EncodeToken(start.End())
}

// progressBarCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type progressBarCounter struct {
	LoadedBytes    uint64
	TotalBytes     uint64
	ProgressBar    bool
	SimpleBarStyle bool
}

func (w *progressBarCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.LoadedBytes += uint64(n)

	if w.ProgressBar {
		w.printProgress()
	}

	return n, nil
}

func (w progressBarCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	if w.TotalBytes > 0 {
		if w.SimpleBarStyle {
			fmt.Printf("\r%.2f%%", float64(w.LoadedBytes)*100.00/float64(w.TotalBytes))
		} else {
			fmt.Printf("\rDownloading... %s of %s complete", humanize.Bytes(w.LoadedBytes), humanize.Bytes(w.TotalBytes))
		}
	} else {
		fmt.Printf("\rDownloading... %s complete", humanize.Bytes(w.LoadedBytes))
	}
}
