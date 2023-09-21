package bench

import (
	"github.com/quasilyte/pathing"
)

type quasilytePathingAStarTester struct {
	tc *testCase

	grid  *pathing.Grid
	astar *pathing.AStar
}

func newQuasilytePathingAStarTester() *quasilytePathingAStarTester {
	return &quasilytePathingAStarTester{}
}

func (t *quasilytePathingAStarTester) Init(tc *testCase) {
	t.tc = tc
	t.astar = pathing.NewAStar(pathing.AStarConfig{
		NumCols: uint(tc.numCols),
		NumRows: uint(tc.numRows),
	})
	width := tc.cellWidth * tc.numCols
	height := tc.cellHeight * tc.numRows
	t.grid = pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(width),
		WorldHeight: uint(height),
	})
}

func (t *quasilytePathingAStarTester) BuildPath() (pathing.GridPath, gridCoord) {
	from := pathing.GridCoord{X: t.tc.start.X, Y: t.tc.start.Y}
	to := pathing.GridCoord{X: t.tc.finish.X, Y: t.tc.finish.Y}
	result := t.astar.BuildPath(t.grid, from, to, pathingLayer)
	return result.Steps, gridCoord{X: result.Finish.X, Y: result.Finish.Y}
}
