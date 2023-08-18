// Package util
package util

// MinInt returns the smaller one of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
