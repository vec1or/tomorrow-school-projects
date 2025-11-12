// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		// helper to create sample input when running without args (for quick manual testing)
		_ = ioutil.WriteFile("sample.txt", []byte("this is well-known (cap)"), 0644)
		fmt.Println("Wrote sample.txt. Run: go run . sample.txt result.txt")
		return
	}

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . <input-file> <output-file>")
		os.Exit(2)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	out := ProcessFullText(string(data))

	if err := ioutil.WriteFile(outputFile, []byte(out), 0644); err != nil {
		panic(err)
	}
}