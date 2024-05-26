package songphapadonationconsumer

type DonatorHeap []*Donator

func (h DonatorHeap) Len() int           { return len(h) }
func (h DonatorHeap) Less(i, j int) bool { return h[i].Amount > h[j].Amount }
func (h DonatorHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DonatorHeap) Push(x any) {
	*h = append(*h, x.(*Donator))
	if h.Len() > x.(*Donator).topN {
		h.Pop()
	}
}

func (h *DonatorHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}