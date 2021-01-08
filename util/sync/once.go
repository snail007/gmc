// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gsync

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"
)

var onceDataMap = sync.Map{}
var onceDoDataMap = sync.Map{}

func OnceDo(uniqueKey string, f func()) {
	once, _ := onceDoDataMap.LoadOrStore(uniqueKey, &sync.Once{})
	once.(*sync.Once).Do(f)
	return
}

func LoadOnce(uniqueKey string) *sync.Once {
	if uniqueKey == "" {
		key := make([]byte, 16)
		io.ReadFull(rand.Reader, key)
		uniqueKey = fmt.Sprintf("%x", key)
	}
	once, _ := onceDataMap.LoadOrStore(uniqueKey, &sync.Once{})
	return once.(*sync.Once)
}

func RemoveOnce(uniqueKey string) {
	onceDataMap.Delete(uniqueKey)
	return
}
