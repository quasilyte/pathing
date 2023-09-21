package pathing

import (
	"math"
)

type coordMap struct {
	elems   []coordMapElem
	gen     uint32
	numRows int
	numCols int
}

type coordMapElem struct {
	value uint32
	gen   uint32
}

func newCoordMap(numCols, numRows int) *coordMap {
	size := numRows * numCols
	return &coordMap{
		elems:   make([]coordMapElem, size),
		gen:     1,
		numRows: numRows,
		numCols: numCols,
	}
}

func (m *coordMap) Contains(k uint) bool {
	if k < uint(len(m.elems)) {
		return m.elems[k].gen == m.gen
	}
	return false
}

func (m *coordMap) Get(k uint) (uint32, bool) {
	if k < uint(len(m.elems)) {
		el := m.elems[k]
		if el.gen == m.gen {
			return el.value, true
		}
	}
	return 0, false
}

func (m *coordMap) Set(k uint, v uint32) {
	if k < uint(len(m.elems)) {
		m.elems[k] = coordMapElem{value: v, gen: m.gen}
	}
}

func (m *coordMap) Reset() {
	if m.gen == math.MaxUint32 {
		// For most users, this will never happen.
		// But to be safe, we need to handle this correctly.
		// m.gen will be 1, element gen will be 0.
		m.clear()
	} else {
		m.gen++
	}
}

// clear does a real array data reset.
// m.gen becomes 1.
// Every element gen becomes 0.
// This is identical to the initial array state.
//
//go:noinline - called on a cold path, therefore it should not be inlined.
func (m *coordMap) clear() {
	m.gen = 1

	// TODO: could use clear() starting from Go 1.21.
	for i := range m.elems {
		m.elems[i] = coordMapElem{}
	}
}

func (s *coordMap) packCoord(c GridCoord) uint {
	return uint((c.Y * s.numCols) + c.X)
}
