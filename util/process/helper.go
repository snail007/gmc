// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package process

import (
	gdaemon "github.com/snail007/gmc/util/process/daemon"
	ghook "github.com/snail007/gmc/util/process/hook"
)
// Daemon make your program ability of 1.daemonize 2.forever 3.logging stdout to file.
// If you used Daemon, you can run your program with arguments `--daemon` `--forever` `--flog foo.log` optionally.
func Daemon(mainFunc func(), waitAfterMain bool) (err error) {
	// daemon will os.Exit in Start()
	err = gdaemon.Start()
	if err != nil {
		return
	}
	// so daemon never will be reach here.
	if gdaemon.CanRun() {
		// no flags of forever, flog and daemon, will reach here.
		mainFunc()
		if waitAfterMain {
			ghook.WaitShutdown()
		}
	} else {
		// forever and flog will reach here.
		ghook.RegistShutdown(func() {
			gdaemon.Clean()
		})
		ghook.WaitShutdown()
	}
	return
}
