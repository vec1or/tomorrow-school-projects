# Forum

A web forum built with Go, SQLite, HTML, and CSS.

Users can register, publish posts, select categories, write comments, react with
likes or dislikes, and filter discussions. The frontend does not use external
UI frameworks.

## Features

- Registration, login, logout, and cookie-based sessions
- Secure password hashing with bcrypt
- One active session per user
- Creating and viewing posts
- Multiple categories per post
- Comments on posts
- Likes and dislikes for posts and comments
- Public reaction and comment counters
- Filtering by category
- Filtering posts created by the current user
- Filtering posts liked by the current user
- Responsive interface
- Persistent SQLite storage
- Docker support

## Technology stack

- Go 1.26.4
- SQLite and `go-sqlite3`
- Go HTML templates
- HTML, CSS, and JavaScript
- Docker

## Project structure

## Project structure

```text
forum/
├── .dockerignore
├── .gitignore
├── Dockerfile
├── README.md
├── go.mod
├── go.sum
├── schema.sql
│
├── cmd/
│   └── web/
│       ├── context.go
│       ├── handlers.go
│       ├── helpers.go
│       ├── main.go
│       ├── middleware.go
│       ├── routes.go
│       └── templates.go
│
├── internal/
│   ├── models/
│   │   ├── categories.go
│   │   ├── comments.go
│   │   ├── errors.go
│   │   ├── posts.go
│   │   ├── reactions.go
│   │   ├── sessions.go
│   │   └── users.go
│   │
│   └── validator/
│       └── validator.go
│
└── ui/
    ├── html/
    │   ├── base.tmpl
    │   │
    │   ├── pages/
    │   │   ├── create.tmpl
    │   │   ├── home.tmpl
    │   │   ├── login.tmpl
    │   │   ├── signup.tmpl
    │   │   └── view.tmpl
    │   │
    │   └── partials/
    │       └── nav.tmpl
    │
    └── static/
        ├── css/
        │   └── main.css
        │
        ├── img/
        │   ├── favicon.ico
        │   └── logo.png
        │
        └── js/
            └── main.js
```

* `cmd/web` contains the HTTP server, application routes, handlers, middleware, and template configuration.
* `internal/models` contains the SQLite models and database queries.
* `internal/validator` contains form validation helpers.
* `ui/html` contains the Go HTML templates.
* `ui/static` contains frontend assets such as CSS, JavaScript, and images.
* `schema.sql` creates the database tables, indexes, and initial categories.
* `Dockerfile` defines the container build and runtime environment.

The `forum.db` database file and compiled binaries are generated locally at runtime and are not stored in the Git repository.

## Run with Docker (recommended)

Make sure Docker is running:

```bash
docker --version
```

### 1. Build the image

```bash
docker build -t forum-app .
```

### 2. Start the forum

```bash
docker run --name forum-container -p 8080:8080 -v forum-data:/app/data forum-app
```

Open the application in a browser:

<http://localhost:8080>

The application runs in the foreground. Press `Ctrl+C` to stop it.

### Start the existing container again

```bash
docker start -a forum-container
```

## Run locally

Local execution requires Go 1.26.4 or newer and a C compiler because the SQLite
driver uses CGO.

Verify Go and CGO:

```bash
go version
go env CGO_ENABLED
```

`CGO_ENABLED` must be `1`. If it is disabled, enable it with:

```bash
go env -w CGO_ENABLED=1
```

### Windows

Requirements:

- Go 1.26.4 or newer
- GCC/MinGW

Verify GCC:

```powershell
gcc --version
```

Download dependencies and run the application:

```powershell
go mod download
go run ./cmd/web
```

Alternatively, build and run a Windows executable:

```powershell
go build -o forum.exe ./cmd/web
.\forum.exe
```

Open <http://localhost:8080>.

### macOS

Requirements:

- Go 1.26.4 or newer
- Xcode Command Line Tools

Install the command-line tools if necessary:

```bash
xcode-select --install
```

If Homebrew is installed, Go can be installed with:

```bash
brew install go
```

Download dependencies and run the application:

```bash
go mod download
go run ./cmd/web
```

Alternatively, build and run a macOS executable:

```bash
go build -o forum ./cmd/web
./forum
```

Open <http://localhost:8080>.

### Linux

Requirements:

- Go 1.26.4 or newer
- GCC

Verify GCC:

```bash
gcc --version
```

Download dependencies and run the application:

```bash
go mod download
go run ./cmd/web
```

Alternatively, build and run a Linux executable:

```bash
go build -o forum ./cmd/web
./forum
```

Open <http://localhost:8080>.

## Authors
#lsabit #almuratov #bkalzhan #dkalzhan