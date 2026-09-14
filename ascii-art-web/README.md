# ASCII Art Web

A web application written in Go that converts user-provided text into ASCII art.

The application provides a browser-based interface where users can enter text, select a banner style, and view the generated ASCII art.

## Features

- Simple web interface
- Converts text into ASCII art
- Supports letters, numbers, spaces, and special characters
- Handles multiple lines
- Provides three banner styles:
  - Standard
  - Shadow
  - Thinkertoy
- Returns appropriate HTTP status codes
- Uses only the Go standard library on the backend

## Technologies

- Go
- HTML
- CSS
- Go HTML templates
- Go `net/http` package

## Requirements

- [Go](https://go.dev/doc/install)
- A modern web browser

Verify that Go is installed:

```bash
go version
```

## Installation

Clone the repository:

```bash
git clone https://github.com/vec1or/ascii-art-web.git
cd ascii-art-web
```

## Running the Application

Start the web server:

```bash
go run .
```

Open the following address in your browser:

```text
http://localhost:8080
```

Stop the server by pressing:

```text
Ctrl+C
```

## Usage

1. Open the application in your browser.
2. Enter text into the input field.
3. Select a banner style.
4. Click the generate button.
5. The generated ASCII art will appear on the page.

## Banner Styles

The application supports the following banner templates:

| Banner | Description |
| --- | --- |
| `standard` | Classic ASCII art style |
| `shadow` | Characters with a shadow effect |
| `thinkertoy` | Thin and compact character style |

## HTTP Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/` | Displays the main page |
| `POST` | `/ascii-art` | Processes the form and generates ASCII art |

## HTTP Status Codes

| Status | Meaning |
| --- | --- |
| `200 OK` | The request was processed successfully |
| `400 Bad Request` | The submitted input or banner is invalid |
| `404 Not Found` | The requested page or resource does not exist |
| `500 Internal Server Error` | An unexpected server error occurred |

## Implementation Details

The application works as follows:

1. The Go HTTP server starts and listens for incoming requests.
2. A `GET` request to `/` renders the main HTML page.
3. The user submits text and a selected banner through an HTML form.
4. A `POST` request sends the form data to the `/ascii-art` handler.
5. The server validates the submitted text and banner.
6. The selected banner file is loaded and parsed.
7. Each input character is mapped to its eight-line ASCII representation.
8. The generated result is rendered inside the HTML page.
9. Invalid requests are handled using the appropriate HTTP status code.

Go templates are used to pass the generated ASCII art from the server to the webpage.

## Testing

Run all tests with:

```bash
go test ./...
```

## Learning Objectives

This project demonstrates:

- Building an HTTP server in Go
- Working with `GET` and `POST` requests
- Processing HTML forms
- Using Go HTML templates
- Serving static files
- Reading and processing banner files
- HTTP status code handling
- Separating application logic from presentation

## Acknowledgements

This project is based on the [01-edu ASCII Art Web subject](https://public.01-edu.org/subjects/ascii-art-web/).
