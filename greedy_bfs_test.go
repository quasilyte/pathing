package pathing_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/quasilyte/pathing"
)

func BenchmarkGreedyBFS(b *testing.B) {
	l := pathing.MakeGridLayer([4]uint8{1, 0, 1, 1})
	for i := range bfsTests {
		test := bfsTests[i]
		if !test.bench {
			continue
		}
		numCols := len(test.path[0])
		numRows := len(test.path)
		b.Run(fmt.Sprintf("%s_%dx%d", test.name, numCols, numRows), func(b *testing.B) {
			parseResult := testParseGrid(b, test.path)
			bfs := pathing.NewGreedyBFS(pathing.GreedyBFSConfig{
				NumCols: uint(parseResult.numCols),
				NumRows: uint(parseResult.numRows),
			})
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bfs.BuildPath(parseResult.grid, parseResult.start, parseResult.dest, l)
			}
		})
	}
}

func TestGreedyBFS(t *testing.T) {
	runTestOnce := func(t *testing.T, test bfsTestCase, m []string, parseResult testGrid, bfs *pathing.GreedyBFS, grid *pathing.Grid) {
		t.Helper()

		l := pathing.MakeGridLayer([4]uint8{1, 0, 1, 1})

		result := bfs.BuildPath(grid, parseResult.start, parseResult.dest, l)
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
			default:
				panic(fmt.Sprintf("unexpected %c marker", marker))
			}
		}

		have := string(bytes.Join(haveRows, []byte("\n")))
		want := strings.Join(m, "\n")

		if have != want {
			t.Fatalf("paths mismatch\nmap:\n%s\nhave (l=%d):\n%s\nwant (l=%d):\n%s",
				strings.Join(m, "\n"), pathLen, have, parseResult.pathLen, want)
		}

		wantPartial := test.partial
		havePartial := pos != parseResult.dest && result.Partial
		if havePartial != wantPartial {
			t.Fatalf("partial flag mismatch\nmap:\n%s\nhave: %v\nwant: %v", strings.Join(m, "\n"), havePartial, wantPartial)
		}
	}

	runTestCase := func(t *testing.T, test bfsTestCase, offset, offset2 pathing.GridCoord) {
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
		bfs := pathing.NewGreedyBFS(pathing.GreedyBFSConfig{
			NumCols: uint(parseResult.numCols),
			NumRows: uint(parseResult.numRows),
		})
		grid := parseResult.grid

		for i := 0; i < 5; i++ {
			runTestOnce(t, test, m, parseResult, bfs, grid)
		}
	}

	for i := range bfsTests {
		test := bfsTests[i]
		t.Run(test.name, func(t *testing.T) {
			runTestCase(t, test, pathing.GridCoord{}, pathing.GridCoord{})
		})

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

type bfsTestCase struct {
	name    string
	path    []string
	partial bool
	bench   bool
}

var bfsTests = []bfsTestCase{
	{
		name: "trivial_short",
		path: []string{
			"..........",
			"...A   $..",
			"..........",
		},
		bench: true,
	},

	{
		name: "trivial_short2",
		path: []string{
			"..........",
			"...A......",
			"... ......",
			"... ......",
			"...  $....",
			"..........",
		},
		bench: true,
	},

	{
		name: "trivial",
		path: []string{
			".A..........",
			". ..........",
			". ..........",
			". ..........",
			". ..........",
			". ..........",
			".          $",
		},
		bench: true,
	},

	{
		name: "trivial_long",
		path: []string{
			".......................x........",
			"                               $",
			"A...............................",
			"..........................x.....",
		},
		bench: true,
	},

	{
		name: "simple_wall1",
		path: []string{
			"........",
			"...A....",
			"...   ..",
			"....x ..",
			"....x $.",
		},
		bench: true,
	},

	{
		name: "simple_wall2",
		path: []string{
			"...   ..",
			"...Ax ..",
			"....x ..",
			"....x ..",
			"....x $.",
		},
		bench: true,
	},

	{
		name: "simple_wall3",
		path: []string{
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			".............   ................",
			"..            x          $......",
			".. ...........x.................",
			"..A...........x.................",
			"....x...........................",
			"....x...........................",
			"....x...........................",
			"....x...........................",
		},
		bench: true,
	},

	{
		name: "simple_wall4",
		path: []string{
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"................................",
			"..............x.................",
			"..............x.................",
			"..A...........x.................",
			".. .x...........................",
			".. .x...........................",
			".. .x...........................",
			".. .x...........................",
			".. .............................",
			".. .............................",
			".. ..................xxxxxxxx...",
			".. .............................",
			".. .............................",
			".. ...........x.................",
			".. ...........x.................",
			"..    ........x.................",
			"....x ..........................",
			"....x                      $....",
			"....x...........................",
			"....x...........................",
		},
		bench: true,
	},

	{
		name: "zigzag1",
		path: []string{
			"........",
			"   A....",
			" xxxxxx.",
			" .......",
			" .xxxxxx",
			" .......",
			" $......",
		},
		bench: true,
	},

	{
		name: "zigzag2",
		path: []string{
			"........",
			"...A    ",
			".xxxxxx ",
			".....   ",
			"..xxx xx",
			"..... ..",
			".....  $",
		},
		bench: true,
	},

	{
		name: "zigzag3",
		path: []string{
			"...   ....x.....",
			"..A x ....x.....",
			"....x ....x.....",
			"....x ....x.....",
			"....x        $..",
			"....x...........",
		},
		bench: true,
	},

	{
		name: "zigzag4",
		path: []string{
			"...   .x.   x...",
			"... x .x. x x...",
			"... x .x. x x...",
			"... x .x. x   ..",
			"..A x  x  x.x  $",
			"....x.   .x.x...",
		},
		bench: true,
	},

	{
		name: "zigzag5",
		path: []string{
			".A     ..",
			"xxxxxx ..",
			"..     ..",
			".. xxxxxx",
			"..   ....",
			"xxxx x...",
			"....    .",
			"...xxxx x",
			".......$.",
		},
		bench: true,
	},

	{
		name: "double_corner1",
		path: []string{
			".   .x  A.",
			". x .x ...",
			"x x .x ...",
			"  x .x ...",
			" xx    ...",
			" .xxxxxxxx",
			"   $......",
		},
		bench: true,
	},

	{
		name: "double_corner2",
		path: []string{
			".    x..A.",
			". x. x.. .",
			"x x. x.. .",
			"  x. x.. .",
			" xx.     .",
			" .xxxxxxxx",
			"        $.",
			"..........",
		},
	},

	{
		name: "double_corner3",
		path: []string{
			"    x..A.",
			" x. x.. .",
			" x. x.. .",
			" x.     .",
			" xxxxxxxx",
			"       $.",
		},
	},

	{
		name: "labyrinth1",
		path: []string{
			".........x.....",
			"xxxxxxxx.x.  $.",
			"x.     x.x. ...",
			"x. xxx x.x. ...",
			"x.   x x.x. ...",
			"x...Ax   xx .xx",
			"x....x.x x  ...",
			"xxxxxx.x x xxxx",
			"x......x x    .",
			"xxxxxxxx xxxx x",
			"........ x    .",
			"........   ....",
		},
		bench: true,
	},

	{
		name: "labyrinth2",
		path: []string{
			".x......x.......x............",
			".x......x.......x............",
			".x......x.......x............",
			".x......x.......xxxxxxxxxx...",
			".x....       ...x.....    ...",
			".x     xxx.x    x.....$.x  xx",
			"   .x..x...xxx. x.......x.  .",
			"A...x..x...x... xxxxxxxxxxx .",
			"..x.x..x.......     x       .",
			"..x.x..x....x...... x .......",
			"..x.x..x..xxxx...x.   .......",
			"..x.x.......x....x...........",
		},
		bench: true,
	},

	{
		name: "labyrinth3",
		path: []string{
			"...x......x........x............",
			"..Ax......x........x............",
			".. x......x........xxxxxxxxxx...",
			".. x...............x............",
			".. x.....xxx..x....x.......x..xx",
			".. ...x..x....xxx..x.......x....",
			".. ...x..x....x....xxxxxxxxxxx..",
			".. .x.x..x.....x...   .x........",
			".. .x.x..x...xxxx.  x         ..",
			"..        x....... xxxxxxxxxx ..",
			"xxxx.....       .. x........  ..",
			"...x.....xxx..x    x.......x .xx",
			"......x..x....xxx..x.......x   .",
			"......x..x....x....xxxxxx.xxxx .",
			"....x.x........x....x..x...$   .",
		},
		bench: true,
	},

	{
		// This is unfortunate.
		// TODO: can we adjust anything to make it better?
		name: "depth1",
		path: []string{
			"........................",
			".xxxxxxxxxxxxxxxxxxxx...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxx.x...",
			"....................x...",
			".x.xxxxxxxxxxxxxxxxxx...",
			"..................A x.B.",
			".x.xxxxxxxxxxxxxxxxxx...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxx.x...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxxxx...",
			"........................",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "depth2",
		path: []string{
			"...................   ..",
			"..                  x ..",
			".x xxxxxxxxxxxxxxxxxx ..",
			"..                A.x $.",
			".x.xxxxxxxxxxxxxxxxxx...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxx.x...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxxxx...",
			"........................",
		},
		bench: true,
	},

	{
		name: "nopath1",
		path: []string{
			"A    x.B",
			".....x..",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath2",
		path: []string{
			"....Ax.B",
			".....x..",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath3",
		path: []string{
			"........",
			".xxxxx..",
			".x...x..",
			".x.A.x..",
			".x.  x..",
			".xxxxx..",
			".......B",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath4",
		path: []string{
			".....x.....x..",
			".xxxxx.   .x..",
			".x...x. x .x..",
			".x.A.x. x  x..",
			".x.     xxxx..",
			".xxxxxxxx.....",
			".............B",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath5",
		path: []string{
			".B...x.....x..",
			".xxxxx.....x..",
			".x  .x..x..x..",
			".x.A.x..x..x..",
			".x......xxxx..",
			".xxxxxxxx.....",
			"..............",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "tricky1",
		path: []string{
			"               $",
			" .xxxxxxxxxxxx..",
			" ............x..",
			" ............x..",
			" ............x..",
			" ............x..",
			" ............x..",
			"A..xxxxxxxxxxx..",
			"................",
		},
		bench: true,
	},
	{
		name: "tricky2",
		path: []string{
			"...............",
			".             .",
			"  xxxxxxxxxxx $",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			"A.xxxxxxxxxxx..",
			"...............",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky3",
		path: []string{
			"...............",
			"...............",
			"..xxxxxxxxxxx A",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"$ xxxxxxxxxxx .",
			".             .",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky4",
		path: []string{
			"...............",
			".             .",
			". xxxxxxxxxxx $",
			".     ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			".....A......x..",
			"..xxxxxxxxxxx..",
			"...............",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky5",
		path: []string{
			"...............",
			"...............",
			"A.xxxxxxxxxxx..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			"  xxxxxxxxxxx $",
			".             .",
			"...............",
		},
	},

	{
		name: "tricky6",
		path: []string{
			"............$ .",
			"............. .",
			"..xxxxxxxxxxx .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..            .",
			"..A............",
		},
	},

	{
		name: "tricky7",
		path: []string{
			"..          A..",
			".  ............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". .............",
			". $............",
		},
	},

	{
		name: "tricky8",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			".            ..",
			"............A..",
		},
	},

	{
		name: "tricky9",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x         x..",
			". x .......Ax..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			".   ...........",
			"...............",
		},
	},

	{
		name: "tricky10",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". xA........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			".   ...........",
			"...............",
		},
	},

	{
		name: "tricky11",
		path: []string{
			".    $.........",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.        x..",
			". x  ...... x..",
			".   .......  ..",
			"............A..",
		},
	},

	{
		name: "tricky12",
		path: []string{
			"..........$   .",
			"............. .",
			"..xxxxxxxxxxx .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"............  .",
			"............A..",
		},
	},

	{
		name: "tricky13",
		path: []string{
			"...............",
			"           $...",
			" .....xxxxxxx..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			"A.xxxxxxxxxxx..",
			"...............",
			"...............",
		},
	},

	{
		name: "distlimit1",
		path: []string{
			"A                                                        ..........B",
		},
		bench:   true,
		partial: true,
	},

	{
		name: "distlimit2",
		path: []string{
			"A.............x......   ....            ......x.....x.....x....",
			" .............x...... x      xxxxxxxxxx ......x..x..x..x..x....",
			" ...xxxxxxxxxxx...... x...............x ......x..x..x..x..x....",
			"                      x...............x       ...x.....x......B",
		},
		bench:   true,
		partial: true,
	},
}
