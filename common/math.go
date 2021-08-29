package common

// MinInt returns the minimum of the given values
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// MaxInt returns the maximum of the given values
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
