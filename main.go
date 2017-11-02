package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Netmap struct {
	AnchorId int
	Nodes    []Node
	Size     int
}

type FtEntry struct {
	Key       int
	Successor int
}

type FingerTable struct {
	Entries []FtEntry
	Size    int
}

type Node struct {
	Id        int
	Successor int
	Active    bool
	Table     FingerTable
}

func main() {

	// Set flags for node creation and output verbosity.
	ml := flag.Bool("manual", false, "Manually enter the active nodes for the network.")
	ms := flag.Bool("m", false, "Manually enter the active nodes for the network.")
	vl := flag.Bool("verbose", false, "Prints the state of the program after each step. Warning: This will add considerable clutter.")
	vs := flag.Bool("v", false, "Prints the state of the program after each step. Warning: This will add considerable clutter.")
	flag.Parse()

	// Initialize the flags.
	man := false
	if *ml {
		man = *ml
	}
	if *ms {
		man = *ms
	}

	verb := false
	if *vl {
		verb = *vl
	}
	if *vs {
		verb = *vs
	}

	r := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter the size of the CHORD network: ")
	st, _ := r.ReadString('\n')
	st = strings.TrimSpace(st)
	s, err := ParseInt32(st)
	fts, err := ComputeFTableSize(s)

	if err != nil {
		fmt.Println("The size of the network must be some exponential of 2 (e.g. 2^5 = 32).")
	}

	if err != nil {
		log.Fatalf("Could not parse the size. Please enter an integer number.")
	}

	chord := InitializeChord(s)
	if verb {
		PrintNetwork(chord)
	}
	if man {
		CreateActiveNodes(&chord, r)
	} else {
		GenerateActiveNodes(&chord, r)
	}
	if verb {
		PrintActiveNodes(chord)
		PrintNetwork(chord)
	}
	DetermineSuccessors(&chord)
	if verb {
		PrintNetwork(chord)
	}
	CreateFingerTables(&chord, fts)

	if verb {
		PrintNetwork(chord)
	}

	FindSuccessor(&chord, chord.AnchorId, 16)
}

// ComputeFTableSize determines the sizes of finger tables based on the size of the network.
// It returns the finger table size and an error if the network size cannot be expressed as s = 2^x.
func ComputeFTableSize(s int) (int, error) {

	r := math.Log(float64(s)) / math.Log(2)

	if !FloatIsDigit(r) {
		return -1, errors.New("the entered number does not conform to n = 2^x")
	}

	return int(r), nil
}

// InitializeChord creates the network and sets the ids of all the nodes.
// It returns a network without active nodes of size n.
func InitializeChord(size int) Netmap {

	chord := Netmap{
		Nodes: make([]Node, size),
		Size:  size,
	}

	for index, node := range chord.Nodes {
		node.Id = index
		chord.Nodes[index] = node
	}

	return chord
}

// GenerateActiveNodes will generate the active nodes for the network based on a PRNG
// X_(n+1) = (aX_n + c)mod(m).
// It doesn't return anything.
func GenerateActiveNodes(network *Netmap, r *bufio.Reader) {

	fmt.Println()
	fmt.Println("Enter the parameters for the PRNG")
	fmt.Println("-----------------------------------")

	fmt.Print("Seed: ")
	st, _ := r.ReadString('\n')
	st = strings.TrimSpace(st)
	seed, err := ParseInt32(st)

	if err != nil || seed < 0 {
		fmt.Println("Invalid seed number. Please enter a positive integer number.")
	}

	fmt.Print("Increment: ")
	it, _ := r.ReadString('\n')
	st = strings.TrimSpace(it)
	increment, err := ParseInt32(it)

	if err != nil || increment < 0 {
		fmt.Println("Invalid increment number. Please enter a positive integer number.")
	}

	fmt.Print("Multiplier: ")
	mt, _ := r.ReadString('\n')
	st = strings.TrimSpace(mt)
	multiplier, err := ParseInt32(mt)

	if err != nil || multiplier < 0 {
		fmt.Println("Invalid multiplier number. Please enter a positive integer number.")
	}

	min := network.Size

	// Set the first node as active.
	i := ((multiplier * seed) + increment) % network.Size
	network.Nodes[i].Active = true

	for true {

		if i < min {
			min = i
		}

		// Pseudo-random number generator.
		i = ((multiplier * i) + increment) % network.Size

		// We have begun repeating, thus we have generated all active nodes.
		if network.Nodes[i].Active == true {
			// Set the first active node (where queries first enter).
			network.AnchorId = min
			break
		}

		network.Nodes[i].Active = true
	}
}

// CreateActiveNodes prompts the user for manual entry of the active nodes in the network
// bounded by 0 <= aNode <= networkSize - 1.
// It doesn't return anything.
func CreateActiveNodes(network *Netmap, r *bufio.Reader) {

	fmt.Println()
	fmt.Println("Enter an active node (type done to stop)")
	fmt.Println("-----------------------------------")

	min := network.Size
	// Set the active nodes for the network.
	for true {
		fmt.Print("Active Node: ")
		it, _ := r.ReadString('\n')
		it = strings.TrimSpace(it)

		if it == "done" {
			break
		}

		i, err := ParseInt32(it)

		switch {
		case err != nil:
			fmt.Println("Could not parse the node number. Please enter an integer number.")
		case i < 0:
			fmt.Println("Please enter a positive integer number.")
		case i > network.Size-1:
			fmt.Println("Please enter a node that is within the size of the network.")
		default:
			if i < min {
				min = i
			}

			// Set the node to active.
			network.Nodes[i].Active = true
		}
	}

	// Set the first active node (where queries first enter).
	network.AnchorId = min
}

