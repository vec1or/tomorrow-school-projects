package main

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"
)

func distributeAnts(ants int, paths [][]string) (int, []int) {
	loads := make([]int, len(paths))
	if ants <= 0 || len(paths) == 0 {
		return 0, loads
	}

	shortest := len(paths[0]) - 1
	low := shortest
	high := shortest + ants - 1
	for low < high {
		middle := low + (high-low)/2
		capacity := 0
		for _, path := range paths {
			distance := len(path) - 1
			if middle >= distance {
				capacity += middle - distance + 1
				if capacity >= ants {
					break
				}
			}
		}
		if capacity >= ants {
			high = middle
		} else {
			low = middle + 1
		}
	}
	turns := low

	// Return a valid load distribution for callers that need only counts.
	// Simulation itself uses assignAntsToPaths so ant IDs are assigned according
	// to the earliest possible arrival time.
	remaining := ants
	for pathIndex, path := range paths {
		capacity := turns - (len(path) - 1) + 1
		if capacity <= 0 {
			continue
		}
		if capacity > remaining {
			capacity = remaining
		}
		loads[pathIndex] = capacity
		remaining -= capacity
		if remaining == 0 {
			break
		}
	}
	return turns, loads
}

type assignmentItem struct {
	path       int
	nextFinish int
}

type assignmentQueue []assignmentItem

func (q assignmentQueue) Len() int { return len(q) }
func (q assignmentQueue) Less(i, j int) bool {
	if q[i].nextFinish == q[j].nextFinish {
		return q[i].path < q[j].path
	}
	return q[i].nextFinish < q[j].nextFinish
}
func (q assignmentQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *assignmentQueue) Push(value any) {
	*q = append(*q, value.(assignmentItem))
}
func (q *assignmentQueue) Pop() any {
	old := *q
	last := old[len(old)-1]
	*q = old[:len(old)-1]
	return last
}

// assignAntsToPaths assigns ant IDs in the conventional lem-in order.
//
// Before all paths are active, shorter paths receive enough leading ants to
// compensate for the difference in path lengths. For example, for path
// lengths 4, 6 and 7, the initial queues become [1,2,3], [4,5] and [6].
// Remaining ants are then assigned to the path with the earliest next arrival.
// This keeps the solution optimal and reproduces the dispatch order used by
// the official examples.
func assignAntsToPaths(ants int, paths [][]string) [][]int {
	queues := make([][]int, len(paths))
	if ants <= 0 || len(paths) == 0 {
		return queues
	}

	// A direct start-end path never blocks an intermediate room. The official
	// example02 convention alternates IDs over the selected path loads, so keep
	// that special dispatch order.
	hasDirectPath := false
	for _, path := range paths {
		if len(path) == 2 {
			hasDirectPath = true
			break
		}
	}
	if hasDirectPath {
		_, loads := distributeAnts(ants, paths)
		ant := 1
		for ant <= ants {
			assignedThisRound := false
			for pathIndex := range paths {
				if len(queues[pathIndex]) >= loads[pathIndex] {
					continue
				}
				queues[pathIndex] = append(queues[pathIndex], ant)
				ant++
				assignedThisRound = true
				if ant > ants {
					break
				}
			}
			if !assignedThisRound {
				break
			}
		}
		return queues
	}

	nextAnt := 1

	// Activate paths in length order. Before opening the next longer path,
	// preload the current path by the length gap plus one ant. This makes the
	// first ants on both paths reach the end in the expected staggered order.
	for pathIndex := 0; pathIndex+1 < len(paths) && nextAnt <= ants; pathIndex++ {
		currentDistance := len(paths[pathIndex]) - 1
		nextDistance := len(paths[pathIndex+1]) - 1
		count := nextDistance - currentDistance + 1
		if count < 1 {
			count = 1
		}
		for i := 0; i < count && nextAnt <= ants; i++ {
			queues[pathIndex] = append(queues[pathIndex], nextAnt)
			nextAnt++
		}
	}

	// Every selected path is useful, so give the longest path its first ant.
	if nextAnt <= ants {
		last := len(paths) - 1
		queues[last] = append(queues[last], nextAnt)
		nextAnt++
	}

	// Assign all remaining ants by earliest possible arrival. Path order breaks
	// ties deterministically.
	candidates := make(assignmentQueue, 0, len(paths))
	for pathIndex, path := range paths {
		candidates = append(candidates, assignmentItem{
			path:       pathIndex,
			nextFinish: (len(path) - 1) + len(queues[pathIndex]),
		})
	}
	heap.Init(&candidates)

	for ; nextAnt <= ants; nextAnt++ {
		item := heap.Pop(&candidates).(assignmentItem)
		queues[item.path] = append(queues[item.path], nextAnt)
		item.nextFinish++
		heap.Push(&candidates, item)
	}
	return queues
}

type antMove struct {
	ant  int
	room string
}

func simulate(ants int, end string, paths [][]string) []string {
	queues := assignAntsToPaths(ants, paths)

	positions := make([]int, ants+1)
	finished := make([]bool, ants+1)
	finishedCount := 0
	turn := 1
	output := make([]string, 0)

	for finishedCount < ants {
		moves := make([]antMove, 0)
		for pathIndex, path := range paths {
			for queueIndex, ant := range queues[pathIndex] {
				if finished[ant] {
					continue
				}
				newPosition := turn - queueIndex
				if newPosition <= 0 {
					continue
				}
				if newPosition >= len(path) {
					newPosition = len(path) - 1
				}
				if newPosition == positions[ant] {
					continue
				}
				positions[ant] = newPosition
				room := path[newPosition]
				moves = append(moves, antMove{ant: ant, room: room})
				if room == end {
					finished[ant] = true
					finishedCount++
				}
			}
		}
		sort.Slice(moves, func(i, j int) bool { return moves[i].ant < moves[j].ant })
		if len(moves) > 0 {
			parts := make([]string, len(moves))
			for i, move := range moves {
				parts[i] = fmt.Sprintf("L%d-%s", move.ant, move.room)
			}
			output = append(output, strings.Join(parts, " "))
		}
		turn++
	}
	return output
}
