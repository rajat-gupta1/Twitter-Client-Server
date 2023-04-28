package queue

import (
	"sync/atomic"
	"unsafe"
)

type Request struct {
	Command   string  `json:"command"`
	Id        float64 `json:"id"`
	Body      string  `json:"body"`
	TimeStamp float64 `json:"timestamp"`
}

type Node struct {
	value Request
	next  unsafe.Pointer
}

// LockfreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// NewQueue creates and initializes a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
	node := unsafe.Pointer(&Node{})
	return &LockFreeQueue{head: node, tail: node}
}

// Enqueue adds a series of Request to the queue
func (queue *LockFreeQueue) Enqueue(task *Request) {
	node := &Node{value: *task}
	for {
		tail := load(&queue.tail)
		next := load(&tail.next)
		if tail == load(&queue.tail) {
			if next == nil {
				if cas(&tail.next, next, node) {
					cas(&queue.tail, tail, node)
					return
				}
			} else {
				cas(&queue.tail, tail, next)
			}
		}
	}
}

// Dequeue removes a Request from the queue
func (queue *LockFreeQueue) Dequeue() *Request {
	for {
		head := load(&queue.head)
		tail := load(&queue.tail)
		next := load(&head.next)
		if head == load(&queue.head) {
			if head == tail {
				if next == nil {
					return nil
				}
				cas(&queue.tail, tail, next)
			} else {
				val := next.value
				if cas(&queue.head, head, next) {
					return &val
				}
			}
		}
	}
}

func load(p *unsafe.Pointer) (node *Node) {
	return (*Node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *Node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
