// Copyright Clayton Brown 2020. See LICENSE file.

package solvers

import (
	"runtime/pprof"
	"runtime"
	"os"
	"unsafe"
	"fmt"
	"time"
	"foxhole/grid"
	"sync"
)

const DEBUG = true

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
	Used for tracking the current depth
*/
var CurrentDepth = 0
var DepthLock = sync.Mutex{}
var Depths = make([]sync.WaitGroup, 2048)

var LastDepthLock = sync.Mutex{}
var LastDepthTime = time.Now()
var LastDepth = -1

var TestCounter = 0
var TestLock = sync.Mutex{}

/*
	Function for resetting meta values
*/
func reset() {

	GridsToProcess = make(chan *grid.Grid)
	Hashes = make(map[int]bool)

	Solution = nil

	CurrentDepth = 0
	DepthLock = sync.Mutex{}
	Depths = make([]sync.WaitGroup, 2048)

	LastDepthLock = sync.Mutex{}
	LastDepthTime = time.Now()
	LastDepth = -1
}

/*
	Helper function for adding resulting grids to the
	GridsToProcess array.
*/
func processGrids(grids []*grid.Grid, gridSize int) {

	TestLock.Lock()
	TestCounter += len(grids)
	TestLock.Unlock()

	Depths[gridSize].Done()
	Depths[gridSize].Wait()

	Depths[gridSize + 1].Add(1)

	// Handle printing and information
	LastDepthLock.Lock()
	if LastDepth != gridSize && DEBUG {
		fmt.Println()
		fmt.Println("Completed Depth", gridSize)
		tNow := time.Now()
		fmt.Println("Time to Complete Depth:", fmt.Sprintf("%.2f", float64(tNow.Sub(LastDepthTime)) / float64(time.Second)), "seconds")

		hashSize := 0
		for i := range Hashes {
			hashSize += int(unsafe.Sizeof(i))
		}

		fmt.Println("Total Hashes:", len(Hashes), hashSize / 1024)
		
		TestLock.Lock()
		fmt.Println("Total Grids:", TestCounter)
		totalStorage := 0
		for _, grid := range grids {
			totalStorage += int(unsafe.Sizeof(grid.Checks))
		}
		fmt.Println("Average Size:", float64(totalStorage) / float64(len(grids)))
		TestLock.Unlock()

		f, err := os.Create("memory.prof")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println(err)
		}

		LastDepthTime = tNow
		LastDepth = gridSize

	}
	LastDepthLock.Unlock()

	/*
		If a solution hase already been found, we don't
		need to do any more processing. Just break out,
		handling the wait groups appropriately.
	*/
	if Solution != nil {
		solverWaitGroup.Done()
		return
	}

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
		hash := hashes[i]
		if !Hashes[hash] {
			Hashes[hash] = true
			toProcess = append(toProcess, grids[index])
		}
	}
	HashLock.Unlock()

	// Queue all the future grids to process
	SolutionLock.Lock()
	if Solution == nil {
		solverWaitGroup.Add(len(toProcess))
		for _, grid := range toProcess {
			Depths[len(grid.Checks)].Add(1)
			GridsToProcess <- grid
		}
	}

	SolutionLock.Unlock()

	Depths[gridSize + 1].Done()
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
	repetition, grids := grid.BaseGrid.RepeatingGrid()

	// Try for each solution type
	for i := 0; i < repetition; i++ {

		// Log how long a solve is taking
		t0 := time.Now()

		baseGrid := grids[i]

		if DEBUG {
			fmt.Println("Base Grid:", baseGrid)
		}

		// Reset parameters
		reset()

		// Start all the solvers
		for i := 0; i < nSolvers; i++ {
			go solveRoutine(solver, checks)
		}

		// Add the baseGrid to the grids to process to get everything started.
		solverWaitGroup.Add(1)
		Depths[0].Add(1)
		Hashes[baseGrid.Hash()] = true
		GridsToProcess <- baseGrid

		// Await for all the processing to complete
		solverWaitGroup.Wait()

		// Once everything is completed, kill all the processing routines.
		for i := 0; i < nSolvers; i++ {
			solverKillChannel <- true
		}

		// Log the total time
		fmt.Println()
		if Solution != nil {
			fmt.Println("Solution", Solution)
			fmt.Println("Solution Length", len(Solution.Checks))
			fmt.Println("Total Hashes:", len(Hashes))
			fmt.Println("Time to Process:", fmt.Sprintf("%.2f", float64(time.Since(t0)) / float64(time.Second)), "seconds")
		} else {
			fmt.Println("No Solutions Found")
			fmt.Println("Time to Process:", fmt.Sprintf("%.2f", float64(time.Since(t0)) / float64(time.Second)), "seconds")
		}

		if DEBUG {
			break
		}
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
		case gridToProcess := <-GridsToProcess:
			newGrids := solver(gridToProcess, checks)
			go processGrids(newGrids, len(gridToProcess.Checks))
		case <-solverKillChannel:
			break mainLoop
		}

	}

}
