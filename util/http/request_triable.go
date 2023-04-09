package ghttp

import (
	"bytes"
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
	body    []byte
}

// NewTriableRequest new a TriableRequest, maxTry is the max retry count when a request error occurred.
// If client is nil, default client will be used.
func NewTriableRequest(req *http.Request, client *http.Client, maxTry int, timeout time.Duration) *TriableRequest {
	if client == nil {
		if client == nil {
			client = defaultClient()
		}
	}
	tr := &TriableRequest{req: req, timeout: timeout, client: client, maxTry: maxTry}
	if req.Body != nil {
		tr.reqBody, _ = ioutil.ReadAll(req.Body)
	}
	return tr
}

func (s *TriableRequest) DoFunc(doFunc func(req *http.Request) (*http.Response, error)) *TriableRequest {
	s.doFunc = doFunc
	return s
}

func (s *TriableRequest) Success() bool {
	return s.resp != nil
}

func (s *TriableRequest) Response() *Response {
	return s.resp
}

func (s *TriableRequest) Err() error {
	for _, v := range s.errs {
		return v
	}
	return nil
}

func (s *TriableRequest) Errs() []error {
	return s.errs
}

func (s *TriableRequest) do(req *http.Request) (*http.Response, error) {
	if s.doFunc != nil {
		return s.doFunc(req)
	}
	return s.client.Do(req)
}

// Execute send request with retrying ability.
func (s *TriableRequest) Execute() *TriableRequest {
	var err error
	var resp *http.Response
	for s.maxTry >= 0 {
		req := withTimeout(s.req, s.timeout)
		if s.reqBody != nil {
			req.Body = ioutil.NopCloser(bytes.NewReader(s.reqBody))
		}
		resp, err = s.do(req)
		if err != nil {
			s.errs = append(s.errs, err)
			s.maxTry--
			continue
		}
		s.resp = NewResponse(resp)
		s.errs = nil
		break
	}
	return s
}
