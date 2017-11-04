package chord

import (
	"fmt"
	"os"
	"testing"
)

var networkA Netmap
var networkB Netmap
var networkC Netmap
var networkD Netmap

func TestMain(m *testing.M) {

	networkSize := 32
	fingerTableSize := 5

	// Initialize a basic network.
	networkA = InitializeTestNetworks(networkSize, fingerTableSize)
	networkB = InitializeTestNetworks(networkSize, fingerTableSize)
	networkC = InitializeTestNetworks(networkSize, fingerTableSize)
	networkD = InitializeTestNetworks(networkSize, fingerTableSize)

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
	for k := range network.Nodes {
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

func TestComputeFTableSize(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Printf("[ComputeFTableSize TEST]\n")
	fmt.Println("-----------------------------------")

	// Create tests
	tests := []struct {
		in       int
		expected int
	}{
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
	}

	// Run the tests.
	for index, test := range tests {
		fmt.Printf("\n[Test %d] Input: %d, Expected: %d\n", index, test.in, test.expected)
		out, err := ComputeFTableSize(test.in)

		if err != nil {
			t.Errorf("[Test %d] Expected: %d, Received: %d AND Error: %s", index, test.expected, out, err)
		}

		if out != test.expected {
			t.Errorf("[Test %d] Expected: %d, Received: %d", index, test.expected, out)
		}

		fmt.Printf("[Test %d] Passed!\n", index)
	}
}

func TestInitializeChord(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Printf("[InitializeChord TEST]\n")
	fmt.Println("-----------------------------------")

	// Create tests
	tests := []struct {
		in       int
		expected int
	}{
		{4, 4},
		{8, 8},
		{16, 16},
		{32, 32},
		{64, 64},
	}

	// Run the tests.
	for index, test := range tests {
		fmt.Printf("\n[Test %d] Input: %d, Expected: %d\n", index, test.in, test.expected)
		out := InitializeChord(test.in)

		if out.Size == 0 {
			t.Errorf("[Test %d] Expected: %d, Received: %d", index, test.in, out.Size)
		}

		for j, node := range out.Nodes {
			if node.ID != j {
				t.Errorf("[Test %d] Expected: %d, Received: %d", index, j, node.ID)
			}
		}
		fmt.Printf("[Test %d] Passed!\n", index)

	}
}

func TestFindSuccessor(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Printf("[FindSuccessor TEST]\n")
	fmt.Println("-----------------------------------")

	netANodeActive := []int{
		1, 2, 3, 4, 5,
	}

	netBNodeActive := []int{
		1, 4, 7, 12, 15, 20, 27,
	}

	netCNodeActive := []int{
		1, 2,
	}

	netDNodeActive := []int{
		1, 31,
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

	// Set network D active nodes.
	for _, node := range netDNodeActive {
		if networkD.Nodes[node].ID == node {
			networkD.Nodes[node].Active = true
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

	dFingerTableSize, _ := ComputeFTableSize(networkD.Size)
	DetermineSuccessors(&networkD)
	CreateFingerTables(&networkD, dFingerTableSize)

	// Create tests
	tests := []struct {
		net       *Netmap
		netID     string
		startNode int
		node      int
		expected  int
	}{
		{&networkA, "A", 3, 3, 3},
		{&networkA, "A", 1, 14, 1},
		{&networkA, "A", 1, 17, 1},
		{&networkA, "A", 1, 18, 1},
		{&networkB, "B", 1, 3, 4},
		{&networkB, "B", 1, 14, 15},
		{&networkB, "B", 1, 16, 20},

		{&networkB, "B", 7, 21, 27},
		{&networkB, "B", 27, 0, 1},
		{&networkB, "B", 1, 0, 1},

		{&networkC, "C", 1, 20, 1},
		{&networkC, "C", 1, 0, 1},
		{&networkC, "C", 2, 0, 1},
		{&networkD, "D", 1, 0, 1},
		{&networkD, "D", 1, 20, 31},
		{&networkD, "D", 1, 1, 1},
	}

	// Run the tests.
	for index, test := range tests {
		fmt.Printf("\n[Test %d] Network: %s, Node: %d, StartingAt: %d, Expected: %d\n", index, test.netID, test.node, test.startNode, test.expected)

		// Suppress the function's print statements.
		o := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
		out := FindSuccessor(test.net, test.startNode, test.node)
		os.Stdout = o

		if out != test.expected {
			t.Errorf("[Test %d] Network: %s failed lookup for %d starting at %d. Received: %d, Expected: %d", index, test.netID, test.node, test.startNode, out, test.expected)
		}

		fmt.Printf("[Test %d] Passed!\n", index)
	}
}
