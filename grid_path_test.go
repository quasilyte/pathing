package pathing_test

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/quasilyte/pathing"
)

func testParsePath(s string) pathing.GridPath {
	s = s[1 : len(s)-1] // Drop "{}"
	if s == "" {
		return pathing.GridPath{}
	}
	var directions []pathing.Direction
	for _, part := range strings.Split(s, ",") {
		switch part {
		case "Right":
			directions = append(directions, pathing.DirRight)
		case "Down":
			directions = append(directions, pathing.DirDown)
		case "Left":
			directions = append(directions, pathing.DirLeft)
		case "Up":
			directions = append(directions, pathing.DirUp)
		default:
			panic("unexpected part: " + part)
		}
	}
	return pathing.MakeGridPath(directions...)
}

func TestGridPathString(t *testing.T) {
	tests := []string{
		"{}",
		"{Left}",
		"{Left,Right}",
		"{Right,Left}",
		"{Down,Down,Down,Up}",
		"{Left,Right,Up,Down}",
		"{Left,Right,Right,Right,Left}",
		"{Up,Up,Down,Down,Left,Left,Right,Right,Down,Down}",
		"{Up,Up,Down,Down,Left,Left,Right,Right,Down,Down,Down,Left,Up,Right}",
		"{Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left}",
		"{Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Right}",
		"{Up,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left,Left}",
	}

	for _, test := range tests {
		p := testParsePath(test)
		if p.String() != test {
			t.Fatalf("results mismatched:\nhave: %q\nwant: %q", p.String(), test)
		}
	}
}

func TestGridPath(t *testing.T) {
	tests := [][]pathing.Direction{
		{},
		{pathing.DirLeft},
		{pathing.DirDown},
		{pathing.DirLeft, pathing.DirRight, pathing.DirUp},
		{pathing.DirLeft, pathing.DirLeft, pathing.DirLeft},
		{pathing.DirDown, pathing.DirDown, pathing.DirDown},
		{pathing.DirDown, pathing.DirUp, pathing.DirLeft, pathing.DirRight, pathing.DirLeft, pathing.DirRight},
		{pathing.DirDown, pathing.DirLeft, pathing.DirLeft, pathing.DirLeft, pathing.DirLeft, pathing.DirDown},
		{pathing.DirRight, pathing.DirRight, pathing.DirRight, pathing.DirRight, pathing.DirRight, pathing.DirRight, pathing.DirRight},
		{pathing.DirDown, pathing.DirRight, pathing.DirRight, pathing.DirDown, pathing.DirRight, pathing.DirUp, pathing.DirRight, pathing.DirLeft},
	}

	for i, directions := range tests {
		p := pathing.MakeGridPath(directions...)
		reconstructed := []pathing.Direction{}
		for p.HasNext() {
			reconstructed = append(reconstructed, p.Next())
		}
		if !reflect.DeepEqual(directions, reconstructed) {
			t.Fatalf("test%d paths mismatch", i)
		}
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 100; i++ {
		size := r.Intn(20) + 10
		directions := []pathing.Direction{}
		for j := 0; j < size; j++ {
			d := r.Intn(4)
			directions = append(directions, pathing.Direction(d))
		}
		p := pathing.MakeGridPath(directions...)
		reconstructed := []pathing.Direction{}
		for p.HasNext() {
			reconstructed = append(reconstructed, p.Next())
		}
		if !reflect.DeepEqual(directions, reconstructed) {
			t.Fatalf("test%d paths mismatch", i)
		}

		p.Rewind()
		reconstructed = reconstructed[:0]
		for p.HasNext() {
			reconstructed = append(reconstructed, p.Next())
		}
		if !reflect.DeepEqual(directions, reconstructed) {
			t.Fatalf("test%d paths mismatch", i)
		}
	}
}