// DetermineSuccessors computes the successor for each node in the network which forms the
// logical path between each node (creating a circle).
// It doesn't return anything.
func DetermineSuccessors(network *Netmap) {

	lBound := 0

	for index, node := range network.Nodes {

		if node.Active == true {
			// Set the successor for all the nodes between two active nodes to be
			// the current node as the successor.
			for lBound <= index {
				network.Nodes[lBound].Successor = node.Id
				lBound++
			}
		}

		// When it reaches the end of the circular structure, there could be nodes
		// which haven't been assigned a successor. So, they should be assigned to
		// the first active node found since they are immediately before it logically.
		if index == network.Size-1 {
			for lBound <= index {
				network.Nodes[lBound].Successor = network.AnchorId
				lBound++
			}
		}
	}
}

// CreateFingerTables creates the finger tables for the network given an appropriate size.
// It doesn't return anything.
func CreateFingerTables(network *Netmap, fingerTableSize int) {

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
}

// FindSuccessor is the lookup algorithm for a particular node to determine which active node
// contains the node's data/information.
// It returns the ID of the node which contains the node's data (it's successor).
func FindSuccessor(network *Netmap, node int, find int) int {

	fmt.Println("Entering node ", node)
	PrintNodeFingerTable(network.Nodes[node])

	min := network.Size
	for _, entries := range network.Nodes[node].Table.Entries {
		switch {
		// We know this by definition of the algorithm.
		case find < node:
			fmt.Println("  > By definition, we know node ", node, " must be who we are looking for.")
			return node
		case entries.Key < find:
			fmt.Println("  > key ", entries.Key, " is less than our value ", find)
			min = entries.Successor
		case entries.Key == find:
			fmt.Println("  > Found ", find, " with node ", entries.Successor, " as it's successor!")
			return entries.Successor
		case entries.Key > find:
			fmt.Println("  > Node ", min, " is the closest preceeding node. Moving to node ", min)
			break
		}
	}

	// If the successor in the node's table is itself, then it is the value's successor.
	if min == node {
		fmt.Println("  > By definition, we know node ", node, " must be who we are looking for.")
		return node
	}

	fmt.Println()
	return FindSuccessor(network, min, find)
}

// ParseInt32 takes a string, trims the space, and parses it into an int.
// It returns the parsed int or an error due to a failure in parsing.
func ParseInt32(s string) (int, error) {

	s = strings.TrimSpace(s)
	st, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return -1, err
	}

	si := int(st)

	return si, nil
}

// FloatIsDigit determines if a given float is a whole number.
// It returns true if it is a whole number, otherwise false.
func FloatIsDigit(n float64) bool {

	return n == float64(int(n))
}

// Pow computes x^y for ints.
// It returns an int value of the exponentiation.
func Pow(x int, y int) int {

	res := x

	switch {
	case y == 0:
		return 1
	case y == 1:
		return x
	default:
		for i := 1; i < y; i++ {
			res *= x
		}
		return res
	}
}

// PrintNetwork displays the network information and all of the node information.
// It doesn't return anything.
func PrintNetwork(network Netmap) {

	fmt.Println("Network Size: ", network.Size)
	fmt.Println("Network Anchor: ", network.AnchorId)

	for _, node := range network.Nodes {
		fmt.Println("Node: ", node.Id)
		fmt.Println("--------------")
		fmt.Println("Active: ", node.Active)
		fmt.Println("Successor: ", node.Successor)
		fmt.Println("--------------")
		fmt.Println("FINGER TABLE")
		fmt.Println("--------------")

		for _, entry := range node.Table.Entries {
			fmt.Print("Key = ", entry.Key)
			fmt.Print(" , Value = ", entry.Successor)
			fmt.Println()
		}
	}
}

// PrintNode displays a particular node's information.
// It doesn't return anything.
func PrintNode(node Node) {

	fmt.Println("Node: ", node.Id)
	fmt.Println("-------------------")
	fmt.Println("Active: ", node.Active)
	fmt.Println("Successor: ", node.Successor)
	fmt.Println("-------------------")
	fmt.Println("FINGER TABLE")
	fmt.Println("-------------------")

	for _, entry := range node.Table.Entries {
		fmt.Print("Key = ", entry.Key)
		fmt.Print(" , Value = ", entry.Successor)
		fmt.Println()
	}
}

// PrintActiveNodes displays a newline separated list of all the active nodes
// for a network.
// It doesn't return anything.
func PrintActiveNodes(network Netmap) {

	fmt.Println("Active Nodes:")
	fmt.Println("-------------------")
	for _, node := range network.Nodes {
		if node.Active == true {
			fmt.Println("Node: ", node.Id)
		}
	}
}

// PrintNodeFingerTable displays only a particular node's finger table.
// It doesn't return anything.
func PrintNodeFingerTable(node Node) {

	fmt.Println("-------------------")

	for _, entry := range node.Table.Entries {
		fmt.Print("Key = ", entry.Key)
		fmt.Print(" , Value = ", entry.Successor)
		fmt.Println()
	}
	fmt.Println("-------------------")
}
