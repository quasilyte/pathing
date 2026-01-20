package pathing

type Grid struct {
	numCols uint
	numRows uint

	bytes []byte

	cellWidth  int
	cellHeight int

	fcellWidth      float64
	fcellHeight     float64
	fcellHalfWidth  float64
	fcellHalfHeight float64
}

// GridConfig is a NewGrid() function parameter.
// See field comments for more details.
type GridConfig struct {
	// WorldWidth and WorldHeight specify the modelled world size.
	// These are used in combination with cell sizes to map positions
	// and grid coordinates.
	// The world size is specified in pixels.
	// Although the positions are expected to be a pair of float64,
	// the world size is a pair of uints, because sizes like 200.5 make no sense.
	WorldWidth  uint
	WorldHeight uint

	// CellWidth and CellHeight specify the grid cell size.
	// If the world size is 320x320 and the cell size is 32x32,
	// that would mean that there are 10x10 (100) cells in total.
	//
	// If left unset (0), the default size will be used (32x32).
	CellWidth  uint
	CellHeight uint

	// DefaultTile controls the default grid fill.
	// Only the 3 lower bits matter as a tile tag value can't exceed a value of 7.
	// This value is a minor option, but it can be used to populate the grid
	// with the most common tile.
	// Although it does fill the grid in an optimized way, it's mostly a convenience option
	// to make the initialization easier.
	DefaultTile uint8
}

// NewGrid creates a Grid object.
// See GridConfig comment to learn more.
func NewGrid(config GridConfig) *Grid {
	if config.CellWidth == 0 {
		config.CellWidth = 32
	}
	if config.CellHeight == 0 {
		config.CellHeight = 32
	}

	g := &Grid{
		cellWidth:  int(config.CellWidth),
		cellHeight: int(config.CellHeight),

		fcellWidth:  float64(config.CellWidth),
		fcellHeight: float64(config.CellHeight),
	}

	g.fcellHalfWidth = float64(config.CellWidth / 2)
	g.fcellHalfHeight = float64(config.CellHeight / 2)

	g.numCols = config.WorldWidth / config.CellWidth
	g.numRows = config.WorldHeight / config.CellHeight

	numCells := g.numCols * g.numRows
	numBytes := numCells / 2
	if numCells%2 != 0 {
		numBytes++
	}
	b := make([]byte, numBytes)

	defaultTileTag := config.DefaultTile
	defaultTileTag &= 0b111
	if defaultTileTag != 0 {
		v := uint8(0)
		switch defaultTileTag {
		case 1:
			v = 0b0001_0001
		case 2:
			v = 0b0010_0010
		case 3:
			v = 0b0011_0011
		case 4:
			v = 0b0100_0100
		case 5:
			v = 0b0101_0101
		case 6:
			v = 0b0110_0110
		case 7:
			v = 0b0111_0111
		}
		for i := range b {
			b[i] = v
		}
	}

	g.bytes = b

	return g
}

// NumCols returns the number of columns this grid has.
func (g *Grid) NumCols() int { return int(g.numCols) }

// NumRows returns the number of rows this grid has.
func (g *Grid) NumRows() int { return int(g.numRows) }

// SetCellTile assigns the tile tag for the given cell coordinate.
//
// You usually do not want to change the tile types after the grid
// is filled, but if your map is dynamic, it is OK to do so
// as it is an O(1) operation.
//
// For dynamic info like "tile is blocked" use the separate bit
// accessible through some of the APIs (e.g. SetCellIsBlocked, GetCellTile2, GetCellIsBlocked)
func (g *Grid) SetCellTile(c GridCoord, tileTag uint8) {
	i := uint(c.Y)*g.numCols + uint(c.X)
	byteIndex := i / 2
	if byteIndex < uint(len(g.bytes)) {
		shift := (i % 2) * 4
		b := g.bytes[byteIndex]
		b &^= 0b0111 << shift            // Clear the lower 3 data bits
		b |= (tileTag & 0b0111) << shift // Mix it with provided bits
		g.bytes[byteIndex] = b
	}
}

