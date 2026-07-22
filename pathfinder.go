package main

import (
	"container/heap"
	"sort"
	"strings"
)

const infinity = int(^uint(0) >> 2)

type flowEdge struct {
	to          int
	rev         int
	cap         int
	originalCap int
	cost        int
	tunnel      bool
	roomTo      string
}

type flowGraph [][]flowEdge

type queueItem struct {
	node int
	dist int
}

type minQueue []queueItem

func (q minQueue) Len() int { return len(q) }
func (q minQueue) Less(i, j int) bool {
	if q[i].dist == q[j].dist {
		return q[i].node < q[j].node
	}
	return q[i].dist < q[j].dist
}
func (q minQueue) Swap(i, j int)   { q[i], q[j] = q[j], q[i] }
func (q *minQueue) Push(value any) { *q = append(*q, value.(queueItem)) }
func (q *minQueue) Pop() any {
	old := *q
	last := old[len(old)-1]
	*q = old[:len(old)-1]
	return last
}

func quickestPaths(colony Colony) ([][]string, error) {
	graph, source, sink := buildFlowNetwork(colony)
	potentials := make([]int, len(graph))
	bestTurns := infinity
	bestTotalLength := infinity
	var bestPaths [][]string

	startDegree := 0
	endDegree := 0
	for _, link := range colony.Links {
		if link.A == colony.Start || link.B == colony.Start {
			startDegree++
		}
		if link.A == colony.End || link.B == colony.End {
			endDegree++
		}
	}
	maxPaths := colony.Ants
	if startDegree < maxPaths {
		maxPaths = startDegree
	}
	if endDegree < maxPaths {
		maxPaths = endDegree
	}

	for flow := 1; flow <= maxPaths; flow++ {
		if !augmentOne(graph, source, sink, potentials) {
			break
		}
		paths, ok := extractFlowPaths(graph, source, sink, colony.Start, flow)
		if !ok {
			return nil, invalidData("could not construct valid paths")
		}
		sortPaths(paths, colony)
		turns, _ := distributeAnts(colony.Ants, paths)
		totalLength := 0
		for _, path := range paths {
			totalLength += len(path) - 1
		}
		if turns < bestTurns || (turns == bestTurns && totalLength < bestTotalLength) {
			bestTurns = turns
			bestTotalLength = totalLength
			bestPaths = clonePaths(paths)
		}
	}

	if len(bestPaths) == 0 {
		return nil, invalidData("no path between ##start and ##end")
	}
	return bestPaths, nil
}

func buildFlowNetwork(colony Colony) (flowGraph, int, int) {
	roomNames := make([]string, 0, len(colony.Rooms))
	for name := range colony.Rooms {
		roomNames = append(roomNames, name)
	}
	sort.Strings(roomNames)

	inNode := make(map[string]int, len(roomNames))
	outNode := make(map[string]int, len(roomNames))
	for i, name := range roomNames {
		inNode[name] = 2 * i
		outNode[name] = 2*i + 1
	}
	graph := make(flowGraph, 2*len(roomNames))

	for _, name := range roomNames {
		if name == colony.Start || name == colony.End {
			continue
		}
		addFlowEdge(graph, inNode[name], outNode[name], 1, 0, false, "")
	}

	// Keep tunnel insertion order. It is part of deterministic tie-breaking:
	// equally good paths should follow the order in which the farm was given,
	// rather than an unrelated alphabetical order.
	for _, link := range colony.Links {
		if link.A != colony.End && link.B != colony.Start {
			addFlowEdge(graph, outNode[link.A], inNode[link.B], 1, 1, true, link.B)
		}
		if link.B != colony.End && link.A != colony.Start {
			addFlowEdge(graph, outNode[link.B], inNode[link.A], 1, 1, true, link.A)
		}
	}

	return graph, outNode[colony.Start], inNode[colony.End]
}

func addFlowEdge(graph flowGraph, from, to, capacity, cost int, tunnel bool, roomTo string) {
	forward := flowEdge{
		to:          to,
		rev:         len(graph[to]),
		cap:         capacity,
		originalCap: capacity,
		cost:        cost,
		tunnel:      tunnel,
		roomTo:      roomTo,
	}
	reverse := flowEdge{
		to:   from,
		rev:  len(graph[from]),
		cost: -cost,
	}
	graph[from] = append(graph[from], forward)
	graph[to] = append(graph[to], reverse)
}

