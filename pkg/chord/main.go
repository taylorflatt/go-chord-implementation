// Package chord allows a user to generate a network with pseudo-random active nodes or specify active nodes.
// Additionally, it will determine the location of data for a particular node starting at any specified active node.
package chord

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

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

	fmt.Println()
	fmt.Println("Enter a node from where to search from: ")
	at, _ := r.ReadString('\n')
	at = strings.TrimSpace(at)
	anchor, err := ParseInt32(at)

	FindSuccessor(&chord, anchor, 18)
}
