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
type SolverFunction func(*grid.Grid, int) []*grid.Grid

/*
	Current hashes. Used for very quickly identifying if a certain
	arrangement of the grid has been reached before.
*/
var Hashes = make(map[int]bool)
var HashLock = sync.Mutex{}

/*
	For tracking if a solution has been found
*/
var Solution *grid.Grid
var SolutionLock = sync.Mutex{}

/*
	Helper function for adding resulting grids to the
	GridsToProcess array.
*/
func processGrids(grids []*grid.Grid) {

	// Pre-compute all the hashes
	indices := []int{}
	hashes := []int{}
	for i, grid := range grids {
		hash := grid.Hash()

		// No possible locations for the fox
		if hash == 0 {
			SolutionLock.Lock()
			Solution = grid
			SolutionLock.Unlock()
		} else {
			hashes = append(hashes, hash)
			indices = append(indices, i)
		}

	}

	// Quickly lock and insert the hashes
	toProcess := []*grid.Grid{}
	HashLock.Lock()
	for i, index := range indices {
		hash := hashes[index]
		if !Hashes[hash] {
			Hashes[hash] = true
			toProcess = append(toProcess, grids[i])
		}
	}
	HashLock.Unlock()

	// Queue all the future grids to process
	SolutionLock.Lock()
	if Solution == nil {
		solverWaitGroup.Add(len(toProcess))
		for _, grid := range toProcess {
			GridsToProcess <- grid
		}
	}
	SolutionLock.Unlock()

	solverWaitGroup.Done()

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

	// Number of check that can be performed per day
	checks int,

	// Number of concurrent threads
	nSolvers int,

) {

	/*
		Create the base case where the fox can
		be anywhere in the grid.
	*/
	baseGrid := grid.CreateBlankGrid()
	for i := range baseGrid.Values {
		baseGrid.Values[i] = true
	}

	// Start all the solvers
	for i := 0; i < nSolvers; i++ {
		go solveRoutine(solver, checks)
	}

	// Add the baseGrid to the grids to process to get everything started.
	solverWaitGroup.Add(1)
	Hashes[baseGrid.Hash()] = true
	GridsToProcess <- &baseGrid

	// Await for all the processing to complete
	solverWaitGroup.Wait()

	// Once everything is completed, kill all the processing routines.
	for i := 0; i < nSolvers; i++ {
		solverKillChannel <- true
	}

}

/*
	An individual solving routine
*/
func solveRoutine(solver SolverFunction, checks int) {

	// Endless looping
mainLoop:
	for {

		select {
		case grid := <-GridsToProcess:
			newGrids := solver(grid, checks)
			go processGrids(newGrids)
			fmt.Println("wg", solverWaitGroup)
		case <-solverKillChannel:
			break mainLoop
		}

	}

	fmt.Println(Solution)

}
