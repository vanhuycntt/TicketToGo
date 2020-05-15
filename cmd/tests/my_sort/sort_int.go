package my_sort

type SortedInt struct {
	items  []int
	byFunc func(i, j int) bool
}

func (si SortedInt) Less(i, j int) bool {
	return si.byFunc(si.items[i], si.items[j])
}

func (si SortedInt) Len() int {
	return len(si.items)
}

func (si SortedInt) Swap(i, j int) {
	si.items[i], si.items[j] = si.items[j], si.items[i]
}

func NewSortedIn(items []int, lessFunc func(i, j int) bool) *SortedInt {
	return &SortedInt{
		items:  items,
		byFunc: lessFunc,
	}
}
