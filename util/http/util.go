package ghttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewRequest(method, URL string, timeout time.Duration, data, header map[string]string) (req *http.Request, cancel context.CancelFunc, err error) {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return NewPost(URL, timeout, data, header)
	default:
		return NewGet(URL, timeout, data, header)
	}
}

func NewRequestWithContext(ctx context.Context, method, URL string, timeout time.Duration, data, header map[string]string) (req *http.Request, err error) {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return NewPostWithContext(ctx, URL, data, header)
	default:
		return NewGetWithContext(ctx, URL, data, header)
	}
}

func NewGet(URL string, timeout time.Duration, queryData, header map[string]string) (req *http.Request, cancel context.CancelFunc, err error) {
	ctx, cancel := getTimeoutContext(timeout)
	req, err = NewGetWithContext(ctx, URL, queryData, header)
	return
}

func NewGetWithContext(ctx context.Context, URL string, queryData, header map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, "GET", appendQuery(URL, queryData), nil)
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
	return NewPostReaderWithContext(ctx, URL, bytes.NewReader([]byte(encodeData(data))), header)
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

func appendQuery(URL string, queryData map[string]string) string {
	if len(queryData) == 0 {
		return URL
	}
	return URL + getConcatChar(URL) + encodeData(queryData)
}

func getConcatChar(URL string) string {
	if strings.Contains(URL, "?") {
		return "&"
	}
	return "?"
}

func encodeData(data map[string]string) string {
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values.Set(k, v)
		}
	}
	return values.Encode()
}
func withTimeout(req *http.Request, timeout time.Duration) *http.Request {
	if timeout == 0 {
		return req
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req = req.WithContext(ctx)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return req
}

func setHeader(req *http.Request, header map[string]string) {
	if header == nil {
		return
	}
	foundContentType := false
	if header != nil {
		for k, v := range header {
			if strings.EqualFold(k, "host") {
				req.Host = v
				continue
			}
			req.Header.Set(k, v)
			if strings.TrimSpace(strings.ToLower(k)) == "content-type" {
				foundContentType = true
			}
		}
	}
	if !foundContentType && (req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

func CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
