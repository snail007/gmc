package gsession

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	cfg := NewRedisStoreConfig()
	cfg.TTL = 1
	store, err := NewRedisStore(cfg)
	assert.Nil(err)
	sess := NewSession()
	sid := sess.SessionID()
	sess.Set("test", "aaa")
	sess.Touch()
	store.Save(sess)
	s0, ok := store.Load(sid)
	assert.True(ok)
	assert.Equal(s0.SessionID(), sid)
	assert.Equal(s0.Get("test"), "aaa")
	time.Sleep(time.Second)
	_, ok = store.Load(sid)
	assert.False(ok)
	store.Delete(sid)
}
