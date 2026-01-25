package ascii

import (
	"errors"
	"os"
	"strings"
)

func Generate(text, banner string) (string, error) {
	if banner == "" {
		banner = "standard"
	}

	path := "banners/" + banner + ".txt"

	content, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("banner not found")
	}

	lines := strings.Split(string(content), "\n")

if len(lines) < 95*8 {
	return "", errors.New("invalid banner format")
}
	asciiMap := make(map[rune][]string)

	for i := 32; i <= 126; i++ {
		start := (i - 32) * 9
		asciiMap[rune(i)] = lines[start : start+8]
	}

	var result strings.Builder
	rows := strings.Split(text, "\n")

	for _, row := range rows {
		for i := 0; i < 8; i++ {
			for _, ch := range row {
				if val, ok := asciiMap[ch]; ok {
					result.WriteString(val[i])
				}
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}
