package pathing

import "testing"

//go:noinline
func iterateAllSteps(steps *GridPath) int {
	n := 0
	for steps.HasNext() {
		x := steps.Next()
		n += int(x)
	}
	return n
}

//go:noinline
func truncateSteps(steps *GridPath, n int) GridPath {
	t := steps.Truncated(n)
	t = t.Truncated(n) // no-op
	t = t.Truncated(3)
	return t
}

func BenchmarkGridPathConstruct(b *testing.B) {
	var parts []Direction
	for i := 0; i < 20; i++ {
		parts = append(parts, DirLeft)
		parts = append(parts, DirUp)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var p GridPath
		for _, dir := range parts {
			p.push(dir)
		}
	}
}

func BenchmarkGridPathTruncate(b *testing.B) {
	var p GridPath
	for i := 0; i < 20; i++ {
		p.push(DirLeft)
		p.push(DirUp)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		truncateSteps(&p, 15)
	}
}

func BenchmarkGridPathIterate(b *testing.B) {
	var p GridPath
	for i := 0; i < 20; i++ {
		p.push(DirLeft)
		p.push(DirUp)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.Rewind()
		iterateAllSteps(&p)
	}
}
