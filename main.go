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

// Программа ожидает один аргумент — строку для вывода (может содержать реальные символы '\n').
// Опционально флаг -font указывает файл баннера (по умолчанию "standard.txt").

const (
        defaultFont = "standard.txt"
        firstASCII  = 32 // пробел
)

func main() {
        fontFile := flag.String("font", defaultFont, "path to banner font file")
        flag.Parse()

        // Входная строка — либо аргумент, либо пустая строка
        var input string
        if flag.NArg() > 0 {
                input = flag.Arg(0)
        } else {
                // Если аргумента нет — читаем всё из stdin (trim trailing newline)
                b, err := io.ReadAll(os.Stdin)
                if err != nil {
                        fmt.Fprintln(os.Stderr, "error reading stdin:", err)
                        os.Exit(1)
                }
                input = string(b)
        }

        // Обрабатываем CRLF и убираем завершающий один перенос строки, если он пришёл в конце и не нужен
        input = strings.ReplaceAll(input, "\r\n", "\n")
        if len(input) > 0 && input[len(input)-1] == '\n' {
                // оставляем, потому что пользователь мог иметь в конце пустую строку — это семантически важно
        }

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

// loadFont читает файл шрифта, формат: подряд идут блоки по 8 строк на символ,
// порядок символов — начиная с ASCII 32 (пробел) и далее. Количество символов определяется длиной файла.
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

        nChars := len(lines) / 8
        fm := make(fontMap)
        for i := 0; i < nChars; i++ {
                var block [8]string
                for r := 0; r < 8; r++ {
                        block[r] = lines[i*8+r]
                }
                runeVal := rune(firstASCII + i)
                fm[runeVal] = block
        }
        return fm, nil
}

// renderString принимает вход который может содержать '\n' и возвращает готовый ASCII-арт.
func renderString(s string, fm fontMap) (string, error) {
        // Разбиваем на логические строки, сохраняя пустые
        lines := strings.Split(s, "\n")
        var outLines []string

        for li, ln := range lines {
                // Для каждой из 8 рядов формируем строку
                // Если ln == "" — выввести пустую строку-высоту (8 пустых строк)
                if ln == "" {
                        for i := 0; i < 8; i++ {
                                outLines = append(outLines, "")
                        }
                        continue
                }

                // Собираем 8 строк
                rowBuf := make([]strings.Builder, 8)
                for _, ch := range ln {
                        block, ok := fm[ch]
                        if !ok {
                                // если символ не найден — попробуем знак вопроса, иначе пробел
                                if q, qok := fm['?']; qok {
                                        block = q
                                } else if sp, spok := fm[' ']; spok {
                                        block = sp
                                } else {
                                        return "", fmt.Errorf("character %q not in font and neither '?' nor ' ' available", ch)
                                }
                        }
                        for r := 0; r < 8; r++ {
                                rowBuf[r].WriteString(block[r])
                        }
                }
                // Добавляем 8 сформированных строк
                for r := 0; r < 8; r++ {
                        outLines = append(outLines, rowBuf[r].String())
                }
                // После каждой логической строки добавляем пустая строк (по примеру в задании — две пустые строки в конце).
                // Но в примерах видно что после вывода высоты 8 добавляются 2 полностью пустые строки.
                // Чтобы соответствовать, добавим 2 пустых строк-пустышки(пустые текстовые строки).
                if li != len(lines)-1 {
                        // при конце ввода — в поведении примера остаются 2 пустые строки; сохраним тоже
                        outLines = append(outLines, "")
                        outLines = append(outLines, "")
                }
        }

        return strings.Join(outLines, "\n") + "\n", nil
}
