package loop

import (
	"github.com/snail007/gmc/util/gpool"
	"sync"
)

func For(count int, f func(idx int)) {
	for idx := 0; idx < count; idx++ {
		f(idx)
	}
}

type BatchExecutor struct {
	workers int
	tasks   []func()
}

func NewBatchExecutor() *BatchExecutor {
	return &BatchExecutor{}
}

func (s *BatchExecutor) SetWorkers(workersCnt int) {
	s.workers = workersCnt
}
func (s *BatchExecutor) AppendTask(tasks ...func()) {
	s.tasks = append(s.tasks, tasks...)
}

func (s *BatchExecutor) Exec() {
	workers := len(s.tasks)
	if s.workers > 0 {
		workers = s.workers
	}
	p := gpool.New(workers)
	defer p.Stop()
	g := sync.WaitGroup{}
	g.Add(len(s.tasks))
	for _, t := range s.tasks {
		task := t
		p.Submit(func() {
			defer g.Done()
			task()
		})
	}
	g.Wait()
}
