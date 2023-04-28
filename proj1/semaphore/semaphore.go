package semaphore

import "sync"

type Semaphore struct {
	mu    sync.Mutex
	cond  *sync.Cond
	value int
}

func NewSemaphore(capacity int) *Semaphore {
	sema := Semaphore{value: capacity}
	sema.cond = sync.NewCond(&sema.mu)
	return &sema
}
func (s *Semaphore) Up() {
	s.mu.Lock()
	s.value++
	s.cond.Signal()
	s.mu.Unlock()
}
func (s *Semaphore) Down() {
	s.mu.Lock()
	for s.value == 0 {
		s.cond.Wait()
	}
	s.value--
	s.mu.Unlock()
}
