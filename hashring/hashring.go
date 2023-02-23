// Package hashring
// 在理想情况下，每个物理节点受影响的数据量
// 为其节点缓存数据最的1/4
// （X/(N+X)）N为原 有物理节点数，X为新加入物理节点数），
// 也就是集群中已经被缓存的数据有75%可以被继续命中，
// 和未使用虚拟节点的一致性Hash算法结果相同，
// 只是解决的负载均衡的问题。
// @author: xwc1125
// @date: 2021/5/27
package hashring

import (
	"crypto/sha1"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/chain5j/chain5j-pkg/util/convutil"
	"github.com/chain5j/logger"
)

const (
	// DefaultVirtualSpots default virtual spots
	DefaultVirtualSpots = 400
)

type node struct {
	nodeKey   uint64 // 节点的key
	spotValue uint32 // 值
}

type nodesArray []node

func (p nodesArray) Len() int           { return len(p) }
func (p nodesArray) Less(i, j int) bool { return p[i].spotValue < p[j].spotValue }
func (p nodesArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p nodesArray) Sort()              { sort.Sort(p) }

// HashRing store nodes and weights
type HashRing struct {
	log          logger.Logger
	virtualSpots int               // 虚拟点
	nodes        nodesArray        // 节点数组
	weights      map[uint64]uint64 // 权重
	mu           sync.RWMutex      // 锁
}

// NewHashRing create a hash ring with virtual spots
func NewHashRing(spots int) *HashRing {
	if spots == 0 {
		spots = DefaultVirtualSpots
	}

	h := &HashRing{
		log:          logger.Log("harshRing"),
		virtualSpots: spots,
		weights:      make(map[uint64]uint64),
	}
	return h
}

// AddNodes add nodes to hash ring
func (h *HashRing) AddNodes(nodeWeight map[uint64]uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for nodeKey, w := range nodeWeight {
		h.weights[nodeKey] = w
	}
	h.generate()
}

// AddNode add node to hash ring
func (h *HashRing) AddNode(nodeKey uint64, weight uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

// RemoveNode remove node
func (h *HashRing) RemoveNode(nodeKey uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.weights, nodeKey)
	h.generate()
}

// UpdateNode update node with weight
func (h *HashRing) UpdateNode(nodeKey uint64, weight uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

func (h *HashRing) generate() {
	var totalW uint64
	for _, w := range h.weights {
		totalW += w
	}

	totalVirtualSpots := h.virtualSpots * len(h.weights)
	h.nodes = nodesArray{}

	for nodeKey, w := range h.weights {
		spots := int(math.Floor(float64(w) * float64(totalVirtualSpots) / float64(totalW)))
		for i := 1; i <= spots; i++ {
			hashBytes, _ := getHashBytes([]byte(convutil.ToString(nodeKey) + ":" + strconv.Itoa(i)))
			n := node{
				nodeKey:   nodeKey,
				spotValue: genValue(hashBytes[6:10]),
			}
			h.nodes = append(h.nodes, n)
		}
	}
	h.nodes.Sort()
}

// Kemata hash计算
func genValue(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) |
		(uint32(bs[2]) << 16) |
		(uint32(bs[1]) << 8) |
		(uint32(bs[0]))
	return v
}

// GetNode get node with key
func (h *HashRing) GetNode(s string) (uint64, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.nodes) == 0 {
		return 0, fmt.Errorf("nodes is empty")
	}

	hashBytes, err := getHashBytes([]byte(s))
	if err != nil {
		h.log.Warn("hashRing sha1 err", "err", err)
		return 0, err
	}

	v := genValue(hashBytes[6:10])
	i := sort.Search(len(h.nodes), func(i int) bool { return h.nodes[i].spotValue >= v })

	if i == len(h.nodes) {
		// h.log.Debug("hashRing sort.Search index the same err", "nodesLen", len(h.nodes), "index", i)
		i = 0
	}
	return h.nodes[i].nodeKey, nil
}

func getHashBytes(input []byte) ([]byte, error) {
	hash := sha1.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (h *HashRing) GetWeight() map[uint64]uint64 {
	return h.weights
}
