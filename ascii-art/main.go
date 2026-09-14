package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)
const (
	defaultFont = "standard.txt"
	firstASCII  = 32
)

func main() {
	fontFile := flag.String("font", defaultFont, "path to banner font file")
	flag.Parse()

	var input string
	if flag.NArg() > 0 {
		input = flag.Arg(0)
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading stdin:", err)
			os.Exit(1)
		}
		input = string(b)
	}

	// Обрабатываем CRLF и реальные символы \n
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, `\n`, "\n")
	input = strings.TrimSuffix(input, "\n") // убираем лишний перенос в конце

	fontPath := *fontFile
	if !filepath.IsAbs(fontPath) {
		cwd, _ := os.Getwd()
		fontPath = filepath.Join(cwd, fontPath)
	}

	font, err := loadFont(fontPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load font:", err)
		os.Exit(2)
	}

	out, err := renderString(input, font)
	if err != nil {
		fmt.Fprintln(os.Stderr, "render error:", err)
		os.Exit(3)
	}

	fmt.Print(out)
}

type fontMap map[rune][8]string
func loadFont(path string) (fontMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, errors.New("empty font file")
	}
	if len(lines)%8 != 0 {
		return nil, fmt.Errorf("font file has %d lines which is not a multiple of 8", len(lines))
	}

	fm := make(fontMap)
	for i := 0; i < len(lines)/8; i++ {
		var block [8]string
		for r := 0; r < 8; r++ {
			block[r] = lines[i*8+r]
		}
		fm[rune(firstASCII+i)] = block
	}
	return fm, nil
}
func renderString(s string, fm fontMap) (string, error) {
	lines := strings.Split(s, "\n")
	var out []string

	for _, ln := range lines {
		if ln == "" {
			out = append(out, "")
			continue
		}

		rows := make([]strings.Builder, 8)
		for _, ch := range ln {
			block, ok := fm[ch]
			if !ok {
				block = fm[' ']
			}
			for i := 0; i < 8; i++ {
				rows[i].WriteString(block[i])
			}
		}

		for i := 0; i < 8; i++ {
			out = append(out, rows[i].String())
		}
	}
	for len(out) > 0 && out[0] == "" {
		out = out[1:]
	}

	return strings.Join(out, "\n") + "\n", nil
}
