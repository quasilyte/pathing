package pathing_test

import (
	"testing"

	"github.com/quasilyte/pathing"
)

func BenchmarkPathgridGetCellCost(b *testing.B) {
	p := pathing.NewGrid(pathing.GridConfig{WorldWidth: 1856, WorldHeight: 1856})
	l := pathing.MakeGridLayer([4]uint8{1, 0, 2, 3})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GetCellCost(pathing.GridCoord{14, 5}, l)
	}
}

func BenchmarkPathgridSetCellTile(b *testing.B) {
	p := pathing.NewGrid(pathing.GridConfig{WorldWidth: 1856, WorldHeight: 1856})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.SetCellTile(pathing.GridCoord{14, 5}, 1)
	}
}
