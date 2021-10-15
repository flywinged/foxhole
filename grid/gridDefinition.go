// Copyright Clayton Brown 2020. See LICENSE file.

package grid

import (
	"fmt"
	"math"
	"github.com/gitchander/permutation"
)


func init() {

	// Original foxhole problem definition
	// BaseGrid = CreateLinearGrid(5)

	BaseGrid = CreatePrismGrid([]int{8, 8})

}

/*
	The default grid definition that will be used everywhere.
	The format for this is a list of connections for each node
	in the grid. The example for the basic foxhole problem would
	look like this:

	Connections: [][]int{
		{2},
		{1, 3},
		{2, 4},
		{3, 5},
		{4},
	}

	Symmetry: [][]int{
		{0,1,2,3,4},
		{4,3,2,1,0},
	}

*/
var BaseGrid *GridDefinition

type GridDefinition struct {
	// List of connections
	Connections [][]int

	// List of orderings which are symettric for the definitions
	Symmetries [][]int
}

/*
	Helper function for creating linear foxhole patterns
*/
func CreateLinearGrid(n int) *GridDefinition {

	// First create the connections
	connections := [][]int{}
	for i := 0; i < n; i += 1 {
		node := []int{}
		if i > 0 {
			node = append(node, i-1)
		}
		if i < n-1 {
			node = append(node, i+1)
		}
		connections = append(connections, node)
	}

	// Then create the symettries
	forward, backward := []int{}, []int{}
	for i := 0; i < n; i += 1 {
		forward = append(forward, i)
		backward = append(backward, n-i-1)
	}
	symmetries := [][]int{forward, backward}

	return &GridDefinition{
		Connections: connections,
		Symmetries:  symmetries,
	}

}

/*
	Function for determining the location of a cell
	in the values array given 3D coordinates.
*/
func getIndex(dimensionSizes []int, location []int) int {
	index := 0
	for i, x := range location {
		index += dimensionSizes[i] * x
	}
	return index
}

/*
	Function for determining the location of a cell in
	vector form given a value index.
*/
func getLocation(dimensionSizes []int, index int) []int {
	location := []int{}
	for i, size := range dimensionSizes {
		remainder := index
		if i < len(dimensionSizes)-1 {
			remainder = index % dimensionSizes[i+1]
		}

		x := int(math.Floor(float64(remainder) / float64(size)))
		location = append(location, x)
	}
	return location
}

/*
	Create n-cube grid.
*/
func CreatePrismGrid(dimensionLengths []int) *GridDefinition {

	/*
		Create the size of each dimension.
		Use this value to index and un-index things.
	*/
	totalCells := 1
	dimensionSizes := []int{}
	for _, l := range dimensionLengths {
		dimensionSizes = append(dimensionSizes, totalCells)
		totalCells *= l
	}

	// Generate the connections array first
	connections := [][]int{}
	for i := 0; i < totalCells; i++ {
		location := getLocation(dimensionSizes, i)

		// Generate each of the modified locations based on the current location
		connectionLocations := []int{}
		for i, x := range location {

			// Move in the negative direction along an axis
			if x > 0 {
				newLocation := make([]int, len(location))
				copy(newLocation, location)
				newLocation[i] = x - 1
				connectionLocations = append(connectionLocations, getIndex(dimensionSizes, newLocation))
			}

			// Move in the positive direction along an axis
			if x < dimensionLengths[i]-1 {
				newLocation := make([]int, len(location))
				copy(newLocation, location)
				newLocation[i] = x + 1
				connectionLocations = append(connectionLocations, getIndex(dimensionSizes, newLocation))
			}

		}

		// Add the generated locations for this cell to connections
		connections = append(connections, connectionLocations)

	}

	// Generate all the symettries. First we need the base symmetry
	symmetryHashes := map[string]bool{}
	symmetries := [][]int{}

	// Now apply all the other ones by dimension
	for _, dimensionInformation := range GetBinaryArrays(len(dimensionLengths)) {

		order := GetOrderedLocations(dimensionLengths, dimensionInformation)
		indexOrder := []int{}
		for _, location := range order {
			indexOrder = append(indexOrder, getIndex(dimensionSizes, location))
		}

		hash := hashSymmetry(indexOrder)
		if _, exists := symmetryHashes[hash]; !exists {
			symmetryHashes[hash] = true
			symmetries = append(symmetries, indexOrder)
		}

	}

	return &GridDefinition{
		Connections: connections,
		Symmetries:  symmetries,
	}
}


