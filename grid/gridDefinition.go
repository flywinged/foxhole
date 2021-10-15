// Copyright Clayton Brown 2020. See LICENSE file.

package grid

import (
	"math"
)

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

	/*
		Function for determining the location of a cell
		in the values array given 3D coordinates.
	*/
	getIndex := func(location []int) int {
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
	getLocation := func(index int) []int {
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

	// Generate the connections array first
	connections := [][]int{}
	for i := 0; i < totalCells; i++ {
		location := getLocation(i)

		// Generate each of the modified locations based on the current location
		connectionLocations := []int{}
		for i, x := range location {

			// Move in the negative direction along an axis
			if x > 0 {
				newLocation := make([]int, len(location))
				copy(newLocation, location)
				newLocation[i] = x - 1
				connectionLocations = append(connectionLocations, getIndex(newLocation))
			}

			// Move in the positive direction along an axis
			if x < dimensionLengths[i]-1 {
				newLocation := make([]int, len(location))
				copy(newLocation, location)
				newLocation[i] = x + 1
				connectionLocations = append(connectionLocations, getIndex(newLocation))
			}

		}

		// Add the generated locations for this cell to connections
		connections = append(connections, connectionLocations)

	}

	// Generate all the symettries
	baseSymmetry := []int{}
	for i := 0; i < totalCells; i++ {
		baseSymmetry = append(baseSymmetry, i)
	}
	symmetries := [][]int{baseSymmetry}

	return &GridDefinition{
		Connections: connections,
		Symmetries:  symmetries,
	}
}

func init() {

	// Original foxhole problem definition
	// BaseGrid = CreateLinearGrid(5)

	BaseGrid = CreatePrismGrid([]int{3, 3})

}
