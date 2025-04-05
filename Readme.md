# File Sharing & Management System (Golang)

This is a backend project built using Go, PostgreSQL, and Redis. It allows users to upload and manage files securely.

## Features

- User registration and login with JWT authentication
- Upload files to local storage and save metadata in PostgreSQL
- View uploaded files, search by name, type, and date
- Public sharing of files via a sharable URL
- Redis caching for file metadata on listing and sharing
- Background worker to delete expired files from storage and database
- Basic test written for register endpoint

## Stack Used

- Go
- PostgreSQL
- Redis

## How to Use

1. Clone the repo
2. Make sure PostgreSQL and Redis are running
3. Create `.env` file with your DB config and JWT secret
4. Run migrations to create `users` and `files` tables
5. Start the server:
   
   ```bash
   go run cmd/main.go

## Note

Files are stored locally under `/uploads`. Redis is used for caching metadata and the system includes a cleanup goroutine that deletes expired files every few minutes.

## Testing

A basic test is added under `internal/auth/handler_test.go` to check the `/register` endpoint.  
It uses Go’s standard `httptest` library and just verifies response code for a simple POST request.  
This is added to show that the project supports basic API testing.
