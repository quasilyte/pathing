package bench

import (
	"image"

	"github.com/s0rg/grid"
)

// Using the benchmark code for the reference.
// See https://github.com/s0rg/grid/blob/13cbded225ee60de458020897aefd7aa972183d6/grid_test.go#L718

type s0rgGridTester struct {
	tc *testCase

	grid     *grid.Map[struct{}]
	dirCross []image.Point
}

func newS0rgGridTester() *s0rgGridTester {
	return &s0rgGridTester{}
}

func (t *s0rgGridTester) Init(tc *testCase) {
	t.tc = tc
	t.grid = grid.New[struct{}](image.Rect(0, 0, t.tc.numCols, t.tc.numRows))
	t.dirCross = grid.Points(grid.DirectionsCardinal...)
}

func (t *s0rgGridTester) moveCost(coord image.Point, dist float64, _ struct{}) (cost float64, walkable bool) {
	return dist, t.tc.layout[coord.Y][coord.X] != 'x'
}

func (t *s0rgGridTester) BuildPath() ([]image.Point, gridCoord) {
	from := image.Pt(t.tc.start.X, t.tc.start.Y)
	to := image.Pt(t.tc.finish.X, t.tc.finish.Y)
	result, _ := t.grid.Path(from, to, t.dirCross, grid.DistanceManhattan, t.moveCost)
	last := result[len(result)-1]
	return result, gridCoord{X: last.X, Y: last.Y}
}
