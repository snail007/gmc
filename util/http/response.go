// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
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
	if s.Response == nil {
		return nil, nil
	}
	if s.bodyErr != nil {
		return nil, s.bodyErr
	}
	if s.body != nil {
		return s.body, nil
	}
	if s.Response.Body != nil {
		s.body, s.bodyErr = ioutil.ReadAll(s.Response.Body)
		s.Response.Body.Close()
		return s.body, s.bodyErr
	}
	return nil, nil
}
