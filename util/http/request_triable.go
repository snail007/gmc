package ghttp

import (
	"net/http"
	"time"
)

// TriableRequest when request fail, retry to request with MaxTryCount.
type TriableRequest struct {
	req    *http.Request
	client *http.Client
	maxTry int
	resp   *http.Response
	errs   []error
	doFunc func(req *http.Request) (*http.Response, error)
}

// NewTriableRequest new a TriableRequest, maxTry is the max retry count when a request error occurred.
// If client is nil, default client will be used.
func NewTriableRequest(req *http.Request, client *http.Client, maxTry int, timeout time.Duration) *TriableRequest {
	if client == nil {
		if client == nil {
			client = defaultClient()
		}
	}
	return &TriableRequest{req: withTimeout(req, timeout), client: client, maxTry: maxTry}
}

func (s *TriableRequest) DoFunc(doFunc func(req *http.Request) (*http.Response, error)) *TriableRequest {
	s.doFunc = doFunc
	return s
}

func (s *TriableRequest) Success() bool {
	return s.resp != nil
}

func (s *TriableRequest) Response() *http.Response {
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
	for s.maxTry >= 0 {
		s.resp, err = s.do(s.req)
		if err != nil {
			s.errs = append(s.errs, err)
			s.maxTry--
			continue
		}
		s.errs = nil
		break
	}
	return s
}
