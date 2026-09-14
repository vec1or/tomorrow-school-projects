package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestExample00UsesAtMostSixTurns(t *testing.T) {
	input := `4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) > 6 {
		t.Fatalf("expected at most 6 turns, got %d\n%s", len(moves), strings.Join(moves, "\n"))
	}
	validateMoves(t, colony, moves)
}

func TestParallelPaths(t *testing.T) {
	input := `3
2 5 0
##start
0 1 2
##end
1 9 2
3 5 4
0-2
0-3
2-1
3-1
2-3`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) > 3 {
		t.Fatalf("expected at most 3 turns, got %d\n%s", len(moves), strings.Join(moves, "\n"))
	}
	validateMoves(t, colony, moves)
}

func TestInvalidInputs(t *testing.T) {
	tests := []string{
		"0\n##start\na 0 0\n##end\nb 1 1\na-b",
		"1\n##start\na 0 0\na-b",
		"1\n##start\na 0 0\n##end\nb x 1\na-b",
		"1\n##start\na 0 0\n##end\nb 1 1\na-c",
		"1\n##start\na 0 0\n##end\nb 1 1\na-a",
		"1\n##start\na 0 0\n##end\nb 1 1\na-b\nb-a",
		"1\n##start\na 0 0\n##end\na 1 1\na-a",
		"1\n##start\nLa 0 0\n##end\nb 1 1\nLa-b",
	}
	for _, input := range tests {
		if _, err := parseColony(strings.Split(input, "\n")); err == nil || !strings.HasPrefix(err.Error(), "ERROR: invalid data format") {
			t.Fatalf("expected invalid-data error for:\n%s\ngot: %v", input, err)
		}
	}
}

func TestUnknownCommandIsIgnored(t *testing.T) {
	input := `1
##start
start 0 0
##unknown
##end
end 1 0
start-end`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if colony.Start != "start" || colony.End != "end" {
		t.Fatalf("wrong start/end: %#v", colony)
	}
}

func TestRunPrintsInputBlankLineAndMoves(t *testing.T) {
	input := "2\n##start\nstart 0 0\n##end\nend 1 0\nstart-end\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "farm.txt")
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	output, err := run(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "start-end\n\nL1-end\nL2-end") {
		t.Fatalf("unexpected output:\n%s", output)
	}
}

func TestDeterministicOutput(t *testing.T) {
	input := `5
##start
s 0 0
a 1 1
b 1 -1
##end
e 2 0
s-a
a-e
s-b
b-e`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	var first string
	for i := 0; i < 10; i++ {
		paths, err := quickestPaths(colony)
		if err != nil {
			t.Fatal(err)
		}
		result := strings.Join(simulate(colony.Ants, colony.End, paths), "\n")
		if i == 0 {
			first = result
		} else if result != first {
			t.Fatalf("output changed between runs\nfirst:\n%s\nnow:\n%s", first, result)
		}
	}
}

func validateMoves(t *testing.T, colony Colony, lines []string) {
	t.Helper()
	position := make([]string, colony.Ants+1)
	for ant := 1; ant <= colony.Ants; ant++ {
		position[ant] = colony.Start
	}
	adjacent := make(map[string]bool)
	for _, link := range colony.Links {
		adjacent[canonicalLink(link.A, link.B)] = true
	}

	for turnIndex, line := range lines {
		moved := make(map[int]bool)
		usedTunnel := make(map[string]bool)
		nextPosition := append([]string(nil), position...)
		for _, token := range strings.Fields(line) {
			var ant int
			var room string
			if _, err := fmtSscanfMove(token, &ant, &room); err != nil {
				t.Fatalf("turn %d: invalid move %q", turnIndex+1, token)
			}
			if ant < 1 || ant > colony.Ants || moved[ant] {
				t.Fatalf("turn %d: ant %d moved more than once or is invalid", turnIndex+1, ant)
			}
			if !adjacent[canonicalLink(position[ant], room)] {
				t.Fatalf("turn %d: ant %d used a nonexistent tunnel %s-%s", turnIndex+1, ant, position[ant], room)
			}
			tunnel := canonicalLink(position[ant], room)
			if usedTunnel[tunnel] {
				t.Fatalf("turn %d: tunnel %s used more than once", turnIndex+1, tunnel)
			}
			usedTunnel[tunnel] = true
			moved[ant] = true
			nextPosition[ant] = room
		}
		occupied := make(map[string]int)
		for ant := 1; ant <= colony.Ants; ant++ {
			room := nextPosition[ant]
			if room == colony.Start || room == colony.End {
				continue
			}
			if other := occupied[room]; other != 0 {
				t.Fatalf("turn %d: ants %d and %d occupy room %s", turnIndex+1, other, ant, room)
			}
			occupied[room] = ant
		}
		position = nextPosition
	}
	for ant := 1; ant <= colony.Ants; ant++ {
		if position[ant] != colony.End {
			t.Fatalf("ant %d did not reach end; final room %s", ant, position[ant])
		}
	}
}

func fmtSscanfMove(token string, ant *int, room *string) (int, error) {
	if !strings.HasPrefix(token, "L") {
		return 0, invalidData("invalid move")
	}
	parts := strings.SplitN(token[1:], "-", 2)
	if len(parts) != 2 {
		return 0, invalidData("invalid move")
	}
	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	*ant = value
	*room = parts[1]
	return 2, nil
}

func TestProvidedExamplesMeetTurnLimits(t *testing.T) {
	tests := []struct {
		file     string
		maxTurns int
	}{
		{file: "example00.txt", maxTurns: 6},
		{file: "example01.txt", maxTurns: 8},
		{file: "example02.txt", maxTurns: 11},
		{file: "example03.txt", maxTurns: 6},
		{file: "example04.txt", maxTurns: 6},
		{file: "example05.txt", maxTurns: 8},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			lines, err := readInput(test.file)
			if err != nil {
				t.Fatal(err)
			}
			colony, err := parseColony(lines)
			if err != nil {
				t.Fatal(err)
			}
			paths, err := quickestPaths(colony)
			if err != nil {
				t.Fatal(err)
			}
			moves := simulate(colony.Ants, colony.End, paths)
			if len(moves) > test.maxTurns {
				t.Fatalf("expected at most %d turns, got %d\n%s", test.maxTurns, len(moves), strings.Join(moves, "\n"))
			}
			validateMoves(t, colony, moves)
		})
	}
}

func TestEqualLengthPathsUseInputTunnelOrder(t *testing.T) {
	input := `3
##start
s 0 0
a 1 0
b 1 1
c 1 2
##end
e 2 0
s-c
c-e
s-a
a-e
s-b
b-e`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) == 0 || moves[0] != "L1-c L2-a L3-b" {
		t.Fatalf("expected input-order tie breaking, got:\n%s", strings.Join(moves, "\n"))
	}
	validateMoves(t, colony, moves)
}

func TestShorterPathReceivesEarlierAnts(t *testing.T) {
	input := `5
##start
s 0 0
short 1 0
long1 1 1
long2 2 1
##end
e 3 0
s-long1
long1-long2
long2-e
s-short
short-e`
	colony, err := parseColony(strings.Split(input, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) == 0 || moves[0] != "L1-short L3-long1" {
		t.Fatalf("expected earliest-arrival ant assignment, got:\n%s", strings.Join(moves, "\n"))
	}
	validateMoves(t, colony, moves)
}

func TestBlankLineIsRejectedInsteadOfSilentlyRemoved(t *testing.T) {
	input := "1\n##start\ns 0 0\n\n##end\ne 1 0\ns-e\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "blank-line.txt")
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	lines, err := readInput(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parseColony(lines); err == nil {
		t.Fatal("expected an empty line inside the farm to be rejected")
	}
}

func TestExample02UsesExpectedDispatchOrder(t *testing.T) {
	lines, err := readInput("example02.txt")
	if err != nil {
		t.Fatal(err)
	}
	colony, err := parseColony(lines)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) != 11 {
		t.Fatalf("expected 11 turns, got %d\n%s", len(moves), strings.Join(moves, "\n"))
	}
	if moves[0] != "L1-3 L2-1" {
		t.Fatalf("unexpected first turn: %s", moves[0])
	}
	if moves[1] != "L2-2 L3-3 L4-1" {
		t.Fatalf("unexpected second turn: %s", moves[1])
	}
	if moves[len(moves)-1] != "L18-3 L20-3" {
		t.Fatalf("unexpected final turn: %s", moves[len(moves)-1])
	}
	validateMoves(t, colony, moves)
}

func TestExample05MatchesExpectedDispatchOrder(t *testing.T) {
	lines, err := readInput("example05.txt")
	if err != nil {
		t.Fatal(err)
	}
	colony, err := parseColony(lines)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := quickestPaths(colony)
	if err != nil {
		t.Fatal(err)
	}
	moves := simulate(colony.Ants, colony.End, paths)
	expected := []string{
		"L1-A0 L4-B0 L6-C0",
		"L1-A1 L2-A0 L4-B1 L5-B0 L6-C1",
		"L1-A2 L2-A1 L3-A0 L4-E2 L5-B1 L6-C2 L9-B0",
		"L1-end L2-A2 L3-A1 L4-D2 L5-E2 L6-C3 L7-A0 L9-B1",
		"L2-end L3-A2 L4-D3 L5-D2 L6-I4 L7-A1 L8-A0 L9-E2",
		"L3-end L4-end L5-D3 L6-I5 L7-A2 L8-A1 L9-D2",
		"L5-end L6-end L7-end L8-A2 L9-D3",
		"L8-end L9-end",
	}
	if strings.Join(moves, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("unexpected example05 output:\n%s", strings.Join(moves, "\n"))
	}
	validateMoves(t, colony, moves)
}
