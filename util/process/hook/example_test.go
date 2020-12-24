// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package ghook

func Example() {
	//the code block should be last in your main function code

	RegistShutdown(func() {
		//do something before program exit
	})
	WaitShutdown()
}
