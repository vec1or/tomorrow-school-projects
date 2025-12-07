package main

import (
        "io/ioutil"
        "os"
        "strings"
        "testing"
)

// Тест создаёт минимальный "шрифт" для символов ' ' (space) и 'A' и проверяет рендер.
func TestRenderSimple(t *testing.T) {
        // Подготовим временный файл
        content := strings.Join([]string{
                " ",
                " ",
                " ",
                " ",
                " ",
                " ",
                " ",
                " ",
                " A ",
                " A A",
                "AAAA",
                "A A",
                "A A",
                " ",
                " ",
                " ",
        }, "\n")

        tmp, err := ioutil.TempFile("", "font-*.txt")
        if err != nil {
                t.Fatal(err)
        }
        defer os.Remove(tmp.Name())

        ioutil.WriteFile(tmp.Name(), []byte(content), 0644)

        fm, err := loadFont(tmp.Name())
        if err != nil {
                t.Fatal(err)
        }

        out, err := renderString("A", fm)
        if err != nil {
                t.Fatal(err)
        }

        if !strings.Contains(out, "AAAA") {
                t.Fatalf("unexpected render output: %q", out)
        }
}
