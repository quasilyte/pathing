# quasilyte/pathing

![Build Status](https://github.com/quasilyte/pathing/workflows/Go/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/quasilyte/pathing)](https://pkg.go.dev/mod/github.com/quasilyte/pathing)

A very fast & zero-allocation, grid-based, pathfinding library for Go.

## Overview

This library has several things that make it much faster than the alternatives. Some of them are fundamental, and some of them are just a side-effect of the goal I was pursuing for myself.

Some of the limitations you may want to know about before using this library:

1. Its max path length per `BuildPath()` is limited
2. Only 4 tile kinds per `Grid` are supported

Both of these limitations can be worked around:

1. Connect the partial results to traverse a bigger map
2. Use different "layers" for different biomes

To learn more about this library and its internals, see this presentation: TODO link.

When to use this library?

* You need a very fast pathfinding
* You can live with the limitations listed above

If you answer "yes" to both, consider using this library.

Some games that use this library:

* [Roboden](https://store.steampowered.com/app/2416030/Roboden/)
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
		// A 5x4 map.
		WorldWidth:  5 * cellSize,
		WorldHeight: 4 * cellSize,
		CellWidth:   cellSize,
		CellHeight:  cellSize,
	})

	// We'll use Greedy BFS pathfinder (A* is also available).
	bfs := pathing.NewGreedyBFS(pathing.GreedyBFSConfig{})

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
	// . . m . . | [m] - mountain
	// . . f . . | [f] - forest
	// . . f . . | [.] - plain
	// . . . . .
	g.SetCellTile(pathing.GridCoord{X: 2, Y: 0}, tileMountain)
	g.SetCellTile(pathing.GridCoord{X: 2, Y: 1}, tileForest)
	g.SetCellTile(pathing.GridCoord{X: 2, Y: 2}, tileForest)

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
	// . . m . . | [m] - mountain
	// . A f B . | [f] - forest
	// . . f . . | [.] - plain
	// . . . . . | [A] - start, [B] - finish
	startPos := pathing.GridCoord{X: 1, Y: 1}
	finishPos := pathing.GridCoord{X: 3, Y: 1}

	// Let's build a normal path first, for a non-flying unit.
	p := bfs.BuildPath(g, startPos, finishPos, normalLayer)

	// The path reads as: Down, Down, Right, Right, Up, Up.
	fmt.Println(p.Steps.String(), "- normal layer path")

	// You can iterate the path.
	for p.Steps.HasNext() {
		fmt.Println("> step:", p.Steps.Next())
	}

	// A flying unit can go in a straight line.
	p = bfs.BuildPath(g, startPos, finishPos, flyingLayer)
	fmt.Println(p.Steps.String(), "- flying layer path")
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
| quasilyte/pathing BFS | 3525 | 2084 | 2688 |
| quasilyte/pathing A* | 20140 | 3415 | 13310 |
| fzipp/astar | 948367 | 1554290 | 1842812 |
| beefsack/go-astar | 453939 | 939300 | 1032581 |
| s0rg/grid | 1816039 | 1154117 | 1189989 |
| SolarLune/paths | 6588751 | 5158604 | 6114856 |

Allocations - **allocs/op**:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing BFS | 0 | 0 | 0 |
| quasilyte/pathing A* | 0 | 0 | 0 |
| fzipp/astar | 2008 | 3677 | 3600 |
| beefsack/go-astar | 529 | 1347 | 1557 |
| s0rg/grid | 2976 | 1900 | 1759 |
| SolarLune/paths | 7199 | 6368 | 7001 |

Allocations -  **bytes/op**:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing BFS | 0 | 0 | 0 |
| quasilyte/pathing A* | 0 | 0 | 0 |
| fzipp/astar | 337336 | 511908 | 722690 |
| beefsack/go-astar | 43653 | 93122 | 130731 |
| s0rg/grid | 996889 | 551976 | 740523 |
| SolarLune/paths | 235168 | 194768 | 230416 |

I hope that my contribution to this lineup will increase the competition, so we get better Go gamedev libraries in the future.

Some of my findings that can make these libraries faster:

* Never use `container/heap`; use a generic non-interface version
* Better yet, try a bucket priority queue instead of minheap
* Do not use `map`, prefer something that allows a memory re-use
* The [sparse-dense](https://research.swtch.com/sparse) is a good structure to consider
* Allocating the result path slice is expensive; consider deltas (2 bits per step)
* Interface method calls are slow for a hot loop
* Try to be cache-friendly; everything that can be packed should be packed
* Not every game needs A*, don't underestimate the power of a simpler (and faster) algorithm

If you want to learn more details, look at my library implementation and/or see TODO talk link.
