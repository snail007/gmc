// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"fmt"
)

func ExampleNew() {
	//we create a poll named "p" with 3 workers
	p := NewGPool(3)
	//after New, you can submit a function as a task, you can repeat Submit() many times anywhere as you need.
	a := make(chan bool)
	p.Submit(func() {
		a <- true
	})
	fmt.Println(<-a)
}
