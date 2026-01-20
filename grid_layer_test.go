package pathing

import (
	"testing"
)

func TestGridLayerBlocked(t *testing.T) {
	tests := []struct {
		data  [2][8]uint8
		index int
		want  uint8
	}{
		{
			data:  [2][8]uint8{{1: 11}, {2: 22}},
			index: 1,
			want:  11,
		},
		{
			data:  [2][8]uint8{{1: 11}, {2: 22}},
			index: 0,
			want:  0,
		},
		{
			data:  [2][8]uint8{{1: 11}, {2: 22}},
			index: 2 | (1 << 3),
			want:  22,
		},
		{
			data:  [2][8]uint8{{1: 11}, {2: 22}},
			index: 1 | (1 << 3),
			want:  0,
		},
		{
			data:  [2][8]uint8{{1: 11}, {}},
			index: 2 | (1 << 3),
			want:  0,
		},
		{
			data:  [2][8]uint8{{1: 11}, {}},
			index: 1 | (1 << 3),
			want:  0,
		},
	}

	for _, test := range tests {
		l := MakeGridLayerWithBlocked(test.data[0], test.data[1])
		want := test.want
		have := l.getFast(uint8(test.index))
		if want != have {
			t.Fatalf("layer %v getFast(%d %04b):\nhave: %v\nwant: %v", test.data, test.index, test.index, have, want)
		}
	}
}

func TestGridLayer(t *testing.T) {
	tests := [][]uint8{
		{0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 1, 0, 0, 2, 0},
		{1, 0, 0, 0, 0, 0, 3, 0},
		{1, 1, 1, 1, 0, 0, 4, 0},
		{10, 0, 10, 0, 0, 0, 0, 0xff},
		{1, 2, 3, 4, 0, 0, 0, 0},
		{4, 3, 2, 1, 0, 0, 0, 0},
		{0xff, 0xff, 0xff, 0xff, 0, 0xff, 0, 0},
		{100, 0xff, 0xff, 100, 0, 0, 0xfe, 0xfa},
		{24, 53, 21, 99, 0, 0, 0, 0},
		{99, 145, 9, 0, 0, 0, 0, 0},
		{10, 0, 20, 30, 0, 0, 0, 0},
	}

	for _, test := range tests {
		l := MakeGridLayer(([8]uint8)(test))
		for i := uint8(0); i <= 7; i++ {
			want := test[i]
			have := l.Get(i)
			if want != have {
				t.Fatalf("(%v).Get(%d): have %v, want %v", test, i, have, want)
			}
			have2 := l.getFast(i)
			if want != have2 {
				t.Fatalf("(%v).getFast(%d): have %v, want %v", test, i, have, want)
			}
		}
	}
}
