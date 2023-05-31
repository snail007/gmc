package gbatch

import (
	"github.com/snail007/gmc/util/gpool"
	"sync"
)

type Executor struct {
	workers int
	tasks   []func()
}

func NewBatchExecutor() *Executor {
	return &Executor{}
}

func (s *Executor) SetWorkers(workersCnt int) {
	s.workers = workersCnt
}
func (s *Executor) AppendTask(tasks ...func()) {
	s.tasks = append(s.tasks, tasks...)
}

func (s *Executor) Exec() {
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
