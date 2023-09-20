package bench

import (
	"github.com/quasilyte/pathing"
)

var pathingLayer = pathing.MakeGridLayer([4]uint8{1, 0, 0, 0})

type quasilytePathingTester struct {
	tc *testCase

	grid *pathing.Grid
	bfs  *pathing.GreedyBFS
}

func newQuasilytePathingTester() *quasilytePathingTester {
	return &quasilytePathingTester{}
}

func (t *quasilytePathingTester) Init(tc *testCase) {
	t.tc = tc
	t.bfs = pathing.NewGreedyBFS(tc.numCols, tc.numRows)
	width := tc.cellWidth * tc.numCols
	height := tc.cellHeight * tc.numRows
	t.grid = pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(width),
		WorldHeight: uint(height),
	})
}

func (t *quasilytePathingTester) BuildPath() (pathing.GridPath, gridCoord) {
	from := pathing.GridCoord{X: t.tc.start.X, Y: t.tc.start.Y}
	to := pathing.GridCoord{X: t.tc.finish.X, Y: t.tc.finish.Y}
	result := t.bfs.BuildPath(t.grid, from, to, pathingLayer)
	return result.Steps, gridCoord{X: result.Finish.X, Y: result.Finish.Y}
}
