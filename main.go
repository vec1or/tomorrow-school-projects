package main

import (
	//"fmt"
	"os"
	"strings"
)

func main() {
	// Создаём тестовый входной файл
	createInputFile()

	// Проверка аргументов
	if len(os.Args) != 3 {
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Чтение входного файла
	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	// Разбиваем текст на слова
	words := strings.Fields(string(data))

	// Обрабатываем текст
	words = ProcessText(words)

	// Собирам текст обратно
	result := strings.Join(words, " ")

	// Записываем результат
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}

// Создание входного файла sample.txt
func createInputFile() {
	file, err := os.Create("sample.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.")
	if err != nil {
		panic(err)
	}
}
