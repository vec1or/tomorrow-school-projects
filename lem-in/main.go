package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(invalidData("expected exactly one input file"))
		return
	}

	result, err := run(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func run(filePath string) (string, error) {
	lines, err := readInput(filePath)
	if err != nil {
		return "", err
	}

	colony, err := parseColony(lines)
	if err != nil {
		return "", err
	}

	paths, err := quickestPaths(colony)
	if err != nil {
		return "", err
	}

	moves := simulate(colony.Ants, colony.End, paths)
	if len(moves) == 0 {
		return "", invalidData("no ant movements generated")
	}

	return strings.Join(colony.Input, "\n") + "\n\n" + strings.Join(moves, "\n"), nil
}
