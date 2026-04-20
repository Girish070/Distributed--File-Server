package p2p

import (
	"crypto/sha1"
	"encoding/binary"
	"sort"
	"sync"
)

type HashRing struct {
	mu      sync.RWMutex
	nodes   []uint32
	nodeMap map[uint32]string
}

func NewHashring() *HashRing {
	return &HashRing{
		nodes:   []uint32{},
		nodeMap: make(map[uint32]string),
	}
}

func (h *HashRing) hash(key string) uint32 {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	sum := hasher.Sum(nil)

	return binary.BigEndian.Uint32(sum[:4])
}

func (h *HashRing) AddNode(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	hash := h.hash(addr)
	if _, existed := h.nodeMap[hash]; !existed {
		h.nodes = append(h.nodes, hash)

		sort.Slice(h.nodes, func(i, j int) bool {
			return h.nodes[i] < h.nodes[j]
		})
		h.nodeMap[hash] = addr
	}
}

func (h *HashRing) RemoveNode(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	hash := h.hash(addr)
	if _, exists := h.nodeMap[hash]; exists {
		delete(h.nodeMap, hash)

		var newNode []uint32
		for _, n := range h.nodes {
			if n != hash {
				newNode = append(newNode, n)
			}
		}
		h.nodes = newNode
	}
}

func (h *HashRing) GetNode(key string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.nodes) == 0 {
		return ""
	}
	hash := h.hash(key)

	idx := sort.Search(len(h.nodes), func(i int) bool {
		return h.nodes[i] >= hash
	})

	if idx == len(h.nodes) {
		idx = 0
	}
	return h.nodeMap[h.nodes[idx]]
}
