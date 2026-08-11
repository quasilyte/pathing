package pathing

var (
	// UnsetCoord is a conventional "non-existing" coord sentinel.
	// It will only work for games that can't have negative coordinates.
	//
	// The zero value of GridCoord is a valid [0, 0] coordinate,
	// so it can't be used as "undefined" value. This value
	// is declared just for that. And it's exported from the package
	// to allow games have this well-defined shared sentinel.
	UnsetCoord = GridCoord{X: -1, Y: -1}

	// UnsetTinyCoord is like UnsetCoord, but for smaller coords.
	UnsetTinyCoord = TinyGridCoord{X: -1, Y: -1}
)

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

// Add performs a '+' operation and returns the result coordinate.
func (c GridCoord) Add(other GridCoord) GridCoord {
	return GridCoord{X: c.X + other.X, Y: c.Y + other.Y}
}

// Sub performs a '-' operation and returns the result coordinate.
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

func (c GridCoord) Midpoint(other GridCoord) GridCoord {
	return GridCoord{
		X: (c.X + other.X) / 2,
		Y: (c.Y + other.Y) / 2,
	}
}

func (c GridCoord) DirTo(other GridCoord) Direction {
	return dirTo(c.X, c.Y, other.X, other.Y)
}

func (c GridCoord) ToTinyCoord() TinyGridCoord {
	return TinyGridCoord{
		X: int8(c.X),
		Y: int8(c.Y),
	}
}

// TinyGridCoord is a compact GridCoord representation
// that is only capable of storing small values that fit in int8 X/Y.
// It can be expanded into GridCoord and vice versa.
type TinyGridCoord struct {
	X int8
	Y int8
}

func (c TinyGridCoord) ToCoord() GridCoord {
	return GridCoord{
		X: int(c.X),
		Y: int(c.Y),
	}
}

// IsZero reports whether the coord is {0, 0}.
func (c TinyGridCoord) IsZero() bool {
	return c.X == 0 && c.Y == 0
}

// Add performs a '+' operation and returns the result coordinate.
func (c TinyGridCoord) Add(other TinyGridCoord) TinyGridCoord {
	return TinyGridCoord{X: c.X + other.X, Y: c.Y + other.Y}
}

// Sub performs a '-' operation and returns the result coordinate.
func (c TinyGridCoord) Sub(other TinyGridCoord) TinyGridCoord {
	return TinyGridCoord{X: c.X - other.X, Y: c.Y - other.Y}
}

func (c TinyGridCoord) Move(d Direction) TinyGridCoord {
	return c.ToCoord().Move(d).ToTinyCoord()
}

// Dist finds a Manhattan distance between the two coordinates.
func (c TinyGridCoord) Dist(other TinyGridCoord) int {
	return int(int8abs(c.X-other.X)) + int(int8abs(c.Y-other.Y))
}

func (c TinyGridCoord) Midpoint(other TinyGridCoord) TinyGridCoord {
	return TinyGridCoord{
		X: (c.X + other.X) / 2,
		Y: (c.Y + other.Y) / 2,
	}
}

func (c TinyGridCoord) DirTo(other TinyGridCoord) Direction {
	return dirTo(int(c.X), int(c.Y), int(other.X), int(other.Y))
}

func int8abs(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

func intabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func dirTo(fromX, fromY, toX, toY int) Direction {
	dx := toX - fromX
	dy := toY - fromY

	if dx == 0 && dy == 0 {
		return DirNone
	}
	if intabs(dx) >= intabs(dy) {
		if dx > 0 {
			return DirRight
		}
		return DirLeft
	}
	if dy > 0 {
		return DirDown
	}
	return DirUp
}
