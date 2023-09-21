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
