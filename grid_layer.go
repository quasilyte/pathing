package pathing

import "unsafe"

type GridLayer uint32

func MakeGridLayer(values [4]uint8) GridLayer {
	v0 := values[0]
	v1 := values[1]
	v2 := values[2]
	v3 := values[3]
	merged := uint32(v0) | uint32(v1)<<8 | uint32(v2)<<16 | uint32(v3)<<24
	return GridLayer(merged)
}

func (l GridLayer) Get(tag uint8) uint8 {
	return uint8(l >> (uint32(tag) * 8))
}

func (l GridLayer) getFast(tag uint8) uint8 {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(&l), tag))
}
