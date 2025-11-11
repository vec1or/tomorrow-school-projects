package main

import (
	"strconv"
	"strings"
	"unicode"
)

// Главная функция обработки текста
func ProcessText(words []string) []string {
	i := 0
	for i < len(words) {
		if i < len(words)-1 {
			next := words[i+1]

			// HEX → DEC
			if next == "(hex)" {
				words[i] = HexToDec(words[i])
				words = removeAt(words, i+1)
				continue
			}

			// BIN → DEC
			if next == "(bin)" {
				words[i] = BinToDec(words[i])
				words = removeAt(words, i+1)
				continue
			}

			// UP → uppercase
			if strings.HasPrefix(next, "(up") {
				count := getCount(next)
				words = ChangeCase(words, i, count, "up")
				words = removeAt(words, i+1)
				continue
			}

			// LOW → lowercase
			if strings.HasPrefix(next, "(low") {
				count := getCount(next)
				words = ChangeCase(words, i, count, "low")
				words = removeAt(words, i+1)
				continue
			}

			// CAP → capitalize
			if strings.HasPrefix(next, "(cap") {
				count := getCount(next)
				words = ChangeCase(words, i, count, "cap")
				words = removeAt(words, i+1)
				continue
			}
		}
		i++
	}

	// Можно здесь добавить обработку пунктуации и кавычек
	words = FixPunctuation(words)
	return words
}

// Удаляет элемент из среза по индексу
func removeAt(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// HEX → DEC
func HexToDec(word string) string {
	n, _ := strconv.ParseInt(word, 16, 64)
	return strconv.Itoa(int(n))
}

// BIN → DEC
func BinToDec(word string) string {
	n, _ := strconv.ParseInt(word, 2, 64)
	return strconv.Itoa(int(n))
}

// Обработка регистра слов
func ChangeCase(words []string, start, count int, mode string) []string {
	if count <= 0 {
		count = 1
	}
	end := start + count
	if end > len(words) {
		end = len(words)
	}

	for i := start; i < end; i++ {
		switch mode {
		case "up":
			words[i] = strings.ToUpper(words[i])
		case "low":
			words[i] = strings.ToLower(words[i])
		case "cap":
			words[i] = Capitalize(words[i])
		}
	}
	return words
}

// Извлечение числа из скобок типа (up, 3)
func getCount(s string) int {
	// Простая реализация: если есть число, возвращаем его, иначе 1
	start := strings.Index(s, ",")
	end := strings.Index(s, ")")
	if start != -1 && end != -1 && start < end {
		numStr := strings.TrimSpace(s[start+1 : end])
		n, err := strconv.Atoi(numStr)
		if err == nil {
			return n
		}
	}
	return 1
}

// Capitalize слово
func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// Исправление пробелов вокруг знаков пунктуации
func FixPunctuation(words []string) []string {
	puncts := []string{".", ",", "!", "?", ":", ";"}
	for i := 0; i < len(words); i++ {
		for _, p := range puncts {
			if words[i] == p && i > 0 {
				words[i-1] += p
				words = removeAt(words, i)
				i--
				break
			}
		}
	}
	return words
}