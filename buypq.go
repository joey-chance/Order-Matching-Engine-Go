package main

import (
	"container/heap"
	"fmt"
	"os"
)

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
	fmt.Fprintf(os.Stderr, "Push BO Count: %v\n", x.(*order).inp.count)
	n := len(*bpq)
	item := x.(*order)
	item.index = n
	*bpq = append(*bpq, item)
	fmt.Fprintf(os.Stderr, "Len: %v\n", len(*bpq))
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

func (bpq *BuyPriorityQueue) Peek() any {
	old := *bpq
	n := len(old)
	item := old[n-1]
	return item
}

func (bpq *BuyPriorityQueue) RemoveOrderId(orderId uint32) any {
	//Find orderId

	removedIdx := 0
	var removedOrderPtr *order = nil
	old := *bpq
	n := len(old)

	if n <= 0 {
		return nil
	}

	for idx, order := range old {
		if order.inp.orderId == orderId {
			removedIdx = idx
			removedOrderPtr = order
		}
	}

	if removedOrderPtr == nil {
		return nil
	}

	//Swap this order with the last (n-1)
	fmt.Fprintf(os.Stderr, "Length: %v\n", n)

	item := old[n-1]
	// fmt.Println("HERE?")
	old[n-1] = removedOrderPtr
	old[removedIdx] = item
	// fmt.Println("HERE?")
	//Get rid of last order (which is the one we want to remove)
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety

	//Reslice
	*bpq = old[0 : n-1]
	//Fix heap invariant
	heap.Init(bpq)
	return item

}
