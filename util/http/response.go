// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	idx int
	req *http.Request
	*http.Response
	body      []byte
	bodyErr   error
	respErr   error
	usedTime  time.Duration
	startTime time.Time
	endTime   time.Time
	cancel    context.CancelFunc
}

func (s *Response) Request() *http.Request {
	return s.req
}

func (s *Response) UsedTime() time.Duration {
	return s.usedTime
}

func (s *Response) StartTime() time.Time {
	return s.startTime
}

func (s *Response) EndTime() time.Time {
	return s.endTime
}

func (s *Response) Idx() int {
	return s.idx
}

func (s *Response) Err() error {
	return s.respErr
}

func (s *Response) Body() []byte {
	b, _ := s.BodyE()
	return b
}

func (s *Response) BodyE() ([]byte, error) {
	if s.Response == nil || s.Response.Body == nil || s.respErr != nil {
		return nil, nil
	}
	if s.bodyErr != nil {
		return nil, s.bodyErr
	}
	if s.body != nil {
		return s.body, nil
	}
	s.body, s.bodyErr = ioutil.ReadAll(s.Response.Body)
	s.Response.Body.Close()
	return s.body, s.bodyErr
}

func (s *Response) Close() {
	if s.Response != nil && s.Response.Body != nil {
		s.Response.Body.Close()
	}
	if s.cancel != nil {
		s.cancel()
	}
}
