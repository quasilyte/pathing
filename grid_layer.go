package pathing

import (
	"unsafe"
)

// GridLayer is a tile-to-cost mapper.
// Every Grid cell has a tile tag value ranging from 0 to 3 (2 bits).
// Layers are used to turn that tag value into an actual traversal cost.
//
// Although you can construct a GridLayer value yourself,
// it's much easier to use a MakeGridLayer() function.
//
// A value of 0 means "the path is blocked".
// A value higher than 0 means "traversing this cell costs X points".
// The pathfinding algorithms will respect that value when finding a solution.
//
// In the simplest situations just use 0 (no path) and 1 (can traverse).
type GridLayer uint32

// MakeGridLayer is a GridLayer constructor function.
// It uses a temporary array to fill the result layer.
//
// The array represents a mapping from a tile tag (the key)
// to a traversal cost (the value).
func MakeGridLayer(values [4]uint8) GridLayer {
	v0 := values[0]
	v1 := values[1]
	v2 := values[2]
	v3 := values[3]
	merged := uint32(v0) | uint32(v1)<<8 | uint32(v2)<<16 | uint32(v3)<<24
	return GridLayer(merged)
}

// Get maps a given tile tag into a traversal score.
// A tile tag is a value in [0-3] range.
func (l GridLayer) Get(tileTag uint8) uint8 {
	return uint8(l >> (uint32(tileTag) * 8))
}

func (l GridLayer) getFast(tag uint8) uint8 {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(&l), tag))
}
