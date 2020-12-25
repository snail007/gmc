package gsession

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"os"
	"testing"
)

var (
	fileStore gcore.SessionStorage
	memoryStore gcore.SessionStorage
	redisStore gcore.SessionStorage
)
func TestMain(m *testing.M) {
	var err error

	cfg := NewFileStoreConfig()
	cfg.GCtime = 1
	cfg.TTL = 1
	fileStore, err = NewFileStore(cfg)
	if err != nil {
		fmt.Println(err)
	}

	cfg0 := NewMemoryStoreConfig()
	cfg0.GCtime = 1
	cfg0.TTL = 1
	memoryStore, err = NewMemoryStore(cfg0)

	os.Exit(m.Run())
}
