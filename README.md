# CHORD Implementation [![GoDoc](https://godoc.org/github.com/taylorflatt/go-chord-implementation?status.svg)](https://godoc.org/github.com/taylorflatt/go-chord-implementation) [![Build Status](https://travis-ci.org/taylorflatt/go-chord-implementation.svg?branch=master)](https://travis-ci.org/taylorflatt/go-chord-implementation)

A basic implementation of the lookup for a P2P CHORD network.

## Usage
Start the program by running `go build main.go chord.go`. There are currently 2 flags available: 
- `-v, --verbose`: Displays all runtime state information at each step.
- `-m, --manual`: Allows manual input of active nodes. 

A build is required due to the package being non-main. Hence, `go run main.go chord.go` will fail. To circumvent this, change the package name to main.

## Known Bugs
None currently. If you run into any problems, please don't hesistate to create an issue

## Future Ideas
Complete the implementation and gravitate towards a real implementation using SHA-1 for key generation. The implementation would also include nodes joining/leaving in a real network.
