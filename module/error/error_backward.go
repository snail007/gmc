// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

//go:build !go1.13
// +build !go1.13

package gerror

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error.
func Is(e error, original error) bool {
	if e == original {
		return true
	}

	if e, ok := e.(*Error); ok {
		return Is(e.Err, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}

	return false
}
