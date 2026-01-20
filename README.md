# quasilyte/pathing

![Build Status](https://github.com/quasilyte/pathing/workflows/Go/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/quasilyte/pathing)](https://pkg.go.dev/mod/github.com/quasilyte/pathing)

A very fast & zero-allocation, grid-based, pathfinding library for Go.

## Overview

This library has several things that make it much faster than the alternatives. Some of them are fundamental, and some of them are just a side-effect of the goal I was pursuing for myself.

Some of the limitations you may want to know about before using this library:

1. Its max path length per `BuildPath()` is limited (56)
2. Only 8 tile kinds per `Grid` are supported

Both of these limitations can be worked around:

1. Connect the partial results to traverse a bigger map
2. Use different "layers" for different biomes

To learn more about this library and its internals, see [this presentation](https://speakerdeck.com/quasilyte/zero-alloc-pathfinding).

When to use this library?

* You need a very fast pathfinding
* You can live with the limitations listed above

If you answer "yes" to both, consider using this library.

Some games that use this library:

* [Roboden](https://store.steampowered.com/app/2416030/Roboden/)
* [Cavebots](https://quasilyte.itch.io/cavebots)
* [Assemblox](https://itch.io/jam/gmtk-2023/rate/2157747)

## Quick Start

```bash
$ go get github.com/quasilyte/pathing
```

This is a simplified example. See the [full example](example_detailed_test.go) if you want to learn more.

```go
func main() {
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
	// a tile enum values that should fit the 3 bits (max 8 tags per Grid).
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
	normalLayer := pathing.MakeGridLayer([8]uint8{
		tilePlain:    1, // passable
		tileMountain: 0, // not passable
		tileForest:   0, // not passable
	})
	flyingLayer := pathing.MakeGridLayer([8]uint8{
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

	// You can also toggle a "blocked" path bit to make it impossible
	// to be traversed (unless a layer with blocked tile costs is used).
	g.SetCellIsBlocked(pathing.GridCoord{X: 2, Y: 1}, true)

	// This blocked our closest route for the flying unit.
	// Note that it is still known to be a forest tile.
	// m m m m m
	// m A X B m
	// m   f   m
	// m       m
	// m m m m m

	p = bfs.BuildPath(g, startPos, finishPos, flyingLayer)
	fmt.Println(p.Steps.String(), "- after blocking a tile")
}
```

Some terminology hints:

* Grid - a compact matrix that holds the "tile tags"
* GridCoord - an `{X,Y}` object that addresses the grid cell
* Pos - a world coordinate that can be mapped to a grid
* Tile (or a tile tag) - a enum-like value that represents a tile kind
* GridLayer - translates the tag into a pathing value cost (where 0 means "blocked")

Note that it's possible to convert between the GridCoord and world positions via the `Grid` type API.

## Greedy BFS paths quality

This library provides both greedy best-first search as well as A* algorithms.

You may be concerned about the Greedy BFS vs A* results. Due to a couple of tricks I used during the implementation, an unexpected thing happened: some of the paths are actually better than you would expect from a Greedy BFS.

<table>
	<tr>
		<td>A* (21 steps)</td>
		<td>Greedy BFS (27 steps)</td>
		<td>This library's BFS (21 steps)</td>
	<tr>
		<td>
			<img src="https://github.com/quasilyte/pathing/assets/6286655/ba657850-8321-4586-80bd-5e466fa3504c">
		</td>
		<td>
			<img src="https://github.com/quasilyte/pathing/assets/6286655/bef9228a-2b0b-4f6d-a5a3-c676c96149e5">
		</td>
		<td>
			<img src="https://github.com/quasilyte/pathing/assets/6286655/b1da357d-5a9c-40b2-a0d0-e8c6a4bbfdea">
		</td>
	</tr>
</table>

This library worked well enough for me even without A*. You still may want to use A* if you need to have different movement costs for tiles.

In general, A* always build an optimal path and can handle cost-based pathfinding. Greedy BFS requires less memory and works faster.

## Benchmarks & Performance

See [_bench](_bench) folder to reproduce the results.

```bash
# If you're using Linux+Intel processor, consider doing this
# to reduce the noise and make your results more stable:
$ echo "1" | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo

# Running and analyzing the benchmarks:
$ cd _bench
$ go test -bench=. -benchmem -count=10 | results.txt
$ benchstat results.txt
```

Time - **ns/op**:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing BFS | 3525 | 6353 | 16927 |
| quasilyte/pathing A* | 20140 | 35846 | 44756 |
| fzipp/astar | 948367 | 1554290 | 1842812 |
| beefsack/go-astar | 453939 | 939300 | 1032581 |
| kelindar/tile | 107632 ns | 169613 ns | 182342 ns |
| s0rg/grid | 1816039 | 1154117 | 1189989 |
| SolarLune/paths | 6588751 | 5158604 | 6114856 |

Allocations - **allocs/op**:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing BFS | 0 | 0 | 0 |
| quasilyte/pathing A* | 0 | 0 | 0 |
| fzipp/astar | 2008 | 3677 | 3600 |
| beefsack/go-astar | 529 | 1347 | 1557 |
| kelindar/tile | 3 | 3 | 3 |
| s0rg/grid | 2976 | 1900 | 1759 |
| SolarLune/paths | 7199 | 6368 | 7001 |

Allocations -  **bytes/op**:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing BFS | 0 | 0 | 0 |
| quasilyte/pathing A* | 0 | 0 | 0 |
| fzipp/astar | 337336 | 511908 | 722690 |
| beefsack/go-astar | 43653 | 93122 | 130731 |
| tile | 123118 | 32950 | 65763 |
| s0rg/grid | 996889 | 551976 | 740523 |
| SolarLune/paths | 235168 | 194768 | 230416 |

I hope that my contribution to this lineup will increase the competition, so we get better Go gamedev libraries in the future.

Some of my findings that can make these libraries faster:

* Never use `container/heap`; use a generic non-interface version
* Better yet, try a bucket priority queue instead of minheap
* Do not use `map`, prefer something that allows a memory re-use
* The [sparse-dense](https://research.swtch.com/sparse) is a good structure to consider
* The [generations array](https://quasilyte.dev/blog/post/gen-map/) is also a good option
* Allocating the result path slice is expensive; consider deltas (2 bits per step)
* Interface method calls are slow for a hot loop
* Try to be cache-friendly; everything that can be packed should be packed
* Not every game needs A*, don't underestimate the power of a simpler (and faster) algorithm

If you want to learn more details, look at my library implementation and/or see [these slides](https://speakerdeck.com/quasilyte/zero-alloc-pathfinding).
