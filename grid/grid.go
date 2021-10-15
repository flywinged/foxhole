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
func CreateBlankGrid() *Grid {

	values := make([]bool, len(BaseGrid.Connections))
	return &Grid{
		Values: values,
	}

}

/*
	Function for copying a grid
*/
func (grid *Grid) Copy() *Grid {

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
func (grid *Grid) Propogate() *Grid {

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

func (grid *Grid) PropogateWithChecksAndAdd(checks map[int]bool) *Grid {
	newGrid := grid.PropgateWithChecks(checks)
	newGrid.AddChecks(checks)
	return newGrid
}

func (grid *Grid) AddChecks(checks map[int]bool) {
	grid.Checks = append(grid.Checks, checks)
}

func (grid *Grid) PropgateWithChecks(checks map[int]bool) *Grid {

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

	return newGrid

}

/*
	Determines what needs to be checked to remove the possiblity of
	a certain tile appearing in the next propogation. Returns an array
	of the indexes which need to be checked to perform that removal.
*/
func (grid *Grid) HowToRemove(checks map[int]bool) []map[int]bool {

	howToRemove := []map[int]bool{}
	for i := 0; i < len(grid.Values); i++ {
		howToRemove = append(howToRemove, map[int]bool{})
	}

	for i, value := range grid.Values {

		// Don't need to worry if the grid doesn't have this value as a possibility
		if !value {
			continue
		}

		for _, conn := range BaseGrid.Connections[i] {
			howToRemove[conn][i] = true
		}
	}

	return howToRemove

}

/*
	Determine how many trues are currently in the array
*/
func (grid *Grid) NFoxes() int {
	count := 0
	for _, value := range grid.Values {
		if value {
			count++
		}
	}
	return count
}

// Used for doing powers way faste
var POWERS = map[int]int{}
func init() {
	for i := 0; i < 1024; i++ {
		POWERS[i] = power2(i)
	}
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
	lowestHashFirstIndex := len(BaseGrid.Symmetries[0])

	// Loop through each of the valid symettric configurations
configurationLoop:
	for _, configuration := range BaseGrid.Symmetries {

		// Construct this has in the order of the symmetry configuration
		hash := 0
		firstIndex := -1

		for power, i := range configuration {

			if firstIndex == -1 && i > lowestHashFirstIndex {
				continue configurationLoop
			}

			if grid.Values[i] {
				if firstIndex == -1 {
					firstIndex = i
				}
				hash += POWERS[power]
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

/*
	Check if a grid is equal to another grid
*/
func (grid *Grid) Equal(other *Grid) bool {

	for i, gridValue := range grid.Values {
		otherValue := other.Values[i]
		if gridValue != otherValue {
			return false
		}
	}

	return true

}