package pathing

import (
	"strings"
)

// GridPath represents a constructed path from point A to point B.
//
// Instead of storing the actual coordinates, it stores deltas in a form of step directions.
//
// You get from point A to point B by taking the steps into the directions
// specified by the path.
// The path object is essentialy an iterator.
//
// The path can be copied by simply assigning it, it has a value semantics.
// You want to pass it around as a value 90% of time,
// but if you want some function to be able to affect the iterator state,
// pass it by the pointer.
type GridPath struct {
	bytes [gridPathBytes]byte
	len   byte
	pos   byte
}

// MakeGridPath construct a path from the given set of steps.
func MakeGridPath(steps ...Direction) GridPath {
	var result GridPath
	for i := len(steps) - 1; i >= 0; i-- {
		result.push(steps[i])
	}
	result.Rewind()
	return result
}

// String returns a debug-print version of the path.
// It's not intended to be used a fast path-to-string method.
func (p GridPath) String() string {
	parts := make([]string, 0, p.len)
	prevPos := p.pos // Restore the pos later
	p.Rewind()
	for p.HasNext() {
		parts = append(parts, p.Next().String())
	}
	p.pos = prevPos
	return "{" + strings.Join(parts, ",") + "}"
}

// Len returns the path length.
// It's not affected by the iterator state; the result is always
// a total path length regardless of the progress.
func (p *GridPath) Len() int {
	return int(p.len)
}

// HasNext reports whether there are more steps inside this path.
// Use Next() to extract the next path segment if there are any.
func (p *GridPath) HasNext() bool {
	return p.pos != 0
}

// Rewind resets the iterator and allows you to traverse it again.
func (p *GridPath) Rewind() {
	p.pos = p.len
}

// Peek returns the next path step without advancing the iterator.
func (p *GridPath) Peek() Direction {
	return p.get(p.pos - 1)
}

// Next returns the next path step and advances the iterator.
func (p *GridPath) Next() Direction {
	d := p.Peek()
	p.pos--
	return d
}

// Skip consumes n next path steps.
func (p *GridPath) Skip(n byte) {
	p.pos -= n
}

// Peek2 is like Peek(), but it returns two next steps instead of just one.
func (p *GridPath) Peek2() (Direction, Direction) {
	// If p.pos is 1, p.pos-2 overflows to 255.
	// byteIndex will not be inside len(p.bytes), so
	// p.get(p.pos-2) will return DirNone as it should.
	// No need to check for that condition here explicitely.
	return p.get(p.pos - 1), p.get(p.pos - 2)
}

func (p *GridPath) push(dir Direction) {
	i := p.pos
	p.pos++
	p.len++
	byteIndex := i / 4
	bitShift := (i % 4) * 2
	if byteIndex < uint8(len(p.bytes)) {
		p.bytes[byteIndex] |= byte(dir << bitShift)
	}
}

func (p *GridPath) get(i byte) Direction {
	byteIndex := i / 4
	bitShift := (i % 4) * 2
	if byteIndex < uint8(len(p.bytes)) {
		return Direction((p.bytes[byteIndex] >> bitShift) & 0b11)
	}
	return DirNone
}

func constructPath(from, to GridCoord, pathmap *coordMap) GridPath {
	// We walk from the finish point towards the start.
	// The directions are pushed in that order and would lead
	// to a reversed path, but since GridPath does its iteration
	// in reversed order itself, we don't need to do any
	// post-build reversal here.
	var result GridPath
	pos := to
	for {
		d, _ := pathmap.Get(pathmap.packCoord(pos))
		if pos == from {
			break
		}
		result.push(Direction(d))
		pos = pos.reversedMove(Direction(d))
	}
	return result
}

func findPathOrigin(from GridCoord) GridCoord {
	origin := GridCoord{}
	if originX := from.X - gridPathMaxLen; originX > 0 {
		origin.X = originX
	}
	if originY := from.Y - gridPathMaxLen; originY > 0 {
		origin.Y = originY
	}
	return origin
}
