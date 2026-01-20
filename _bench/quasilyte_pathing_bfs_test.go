package bench

import (
	"github.com/quasilyte/pathing"
)

var pathingLayer = pathing.MakeGridLayer([8]uint8{1, 0, 0, 0, 0, 0, 0, 0})

type quasilytePathingBFSTester struct {
	tc *testCase

	grid *pathing.Grid
	bfs  *pathing.GreedyBFS
}

func newQuasilytePathingBFSTester() *quasilytePathingBFSTester {
	return &quasilytePathingBFSTester{}
}

func (t *quasilytePathingBFSTester) Init(tc *testCase) {
	t.tc = tc
	t.bfs = pathing.NewGreedyBFS(pathing.GreedyBFSConfig{
		NumCols: uint(tc.numCols),
		NumRows: uint(tc.numRows),
	})
	width := tc.cellWidth * tc.numCols
	height := tc.cellHeight * tc.numRows
	t.grid = pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(width),
		WorldHeight: uint(height),
	})
	fillPathingGrid(t.grid, tc)
}

func (t *quasilytePathingBFSTester) BuildPath() (pathing.GridPath, gridCoord) {
	from := pathing.GridCoord{X: t.tc.start.X, Y: t.tc.start.Y}
	to := pathing.GridCoord{X: t.tc.finish.X, Y: t.tc.finish.Y}
	result := t.bfs.BuildPath(t.grid, from, to, pathingLayer)
	return result.Steps, gridCoord{X: result.Finish.X, Y: result.Finish.Y}
}

func fillPathingGrid(g *pathing.Grid, tc *testCase) {
	for y, row := range tc.layout {
		for x, col := range row {
			if col != 'x' {
				continue
			}
			g.SetCellTile(pathing.GridCoord{X: x, Y: y}, 1)
		}
	}
}
