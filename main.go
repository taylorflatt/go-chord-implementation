package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Netmap struct {
	Nodes []Node
	Size  int
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
	Id          int
	Successor   int
	Predecessor int
	Active      bool
	Table       FingerTable
}

// ParseInt32 custom method.
func ParseInt32(s string) (int, error) {

	s = strings.TrimSpace(s)
	st, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return -1, err
	}

	si := int(st)

	return si, nil
}

func FloatIsDigit(n float64) bool {

	return n == float64(int(n))
}

func Pow(x int, y int) int {

	res := x

	switch {
	case y == 0:
		return 1
	case y == 1:
		return x
	default:
		for y > 0 {
			res *= x
			y--
		}
		return res
	}
}

func ComputeFTableSize(s int) (int, error) {

	r := math.Log(float64(s)) / math.Log(2)

	if !FloatIsDigit(r) {
		return -1, errors.New("the entered number does not work with 2")
	}

	return int(r), nil
}

func InitializeChord(size int) Netmap {

	chord := Netmap{
		Nodes: make([]Node, size),
		Size:  size,
	}

	return chord
}

func CreateActiveNodes(network Netmap, r *bufio.Reader) {

	fmt.Println()
	fmt.Println("Enter an active node (type done to stop)")
	fmt.Println("-----------------------------------")

	// Set the active nodes for the network.
	for true {
		fmt.Print("Active Node: ")
		it, _ := r.ReadString('\n')
		it = strings.TrimSpace(it)
		fmt.Print(it)

		if it == "done" {
			break
		}

		i, err := ParseInt32(it)

		switch {
		case err != nil:
			fmt.Println("Could not parse the node number. Please enter an integer number.")
		case i > network.Size-1:
			fmt.Println("Please enter a node that is within the size of the network.")
		default:
			// Create the node and set it to active.
			n := Node{
				Id:     i,
				Active: true,
			}

			// Add the node to the chord network.
			network.Nodes[i] = n
		}
	}
}

func DetermineSuccessors(network Netmap) {

	lBound := 0
	first := false
	firstActive := 0

	for index, node := range network.Nodes {

		if node.Active == true {
			// Set the successor for all the nodes between two active nodes to be
			// the current node as the successor.
			for lBound <= index {
				network.Nodes[lBound].Successor = node.Id
				lBound++
			}

			// Store the first active node for later use.
			if first {
				firstActive = node.Id
				first = false
			}

			// The active node is its own successor, hence + 1.
			lBound++
		}

		// When it reaches the end of the circular structure, there could be nodes
		// which haven't been assigned a successor. So, they should be assigned to
		// the first active node found since they are immediately before it logically.
		if index == network.Size-1 {
			for lBound <= index {
				network.Nodes[lBound].Successor = firstActive
				lBound++
			}
		}
	}
}

func CreateFingerTables(network Netmap, fingerTableSize int) {

	for k, node := range network.Nodes {
		table := FingerTable{
			Entries: make([]FtEntry, fingerTableSize),
			Size:    fingerTableSize,
		}

		// Generate an entry into the node's finger table.
		for i := 0; i < fingerTableSize; i++ {
			key := k + Pow(2, i)
			successor := network.Nodes[key].Successor

			table.Entries[i] = FtEntry{
				Key:       key,
				Successor: successor,
			}
		}

		// Set the node's completed finger table.
		node.Table = table
	}
}

func main() {

	r := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter the size of the CHORD network: ")
	st, _ := r.ReadString('\n')
	st = strings.TrimSpace(st)
	s, err := ParseInt32(st)
	fts, err := ComputeFTableSize(s)

	if err != nil {
		fmt.Println("The size of the network must be some exponential of 2 (e.g. 2^5 = 32).")
	}

	//fmt.Print(fts)

	if err != nil {
		log.Fatalf("Could not parse the size. Please enter an integer number.")
	}

	chord := InitializeChord(s)
	CreateActiveNodes(chord, r)
	DetermineSuccessors(chord)
	CreateFingerTables(chord, fts)

	for _, node := range chord.Nodes {
		fmt.Println("Node: ", node.Id)
		fmt.Println("--------------")
		fmt.Println("Active: ", node.Active)
		fmt.Println("Predecessor: ", node.Predecessor)
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
