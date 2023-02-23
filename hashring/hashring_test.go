// Package hashring
//
// @author: xwc1125
// @date: 2021/8/18
package hashring

import (
	"fmt"
	"testing"
)

func TestHashRing_AddNode(t *testing.T) {
	nodeWeight := map[uint64]uint64{
		0: 100,
		1: 100,
		// 2: 100,
	}
	hashRing := NewHashRing(10)
	hashRing.AddNodes(nodeWeight)

	for i := 0; i < 10; i++ {
		fmt.Println(hashRing.GetNode(fmt.Sprintf("%d", i)))
	}
}
