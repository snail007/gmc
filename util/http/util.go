package ghttp

import (
	"bytes"
	"context"
	gurl "github.com/snail007/gmc/util/url"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type BeforeDoFunc func(idx int, req *http.Request)
type AfterDoFunc func(resp *Response)

func NewRequest(method, URL string, timeout time.Duration, data, header map[string]string) (req *http.Request, cancel context.CancelFunc, err error) {
	ctx, cancel := getTimeoutContext(timeout)
	req, err = NewRequestWithContext(ctx, method, URL, data, header)
	return
}

func NewRequestWithContext(ctx context.Context, method, URL string, data, header map[string]string) (req *http.Request, err error) {
	if IsFormMethod(method) {
		return NewPostWithContext(ctx, URL, data, header)
	}
	return NewGetWithContext(ctx, URL, data, header)
}

func NewGet(URL string, timeout time.Duration, queryData, header map[string]string) (req *http.Request, cancel context.CancelFunc, err error) {
	ctx, cancel := getTimeoutContext(timeout)
	req, err = NewGetWithContext(ctx, URL, queryData, header)
	return
}

func NewGetWithContext(ctx context.Context, URL string, queryData, header map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, "GET", gurl.AppendQuery(URL, queryData), nil)
	if err != nil {
		return
	}
	setHeader(req, header)
	return
}

func NewPost(URL string, timeout time.Duration, data, header map[string]string) (req *http.Request, cancel context.CancelFunc, err error) {
	ctx, cancel := getTimeoutContext(timeout)
	req, err = NewPostWithContext(ctx, URL, data, header)
	return
}

func NewPostWithContext(ctx context.Context, URL string, data, header map[string]string) (req *http.Request, err error) {
	return NewPostReaderWithContext(ctx, URL, bytes.NewReader([]byte(gurl.EncodeData(data))), header)
}

func NewPostReaderWithContext(ctx context.Context, URL string, r io.Reader, header map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, "POST", URL, r)
	if err != nil {
		return
	}
	setHeader(req, header)
	return
}

func getTimeoutContext(timeout time.Duration) (ctx context.Context, cancel context.CancelFunc) {
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx = context.Background()
		cancel = func() {}
	}
	return
}

func IsFormRequest(req *http.Request) bool {
	return IsFormMethod(req.Method)
}

func IsFormMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

func setHeader(req *http.Request, header map[string]string) {
	defer func() {
		if IsFormRequest(req) && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}()
	if header == nil {
		return
	}
	for k, v := range header {
		if strings.EqualFold(k, "host") {
			req.Host = v
			continue
		}
		req.Header.Set(k, v)
	}
}

func CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
func GetResponseBodyE(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, nil
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return b, nil
}
func GetResponseBody(resp *http.Response) []byte {
	b, _ := GetResponseBodyE(resp)
	return b
}

func GetContent(_url string, timeout time.Duration) []byte {
	d, _ := GetContentE(_url, timeout)
	return d
}

func GetContentE(_url string, timeout time.Duration) ([]byte, error) {
	d, _, e := Download(_url, timeout, nil, nil)
	return d, e
}
