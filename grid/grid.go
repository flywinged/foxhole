// Copyright Clayton Brown 2020. See LICENSE file.

package grid

/*
	A grid is an array of bools which indicate whether or
	not a fox can be there or not.
*/
type Grid struct {
	Values []bool
	Checks []map[int]bool
}

/*
	Function for generating a new grid with the appropriate
	preallocated space.
*/
func CreateBlankGrid() Grid {

	values := make([]bool, len(BaseGrid.Connections))
	return Grid{
		Values: values,
	}

}

/*
	Function for copying a grid
*/
func (grid *Grid) Copy() Grid {

	newGrid := CreateBlankGrid()
	for i, value := range grid.Values {
		newGrid.Values[i] = value
	}
	newChecks := make([]map[int]bool, len(grid.Checks))
	for i, value := range grid.Checks {
		newChecks[i] = value
	}
	newGrid.Checks = newChecks

	return newGrid

}

/*
	Function for propogating a grid
*/
func (grid *Grid) Propogate() Grid {

	newGrid := CreateBlankGrid()
	for i, value := range grid.Values {
		for _, j := range BaseGrid.Connections[i] {
			newGrid.Values[j] = newGrid.Values[j] || value
		}
	}

	newChecks := make([]map[int]bool, len(grid.Checks))
	for i, value := range grid.Checks {
		newChecks[i] = value
	}
	newGrid.Checks = newChecks

	return newGrid

}

func (grid *Grid) PropogateWithChecks(checks map[int]bool) Grid {

	newGrid := CreateBlankGrid()
	for i, value := range grid.Values {
		checkedValue := checks[i]
		for _, j := range BaseGrid.Connections[i] {
			propogatedValue := !checkedValue && value
			newGrid.Values[j] = newGrid.Values[j] || propogatedValue
		}
	}

	newChecks := make([]map[int]bool, len(grid.Checks))
	for i, value := range grid.Checks {
		newChecks[i] = value
	}
	newGrid.Checks = newChecks
	newGrid.Checks = append(newGrid.Checks, checks)

	return newGrid

}

// Helper function for large int powers
func power2(p int) int {
	n := 1
	for i := 0; i < p; i++ {
		n *= 2
	}
	return n
}

/*
	Create a hash for the grid
*/
func (grid *Grid) Hash() int {

	// Because of all the symmetries, only grab the lowest hash
	lowestHash := -1

	// Loop through each of the valid symettric configurations
	for _, configuration := range BaseGrid.Symmetries {

		// Construct this has in the order of the symmetry configuration
		hash := 0
		for power, i := range configuration {
			if grid.Values[i] {
				hash += power2(power)
			}
		}

		// Replace the lowest hash if applicable
		if lowestHash == -1 {
			lowestHash = hash
		} else if hash < lowestHash {
			lowestHash = hash
		}

	}

	return lowestHash
}
