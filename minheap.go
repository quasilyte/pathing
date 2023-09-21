package pathing

type minheap[T any] struct {
	elems []minheapElem[T]
}

type minheapElem[T any] struct {
	Value    T
	Priority int
}

func newMinheap[T any](size int) *minheap[T] {
	h := &minheap[T]{
		elems: make([]minheapElem[T], 0, size),
	}
	return h
}

func (h *minheap[T]) IsEmpty() bool {
	return len(h.elems) == 0
}

func (h *minheap[T]) Reset() {
	h.elems = h.elems[:0]
}

func (q *minheap[T]) Push(priority int, value T) {
	q.elems = append(q.elems, minheapElem[T]{
		Priority: priority,
		Value:    value,
	})

	elems := q.elems
	i := uint(len(elems) - 1)
	for {
		j := (i - 1) / 2
		if i <= j || i >= uint(len(elems)) || elems[i].Priority >= elems[j].Priority {
			break
		}
		elems[i], elems[j] = elems[j], elems[i]
		i = j
	}
}

func (q *minheap[T]) Pop() T {
	if len(q.elems) == 0 {
		var zero T
		return zero
	}
	elems := q.elems
	size := len(elems) - 1
	elems[0], elems[size] = elems[size], elems[0]

	// down(0)
	i := 0
	for {
		j := 2*i + 1
		if j >= size {
			break
		}
		if j2 := j + 1; j2 < size && elems[j2].Priority < elems[j].Priority {
			j = j2
		}
		if elems[i].Priority < elems[j].Priority {
			break
		}
		elems[i], elems[j] = elems[j], elems[i]
		i = j
	}

	value := elems[size].Value
	q.elems = elems[:size]
	return value
}
