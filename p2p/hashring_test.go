package p2p

import (
	"testing"
)

func TestHashRing_AddAndGet(t *testing.T) {
	ring := NewHashring()

	nodes := []string{"192.168.1.1:3000", "192.168.1.2:3000", "192.168.1.3:3000"}
	for _, n := range nodes {
		ring.AddNode(n)
	}

	if len(ring.nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(ring.nodes))
	}

	// Verify that a specific key consistently maps to the exact same node
	key := "my_test_file.txt"
	expectedNode := ring.GetNode(key)

	if expectedNode == "" {
		t.Fatal("expected a valid node address, got empty string")
	}

	// Do it 100 times to ensure the mutexes and mapping remain stable
	for i := 0; i < 100; i++ {
		node := ring.GetNode(key)
		if node != expectedNode {
			t.Fatalf("consistent hash failed: expected %s, got %s", expectedNode, node)
		}
	}
}

func TestHashRing_RemoveNode(t *testing.T) {
	ring := NewHashring()
	ring.AddNode("nodeA")
	ring.AddNode("nodeB")

	key := "test_file"
	firstOwner := ring.GetNode(key)

	// Remove the node that owns the key
	ring.RemoveNode(firstOwner)

	if len(ring.nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(ring.nodes))
	}

	// The key should now smoothly map to the only remaining node
	secondOwner := ring.GetNode(key)
	if secondOwner == firstOwner {
		t.Fatalf("expected owner to change after removal")
	}
}