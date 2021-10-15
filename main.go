// Copyright Clayton Brown 2020. See LICENSE file.

package main

import (
	// "fmt"
	// "foxhole/grid"
	"foxhole/solvers"
)

func main() {

	// fmt.Println(len(solvers.Combinations(5, 2)))

	// fmt.Println(solvers.Hashes)
	solvers.Solve(solvers.Brute, 5, 12)

	// g := grid.CreateBlankGrid()
	// g.Values[0] = true
	// g.Values[1] = true
	// g.Values[2] = true
	// fmt.Println(g.Values)

	// g = g.PropogateWithCheck(1)
	// fmt.Println(g.Values)

}
