package course

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func modulo(a, b int) int {
	r := a % b
	if r < 0 {
		return r + b
	}
	return r
}
