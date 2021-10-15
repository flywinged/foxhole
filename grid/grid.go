// Copyright Clayton Brown 2020. See LICENSE file.

package grid

/*
	A grid is an array of bools which indicate whether or
	not a fox can be there or not.
*/
type Grid []bool

/*
	Function for generating a new grid with the appropriate
	preallocated space.
*/
func CreateBlankGrid() Grid {

	grid := make([]bool, len(BaseGrid.Connections))
	return grid

}

/*
	Function for copying a grid
*/
func (grid *Grid) Copy() Grid {

	newGrid := CreateBlankGrid()
	for i, value := range *grid {
		newGrid[i] = value
	}

	return newGrid

}

/*
	Function for propogating a grid
*/
func (grid *Grid) Propogate() Grid {

	newGrid := CreateBlankGrid()
	for i, value := range *grid {
		for _, j := range BaseGrid.Connections[i] {
			newGrid[j] = newGrid[j] || value
		}
	}

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
			if (*grid)[i] {
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
