package pcfg

import (
	"fmt"
	"github.com/emirpasic/gods/trees/binaryheap"
	"sync"
)

type PrintJob struct {
	line  string
	count int
}

type Printer struct {
	count int
	heap  *binaryheap.Heap
	mutex sync.Mutex
}

func NewPrintHeap() *Printer {
	printJobComparator := func(a, b interface{}) int {
		aAsserted := a.(*PrintJob)
		bAsserted := b.(*PrintJob)

		switch {
		case aAsserted.count > bAsserted.count:
			return 1
		case aAsserted.count < bAsserted.count:
			return -1
		default:
			return 0
		}
	}

	return &Printer{
		heap:  binaryheap.NewWith(printJobComparator),
		mutex: sync.Mutex{},
	}
}

func (ph *Printer) Push(job *PrintJob) {
	ph.mutex.Lock()
	defer ph.mutex.Unlock()

	ph.heap.Push(job)
}

func (ph *Printer) Print() {
	ph.mutex.Lock()
	defer ph.mutex.Unlock()

	for {
		val, ok := ph.heap.Peek()

		if !ok {
			return
		}

		job, ok := val.(*PrintJob)

		if !ok {
			return
		}

		if job.count != ph.count {
			return
		}

		fmt.Println(job.line)

		ph.heap.Pop()
		ph.count++
	}
}
