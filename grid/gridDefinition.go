// Copyright Clayton Brown 2020. See LICENSE file.

package grid

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
func createLinearGrid(n int) *GridDefinition {

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

func init() {

	// Original foxhole problem definition
	BaseGrid = createLinearGrid(5)
}
