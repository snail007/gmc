package ghttp

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// TriableRequest when request fail, retry to request with MaxTryCount.
type TriableRequest struct {
	req     *http.Request
	client  *http.Client
	reqBody []byte
	maxTry  int
	timeout time.Duration
	resp    *Response
	errs    []error
	doFunc  func(req *http.Request) (*http.Response, error)
	afterDo []AfterDoFunc
	body    []byte
}

// NewTriableRequest new a TriableRequest by *http.Request, maxTry is the max retry count when a request error occurred.
// If client is nil, default client will be used.
func NewTriableRequest(req *http.Request, client *http.Client, maxTry int, timeout time.Duration) *TriableRequest {
	if client == nil {
		if client == nil {
			client = defaultClient()
		}
	}
	tr := &TriableRequest{req: req, timeout: timeout, client: client, maxTry: maxTry}
	return tr
}

// NewTriableURL new a TriableRequest by URL.
func NewTriableURL(client *http.Client, method string, URL string, maxTry int, timeout time.Duration, data map[string]string, header map[string]string) (tr *TriableRequest, err error) {
	var req *http.Request
	var cancel context.CancelFunc
	req, cancel, err = NewRequest(method, URL, timeout, data, header)
	if err != nil {
		return
	}
	return NewTriableRequest(req, client, maxTry, timeout).
		AfterDo(func(resp *Response) {
			cancel()
		}), nil
}

// AfterDo add a callback call after request sent.
func (s *TriableRequest) AfterDo(afterDo AfterDoFunc) *TriableRequest {
	s.afterDo = append(s.afterDo, afterDo)
	return s
}

func (s *TriableRequest) callAfterDo(resp *Response) {
	for _, f := range s.afterDo {
		f(resp)
	}
}

// DoFunc sets a request sender.
func (s *TriableRequest) DoFunc(doFunc func(req *http.Request) (*http.Response, error)) *TriableRequest {
	s.doFunc = doFunc
	return s
}

// ErrAll returns all requests error.
func (s *TriableRequest) ErrAll() []error {
	return s.errs
}

func (s *TriableRequest) do(req *http.Request) (*http.Response, error) {
	if s.doFunc != nil {
		return s.doFunc(req)
	}
	return s.client.Do(req)
}

func (s *TriableRequest) init() *TriableRequest {
	s.resp = nil
	s.errs = nil
	s.reqBody = nil
	if len(s.reqBody) == 0 && s.req.Body != nil {
		s.reqBody, _ = ioutil.ReadAll(s.req.Body)
	}
	if s.client == nil && s.doFunc == nil {
		s.client = defaultClient()
	}
	return s
}

func (s *TriableRequest) forDo() *http.Request {
	req := withTimeout(s.req, s.timeout)
	if s.reqBody != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(s.reqBody))
	}
	return req
}

// Execute send request with retrying ability.
func (s *TriableRequest) Execute() *Response {
	s.init()
	var err error
	var resp *http.Response
	maxTry := s.maxTry
	tryCount := 0
	for tryCount <= maxTry {
		req := s.forDo()
		startTime := time.Now()
		resp, err = s.do(req)
		endTime := time.Now()
		s.resp = &Response{
			idx:       tryCount,
			req:       req,
			Response:  resp,
			respErr:   err,
			usedTime:  endTime.Sub(startTime),
			startTime: startTime,
			endTime:   endTime,
		}
		s.callAfterDo(s.resp)
		if err != nil {
			s.errs = append(s.errs, err)
			tryCount++
			continue
		}
		s.errs = nil
		break
	}
	return s.resp
}
