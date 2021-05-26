// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttpserver_test

import (
	_ "github.com/snail007/gmc"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type APITestSuite struct {
	suite.Suite
	apiAddr string
	assert  *assert2.Assertions
}

func (s *APITestSuite) SetupTest() {
	s.assert = s.Assertions
	var err error
	s.apiAddr, err = apiServer()
	s.assert.Nil(err)
}

func (s *APITestSuite) TestAPI() {

}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
