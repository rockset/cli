package sort

import "sort"

// Multi is used to sort a slice of T using the list of LessFuncs, so you can order a table by different columns
type Multi[T any] struct {
	changes   []T
	LessFuncs []func(p1, p2 *T) bool
}

func (ms *Multi[T]) Sort(changes []T) {
	ms.changes = changes
	sort.Sort(ms)
}

func (ms *Multi[T]) Len() int {
	return len(ms.changes)
}

func (ms *Multi[T]) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *Multi[T]) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.LessFuncs)-1; k++ {
		less := ms.LessFuncs[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.LessFuncs[k](p, q)
}
