package bench

import (
	"fmt"
	"strings"

	"github.com/beefsack/go-astar"
)

// This implementation comes from the provided benchmark.
// See https://github.com/beefsack/go-astar/blob/4ecf9e3044829f1e2d9eba4ff26295a2bbb9b2cd/path_test.go#L101
//
// It doesn't look good performance-wise: there are several things that can be improved.
// But I would guess that it's something that most users will start with.
//
// It is interesting, hovewer, to see the real limits of this library as opposed to
// limits of this particular interface implmenetation limits.
// Maybe some other day?

type beefsackAstarTester struct {
	tc *testCase

	w astarWorld
}

func newBeefsackAstarTester() *beefsackAstarTester {
	return &beefsackAstarTester{}
}

func (t *beefsackAstarTester) Init(tc *testCase) {
	t.tc = tc
	t.w = t.ParseWorld(strings.Join(tc.rawLayout, "\n"))
}

func (t *beefsackAstarTester) BuildPath() ([]astar.Pather, gridCoord) {
	p, _, _ := astar.Path(t.w.From(), t.w.To())
	last := p[0].(*astarTile)
	return p, last.Coord
}

type astarTile struct {
	Kind  int
	Coord gridCoord
	World astarWorld
}

func (t *astarTile) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.World.Tile(t.Coord.X+offset[0], t.Coord.Y+offset[1]); n != nil &&
			n.Kind != KindBlocker {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (t *astarTile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*astarTile)
	return KindCosts[toT.Kind]
}

func (t *astarTile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*astarTile)
	absX := toT.Coord.X - t.Coord.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Coord.Y - t.Coord.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

type astarWorld map[int]map[int]*astarTile

func (w astarWorld) Tile(x, y int) *astarTile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

func (w astarWorld) RenderPath(path []astar.Pather) string {
	width := len(w)
	if width == 0 {
		return ""
	}
	height := len(w[0])
	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*astarTile)
		pathLocs[fmt.Sprintf("%d,%d", pT.Coord.X, pT.Coord.Y)] = true
	}
	rows := make([]string, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t := w.Tile(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = KindRunes[KindPath]
			} else if t != nil {
				r = KindRunes[t.Kind]
			}
			rows[y] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

func (w astarWorld) SetTile(t *astarTile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*astarTile{}
	}
	w[x][y] = t
	t.Coord.X = x
	t.Coord.Y = y
	t.World = w
}

func (w astarWorld) FirstOfKind(kind int) *astarTile {
	for _, row := range w {
		for _, t := range row {
			if t.Kind == kind {
				return t
			}
		}
	}
	return nil
}

func (w astarWorld) From() *astarTile {
	return w.FirstOfKind(KindFrom)
}

func (w astarWorld) To() *astarTile {
	return w.FirstOfKind(KindTo)
}

// Kind* constants refer to tile kinds for input and output.
const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota

	// KindBlocker (X) is a tile which blocks movement.
	KindBlocker

	// KindFrom (F) is a tile which marks where the path should be calculated
	// from.
	KindFrom

	// KindTo (T) is a tile which marks the goal of the path.
	KindTo

	// KindPath (●) is a tile to represent where the path is in the output.
	KindPath
)

var KindRunes = map[int]rune{
	KindPlain:   ' ',
	KindBlocker: 'x',
	KindFrom:    'S',
	KindTo:      'F',
	KindPath:    '●',
}

var RuneKinds = map[rune]int{
	' ': KindPlain,
	'x': KindBlocker,
	'S': KindFrom,
	'F': KindTo,
}

var KindCosts = map[int]float64{
	KindPlain: 1.0,
	KindFrom:  1.0,
	KindTo:    1.0,
}

func (t *beefsackAstarTester) ParseWorld(input string) astarWorld {
	w := astarWorld{}
	for y, row := range strings.Split(input, "\n") {
		for x, raw := range row {
			kind, ok := RuneKinds[raw]
			if !ok {
				kind = KindBlocker
			}
			w.SetTile(&astarTile{
				Kind: kind,
			}, x, y)
		}
	}
	return w
}
