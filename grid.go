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

type GridConfig struct {
	WorldWidth  uint
	WorldHeight uint

	CellWidth  uint
	CellHeight uint

	DefaultTile uint8
}

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
	numBytes := numCells / 4
	if numCells%4 != 0 {
		numBytes++
	}
	b := make([]byte, numBytes)

	defaultTileTag := config.DefaultTile
	defaultTileTag &= 0b11
	if defaultTileTag != 0 {
		v := uint8(0)
		switch defaultTileTag {
		case 1:
			v = 0b01010101
		case 2:
			v = 0b10101010
		case 3:
			v = 0b11111111
		}
		for i := range b {
			b[i] = v
		}
	}

	g.bytes = b

	return g
}

func (g *Grid) NumCols() int { return int(g.numCols) }

func (g *Grid) NumRows() int { return int(g.numRows) }

func (g *Grid) SetCellTile(c GridCoord, tileTag uint8) {
	i := uint(c.Y)*g.numCols + uint(c.X)
	byteIndex := i / 4
	if byteIndex < uint(len(g.bytes)) {
		shift := (i % 4) * 2
		b := g.bytes[byteIndex]
		b &^= 0b11 << shift            // Clear the two data bits
		b |= (tileTag & 0b11) << shift // Mix it with provided bits
		g.bytes[byteIndex] = b
	}
}

func (g *Grid) GetCellValue(c GridCoord, l GridLayer) uint8 {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.numCols || y >= g.numRows {
		// Consider out of bound cells as blocked.
		return 0
	}
	return g.getCellValue(x, y, l)
}

func (g *Grid) getCellValue(x, y uint, l GridLayer) uint8 {
	i := y*g.numCols + x
	byteIndex := i / 4
	shift := (i % 4) * 2
	tileTag := ((readByte(g.bytes, byteIndex)) >> shift) & 0b11
	return l.getFast(tileTag)
}

func (g *Grid) AlignPos(x, y float64) (float64, float64) {
	return g.CoordToPos(g.PosToCoord(x, y))
}

func (g *Grid) PosToCoord(x, y float64) GridCoord {
	return GridCoord{
		X: int(x) / g.cellWidth,
		Y: int(y) / g.cellHeight,
	}
}

func (g *Grid) CoordToPos(cell GridCoord) (float64, float64) {
	x := (float64(cell.X) * g.fcellWidth) + g.fcellHalfWidth
	y := (float64(cell.Y) * g.fcellHeight) + g.fcellHalfHeight
	return x, y
}

func (g *Grid) UnpackCoord(v uint32) GridCoord {
	u32 := uint32(v)
	x := int(u32 & 0xffff)
	y := int(u32 >> 16)
	return GridCoord{X: x, Y: y}
}

func (g *Grid) PackCoord(cell GridCoord) uint32 {
	return uint32(cell.X) | uint32(cell.Y<<16)
}
