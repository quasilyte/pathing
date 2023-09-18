package pathing

import (
	"fmt"
	"testing"
)

func BenchmarkCoordMapReset(b *testing.B) {
	sizes := []int{32, 256, 2048}
	for i := range sizes {
		size := sizes[i]
		b.Run(fmt.Sprintf("size%d", size), func(b *testing.B) {
			m := newCoordMap(size, size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Reset()
			}
		})
	}
}

func BenchmarkCoordMapSet(b *testing.B) {
	m := newCoordMap(8, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Reset()
		for j := 0; j < 8; j++ {
			m.Set(uint(j), DirUp)
		}
	}
}

func BenchmarkCoordMapGet(b *testing.B) {
	m := newCoordMap(8, 8)
	for j := 0; j < 8; j++ {
		m.Set(uint(j), DirUp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 8; j++ {
			_ = m.Get(uint(j))
		}
	}
}

func TestEmptyCoordMap(t *testing.T) {
	m := newCoordMap(0, 0)

	coords := []GridCoord{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},

		{0, -1},
		{-1, 0},
		{-1, -1},

		{0, 10},
		{10, 0},
		{10, 10},

		{100, 100},
		{-100, -100},
	}

	for _, coord := range coords {
		if m.Get(m.packCoord(coord)) != DirNone {
			t.Fatalf("empty coord map returns invalid result for %v", coord)
		}
	}
}

func TestCoordMap(t *testing.T) {
	m := newCoordMap(32, 32)

	coords := []GridCoord{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},

		{0, 10},
		{10, 0},
		{10, 10},
		{10, 30},

		{31, 31},
	}

	for i, coord := range coords {
		if m.Get(m.packCoord(coord)) != DirNone {
			t.Fatalf("Get(%v) expected to give None before insertion", coord)
		}
		dir := Direction(i % 4)
		m.Set(m.packCoord(coord), dir)
		for j := 0; j < 3; j++ {
			if got := m.Get(m.packCoord(coord)); got != dir {
				t.Fatalf("Get(%v) gives %s, expected %s", coord, got, dir)
			}
		}
		dir = Direction(3 - (i % 4))
		m.Set(m.packCoord(coord), dir)
		for j := 0; j < 3; j++ {
			if got := m.Get(m.packCoord(coord)); got != dir {
				t.Fatalf("Get(%v) gives %s, expected %s", coord, got, dir)
			}
		}
		for _, otherCoord := range coords[i:] {
			if coord == otherCoord {
				continue
			}
			if got := m.Get(m.packCoord(otherCoord)); got != DirNone {
				t.Fatalf("unrelated Get(%v) after Set(%v) gives %s, expected None", otherCoord, coord, got)
			}
		}
	}
}
