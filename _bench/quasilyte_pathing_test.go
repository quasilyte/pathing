package bench

import (
	"github.com/quasilyte/pathing"
)

var pathingLayer = pathing.MakeGridLayer(1, 0, 0, 0)

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
	width := float64(tc.cellWidth) * float64(tc.numCols)
	height := float64(tc.cellHeight) * float64(tc.numRows)
	t.grid = pathing.NewGrid(width, height, 0)
}

func (t *quasilytePathingTester) BuildPath() (pathing.GridPath, gridCoord) {
	from := pathing.GridCoord{X: t.tc.start.X, Y: t.tc.start.Y}
	to := pathing.GridCoord{X: t.tc.finish.X, Y: t.tc.finish.Y}
	result := t.bfs.BuildPath(t.grid, from, to, pathingLayer)
	return result.Steps, gridCoord{X: result.Finish.X, Y: result.Finish.Y}
}
