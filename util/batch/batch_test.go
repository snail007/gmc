package gbatch

import (
	gloop "github.com/snail007/gmc/util/loop"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBatchExecutor(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(10)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func() {
			i.Increase(idx)
			time.Sleep(time.Second)
		})
	})
	be.Exec()
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	assert.True(t, diff >= time.Second)
	assert.True(t, diff < time.Millisecond*1500)
}

func TestNewBatchExecutor2(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(5)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func() {
			i.Increase(idx)
			time.Sleep(time.Second)
		})
	})
	be.Exec()
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	assert.True(t, diff >= time.Second*2)
	assert.True(t, diff < time.Millisecond*2500)
}
