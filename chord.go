package main

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
