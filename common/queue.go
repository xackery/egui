package common

// Item represents a node item
type Item struct {
	node     *Node
	priority float64
	index    int
}

// PriorityQueue is an array of items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	if len(pq)-1 < i || len(pq)-1 < j {
		return
	}
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push add an item to a queue
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an item from a queue
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	var item *Item
	if len(old) == 0 {
		return nil
	}
	item = old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]

	return item
}
