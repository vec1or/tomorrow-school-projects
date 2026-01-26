package ascii

import (
    "fmt"
    "os"
    "strings"
)

func Generate(text, banner string) (string, error) {
    file := "banners/" + banner + ".txt"

    data, err := os.ReadFile(file)
    if err != nil {
        return "", err
    }

    lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
    if len(lines) != 760 {
        return "", fmt.Errorf("invalid banner format: %d lines", len(lines))
    }

    // 🔥 ВАЖНО: Windows fix
    text = strings.ReplaceAll(text, "\r", "")
    textLines := strings.Split(text, "\n")

    var result strings.Builder

    for _, line := range textLines {
        if line == "" {
            result.WriteString("\n")
            continue
        }

        for row := 0; row < 8; row++ {
            for _, r := range line {
                if r < 32 || r > 126 {
                    return "", fmt.Errorf("unsupported character: %q", r)
                }
                index := int(r-32)*8 + row
                result.WriteString(lines[index])
            }
            result.WriteString("\n")
        }
    }

    return result.String(), nil
}
