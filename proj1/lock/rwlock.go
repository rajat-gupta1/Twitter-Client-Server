// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import (
	"sync"
)

type RWMutex struct {
	cond                *sync.Cond
	mu                  sync.Mutex
	NumWritersWaiting int
	WriterActive       bool
	NumReadersActive  int
}

func NewRWMutex() *RWMutex {
	newRWMutex := RWMutex{};
	newRWMutex.cond = sync.NewCond(&newRWMutex.mu)
	newRWMutex.NumReadersActive = 0
	newRWMutex.NumWritersWaiting = 0
	newRWMutex.WriterActive = false
	return &newRWMutex
}

func (rw *RWMutex) Lock() {
	rw.mu.Lock()
	rw.NumWritersWaiting += 1
	for rw.NumReadersActive > 0 || rw.WriterActive {
		rw.cond.Wait()
	}
	rw.NumWritersWaiting -= 1
	rw.WriterActive = true
	rw.mu.Unlock()
}

func (rw *RWMutex) Unlock() {
	rw.mu.Lock()
	rw.WriterActive = false
	rw.cond.Broadcast()
	rw.mu.Unlock()
}

func (rw *RWMutex) RLock() {
	rw.mu.Lock()
	for rw.NumWritersWaiting > 0 || rw.WriterActive || rw.NumReadersActive >= 32 {
		rw.cond.Wait()
	}
	rw.NumReadersActive += 1
	rw.mu.Unlock()
}

func (rw *RWMutex) RUnlock() {
	rw.mu.Lock()
	rw.NumReadersActive -= 1
	if rw.NumReadersActive == 0 {
		rw.cond.Broadcast()
	}
	rw.mu.Unlock()
}