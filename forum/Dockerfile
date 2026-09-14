FROM golang:1.26.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /forum ./cmd/web


FROM alpine:3.23

RUN addgroup -S forum && adduser -S -G forum forum

WORKDIR /app

COPY --from=builder /forum ./forum
COPY schema.sql ./schema.sql
COPY ui ./ui

RUN mkdir -p /app/data && chown -R forum:forum /app

USER forum

EXPOSE 8080

CMD ["./forum", "-addr", ":8080", "-dsn", "file:/app/data/forum.db?_foreign_keys=on"]