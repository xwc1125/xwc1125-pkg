// Package shift
//
// @author: xwc1125
package shift

import (
	"fmt"
	"testing"
)

func TestShift(t *testing.T) {
	// 1c2d3e4f5g6789ab
	shift := Shift("123456789abcdefg", 10)
	fmt.Println(shift)
}
