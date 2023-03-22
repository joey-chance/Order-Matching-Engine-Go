package main

type BuyPriorityQueue []*order

func (bpq BuyPriorityQueue) Len() int { return len(bpq) }

func (bpq BuyPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	if bpq[i].inp.price == bpq[j].inp.price {
		return bpq[i].timestamp < bpq[j].timestamp
	}
	// The buy order we want first is the higher price one
	return bpq[i].inp.price > bpq[j].inp.price
}

func (bpq BuyPriorityQueue) Swap(i, j int) {
	bpq[i], bpq[j] = bpq[j], bpq[i]
	bpq[i].index = i
	bpq[j].index = j
}

func (bpq *BuyPriorityQueue) Push(x any) {
	n := len(*bpq)
	item := x.(*order)
	item.index = n
	*bpq = append(*bpq, item)
}

func (bpq *BuyPriorityQueue) Pop() any {
	old := *bpq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*bpq = old[0 : n-1]
	return item
}
