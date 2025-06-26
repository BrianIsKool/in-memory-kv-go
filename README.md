# go-cache-server

A simple in-memory key-value cache server written in Go with TTL support and automatic expiration.

## 🚀 Features

- `POST /set` — Add a key-value pair with TTL (in seconds)
- `GET /get?key=` — Retrieve value by key
- `POST /remove?key=` — Remove entry by key
- ⏱ Automatic removal of expired entries via background goroutine

## 🧪 Getting Started

### Run the server

```bash
go run main.go
