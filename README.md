# Groupie Tracker

A web application that displays information about music artists and bands using data from the provided https://groupietrackers.herokuapp.com/api.

You can browse all artists, check their details: members, creation date, first album and see where and when they performed.

## How to Run

```bash
go run ./cmd/web/
```

The server starts on `http://localhost:8080` by default.

You can also specify a custom port:

```bash
go run ./cmd/web/ -addr=:3000
```

## How to Test

```bash
go test ./internal/server/
```

## Project Structure

```
.
├── cmd/web/                    # Entry point (main.go)
├── internal/
│   ├── app/                    # Application setup, template cache, graceful shutdown
│   ├── config/                 # Application struct and service interface
│   ├── middleware/              # Recovery, logging, allowed methods, secure headers
│   ├── models/                 # Data models
│   ├── repository/             # Fetching data from external API
│   ├── server/                 # HTTP handlers, routes, tests
│   ├── service/                # Business logic layer
│   └── assert/                 # Test helper
└── ui/
    ├── html/                   # Go templates (base layout + pages)
    └── static/css/             # Stylesheet
```

## Architecture

The project follows a layered architecture:

- **Repository** - fetches raw JSON data from the Groupie Trackers API
- **Service** - stores and provides artist data to handlers
- **Handlers** - process HTTP requests, call the service, render templates


## Features

- Browse all artists on the main page
- View individual artist page with members, dates and concert locations
- Graceful server shutdown on Ctrl+C
- Middleware: panic recovery, request logging, method filtering, security headers
- Widget (handler-level) tests with mock service

## Built With

- Go (standard library only)
- HTML templates
- CSS

## How to test

```bash
go test ./internal/server/
go test ./internal/service/
go test ./internal/repository/
```
Or all in one command:
```bash
go test ./...
```