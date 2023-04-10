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
	return tr
}

func (s *TriableRequest) DoFunc(doFunc func(req *http.Request) (*http.Response, error)) *TriableRequest {
	s.doFunc = doFunc
	return s
}

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