// SetCellIsBlocked writes to a special tile "blocked" bit (0 or 1).
// You can read that bit with GetCellIsBlocked, it is also reported with GetCellTile2.
// Blocked cells usually can not be traversed, but it can be changed with the layer
// with a dedicated blocked tile map.
func (g *Grid) SetCellIsBlocked(c GridCoord, blocked bool) {
	bit := uint8(0)
	if blocked {
		bit = 0b1000
	}
	i := uint(c.Y)*g.numCols + uint(c.X)
	byteIndex := i / 2
	if byteIndex < uint(len(g.bytes)) {
		shift := (i % 2) * 4
		b := g.bytes[byteIndex]
		b &^= 0b1000 << shift // Clear the bit
		b |= bit << shift     // Mix it in
		g.bytes[byteIndex] = b
	}
}

// GetCellTile returns the cell tile tag.
// This operation is only useful for the Grid debugging as
// for the pathfinding tasks you would want to use GetCellCost() method instead.
//
// An out-of-bounds access returns 0.
//
// This function never reports whether a tile is blocked or not,
// it only returns the tag associated with it.
func (g *Grid) GetCellTile(c GridCoord) uint8 {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.numCols || y >= g.numRows {
		return 0
	}
	i := y*g.numCols + x
	byteIndex := i / 2
	shift := (i % 2) * 4
	return ((readByte(g.bytes, byteIndex)) >> shift) & 0b111
}

// GetCellTile2 is like GetCellTile, but also reports
// whether the tile was marked as blocked (occupied).
//
// An out-of-bounds access is not marked as occupied.
func (g *Grid) GetCellTile2(c GridCoord) (uint8, bool) {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.numCols || y >= g.numRows {
		return 0, false
	}
	i := y*g.numCols + x
	byteIndex := i / 2
	shift := (i % 2) * 4
	bits := ((readByte(g.bytes, byteIndex)) >> shift)
	return bits & 0b111, bits&0b1000 != 0
}

// GetCellIsBlocked reports whether the given cell was marked as blocked.
// By default, no cells are marked as such, until SetCellIsBlocked is called.
//
// An out-of-bounds cell is never blocked.
func (g *Grid) GetCellIsBlocked(c GridCoord) bool {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.numCols || y >= g.numRows {
		return false
	}
	i := y*g.numCols + x
	byteIndex := i / 2
	shift := (i % 2) * 4
	return ((readByte(g.bytes, byteIndex))>>shift)&0b1000 != 0
}

// GetCellCost returns a travelling cost for a given cell as specified in the layer.
// The return value interpreted as this: 0 is a blocked path while any other value
// is a travelling cost.
//
// An out-of-bounds access returns 0 (interpreted as blocked).
func (g *Grid) GetCellCost(c GridCoord, l GridLayer) uint8 {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.numCols || y >= g.numRows {
		// Consider out of bound cells as blocked.
		return 0
	}
	return g.getCellCost(x, y, l)
}

func (g *Grid) getCellCost(x, y uint, l GridLayer) uint8 {
	i := y*g.numCols + x
	byteIndex := i / 2
	shift := (i % 2) * 4
	tileTag := ((readByte(g.bytes, byteIndex)) >> shift) & 0b1111
	return l.getFast(tileTag)
}

// AlignPos is an easy way to center the world position inside a grid cell.
// For instance, with a cell size of 32x32, a {10,10} pos would become {16,16}.
func (g *Grid) AlignPos(x, y float64) (float64, float64) {
	return g.CoordToPos(g.PosToCoord(x, y))
}

// PosToCoord converts a world position into a grid coord.
func (g *Grid) PosToCoord(x, y float64) GridCoord {
	return GridCoord{
		X: int(x) / g.cellWidth,
		Y: int(y) / g.cellHeight,
	}
}

// CoordToPos converts a grid coord into a world position.
func (g *Grid) CoordToPos(c GridCoord) (float64, float64) {
	x := (float64(c.X) * g.fcellWidth) + g.fcellHalfWidth
	y := (float64(c.Y) * g.fcellHeight) + g.fcellHalfHeight
	return x, y
}

// PackCoord returns a packed version of a grid coordinate.
// It can be useful to get an efficient map key.
// A packed coordinate can later be unpacked with UnpackCoord() method.
func (g *Grid) PackCoord(c GridCoord) uint32 {
	return uint32(c.X) | uint32(c.Y<<16)
}

// UnpackCoord takes a packed coord and returns its unpacked version.
func (g *Grid) UnpackCoord(v uint32) GridCoord {
	u32 := uint32(v)
	x := int(u32 & 0xffff)
	y := int(u32 >> 16)
	return GridCoord{X: x, Y: y}
}
