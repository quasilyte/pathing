package pathing_test

import (
	"fmt"

	"github.com/quasilyte/pathing"
)

func Example() {
	// Grid is a "map" that stores cell info.
	const cellSize = 40
	g := pathing.NewGrid(pathing.GridConfig{
		// A 5x5 map.
		WorldWidth:  5 * cellSize,
		WorldHeight: 5 * cellSize,
		CellWidth:   cellSize,
		CellHeight:  cellSize,
	})

	// We'll use Greedy BFS pathfinder.
	// Re-use it, don't create a new BFS every time.
	bfs := pathing.NewGreedyBFS(pathing.GreedyBFSConfig{
		NumCols: uint(g.NumCols()),
		NumRows: uint(g.NumRows()),
	})

	// Tile kinds are needed to interpret the cell values.
	// Let's define some.
	const (
		tilePlain = iota
		tileForest
		tileMountain
	)

	// Grid map cells contain "tile tags"; these are basically
	// a tile enum values that should fit the 2 bits (max 4 tags per Grid).
	// The default tag is 0 (tilePlain).
	// Let's add some forests and mountains.
	//
	// The result map layout will look like this:
	// m m m m m | [m] - mountain
	// m   f   m | [f] - forest
	// m   f   m | [ ] - plain
	// m       m
	// m m m m m
	g.SetCellTile(pathing.GridCoord{X: 2, Y: 1}, tileForest)
	g.SetCellTile(pathing.GridCoord{X: 2, Y: 2}, tileForest)
	for y := 0; y < g.NumRows(); y++ {
		for x := 0; x < g.NumCols(); x++ {
			if !(y == 0 || y == g.NumRows()-1 || x == 0 || x == g.NumCols()-1) {
				continue
			}
			g.SetCellTile(pathing.GridCoord{X: x, Y: y}, tileMountain)
		}
	}

	// Now we need to tell the pathfinding library how to interpret
	// these tiles. For instance, which tiles are passable and not.
	// We do that by using layers. I'll define two layers here
	// to show you how it's possible to interpret the grid differently
	// depending on the layer.
	normalLayer := pathing.MakeGridLayer([4]uint8{
		tilePlain:    1, // passable
		tileMountain: 0, // not passable
		tileForest:   0, // not passable
	})
	flyingLayer := pathing.MakeGridLayer([4]uint8{
		tilePlain:    1,
		tileMountain: 1,
		tileForest:   1,
	})

	// Our map with markers will look like this:
	// m m m m m | [m] - mountain
	// m A f B m | [f] - forest
	// m   f   m | [ ] - plain
	// m       m | [A] - start
	// m m m m m | [B] - finish
	startPos := pathing.GridCoord{X: 1, Y: 1}
	finishPos := pathing.GridCoord{X: 3, Y: 1}

	// Let's build a normal path first, for a non-flying unit.
	p := bfs.BuildPath(g, startPos, finishPos, normalLayer)
	fmt.Println(p.Steps.String(), "- normal layer path")

	// You can iterate the path.
	for p.Steps.HasNext() {
		fmt.Println("> step:", p.Steps.Next())
	}

	// A flying unit can go in a straight line.
	p = bfs.BuildPath(g, startPos, finishPos, flyingLayer)
	fmt.Println(p.Steps.String(), "- flying layer path")

	// A path building result has some extra information bits you might be interested in.
	// Usually, you only need the Steps part, so you can pass it around instead of the
	// entire result object
	fmt.Println(p.Finish, p.Partial)

	// Output:
	// {Down,Down,Right,Right,Up,Up} - normal layer path
	// > step: Down
	// > step: Down
	// > step: Right
	// > step: Right
	// > step: Up
	// > step: Up
	// {Right,Right} - flying layer path
	// {3 1} false
}
