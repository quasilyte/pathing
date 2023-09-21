package pathing

// GridCoord represents a grid-local coordinate.
// You can translate it to a world coordinate using a grid.
//
// If the grid cell size is 32x32, then this table can explain the mapping:
//
//   - pos{0, 0}   => coord{0, 0}
//   - pos{16, 16} => coord{0, 0}
//   - pos{20, 20} => coord{0, 0}
//   - pos{35, 10} => coord{1, 0}
//   - pos{50, 50} => coord{1, 1}
//   - pos{90, 90} => coord{2, 2}
type GridCoord struct {
	X int
	Y int
}

// IsZero reports whether the coord is {0, 0}.
func (c GridCoord) IsZero() bool {
	return c.X == 0 && c.Y == 0
}

// Add performs a + operation and returns the result coordinate.
func (c GridCoord) Add(other GridCoord) GridCoord {
	return GridCoord{X: c.X + other.X, Y: c.Y + other.Y}
}

// Sub performs a - operation and returns the result coordinate.
func (c GridCoord) Sub(other GridCoord) GridCoord {
	return GridCoord{X: c.X - other.X, Y: c.Y - other.Y}
}

func (c GridCoord) reversedMove(d Direction) GridCoord {
	switch d {
	case DirRight:
		return GridCoord{X: c.X - 1, Y: c.Y}
	case DirDown:
		return GridCoord{X: c.X, Y: c.Y - 1}
	case DirLeft:
		return GridCoord{X: c.X + 1, Y: c.Y}
	case DirUp:
		return GridCoord{X: c.X, Y: c.Y + 1}
	default:
		return c
	}
}

// Move translates the coordinate one step towards the direction.
//
// Note that the coordinates are not validated.
// It's possible to get an out-of-bounds coordinate that
// will not belong to a Grid.
//
//   - {2,2}.Move(DirLeft) would give {1,2}
//   - {2,2}.Move(DirDown) would give {2,3}
func (c GridCoord) Move(d Direction) GridCoord {
	switch d {
	case DirRight:
		return GridCoord{X: c.X + 1, Y: c.Y}
	case DirDown:
		return GridCoord{X: c.X, Y: c.Y + 1}
	case DirLeft:
		return GridCoord{X: c.X - 1, Y: c.Y}
	case DirUp:
		return GridCoord{X: c.X, Y: c.Y - 1}
	default:
		return c
	}
}

// Dist finds a Manhattan distance between the two coordinates.
func (c GridCoord) Dist(other GridCoord) int {
	return intabs(c.X-other.X) + intabs(c.Y-other.Y)
}

func intabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
