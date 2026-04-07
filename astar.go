package pathing

// AStar implements an A* search pathfinding algorithm.
// You must use NewAStar() function to obtain an instance of this type.
//
// AStar is a bit slower than GreedyBFS, but its results can be more optimal.
// It also supports a proper weight/cost based pathfinding.
//
// Once created, you should re-use it to build paths.
// Do not throw the instance away after building the path once.
type AStar struct {
	frontier *minheap[astarCoord]
	costmap  *coordMap
	pathmap  *coordMap
}

type AStarConfig struct {
	// NumCols and NumRows are size hints for the AStar constructor.
	// Grid.NumCols() and Grid.NumRows() methods will come in handy to initialize these.
	// If you keep them at 0, the max amount of the working space will be allocated.
	// It's like a size hint: the constructor may allocate a smaller working area
	// if the grids you're going operate on are small.
	NumCols uint
	NumRows uint
}

type astarCoord struct {
	Coord  GridCoord
	Weight int32
	Cost   int32
}

// NewAStar creates a ready-to-use AStar object.
func NewAStar(config AStarConfig) *AStar {
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

	astar := &AStar{
		frontier: newMinheap[astarCoord](32),
		pathmap:  newCoordMap(coordMapCols, coordMapRows),
		costmap:  newCoordMap(coordMapCols, coordMapRows),
	}

	return astar
}

// WalkCosts iterates over all reachable tiles, going from the given pos.
// It's like pathfinding, but in all directions, but without the need
// to actually build any paths - therefore it is much faster than building
// paths from pos to every coord in the area to calculate tile costs around the pos.
//
// The provided f callback is called for every such reachable tile, given
// the current coord and the cost associated with a path from pos to that coord.
// The function is not called for the pos itself.
// The return value of "true" means the value was found and the iteration
// will be stopped.
//
// There are many use cases for this, one being showing the reachable area
// for the given movement budget.
func (astar *AStar) WalkCosts(g *Grid, pos GridCoord, l GridLayer, maxCost int, f func(c GridCoord, cost int) bool) {
	// Almost identical to the BuildPath, with a few key differences.
	// It doesn't need a pathmap as paths are never reconstructed.
	// It can't stop as soon as it finds the solution, as it needs as many of them
	// as possible - but the user can stop the process by returning true,
	// so even with a large maxCost it should terminate easily.

	origin := findPathOrigin(pos)

	localStart := pos.Sub(origin)

	frontier := astar.frontier
	frontier.Reset()

	costmap := astar.costmap
	costmap.Reset()

	frontier.Push(0, astarCoord{Coord: localStart, Cost: 0})

	for !frontier.IsEmpty() {
		current := frontier.Pop()

		// Accept the tile path: the first pop wins, subsequent pops for the
		// same tile are stale duplicates produced by the lazy-deletion
		// approach and must be ignored.
		k := costmap.packCoord(current.Coord)
		if costmap.Contains(k) {
			continue
		}
		costmap.Set(k, uint32(current.Cost))

		// For now we have an explicit start coord check here,
		// but ideally this coord should be ignored as a side effect
		// of the algorithm (just need to figure out how).
		if current.Coord != localStart {
			if f(current.Coord.Add(origin), int(current.Cost)) {
				break
			}
		}

		for _, offset := range &neighborOffsets {
			next := current.Coord.Add(offset)

			cx := uint(next.X) + uint(origin.X)
			cy := uint(next.Y) + uint(origin.Y)
			if cx >= g.numCols || cy >= g.numRows {
				continue
			}

			nextCellCost := g.getCellCost(cx, cy, l)
			if nextCellCost == 0 {
				continue
			}

			newCost := int(current.Cost) + int(nextCellCost)
			if newCost > maxCost {
				continue
			}

			k := costmap.packCoord(next)
			if costmap.Contains(k) {
				continue
			}

			frontier.Push(newCost, astarCoord{Coord: next, Cost: int32(newCost)})
		}
	}
}

// BuildPath attempts to find a path between the two coordinates.
// It will use a provided Grid in combination with a GridLayer.
// The Grid is expected to store the tile tags and the GridLayer is
// used to interpret these tags.
func (astar *AStar) BuildPath(g *Grid, from, to GridCoord, l GridLayer) BuildPathResult {
	var result BuildPathResult
	if from == to {
		result.Finish = to
		return result
	}

	origin := findPathOrigin(from)

	localStart := from.Sub(origin)
	localGoal := to.Sub(origin)

	frontier := astar.frontier
	frontier.Reset()

	pathmap := astar.pathmap
	pathmap.Reset()

	costmap := astar.costmap
	costmap.Reset()

	frontier.Push(0, astarCoord{Coord: localStart})

	shortestDist := 0xffffffff
	var fallbackCoord GridCoord
	var fallbackCost int
	foundPath := false
	for !frontier.IsEmpty() {
		current := frontier.Pop()

		if current.Coord == localGoal {
			result.Steps = constructPath(localStart, localGoal, pathmap)
			result.Finish = to
			result.Cost = int(current.Cost)
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
			fallbackCost = int(current.Cost)
		}

		currentCost, _ := costmap.Get(costmap.packCoord(current.Coord))
		for dir, offset := range &neighborOffsets {
			next := current.Coord.Add(offset)
			cx := uint(next.X) + uint(origin.X)
			cy := uint(next.Y) + uint(origin.Y)
			if cx >= g.numCols || cy >= g.numRows {
				continue
			}
			nextCellCost := g.getCellCost(cx, cy, l)
			if nextCellCost == 0 {
				continue
			}
			newNextCost := currentCost + uint32(nextCellCost)
			k := costmap.packCoord(next)
			oldNextCost, ok := costmap.Get(k)
			if ok && newNextCost >= oldNextCost {
				continue
			}
			costmap.Set(k, newNextCost)
			priority := newNextCost + uint32(localGoal.Dist(next))
			nextWeighted := astarCoord{
				Coord:  next,
				Cost:   int32(newNextCost),
				Weight: int32(current.Weight + 1),
			}
			frontier.Push(int(priority), nextWeighted)
			pathmap.Set(k, uint32(dir))
		}
	}

	if !foundPath {
		result.Steps = constructPath(localStart, fallbackCoord, pathmap)
		result.Finish = fallbackCoord.Add(origin)
		result.Cost = fallbackCost
		result.Partial = true
	}

	return result
}
