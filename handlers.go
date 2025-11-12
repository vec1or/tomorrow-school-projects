package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ProcessFullText reads the whole text and applies transformations
func ProcessFullText(s string) string {
	s = strings.TrimSpace(s)

	// Tokenize: words (including hyphens), commands, punctuation, quotes
	re := regexp.MustCompile(`\([^)]*\)|[\p{L}0-9-]+|[.,!:;?]+|'`)
	tokens := re.FindAllString(s, -1)

	// Apply commands like (hex), (bin), (cap, n), (up, n), (low, n)
	tokens = applyCommands(tokens)

	// Assemble final text with punctuation and quotes rules
	return assemble(tokens)
}

// applyCommands handles commands modifying previous words
func applyCommands(tokens []string) []string {
	cmdRE := regexp.MustCompile(`^\(\s*(?i)(hex|bin|up|low|cap)\s*(?:,\s*(\d+)\s*)?\)$`)

	i := 0
	for i < len(tokens) {
		m := cmdRE.FindStringSubmatch(tokens[i])
		if m != nil {
			cmd := strings.ToLower(m[1])
			count := 1
			if m[2] != "" {
				if n, err := strconv.Atoi(m[2]); err == nil && n > 0 {
					count = n
				}
			}

			switch cmd {
			case "hex":
				if j := findPrevWordIndex(tokens, i); j != -1 {
					tokens[j] = hexToDec(tokens[j])
				}
			case "bin":
				if j := findPrevWordIndex(tokens, i); j != -1 {
					tokens[j] = binToDec(tokens[j])
				}
			case "up", "low", "cap":
				// collect previous `count` "word-like" tokens (not commands or punctuation)
				indices := []int{}
				j := i - 1
				for j >= 0 && len(indices) < count {
					if tokens[j] != "" && !strings.HasPrefix(tokens[j], "(") && !isPunct(tokens[j]) && tokens[j] != "'" {
						indices = append([]int{j}, indices...) // prepend to preserve order
					}
					j--
				}
				for _, idx := range indices {
					switch cmd {
					case "up":
						tokens[idx] = strings.ToUpper(tokens[idx])
					case "low":
						tokens[idx] = strings.ToLower(tokens[idx])
					case "cap":
						tokens[idx] = Capitalize(tokens[idx])
					}
				}
			}

			// Remove the command token
			tokens = append(tokens[:i], tokens[i+1:]...)
			continue // do not increment i
		}
		i++
	}

	return tokens
}

// findPrevWordIndex returns the index of previous "word-like" token
func findPrevWordIndex(tokens []string, from int) int {
	for j := from - 1; j >= 0; j-- {
		if tokens[j] != "" && !strings.HasPrefix(tokens[j], "(") && !isPunct(tokens[j]) && tokens[j] != "'" {
			return j
		}
	}
	return -1
}

func hexToDec(s string) string {
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return s
	}
	return strconv.FormatInt(n, 10)
}

func binToDec(s string) string {
	n, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s
	}
	return strconv.FormatInt(n, 10)
}

// Capitalize first letter, lowercase rest, preserve hyphened parts
func Capitalize(word string) string {
	if word == "" {
		return word
	}

	if strings.Contains(word, "-") {
		parts := strings.SplitN(word, "-", 2)
		runes := []rune(parts[0])
		runes[0] = unicode.ToUpper(runes[0])
		for i := 1; i < len(runes); i++ {
			runes[i] = unicode.ToLower(runes[i])
		}
		return string(runes) + "-" + parts[1]
	}

	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// assemble builds final string with proper punctuation, quotes, and articles
func assemble(tokens []string) string {
	var b strings.Builder
	inSingleQuote := false

	isWord := regexp.MustCompile(`^[\p{L}0-9-]+$`)

	for i, tok := range tokens {
		if tok == "'" {
			if !inSingleQuote {
				b.WriteString("'")
				inSingleQuote = true
				continue
			} else {
				b.WriteString("'")
				inSingleQuote = false
				continue
			}
		}

		if isPunct(tok) {
			b.WriteString(tok)
			if i+1 < len(tokens) && !isPunct(tokens[i+1]) && tokens[i+1] != "'" {
				b.WriteString(" ")
			}
			continue
		}

		// handle article 'a' -> 'an'
		if (tok == "a" || tok == "A") && i+1 < len(tokens) && isWord.MatchString(tokens[i+1]) {
			next := strings.ToLower(tokens[i+1])
			if len(next) > 0 && strings.ContainsRune("aeiouh", rune(next[0])) {
				if tok == "A" {
					tok = "An"
				} else {
					tok = "an"
				}
			}
		}

		if b.Len() > 0 {
			last := b.String()[b.Len()-1]
			if last != ' ' && last != '\'' {
				b.WriteString(" ")
			}
		}

		b.WriteString(tok)
	}

	return strings.TrimSpace(b.String())
}

// isPunct returns true if token is punctuation
func isPunct(tok string) bool {
	return regexp.MustCompile(`^[.,!:;?]+$`).MatchString(tok)
}
