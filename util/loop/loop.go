package gloop

type StepHandler func(loopIndex, loopValue int)
type Handler func(loopIndex int)
type Condition func(loopIndex int) bool

func For(count int, f Handler) {
	ForBy(count, 1, func(_, loopValue int) {
		f(loopValue)
	})
}

func ForBy(max, step int, f StepHandler) {
	if step > 0 {
		idx := 0
		for value := 0; value < max; value += step {
			f(idx, value)
			idx++
		}
	} else {
		idx := 0
		for value := max; value >= 0; value += step {
			f(idx, value)
			idx++
		}
	}
}

type doWhileUntil struct {
	do Handler
}

func Do(f Handler) *doWhileUntil {
	return &doWhileUntil{do: f}
}

func (s *doWhileUntil) While(f Condition) {
	s.do(0)
	idx := 1
	for f(idx) {
		s.do(idx)
		idx++
	}
}

func (s *doWhileUntil) Until(f Condition) {
	s.do(0)
	idx := 1
	for !f(idx) {
		s.do(idx)
		idx++
	}
}
