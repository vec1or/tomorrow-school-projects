# GO-RELOADED

## For Audit Questions
- https://github.com/01-edu/public/tree/master/subjects/go-reloaded/audit

## Overview
Go-Reloaded is a command-line tool that reads text files and applies specific transformations using various commands within the text. Finally, it outputs the transformed text.

## Features
- Converting a heximal number to decimal // command (hex)
- Converting a binary number to decimal // command (bin)
- Сonverting words from lowercase to uppercase and vice versa // command (low)/(up)
- Capitalizing first letter // command (cap)
- Punctuation fixing
- Article manipulations

## Pay Attention
For (low), (up), (cap) if a number appears next to it, like so: (low, <number>) it turns the previously specified number of words in lowercase, uppercase or capitalized accordingly. (Ex: "This is so exciting (up, 2)" -> "This is SO EXCITING")

## Project Structure

```
go-reloaded/
├── testdata/
│   ├── tests.txt //test_file
├── go.mod
├── handlers.go //functions
├── main.go
├── result.txt //output_text
├── sample.txt //input_text
└── processor_test.go //running_test_file
```

## Usage
Write text with commands in sample.txt file. Then type:
```go run . sample.txt result.txt```
Check result.txt for final text.

### Example
```
$ cat sample.txt
it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.
$ go run . sample.txt result.txt
$ cat result.txt
It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.
$ cat sample.txt
Simply add 42 (hex) and 10 (bin) and you will see the result is 68.
$ go run . sample.txt result.txt
$ cat result.txt
Simply add 66 and 2 and you will see the result is 68.
$ cat sample.txt
There is no greater agony than bearing a untold story inside you.
$ go run . sample.txt result.txt
$ cat result.txt
There is no greater agony than bearing an untold story inside you.
$ cat sample.txt
Punctuation tests are ... kinda boring ,what do you think ?
$ go run . sample.txt result.txt
$ cat result.txt
Punctuation tests are... kinda boring, what do you think?
```

## For Testing
Type: ```go test -v```