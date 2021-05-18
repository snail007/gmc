// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtest

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func newCmdFromEnv(runName string) string {
	if strings.Contains(runName, "/") {
		rs := strings.Split(runName, "/")
		runName = "^" + strings.Join(rs, "$/^") + "$"
	}
	rb := make([]byte, 16)
	io.ReadFull(rand.Reader, rb)
	d, _ := os.Getwd()
	coverfile := filepath.Join(d, fmt.Sprintf("%x", rb)) + ".gocc.tmp"
	race := os.Getenv("GMCT_COVER_RACE")
	packages := os.Getenv("GMCT_COVER_PACKAGES")
	pkg := strings.TrimPrefix(d, filepath.Join(os.Getenv("GOPATH"), "src"))[1:]
	if race == "true" {
		race = "-v"
	} else {
		race = ""
	}
	cover := ""
	if packages != "" {
		cover = fmt.Sprintf("-covermode=atomic -coverprofile=%s -coverpkg=%s", coverfile, packages)
	}
	return fmt.Sprintf(`go test -v -run=%s %s %s %s`,
		runName, cover, race, pkg)
}

func ExecTestFunc(testFuncName string) (out string, err error) {
	isVerbose := os.Getenv("GMCT_COVER_VERBOSE") == "true"
	defer func() {
		if isVerbose {
			fmt.Printf("output: %s", out)
			fmt.Printf(">>> end child testing process %s\n", testFuncName)
		}
	}()
	if isVerbose {
		fmt.Printf(">>> start child testing process %s\n", testFuncName)
	}
	cmdStr := newCmdFromEnv(testFuncName)
	if isVerbose {
		fmt.Println(cmdStr)
	}
	c := exec.Command("bash", "-c", cmdStr)
	c.Env = append(c.Env, os.Environ()...)
	b, err := c.CombinedOutput()
	out = string(b)
	return
}
