package snow

import (
	"container/heap"
	"sync"
)

// An pqItem is something we manage in a priority queue.
type pqItem struct {
	value    interface{} // The value of the item; arbitrary.
	priority int64       // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A myPriorityQueue implements heap.Interface and holds Items.
type myPriorityQueue []*pqItem

type PriorityQueue struct {
	list myPriorityQueue
	lock sync.Mutex
}

func (pq myPriorityQueue) Len() int { return len(pq) }

func (pq myPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq myPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *myPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *myPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// 优先队列
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

func (pq *PriorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return pq.list.Len()
}

func (pq *PriorityQueue) Push(value interface{}, priority int64) {
	item := &pqItem{
		value:    value,
		priority: priority,
		index:    0,
	}

	pq.lock.Lock()
	defer pq.lock.Unlock()

	heap.Push(&pq.list, item)
	heap.Fix(&pq.list, item.index)
}

func (pq *PriorityQueue) Pop() (value interface{}, priority int64, ok bool) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.list.Len() > 0 {
		item := heap.Pop(&pq.list).(*pqItem)
		return item.value, item.priority, true
	}
	return nil, 0, false
}

func (pq *PriorityQueue) Peek() (value interface{}, priority int64, ok bool) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.list.Len() > 0 {
		item := heap.Pop(&pq.list).(*pqItem)
		heap.Push(&pq.list, item)
		heap.Fix(&pq.list, item.index)
		return item.value, item.priority, true
	}
	return
}
