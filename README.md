# URL Shortener

A simple URL shortener service built with Go. This project provides an API to shorten URLs and redirect users to the original URLs using a short identifier.

## Features

- Shorten long URLs into short, unique identifiers.
- Redirect users from the short URL to the original URL.
- Caching with Redis for faster lookups.
- Persistent storage using SQLite.

## Prerequisites

Before running the project, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.24 or higher)
- [SQLite](https://www.sqlite.org/download.html)
- [Redis](https://redis.io/download)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/adi-253/url_shortener.git
   cd url_shortener
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Start Redis:

   ```bash
   redis-server
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### 1. Shorten URL

**Endpoint:** `POST /post_url`  
**Request Body:**

```json
{
  "url": "https://example.com"
}
```

**Response:**

```json
{
  "short_url": "http://localhost:8080/<short_id>"
}
```

### 2. Redirect to Original URL

**Endpoint:** `GET /{short_id}`  
**Response:** Redirects to the original URL.

## Project Structure

- `main.go`: Entry point of the application.
- `api/server.go`: Contains the HTTP server and route handlers.
- `database/db.go`: Handles SQLite database operations.
- `cache/redis.go`: Manages Redis caching.

## How It Works

1. A user sends a POST request to `/post_url` with a long URL.
2. The server generates a short identifier using the `shortid` library.
3. The long URL and short identifier are stored in SQLite.
4. The short URL is returned to the user.
5. When a user accesses the short URL, the server:
   - Checks Redis for the original URL.
   - If not found, fetches it from SQLite and caches it in Redis.
   - Redirects the user to the original URL.

## Logging

The application uses structured logging with `slog` to log important events, such as URL creation and redirection.

