package main

import (
	"os"
	"testing"
)

var networkA Netmap
var networkB Netmap
var networkC Netmap

func TestMain(m *testing.M) {

	networkSize := 32
	fingerTableSize := 5

	// Initialize a basic network.
	networkA = InitializeTestNetworks(networkSize, fingerTableSize)
	networkB = InitializeTestNetworks(networkSize, fingerTableSize)
	networkC = InitializeTestNetworks(networkSize, fingerTableSize)

	os.Exit(m.Run())
}

func InitializeTestNetworks(networkSize int, fingerTableSize int) Netmap {

	// Initialize a basic network.
	network := Netmap{
		Nodes: make([]Node, networkSize),
		Size:  networkSize,
	}

	// Set the node basic properties.
	for index, node := range network.Nodes {
		node.ID = index
		network.Nodes[index] = node
	}

	// Create the finger tables for each node.
	for k, _ := range network.Nodes {
		table := FingerTable{
			Entries: make([]FtEntry, fingerTableSize),
			Size:    fingerTableSize,
		}

		// Generate an entry into the node's finger table.
		for i := 0; i < fingerTableSize; i++ {
			key := (k + Pow(2, i)) % network.Size
			successor := network.Nodes[key].Successor

			table.Entries[i] = FtEntry{
				Key:       key,
				Successor: successor,
			}
		}

		// Set the node's completed finger table.
		network.Nodes[k].Table = table
	}

	return network
}

func TestFindSuccessor(t *testing.T) {

	netANodeActive := []int{
		1, 2, 3, 4, 5,
	}

	netBNodeActive := []int{
		1, 4, 7, 12, 15, 20, 27,
	}

	netCNodeActive := []int{
		1, 2,
	}

	// Set network A active nodes.
	for _, node := range netANodeActive {
		if networkA.Nodes[node].ID == node {
			networkA.Nodes[node].Active = true
		}
	}

	// Set network B active nodes.
	for _, node := range netBNodeActive {
		if networkB.Nodes[node].ID == node {
			networkB.Nodes[node].Active = true
		}
	}

	// Set network C active nodes.
	for _, node := range netCNodeActive {
		if networkC.Nodes[node].ID == node {
			networkC.Nodes[node].Active = true
		}
	}

	// Create finger tables
	// Note: Should probably have nested testing where each parent function is tested first.
	aFingerTableSize, _ := ComputeFTableSize(networkA.Size)
	DetermineSuccessors(&networkA)
	CreateFingerTables(&networkA, aFingerTableSize)

	bFingerTableSize, _ := ComputeFTableSize(networkB.Size)
	DetermineSuccessors(&networkB)
	CreateFingerTables(&networkB, bFingerTableSize)

	cFingerTableSize, _ := ComputeFTableSize(networkC.Size)
	DetermineSuccessors(&networkC)
	CreateFingerTables(&networkC, cFingerTableSize)

	// Create tests
	tests := []struct {
		net       *Netmap
		startNode int
		node      int
		expected  int
	}{
		{&networkA, 1, 14, 1},
		{&networkA, 1, 17, 1},
		{&networkA, 1, 18, 1},
		{&networkB, 1, 3, 4},
		{&networkB, 1, 14, 15},
		{&networkB, 1, 16, 20},
		{&networkC, 1, 20, 1},
		{&networkC, 1, 0, 1},
		{&networkC, 2, 0, 1},
	}

	// Run the tests.
	for _, test := range tests {
		out := FindSuccessor(test.net, test.startNode, test.node)

		if out != test.expected {
			t.Errorf("Failed to successfully lookup the successor for %d from node %d. Got: %d, expected: %d", test.node, test.startNode, out, test.expected)
		}
	}
}