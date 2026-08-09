package pathing

type uint128 struct {
	lo uint64
	hi uint64
}

func (u uint128) ShiftRight(n uint) uint128 {
	if n > 64 {
		u.lo = u.hi >> (n - 64)
		u.hi = 0
	} else {
		u.lo = u.lo>>n | u.hi<<(64-n)
		u.hi = u.hi >> n
	}
	return u
}
