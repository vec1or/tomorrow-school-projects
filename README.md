# Forum

A Go and SQLite discussion forum with authentication, categories, comments,
likes/dislikes, and post filters.

## Run on Windows (PowerShell)

Requirements: Go, GCC/MinGW, and `CGO_ENABLED=1`.

```powershell
go version
gcc --version
go env CGO_ENABLED
go mod download
go build -o forum.exe ./cmd/web
.\forum.exe
```

Keep that PowerShell window open. After the log line `Starting server on: :8080`,
open <http://127.0.0.1:8080/> in a browser.

To use a different port:

```powershell
.\forum.exe -addr :4000
```

If you previously ran an older project version and do not need its demo data,
start with a clean schema:

```powershell
Remove-Item .\forum.db -ErrorAction SilentlyContinue
.\forum.exe
```

The database and its tables are created automatically from `schema.sql`.

## Implemented features

- Registration, login, logout, and one active session per user
- Creating and viewing posts
- One or more categories per post
- Comments for authenticated users
- Toggleable like/dislike reactions on posts and comments
- Public reaction and comment counts
- Filters by category, posts created by the current user, and posts liked by the current user
- Responsive interface and standard HTTP error responses

Docker is intentionally deferred until the final project stage.
