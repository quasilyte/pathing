package pathing_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/quasilyte/pathing"
)

type pathBuilder interface {
	BuildPath(g *pathing.Grid, from, to pathing.GridCoord, l pathing.GridLayer) pathing.BuildPathResult
}

type testGrid struct {
	start    pathing.GridCoord
	dest     pathing.GridCoord
	grid     *pathing.Grid
	pathLen  int
	numCols  int
	numRows  int
	haveRows [][]byte
}

func runPathfindTest(t *testing.T, test pathfindTestCase, constructor func(uint, uint) pathBuilder) {
	t.Helper()

	runTestOnce := func(t *testing.T, test pathfindTestCase, m []string, parseResult testGrid, impl pathBuilder, grid *pathing.Grid) {
		t.Helper()

		l := test.layer
		if l == 0 {
			l = pathing.MakeGridLayer([4]uint8{1, 0, 2, 3})
		}

		result := impl.BuildPath(grid, parseResult.start, parseResult.dest, l)
		path := result.Steps

		haveRows := make([][]byte, len(parseResult.haveRows))
		for i, row := range parseResult.haveRows {
			haveRows[i] = make([]byte, len(row))
			copy(haveRows[i], row)
		}

		pos := parseResult.start
		pathLen := 0
		for path.HasNext() {
			pathLen++
			d := path.Next()
			pos = pos.Move(d)
			marker := haveRows[pos.Y][pos.X]
			switch marker {
			case 'A':
				haveRows[pos.Y][pos.X] = 'A'
			case 'B':
				haveRows[pos.Y][pos.X] = '$'
			case ' ':
				t.Fatal("visited one cell more than once")
			case '.':
				haveRows[pos.Y][pos.X] = ' '
			case 'o':
				haveRows[pos.Y][pos.X] = 'O'
			case 'w':
				haveRows[pos.Y][pos.X] = 'W'
			case 'O', 'W':
				haveRows[pos.Y][pos.X] = marker
			default:
				panic(fmt.Sprintf("unexpected %c marker", marker))
			}
		}

		have := string(bytes.Join(haveRows, []byte("\n")))
		want := strings.Join(m, "\n")

		haveCost := result.Cost
		wantCost := test.cost
		if wantCost == 0 {
			wantCost = result.Steps.Len()
		}
		if haveCost != wantCost {
			t.Fatalf("costs mismatch\nmap:\n%s\nhave (l=%d c=%d):\n%s\nwant (l=%d c=%d):\n%s",
				strings.Join(m, "\n"), pathLen, haveCost, have, parseResult.pathLen, wantCost, want)
		}

		if have != want {
			t.Fatalf("paths mismatch\nmap:\n%s\nhave (l=%d c=%d):\n%s\nwant (l=%d):\n%s",
				strings.Join(m, "\n"), pathLen, result.Cost, have, parseResult.pathLen, want)
		}

		wantPartial := test.partial
		havePartial := pos != parseResult.dest && result.Partial
		if havePartial != wantPartial {
			t.Fatalf("partial flag mismatch\nmap:\n%s\nhave: %v\nwant: %v", strings.Join(m, "\n"), havePartial, wantPartial)
		}
	}

	runTestCase := func(t *testing.T, test pathfindTestCase, offset, offset2 pathing.GridCoord) {
		t.Helper()

		m := make([]string, len(test.path))
		copy(m, test.path)
		if offset.X != 0 {
			for y := range m {
				m[y] = strings.Repeat("x", offset.X) + m[y]
			}
		}
		if offset2.X != 0 {
			for y := range m {
				m[y] = m[y] + strings.Repeat("x", offset2.X)
			}
		}
		if offset.Y != 0 {
			row := strings.Repeat("x", len(m[0]))
			extraRows := make([]string, offset.Y)
			for i := range extraRows {
				extraRows[i] = row
			}
			m = append(extraRows, m...)
		}
		if offset2.Y != 0 {
			row := strings.Repeat("x", len(m[0]))
			for i := 0; i < offset2.Y; i++ {
				m = append(m, row)
			}
		}

		parseResult := testParseGrid(t, m)
		impl := constructor(uint(parseResult.numCols), uint(parseResult.numRows))
		grid := parseResult.grid

		for i := 0; i < 5; i++ {
			runTestOnce(t, test, m, parseResult, impl, grid)
		}
	}

	t.Run(test.name, func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{}, pathing.GridCoord{})
	})
	if t.Failed() {
		return
	}

	t.Run(test.name+"with_offset_x", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 8}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_x2", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 500}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_y", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{Y: 24}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_y2", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{Y: 600}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_xy", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 32, Y: 120}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_xy2", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 64, Y: 32}, pathing.GridCoord{})
	})
	t.Run(test.name+"with_offset_xy3", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 150, Y: 150}, pathing.GridCoord{})
	})

	t.Run(test.name+"with_offset2", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 8, Y: 8}, pathing.GridCoord{X: 8, Y: 8})
	})
	t.Run(test.name+"with_offset2_2", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{}, pathing.GridCoord{X: 150, Y: 150})
	})
	t.Run(test.name+"with_offset2_x", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{X: 150}, pathing.GridCoord{X: 150})
	})
	t.Run(test.name+"with_offset2_y", func(t *testing.T) {
		runTestCase(t, test, pathing.GridCoord{Y: 150}, pathing.GridCoord{Y: 150})
	})
}

func testParseGrid(tb testing.TB, m []string) testGrid {
	tb.Helper()

	numCols := len(m[0])
	numRows := len(m)

	grid := pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  32 * uint(numCols),
		WorldHeight: 32 * uint(numRows),
	})

	pathLen := 0
	var startPos pathing.GridCoord
	var destPos pathing.GridCoord
	haveRows := make([][]byte, numRows)
	for row := 0; row < numRows; row++ {
		haveRows[row] = make([]byte, numCols)
		for col := 0; col < numCols; col++ {
			marker := m[row][col]
			cell := pathing.GridCoord{X: col, Y: row}
			haveRows[row][col] = marker
			switch marker {
			case 'x':
				grid.SetCellTile(cell, 1)
			case 'o', 'O':
				grid.SetCellTile(cell, 2)
			case 'w', 'W':
				grid.SetCellTile(cell, 3)
			case 'A':
				startPos = cell
			case 'B', '$':
				if marker == '$' {
					pathLen++
				}
				destPos = cell
				haveRows[row][col] = 'B'
			case ' ':
				pathLen++
				haveRows[row][col] = '.'
			}
		}
	}

	return testGrid{
		pathLen:  pathLen,
		start:    startPos,
		dest:     destPos,
		numRows:  numRows,
		numCols:  numCols,
		haveRows: haveRows,
		grid:     grid,
	}
}

type pathfindTestCase struct {
	name    string
	path    []string
	cost    int
	layer   pathing.GridLayer
	partial bool
	bench   bool
}
