// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

/*
gpool is a goroutine pool, it is concurrency safed, it can using a few goroutine to run a huge tasks.

a worker is a goroutine, a task is a function, gpool using a few workers to run a huage tasks.
*/
package gpool
