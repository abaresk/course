package course

func modulo(a, b int) int {
	r := a % b
	if r < 0 {
		return r + b
	}
	return r
}
