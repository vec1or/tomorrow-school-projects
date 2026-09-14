# ASCII Art

A command-line application written in Go that converts text into a graphical representation using ASCII characters.

## Features

- Converts text into ASCII art
- Supports letters, numbers, spaces, and special characters
- Handles newline sequences (`\n`)
- Uses banner templates to render characters
- Uses only the Go standard library

## Requirements

- [Go](https://go.dev/doc/install)

Verify that Go is installed:

```bash
go version
```

## Installation

Clone the repository:

```bash
git clone https://github.com/vec1or/ascii-art.git
cd ascii-art
```

## Usage

Pass the text you want to convert as a command-line argument:

```bash
go run . "Hello"
```

Example output:

```text
 _    _          _   _
| |  | |        | | | |
| |__| |   ___  | | | |   ___
|  __  |  / _ \ | | | |  / _ \
| |  | | |  __/ | | | | | (_) |
|_|  |_|  \___| |_| |_|  \___/

```

### Numbers and special characters

```bash
go run . "Hello 123!"
```

### Multiple lines

Use `\n` to create multiple lines:

```bash
go run . "Hello\nWorld"
```

### Empty input

```bash
go run . ""
```

## Implementation Details

The program works as follows:

1. Reads the input string from the command-line arguments.
2. Loads the ASCII banner template from a text file.
3. Splits the input into lines using `\n`.
4. Finds the eight-line graphical representation of every character.
5. Combines the corresponding rows of all characters.
6. Prints the generated ASCII art to the terminal.

The banner contains representations of printable ASCII characters. Each character is eight lines high and is separated from the next character by an empty line.

## Testing

Run all tests with:

```bash
go test ./...
```

## Learning Objectives

This project demonstrates:

- Command-line argument processing
- File reading in Go
- String and rune manipulation
- Working with ASCII character indexes
- Data formatting and terminal output
- Basic Go project organization

## Acknowledgements

This project is based on the [01-edu ASCII Art subject](https://public.01-edu.org/subjects/ascii-art/).
