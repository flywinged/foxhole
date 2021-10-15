// Copyright Clayton Brown 2020. See LICENSE file.

package solvers

import (
	"fmt"
	"foxhole/grid"
)

/*
	Brute force solver which will just try every possible option
	until everything is exhausted.
*/
func Brute(originalGrid *grid.Grid, checks int) []*grid.Grid {

	fmt.Println()
	fmt.Println(originalGrid.Values)

	nextGrids := []*grid.Grid{}

	for i, value := range originalGrid.Values {
		if value {
			newGrid := originalGrid.PropogateWithCheck(i)
			fmt.Println(i, newGrid.Checks, newGrid.Values)
			nextGrids = append(nextGrids, &newGrid)
		}
	}

	return nextGrids

}
