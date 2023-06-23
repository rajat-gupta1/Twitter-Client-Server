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

// Function to do write lock the queue
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

// Function to do write unlock the queue
func (rw *RWMutex) Unlock() {
	rw.mu.Lock()
	rw.WriterActive = false
	rw.cond.Broadcast()
	rw.mu.Unlock()
}

// Function to do Read Lock
func (rw *RWMutex) RLock() {
	rw.mu.Lock()
	for rw.NumWritersWaiting > 0 || rw.WriterActive || rw.NumReadersActive >= 32 {
		rw.cond.Wait()
	}
	rw.NumReadersActive += 1
	rw.mu.Unlock()
}

// Function to do Read Unlock
func (rw *RWMutex) RUnlock() {
	rw.mu.Lock()
	rw.NumReadersActive -= 1
	if rw.NumReadersActive == 0 {
		rw.cond.Broadcast()
	}
	rw.mu.Unlock()
}