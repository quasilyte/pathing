package bench

import (
	"github.com/solarlune/paths"
)

type solarLunePathsTester struct {
	grid *paths.Grid
	tc   *testCase
}

func newSolarLunePathsTester() *solarLunePathsTester {
	return &solarLunePathsTester{}
}

func (t *solarLunePathsTester) Init(tc *testCase) {
	t.tc = tc
	t.grid = paths.NewGridFromStringArrays(tc.layout, tc.cellWidth, tc.cellHeight)
	t.grid.SetWalkable('x', false)
}

func (t *solarLunePathsTester) BuildPath() (*paths.Path, gridCoord) {
	from := t.grid.Get(t.tc.start.X, t.tc.start.Y)
	to := t.grid.Get(t.tc.finish.X, t.tc.finish.Y)
	result := t.grid.GetPathFromCells(from, to, false, false)
	finish := result.Get(result.Length() - 1)
	return result, gridCoord{X: finish.X, Y: finish.Y}
}
