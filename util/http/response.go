// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	body    []byte
	bodyErr error
}

func NewResponse(response *http.Response) *Response {
	if response == nil {
		return nil
	}
	return &Response{Response: response}
}

func (s *Response) Body() []byte {
	b, _ := s.BodyE()
	return b
}

func (s *Response) BodyE() ([]byte, error) {
	if s.bodyErr != nil {
		return nil, s.bodyErr
	}
	if s.body != nil {
		return s.body, nil
	}
	if s.Response.Body != nil {
		var err error
		s.body, err = ioutil.ReadAll(s.Response.Body)
		s.Response.Body.Close()
		return s.body, err
	}
	return nil, nil
}
