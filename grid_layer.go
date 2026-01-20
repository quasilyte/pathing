package pathing

import (
	"unsafe"
)

// GridLayer is a tile-to-cost mapper.
// Every Grid cell has a tile tag value ranging from 0 to 7 (3 bits).
// Layers are used to turn that tag value into an actual traversal cost.
//
// Although you can construct a GridLayer value yourself,
// it's much easier to use a MakeGridLayer() function.
//
// A byte value of 0 means "the cell can't be traversed".
// A value higher than 0 means "traversing this cell costs X points".
// The pathfinding algorithms will respect that value when finding a solution.
//
// In the simplest situations just use 0 (no path) and 1 (can traverse).
type GridLayer [2]uint64

// MakeGridLayer is a GridLayer constructor function.
// It uses a temporary array to fill the result layer.
//
// The array represents a mapping from a tile tag (the key)
// to a traversal cost (the value).
//
// If you want blocked (occupied) tiles to be traversible,
// see MakeGridLayerWithBlocked constructor.
func MakeGridLayer(values [8]uint8) GridLayer {
	return MakeGridLayerWithBlocked(values, [8]uint8{})
}

// MakeGridLayerWithBlocked is like MakeGridLayer, but allows
// a custom movement cost per a blocked tile.
// Setting a movement cost of 0 will keep that cell impossible to move through.
func MakeGridLayerWithBlocked(values [8]uint8, blocked [8]uint8) GridLayer {
	tileMapping := uint64(0)
	for i := range values {
		tileMapping |= uint64(values[i]) << (i * 8)
	}
	blockedMapping := uint64(0)
	for i := range blocked {
		blockedMapping |= uint64(blocked[i]) << (i * 8)
	}
	return GridLayer([2]uint64{tileMapping, blockedMapping})
}

// Get maps a given tile tag into a traversal score.
// A tile tag is a value in [0-7] range.
func (l GridLayer) Get(tileTag uint8) uint8 {
	return uint8(l[0] >> (uint64(tileTag) * 8))
}

func (l GridLayer) getFast(tag uint8) uint8 {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(&l), tag))
}
