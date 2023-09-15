package bench

import (
	"image"
	"math"

	"github.com/fzipp/astar"
)

type fzippAstarTester struct {
	tc *testCase

	m fzippGraph
}

type fzippGraph []string

func newFzippAstarTester() *fzippAstarTester {
	return &fzippAstarTester{}
}

func (t *fzippAstarTester) Init(tc *testCase) {
	t.tc = tc
	t.m = tc.layout
}

func (t *fzippAstarTester) BuildPath() (astar.Path[image.Point], gridCoord) {
	from := image.Pt(t.tc.start.X, t.tc.start.Y)
	to := image.Pt(t.tc.finish.X, t.tc.finish.Y)
	result := astar.FindPath[image.Point](t.m, from, to, t.nodeDist, t.nodeDist)
	last := result[len(result)-1]
	return result, gridCoord{X: last.X, Y: last.Y}
}

func (t *fzippAstarTester) nodeDist(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}

// Neighbours implements the astar.Graph[Node] interface (with Node = image.Point).
func (g fzippGraph) Neighbours(p image.Point) []image.Point {
	offsets := []image.Point{
		image.Pt(0, -1), // North
		image.Pt(1, 0),  // East
		image.Pt(0, 1),  // South
		image.Pt(-1, 0), // West
	}
	res := make([]image.Point, 0, 4)
	for _, off := range offsets {
		q := p.Add(off)
		if g.isFreeAt(q) {
			res = append(res, q)
		}
	}
	return res
}

func (g fzippGraph) isFreeAt(p image.Point) bool {
	return g.isInBounds(p) && g[p.Y][p.X] == ' '
}

func (g fzippGraph) isInBounds(p image.Point) bool {
	return (0 <= p.X && p.X < len(g[0])) &&
		(0 <= p.Y && p.Y < len(g))
}
