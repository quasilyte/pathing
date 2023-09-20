package pathing_test

import (
	"testing"

	"github.com/quasilyte/pathing"
)

func BenchmarkPathgridGetCellValue(b *testing.B) {
	p := pathing.NewGrid(pathing.GridConfig{WorldWidth: 1856, WorldHeight: 1856})
	l := pathing.MakeGridLayer(1, 0, 2, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GetCellValue(pathing.GridCoord{14, 5}, l)
	}
}

func BenchmarkPathgridSetCellTag(b *testing.B) {
	p := pathing.NewGrid(pathing.GridConfig{WorldWidth: 1856, WorldHeight: 1856})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.SetCellTag(pathing.GridCoord{14, 5}, 1)
	}
}
