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
	req            *http.Request
	client         *http.Client
	reqBody        []byte
	maxTry         int
	timeout        time.Duration
	resp           *Response
	errs           []error
	doFunc         func(req *http.Request) (*http.Response, error)
	beforeDo       []BeforeDoFunc
	afterDo        []AfterDoFunc
	checkErrorFunc func(int, *http.Request, *http.Response) error
	body           []byte
	keepalive      bool
}

// NewTriableRequest new a TriableRequest by *http.Request, maxTry is the max retry count when a request error occurred.
// If client is nil, default client will be used.
func NewTriableRequest(client *http.Client, req *http.Request, maxTry int, timeout time.Duration) *TriableRequest {
	tr := &TriableRequest{req: req, timeout: timeout, client: client, maxTry: maxTry, keepalive: true}
	return tr
}

// NewTriableRequestByURL new a TriableRequest by URL.
func NewTriableRequestByURL(client *http.Client, method string, URL string, maxTry int,
	timeout time.Duration, data map[string]string, header map[string]string) (tr *TriableRequest, err error) {
	var req *http.Request
	var cancel context.CancelFunc
	req, cancel, err = NewRequest(method, URL, timeout, data, header)
	if err != nil {
		return
	}
	return NewTriableRequest(client, req, maxTry, timeout).
		AfterDo(func(resp *Response) {
			cancel()
		}), nil
}

func (s *TriableRequest) init() *TriableRequest {
	s.resp = nil
	s.errs = nil
	s.reqBody = nil
	if len(s.reqBody) == 0 && s.req.Body != nil {
		s.reqBody, _ = ioutil.ReadAll(s.req.Body)
		s.req.Body = ioutil.NopCloser(bytes.NewReader(s.reqBody))
	}
	if s.client == nil && s.doFunc == nil {
		s.client = defaultClient()
	}
	return s
}

// Keepalive sets enable or disable for request keepalive
func (s *TriableRequest) Keepalive(keepalive bool) *TriableRequest {
	s.keepalive = keepalive
	return s
}

// CheckErrorFunc if returns an error, the request treat as fail.
func (s *TriableRequest) CheckErrorFunc(checkErrorFunc func(idx int, req *http.Request, resp *http.Response) error) *TriableRequest {
	s.checkErrorFunc = checkErrorFunc
	return s
}

// SetBeforeDo sets callback call before request sent.
func (s *TriableRequest) SetBeforeDo(beforeDo BeforeDoFunc) *TriableRequest {
	return s.setBeforeDo(beforeDo, true)
}

// AppendBeforeDo add a callback call before request sent.
func (s *TriableRequest) AppendBeforeDo(beforeDo BeforeDoFunc) *TriableRequest {
	return s.setBeforeDo(beforeDo, false)
}

func (s *TriableRequest) setBeforeDo(beforeDo BeforeDoFunc, isSet bool) *TriableRequest {
	if isSet {
		if beforeDo != nil {
			s.beforeDo = []BeforeDoFunc{beforeDo}
		} else {
			s.beforeDo = []BeforeDoFunc{}
		}
	} else if beforeDo != nil {
		s.beforeDo = append(s.beforeDo, beforeDo)
	}
	return s
}

// AfterDo add a callback call after request sent.
func (s *TriableRequest) AfterDo(afterDo AfterDoFunc) *TriableRequest {
	s.afterDo = append(s.afterDo, afterDo)
	return s
}

func (s *TriableRequest) callBeforeDo(idx int, req *http.Request) {
	for _, f := range s.beforeDo {
		f(idx, req)
	}
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

// Err returns first error of fail requests.
func (s *TriableRequest) Err() error {
	if len(s.errs) > 0 {
		return s.errs[0]
	}
	return nil
}

// Close close all response body can context cancel  func.
func (s *TriableRequest) Close() *TriableRequest {
	if s.resp != nil {
		s.resp.Close()
	}
	return s
}

func (s *TriableRequest) do(tryCount int, req *http.Request) (*http.Response, error) {
	req.Close = !s.keepalive
	var resp *http.Response
	var err error
	if s.doFunc != nil {
		resp, err = s.doFunc(req)
	} else {
		resp, err = s.client.Do(req)
	}
	if err != nil {
		return nil, err
	}
	if s.checkErrorFunc != nil {
		err = s.checkErrorFunc(tryCount, req, resp)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
	}
	return resp, nil
}

func (s *TriableRequest) forDo() (req *http.Request, cancel context.CancelFunc) {
	if s.timeout == 0 {
		return s.req, func() {}
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	req = s.req.WithContext(ctx)
	if s.reqBody != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(s.reqBody))
	}
	return req, cancel
}

// Execute send request with retrying ability.
func (s *TriableRequest) Execute() *Response {
	s.init()
	var err error
	var resp *http.Response
	maxTry := s.maxTry
	tryCount := 0
	for tryCount <= maxTry {
		req, cancel := s.forDo()
		s.callBeforeDo(tryCount, req)
		startTime := time.Now()
		resp, err = s.do(tryCount, req)
		endTime := time.Now()
		if err != nil {
			cancel()
		}
		s.resp = &Response{
			idx:       tryCount,
			req:       req,
			Response:  resp,
			respErr:   err,
			usedTime:  endTime.Sub(startTime),
			startTime: startTime,
			endTime:   endTime,
			cancel:    cancel,
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
