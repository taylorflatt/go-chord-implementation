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

func CreateActiveNodes(network Netmap, r *Reader) {

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

func DetermineSuccessor() {

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

	fmt.Println()
	fmt.Println("Enter an active node (type done to stop)")
	fmt.Println("-----------------------------------")

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
		case i > s-1:
			fmt.Println("Please enter a node that is within the size of the network.")
		default:
			successor := (i + 1) % s
			predecessor := -1

			if i == 0 {
				predecessor = s - 1
			} else {
				predecessor = i - 1
			}

			ftb := FingerTable{
				Entries: make([]FtEntry, fts),
				Size:    fts,
			}

			// Create the active node's finger table.
			for k := 0; k < fts; k++ {
				key := (i + Pow(2, k)) % s

				ftb.Entries[k] = FtEntry{
					Key:       key,
					Successor: i,
				}
			}

			// Create the node and set it to active.
			n := Node{
				Id:          i,
				Successor:   successor,
				Predecessor: predecessor,
				Active:      true,
				Table:       ftb,
			}

			// Add the node to the chord network.
			chord.Nodes[i] = n
		}

		//chord.nodes = append(chord.nodes, n)
	}

	// Initialize the other (non-active) nodes within the structure.
	for index, node := range chord.Nodes {
		if node.Active == false {
			node.Id = index

			for j := 0; j < s; j++ {
				if chord.Nodes[j].Active == true {
					node.Successor = j
					break
				}
			}

			// Check for wrapping at the origin.
			if index == 0 {
				node.Predecessor = s - 1
			} else {
				node.Predecessor = index - 1
			}

			node.Successor = (index + 1) % s

			ftb := FingerTable{
				Entries: make([]FtEntry, fts),
				Size:    fts,
			}

			nextSuccessor := -1

			for i, node := range chord.Nodes {
				if i > index && node.Active == true {
					nextSuccessor = i
				}

			}

			for k := 0; k < fts; k++ {
				key := (index + Pow(2, k)) % s
				ftb.Entries[k] = FtEntry{
					Key:       key,
					Successor: nextSuccessor,
				}
			}
		}

		chord.Nodes[index] = node
	}

	//fmt.Printf("%v", chord.nodes)

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
