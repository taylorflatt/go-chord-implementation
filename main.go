package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var chord = make([]node, 0)

type netmap struct {
	nodes []node
	size  int
}

type ftentry struct {
	key int
	s   int
}

type ft struct {
	entry []ftentry
}

type node struct {
	id int
	s  int
	p  int
}

func addNode(i int) {

	n := node{
		id:          i,
		successor:   i + 1,
		predecessor: i - 1,
	}

	chord = append(chord, n)
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

func main() {

	r := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter the size of the CHORD network: ")
	st, _ := r.ReadString('\n')
	st = strings.TrimSpace(st)
	s, err := ParseInt32(st)

	if err != nil {
		log.Fatalf("Could not parse the size. Please enter an integer number.")
	}

	chord := netmap{
		nodes: make([]node, s),
		size:  s,
	}

	fmt.Println()
	fmt.Println("Enter an active node (type done to stop)")
	fmt.Println("-----------------------------------")
	node := ""
	for true {
		fmt.Print("Active Node: ")
		i, _ := r.ReadString('\n')

		r.Scan()
		i := strconv.ParseInt(r.Text(), 10, 32)

		if i == "done" {
			break
		}

		successor := (i + 1) % s
		predecessor := (i - 1) % s

		n := node{
			id: i,
			s:  successor,
			p:  predecessor,
		}

		chord.nodes = append(chord.nodes, n)
	}

}
