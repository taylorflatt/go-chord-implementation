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

type netmap struct {
	nodes []node
	size  int
}

type ftentry struct {
	key int
	s   int
}

type ft struct {
	entries []ftentry
	size    int
}

type node struct {
	id  int
	s   int
	p   int
	act bool
	ft  ft
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

func ComputeFTableSize(s int) (int, error) {

	r := math.Log(float64(s)) / math.Log(2)

	if !FloatIsDigit(r) {
		return -1, errors.New("the entered number does not work with 2")
	}

	return int(r), nil
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

	chord := netmap{
		nodes: make([]node, s),
		size:  s,
	}

	fmt.Println(Pow(2, 0))
	fmt.Println(Pow(2, 1))
	fmt.Println(Pow(2, 2))
	fmt.Println(Pow(2, 3))
	fmt.Println(Pow(2, 4))
	fmt.Println(Pow(2, 5))
	fmt.Println(Pow(2, 6))

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

			ftb := ft{
				entries: make([]ftentry, fts),
				size:    fts,
			}

			n := node{
				id:  i,
				s:   successor,
				p:   predecessor,
				act: true,
				ft:  ftb,
			}

			chord.nodes[i] = n
		}

		//chord.nodes = append(chord.nodes, n)
	}

	// Initialize the other nodes within the structure.
	for index, node := range chord.nodes {
		if node.act == false {
			node.id = index

			for j := 0; j < s; j++ {
				if chord.nodes[j].act == true {
					node.s = j
					break
				}
			}

			// Check for wrapping at the origin.
			if index == 0 {
				node.p = s - 1
			} else {
				node.p = index - 1
			}

			node.s = index + 1%s
		}
	}

	fmt.Printf("%v", chord.nodes)

}
