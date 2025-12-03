package main

import (
    "bufio"
    "os"
    "strings"
    "testing"
)

func readTests(path string) ([][2]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var tests [][2]string
    var lines []string

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    var block []string
    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            if len(block) >= 2 {
                input := strings.TrimSpace(block[0])
                output := strings.TrimSpace(block[1])
                tests = append(tests, [2]string{input, output})
            }
            block = nil
        } else {
            block = append(block, line)
        }
    }

    if len(block) >= 2 {
        input := strings.TrimSpace(block[0])
        output := strings.TrimSpace(block[1])
        tests = append(tests, [2]string{input, output})
    }

    return tests, nil
}

func TestFullTxtFile(t *testing.T) {
    tests, err := readTests("testdata/tests.txt")
    if err != nil {
        t.Fatalf("Failed to load tests: %v", err)
    }

    for idx, pair := range tests {
        input := pair[0]
        want := pair[1]

        got := ProcessFullText(input)

        if got != want {
            t.Errorf("\nTEST #%d FAILED\nInput:\n%s\n\nGot:\n%s\n\nExpected:\n%s\n",
                idx+1, input, got, want,
            )
        }
    }
}
