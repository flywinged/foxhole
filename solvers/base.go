// Copyright Clayton Brown 2020. See LICENSE file.

package solvers

import (
	"fmt"
	"foxhole/grid"
	"sync"
)

/*
	Current grids which need to be processed.
*/
var GridsToProcess = make(chan *grid.Grid)

/*
	Requirement for creating a foxhole solver
*/
type SolverFunction func(*grid.Grid) []*grid.Grid

/*
	Current hashes. Used for very quickly identifying if a certain
	arrangement of the grid has been reached before.
*/
var Hashes = make(map[int]bool)
var Lock = sync.Mutex{}

/*
	Helper function for adding resulting grids to the
	GridsToProcess array.
*/
func processGrids(grids []*grid.Grid) {

	// Pre-compute all the hashes
	hashes := []int{}
	for _, grid := range grids {
		hashes = append(hashes, grid.Hash())
	}

	// Quickly lock and insert the hashes
	toProcess := []*grid.Grid{}
	Lock.Lock()
	for i, hash := range hashes {
		if !Hashes[hash] {
			Hashes[hash] = true
			toProcess = append(toProcess, grids[i])
		}
	}
	Lock.Unlock()

	// Queue all the future grids to process
	for _, grid := range toProcess {
		GridsToProcess <- grid
	}

}

// Used to track how many solvers are currently processing
var solverWaitGroup = sync.WaitGroup{}
var solverKillChannel = make(chan bool)

/*
	Base solve function for handling
*/
func Solve(

	// The solving function to use
	solver SolverFunction,

	// Number of concurrent threads
	nSolvers int,

) {

	/*
		Create the base case where the fox can
		be anywhere in the grid.
	*/
	baseGrid := grid.CreateBlankGrid()
	for i := range baseGrid {
		baseGrid[i] = true
	}

	// Start all the solvers
	for i := 0; i < nSolvers; i++ {
		go solveRoutine(solver)
	}

	// Add the baseGrid to the grids to process to get everything started.
	GridsToProcess <- &baseGrid

}

/*
	An individual solving routine
*/
func solveRoutine(solver SolverFunction) {

	// Endless looping
	for {

		select {
		case grid := <-GridsToProcess:
			fmt.Println(grid)
		case <-solverKillChannel:
			return
		}

	}

}
