package bench

import (
	"strings"
	"testing"
)

func BenchmarkQuasilytePathingBFS(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newQuasilytePathingBFSTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func BenchmarkQuasilytePathingAStar(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newQuasilytePathingAStarTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func BenchmarkFzippAstar(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newFzippAstarTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func BenchmarkS0rgGrid(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newS0rgGridTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func BenchmarkBeefsackAstar(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newBeefsackAstarTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func BenchmarkSolarLunePaths(b *testing.B) {
	for i := range testCaseList {
		tc := testCaseList[i]
		b.Run(tc.name, func(b *testing.B) {
			lib := newSolarLunePathsTester()
			lib.Init(tc)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, finish := lib.BuildPath()
				validateResult(b, tc, finish)
			}
		})
	}
}

func validateResult(tb testing.TB, tc *testCase, finish gridCoord) {
	tb.Helper()
	if finish.X != tc.finish.X || finish.Y != tc.finish.Y {
		tb.Fatalf("incorrect results\nhave: %v\nwant: %v", finish, tc.finish)
	}
}

type testCase struct {
	name      string
	layout    []string
	rawLayout []string

	cellWidth  int
	cellHeight int
	numCols    int
	numRows    int
	start      gridCoord
	finish     gridCoord
}

type gridCoord struct {
	X int
	Y int
}

func initTestCase(c *testCase) *testCase {
	c.numCols = len(c.layout[0])
	c.numRows = len(c.layout)
	c.cellWidth = 32
	c.cellHeight = 32

	// rawLayout keeps S and F markers.
	c.rawLayout = make([]string, len(c.layout))
	copy(c.rawLayout, c.layout)

	for y, row := range c.layout {
		x := strings.IndexByte(row, 'S')
		if x == -1 {
			continue
		}
		c.start = gridCoord{X: x, Y: y}
		c.layout[y] = strings.Replace(row, "S", " ", 1)
		break
	}
	for y, row := range c.layout {
		x := strings.IndexByte(row, 'F')
		if x == -1 {
			continue
		}
		c.finish = gridCoord{X: x, Y: y}
		c.layout[y] = strings.Replace(row, "F", " ", 1)
		break
	}

	if c.start == (gridCoord{}) {
		panic("missing S marker")
	}
	if c.finish == (gridCoord{}) {
		panic("missing F marker")
	}

	return c
}

// All maps are 50x50 cells.
// S is start (free cell).
// F is finish (also a free cell).
var testCaseList = []*testCase{
	initTestCase(&testCase{
		name: "no_walls",
		layout: []string{
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"   S                                              ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                     F            ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
		},
	}),

	initTestCase(&testCase{
		name: "simple_wall",
		layout: []string{
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"   S           x    F                             ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"               x                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
		},
	}),

	initTestCase(&testCase{
		name: "pocket_wall",
		layout: []string{
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"  xxxxxxxxxxxxxxxxxxxxxxxxxxxxx                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x  F                ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"        S                     x                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"                              x                   ",
			"  xxxxxxxxxxxxxxxxxxxxxxxxxxxxx                   ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
		},
	}),

	initTestCase(&testCase{
		name: "multi_wall",
		layout: []string{
			"         x                             x          ",
			"    S    x                             x          ",
			"                                       x          ",
			"         x                             x          ",
			"xxxxx    x                             x          ",
			"         xxxxxxxxx                     x          ",
			"         x       x                     x          ",
			"  xxxxxxxx       x                     x          ",
			"         x       x                     x          ",
			"         x       xxxxxxxxxxxxxxxxxxxxxxx          ",
			"         x       x                                ",
			"         x       x                                ",
			"                 x F                              ",
			"                 x                                ",
			"                 x                                ",
			"                 x          x                     ",
			"                 x          x                     ",
			"                            x                     ",
			"         xxxxxxxxx          x                     ",
			"                 x          x                     ",
			"xxxx             x          x                     ",
			"                            x                     ",
			"                            x                     ",
			"                            x                     ",
			"                            x                     ",
			"xxxxxxxxxx                  x                     ",
			"                            x                     ",
			"                  x         x                     ",
			"                  x         x                     ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                  x                               ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
		},
	}),
}
