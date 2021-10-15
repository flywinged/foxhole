// Copyright Clayton Brown 2020. See LICENSE file.

package solvers

import (
	// "fmt"
	"foxhole/grid"
)

/*
	Brute force solver which will just try every possible option
	until everything is exhausted.
*/
func Brute(originalGrid *grid.Grid, checks int) []*grid.Grid {
	return recursiveBrute(originalGrid, checks, map[int]bool{}, map[int]bool{})
}

/*
	Recursive function used to handle removing multiple foxholes from contention
*/
func recursiveBrute(originalGrid *grid.Grid, checksLeft int, checksMade, hashes map[int]bool) []*grid.Grid {

	resultingGrids := []*grid.Grid{}

	// fmt.Println()
	// fmt.Println("Original Grid:", originalGrid.Values)
	// fmt.Println("Checks Left:", checksLeft)
	// fmt.Println("Checks Made:", checksMade)

	// Can't remove anything if here
	if checksLeft <= 0 {
		newGrid := originalGrid.PropogateWithChecksAndAdd(checksMade)
		resultingGrids = append(resultingGrids, newGrid)
		return resultingGrids
	}

	// Determine what holes can be removed
	removalOptions := originalGrid.HowToRemove(checksMade)

	// fmt.Println("Removals:", removalOptions)
	for _, option := range removalOptions {

		/*
			We can't remove this hole if we don't have enough checks or
			if removing this hole doesn't do anything.
		*/
		removalRequirement := len(option)
		if removalRequirement > checksLeft || removalRequirement == 0 {
			continue
		}

		/*
			Otherwise determine if this series of checks has been made before.
			Union will modify the first argument.
		*/
		SetUnion(option, checksMade)

		optionHash := SetHash(option)
		if _, exists := hashes[optionHash]; !exists {
			hashes[optionHash] = true
			resultingGrids = append(resultingGrids, recursiveBrute(
				originalGrid,
				checksLeft - removalRequirement,
				option,
				hashes,
			)...)
		}
		
	}

	if len(resultingGrids) == 0 {
		newGrid := originalGrid.PropogateWithChecksAndAdd(checksMade)
		resultingGrids = append(resultingGrids, newGrid)
		return resultingGrids
	}
	
	return resultingGrids

}