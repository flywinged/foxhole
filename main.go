// Copyright Clayton Brown 2020. See LICENSE file.

package main

import (
	"fmt"
	"foxhole/grid"
)

func main() {

	g := grid.CreateBlankGrid()
	g[1] = true

	fmt.Println(g, g.Hash())

	g = g.Propogate()
	fmt.Println(g, g.Hash())

	g = g.Propogate()
	fmt.Println(g, g.Hash())

	g = g.Propogate()
	fmt.Println(g, g.Hash())

	fmt.Println(grid.BaseGrid)

}
