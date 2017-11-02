package main

import "fmt"

func PowTest(x int, y int) int {

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

func main() {

	for i := 0; i < 5; i++ {
		fmt.Println(PowTest(2, i))
	}
}
