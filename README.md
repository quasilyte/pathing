# quasilyte/pathing

A very efficient & zero-allocation, grid-based, pathfinding library for Go.

## Overview

This library has several things that make it much faster than the alternatives. Some of them are fundamental, and some of them are just a side-effect of the goal I was pursuing for myself.

Some of the limitations you may want to know about before using this library:

1. Its max path length per `BuildPath()` is limited.
2. Only 4 tile kinds per `Grid` are supported.

Both of these limitations can be worked around:

1. Connect the partial results to traverse a bigger map.
2. Use different "layers" for different biomes.

To learn more about this library and its internals, see this presentation: TODO link.

When to use this library?

* You need a very fast pathfinding
* You can live with the limitations listed above

If you answer "yes" to both, consider using this library.

Some games that use this library:

* [Roboden](https://store.steampowered.com/app/2416030/Roboden/)
* [Assemblox](https://itch.io/jam/gmtk-2023/rate/2157747)

## Quick Start

TODO

## Benchmarks & Performance

See [_bench](_bench) folder to reproduce the results.

Benchmark results (as of 13 Sep 2023), time ns/op:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing | 3525 | 2084 | 2688 |
| fzipp/astar | 948367 | 1554290 | 1842812 |
| beefsack/go-astar | 453939 | 939300 | 1032581 |
| s0rg/grid | 1816039 | 1154117 | 1189989 |
| SolarLune/paths | 6588751 | 5158604 | 6114856 |

Allocations/op:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing | 0 | 0 | 9 |
| fzipp/astar | 2008 | 3677 | 3600 |
| beefsack/go-astar | 529 | 1347 | 1557 |
| s0rg/grid | 2976 | 1900 | 1759 |
| SolarLune/paths | 7199 | 6368 | 7001 |

Allocations bytes/op:

| Library | no_wall | simple_wall | multi_wall |
|---|---|---|---|
| quasilyte/pathing | 0 | 0 | 9 |
| fzipp/astar | 337336 | 511908 | 722690 |
| beefsack/go-astar | 43653 | 93122 | 130731 |
| s0rg/grid | 996889 | 551976 | 740523 |
| SolarLune/paths | 235168 | 194768 | 230416 |

I hope that my contribution to this lineup will increase the competition to get better Go gamedev libraries.

Some of my findings that can make these libraries faster:

* Never use `container/heap`; use a generic non-interface version
* Better yet, try bucket priority queue instead of minheap
* Do not use `map`, prefer something that allows a memory re-use
  * The [sparse-dense](https://research.swtch.com/sparse) map could be an option here
* Allocating the result path slice is expensive; consider deltas (2 bits per step)
* Interface method calls are slow for a hot loop

If you want to learn more details, look at my library implementation and/or see TODO talk link.
