// Copyright Clayton Brown 2020. See LICENSE file.

package solvers

import (
	"foxhole/grid"
)

/*
	Returns int combinations
*/
func Combinations(n, c int) []map[int]bool {

	// Create the return array
	combinations := []map[int]bool{}

	// Base case
	if c == 1 {
		for i := 0; i < n; i++ {
			combinations = append(combinations, map[int]bool{i: true})
		}
		return combinations
	}

	// Otherwise, we're in a recursive case for combinations
	for i := n - 1; i >= c - 1; i-- {
		recursiveCombinations := Combinations(i, c - 1)
		for _, comb := range recursiveCombinations {
			comb[i] = true
			combinations = append(combinations, comb)
		}
	}

	return combinations
}

// Helper for determining if two sets are equal
func SetHash(s map[int]bool) int {

	hash := 0
	for power := range s {
		hash += grid.POWERS[power]
	}

	return hash

}

// Set copy
func SetCopy(s map[int]bool) map[int]bool {

	newS := map[int]bool{}
	for k := range s {
		newS[k] = true
	}

	return newS

}

// Modifies the first set
func SetUnion(s1, s2 map[int]bool) {
	for k := range s2 {
		s1[k] = true
	}
}