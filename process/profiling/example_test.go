// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package profiling

func ExampleStartArg() {
	//StartArg will search command line arguments "profiling" , using it's value as store data folder.
	StartArg("profiling")
	defer Stop()
}

func ExampleStart() {
	//Start will using the argument as store data folder
	Start("debug")
	defer Stop()
}
