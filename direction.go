package pathing

// Direction is a simple enumeration of axial movement directions.
type Direction int

//go:generate stringer -type=Direction -trimprefix=Dir
const (
	DirRight Direction = iota
	DirDown
	DirLeft
	DirUp
	DirNone // A special sentinel value
)

// Reversed returns an opposite direction.
// For instance, DirRight would become DirLeft.
func (d Direction) Reversed() Direction {
	switch d {
	case DirRight:
		return DirLeft
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirUp:
		return DirDown
	default:
		return DirNone
	}
}