/*
	Returns a list of binary arrays
*/
func GetBinaryArrays(n int) [][]int {

	// Create the array list and the base array
	baseArrays := [][]int{}
	baseArray := []int{}
	for i := 0; i < n; i++ {
		baseArray = append(baseArray, i + 1)
	}
	baseArrays = append(baseArrays, baseArray)

	// Now sequentiall apply flips to everything
	for i := 0; i < n; i++ {

		newArrays := [][]int{}
		for _, array := range baseArrays {
			newArray := make([]int, len(array))
			copy(newArray, array)
			newArray[i] = -array[i]
			newArrays = append(newArrays, newArray)
		}
		baseArrays = append(baseArrays, newArrays...)

	}

	// Now permute each of the int options
	finalArrays := [][]int{}
	for _, array := range baseArrays {
		p := permutation.New(permutation.IntSlice(array))
		for p.Next() {
			newArray := make([]int, len(array))
			copy(newArray, array)
			finalArrays = append(finalArrays, newArray)
		}
	}

	// Return a list of the generated baseArrays
	return finalArrays

}

/*
	Creates an ordering based on dimensions
*/
func GetOrderedLocations(dimensionLengths []int, dimensionInformation []int) [][]int {

	// Current position. Start at whatever the appropriate corner is
	currentLocation := make([]int, len(dimensionLengths))
	dimensionDirections := []bool{}
	dimensionOrder := []int{}

	for i, dimensionInfo := range dimensionInformation {

		// Used to determine if the dimension is increasing or decreasing
		if dimensionInfo > 0 {
			dimensionDirections = append(dimensionDirections, true)
		} else {
			dimensionDirections = append(dimensionDirections, false)
			dimensionInfo = -dimensionInfo
		}
		dimensionInfo--

		// Fix the current position into the correct location
		currentLocation[dimensionInfo] = 0
		if !dimensionDirections[i] {
			currentLocation[dimensionInfo] = dimensionLengths[dimensionInfo] - 1
		}

		dimensionOrder = append(dimensionOrder, dimensionInfo)

	}

	/*
		Create a zig-zag boi starting at the current location and changing
		according to the dimension direction and ordering
	*/
	locations := [][]int{}

locationLoop:
	for {

		// Make a copy of the current location and add it to the orderings
		newLocation := make([]int, len(currentLocation))
		copy(newLocation, currentLocation)
		locations = append(locations, newLocation)

		/*
			Adjust the current location to match the orderings and
			direction of the dimensions
		*/
		for i, d := range dimensionOrder {
			direction := dimensionDirections[i]
			length := dimensionLengths[d]

			if direction {
				currentLocation[d]++
				if currentLocation[d] == length {
					currentLocation[d] = 0
				} else {
					continue locationLoop
				}
			} else {

				if currentLocation[d] != 0 {
					currentLocation[d]--
					continue locationLoop
				} else {
					currentLocation[d] = length - 1
				}

			}

		}

		// If we are here, then we didn't hit a continue so it's time to break out
		break

	}

	return locations

}

/*
	Helper function for grabbing symmetry hashes
*/
func hashSymmetry(s []int) string {
	hash := ""
	for i, v := range s {
		hash += fmt.Sprint(v)
		if i != len(s) - 1 {
			hash += ","
		}
	}
	return hash
}

/*
	Function for determining if a grid definition
	is binary or not.
*/
func (d *GridDefinition) RepeatingGrid() (int, []*Grid) {

	// Create an initial grid with only one cell shaded
	grids := []*Grid{}
	grid := CreateBlankGrid()
	grid.Values[0] = true

	// Loop until the last 
	for {

		for i, otherGrid := range grids {
			if grid.Equal(otherGrid) {
				return len(grids) - i, grids[i:]
			}
		}
		
		grids = append(grids, grid)
		grid = grid.Propogate()

	}

}