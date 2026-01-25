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

    // Убираем лишние переводы строк в конце
    lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
    if len(lines) != 760 {
        return "", fmt.Errorf("invalid banner format: %d lines", len(lines))
    }

    var result strings.Builder

    // Проходим по блокам строк (0..7, 8..15, 16..23 ...)
    for row := 0; row < 8; row++ {            // 8 строк на один ряд текста
        for _, r := range text {
            if r < 32 || r > 126 {
                return "", fmt.Errorf("unsupported character: %q", r)
            }
            index := int(r-32)*8 + row       // индекс нужной строки символа
            result.WriteString(lines[index])
        }
        result.WriteString("\n")               // переход на следующую строку баннера
    }

    return result.String(), nil
	
}
