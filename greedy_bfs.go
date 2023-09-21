package pathing

var neighborOffsets = [4]GridCoord{
	{X: 1},
	{Y: 1},
	{X: -1},
	{Y: -1},
}

// GreedyBFS implements a greedy best-first search pathfinding algorithm.
// You must use NewGreedyBFS() function to obtain an instance of this type.
//
// GreedyBFS is a faster pathfinder with a lower memory costs as compared to an AStar.
// It can't handle different movements costs though, so it will treat any non-zero
// value returned by GridLayer identically.
//
// Once created, you should re-use it to build paths.
// Do not throw the instance away after building the path once.
type GreedyBFS struct {
	pqueue     *priorityQueue[weightedGridCoord]
	coordSlice []weightedGridCoord
	coordMap   *coordMap
}

// BuildPathResult is a BuildPath() method return value.
type BuildPathResult struct {
	// Steps is an actual path that was constructed.
	Steps GridPath

	// Finish is where the constructed path ends.
	// It's mostly needed in case of a partial result,
	// since you can build another path from this coord right away.
	Finish GridCoord

	// Whether this is a partial path result.
	// This happens if the destination can't be reached
	// or if it's too far away.
	Partial bool
}

type weightedGridCoord struct {
	Coord  GridCoord
	Weight int
}

type GreedyBFSConfig struct {
	// NumCols and NumRows are size hints for the GreedyBFS constructor.
	// Grid.NumCols() and Grid.NumRows() methods will come in handy to initialize these.
	// If you keep them at 0, the max amount of the working space will be allocated.
	// It's like a size hint: the constructor may allocate a smaller working area
	// if the grids you're going operate on are small.
	NumCols uint
	NumRows uint
}

// NewGreedyBFS creates a ready-to-use GreedyBFS object.
//
// numCols and numRows should be the same as in the Grid
// objects that are going to be used with this pathfinder.
// Grid.NumCols() and Grid.NumRows() methods will come in handy.
//
// Note that you can use different grids with the same
// pathfinder, but they should be of the same size.
func NewGreedyBFS(config GreedyBFSConfig) *GreedyBFS {
	if config.NumCols == 0 {
		config.NumCols = gridMapSide
	}
	if config.NumRows == 0 {
		config.NumRows = gridMapSide
	}

	coordMapCols := gridMapSide
	if int(config.NumCols) < coordMapCols {
		coordMapCols = int(config.NumCols)
	}
	coordMapRows := gridMapSide
	if int(config.NumRows) < coordMapRows {
		coordMapRows = int(config.NumRows)
	}

	bfs := &GreedyBFS{
		pqueue:     newPriorityQueue[weightedGridCoord](),
		coordMap:   newCoordMap(coordMapCols, coordMapRows),
		coordSlice: make([]weightedGridCoord, 0, 40),
	}

	return bfs
}

// BuildPath attempts to find a path between the two coordinates.
// It will use a provided Grid in combination with a GridLayer.
// The Grid is expected to store the tile tags and the GridLayer is
// used to interpret these tags.
func (bfs *GreedyBFS) BuildPath(g *Grid, from, to GridCoord, l GridLayer) BuildPathResult {
	var result BuildPathResult
	if from == to {
		result.Finish = to
		return result
	}

	// Find a search box origin pos. We need these to translate the local coordinates later.
	origin := GridCoord{}
	if originX := from.X - gridPathMaxLen; originX > 0 {
		origin.X = originX
	}
	if originY := from.Y - gridPathMaxLen; originY > 0 {
		origin.Y = originY
	}

	// These will be in local coordinates.
	localStart := from.Sub(origin)
	localGoal := to.Sub(origin)

	frontier := bfs.pqueue
	frontier.Reset()

	hotFrontier := bfs.coordSlice[:0]
	hotFrontier = append(hotFrontier, weightedGridCoord{Coord: localStart})

	pathmap := bfs.coordMap
	pathmap.Reset()

	shortestDist := 0xff
	var fallbackCoord GridCoord
	foundPath := false
	for len(hotFrontier) != 0 || !frontier.IsEmpty() {
		var current weightedGridCoord
		if len(hotFrontier) != 0 {
			current = hotFrontier[len(hotFrontier)-1]
			hotFrontier = hotFrontier[:len(hotFrontier)-1]
		} else {
			current = frontier.Pop()
		}

		if current.Coord == localGoal {
			result.Steps = bfs.constructPath(localStart, localGoal, pathmap)
			result.Finish = to
			foundPath = true
			break
		}
		if current.Weight > gridPathMaxLen {
			break
		}

		dist := localGoal.Dist(current.Coord)
		if dist < shortestDist {
			shortestDist = dist
			fallbackCoord = current.Coord
		}
		for dir, offset := range &neighborOffsets {
			next := current.Coord.Add(offset)
			cx := uint(next.X) + uint(origin.X)
			cy := uint(next.Y) + uint(origin.Y)
			if cx >= g.numCols || cy >= g.numRows {
				continue
			}
			if g.getCellCost(cx, cy, l) == 0 {
				continue
			}
			pathmapKey := pathmap.packCoord(next)
			if pathmap.Get(pathmapKey) != DirNone {
				continue
			}
			pathmap.Set(pathmapKey, Direction(dir))
			nextDist := localGoal.Dist(next)
			nextWeighted := weightedGridCoord{
				Coord: next,
				// This is used to determine the out-of-scope coordinates.
				// It's not a distance score; therefore, we're not using nextDist here.
				Weight: current.Weight + 1,
			}
			if nextDist < dist {
				hotFrontier = append(hotFrontier, nextWeighted)
			} else {
				frontier.Push(nextDist, nextWeighted)
			}
		}
	}

	if !foundPath {
		result.Steps = bfs.constructPath(localStart, fallbackCoord, pathmap)
		result.Finish = fallbackCoord.Add(origin)
		result.Partial = true
	}

	// In case if that slice was growing due to appends,
	// save that extra capacity for later.
	bfs.coordSlice = hotFrontier[:0]

	return result
}

func (bfs *GreedyBFS) constructPath(from, to GridCoord, pathmap *coordMap) GridPath {
	// We walk from the finish point towards the start.
	// The directions are pushed in that order and would lead
	// to a reversed path, but since GridPath does its iteration
	// in reversed order itself, we don't need to do any
	// post-build reversal here.
	var result GridPath
	pos := to
	for {
		d := pathmap.Get(pathmap.packCoord(pos))
		if pos == from {
			break
		}
		result.push(d)
		pos = pos.reversedMove(d)
	}
	return result
}
