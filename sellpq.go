package main

type SellPriorityQueue []*order

func (spq SellPriorityQueue) Len() int { return len(spq) }

func (spq SellPriorityQueue) Less(i, j int) bool {
	if spq[i].inp.price == spq[j].inp.price {
		return spq[i].timestamp < spq[j].timestamp
	}
	// The sell order we want first is the lower price one
	return spq[i].inp.price < spq[j].inp.price
}

func (spq SellPriorityQueue) Swap(i, j int) {
	spq[i], spq[j] = spq[j], spq[i]
	spq[i].index = i
	spq[j].index = j
}

func (spq *SellPriorityQueue) Push(x any) {
	n := len(*spq)
	item := x.(*order)
	item.index = n
	*spq = append(*spq, item)
}

func (spq *SellPriorityQueue) Pop() any {
	old := *spq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*spq = old[0 : n-1]
	return item
}
