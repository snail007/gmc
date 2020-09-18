// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package session

type Store interface {
	Load(sessionID string) (session *Session, isExits bool)
	Save(session *Session) (err error)
	Delete(sessionID string) (err error)
}