func augmentOne(graph flowGraph, source, sink int, potentials []int) bool {
	dist := make([]int, len(graph))
	parentNode := make([]int, len(graph))
	parentEdge := make([]int, len(graph))
	for i := range dist {
		dist[i] = infinity
		parentNode[i] = -1
		parentEdge[i] = -1
	}
	dist[source] = 0
	queue := &minQueue{{node: source, dist: 0}}
	heap.Init(queue)

	for queue.Len() > 0 {
		item := heap.Pop(queue).(queueItem)
		if item.dist != dist[item.node] {
			continue
		}
		for edgeIndex := range graph[item.node] {
			edge := graph[item.node][edgeIndex]
			if edge.cap <= 0 {
				continue
			}
			reducedCost := edge.cost + potentials[item.node] - potentials[edge.to]
			candidate := item.dist + reducedCost
			if candidate < dist[edge.to] {
				dist[edge.to] = candidate
				parentNode[edge.to] = item.node
				parentEdge[edge.to] = edgeIndex
				heap.Push(queue, queueItem{node: edge.to, dist: candidate})
			}
		}
	}

	if parentNode[sink] == -1 {
		return false
	}
	for node := range potentials {
		if dist[node] < infinity {
			potentials[node] += dist[node]
		}
	}

	for node := sink; node != source; node = parentNode[node] {
		from := parentNode[node]
		edgeIndex := parentEdge[node]
		reverseIndex := graph[from][edgeIndex].rev
		graph[from][edgeIndex].cap--
		graph[node][reverseIndex].cap++
	}
	return true
}

func extractFlowPaths(graph flowGraph, source, sink int, start string, expected int) ([][]string, bool) {
	remaining := make([][]int, len(graph))
	for node := range graph {
		remaining[node] = make([]int, len(graph[node]))
		for edgeIndex, edge := range graph[node] {
			if edge.originalCap > 0 {
				flow := edge.originalCap - edge.cap
				if flow > 0 {
					remaining[node][edgeIndex] = flow
				}
			}
		}
	}

	paths := make([][]string, 0, expected)
	for sourceEdgeIndex := range graph[source] {
		if remaining[source][sourceEdgeIndex] <= 0 {
			continue
		}

		path := []string{start}
		current := source
		visited := make([]bool, len(graph))
		for current != sink {
			if visited[current] {
				return nil, false
			}
			visited[current] = true

			chosen := -1
			for edgeIndex := range graph[current] {
				if remaining[current][edgeIndex] > 0 {
					chosen = edgeIndex
					break
				}
			}
			if chosen == -1 {
				return nil, false
			}

			remaining[current][chosen]--
			edge := graph[current][chosen]
			if edge.tunnel {
				path = append(path, edge.roomTo)
			}
			current = edge.to
		}
		if len(path) < 2 {
			return nil, false
		}
		paths = append(paths, path)
	}

	if len(paths) != expected {
		return nil, false
	}
	return paths, true
}

func sortPaths(paths [][]string, colony Colony) {
	linkOrder := make(map[string]int, len(colony.Links))
	for index, link := range colony.Links {
		linkOrder[canonicalLink(link.A, link.B)] = index
	}

	pathOrder := func(path []string) []int {
		order := make([]int, 0, len(path)-1)
		for index := 0; index+1 < len(path); index++ {
			order = append(order, linkOrder[canonicalLink(path[index], path[index+1])])
		}
		return order
	}

	sort.SliceStable(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}

		left := pathOrder(paths[i])
		right := pathOrder(paths[j])
		for index := range left {
			if left[index] != right[index] {
				return left[index] < right[index]
			}
		}
		return strings.Join(paths[i], "\x00") < strings.Join(paths[j], "\x00")
	})
}

func clonePaths(paths [][]string) [][]string {
	copyOfPaths := make([][]string, len(paths))
	for i := range paths {
		copyOfPaths[i] = append([]string(nil), paths[i]...)
	}
	return copyOfPaths
}
