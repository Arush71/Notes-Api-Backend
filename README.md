# Notes API with Versioning (Go)

A backend API built in Go that supports authenticated note management with full version history and rollback support.

This project was built to learn some backend engineering concepts including transactions, concurrency-safe state changes, SQL modeling, and authorization.

---

## 🚀 Features

- User authentication (JWT-based)
- Create / get / update / delete notes
- Protected routes with ownership checks
- Automatic versioning on every update
- List all versions of a note
- Fetch a specific version
- Rollback to any previous version (creates a new version)
- Atomic transactions for consistency

---

## 🛠 Tech Stack

- Go
- PostgreSQL
- sqlc(for sql to go)
- net/http

---

## 🌐 API Routes

All note routes are protected using authentication middleware.
Kindly check out cmd/server/routes.go if you wanna see the routes for yourself.

### Auth
POST /app/auth/...

### Notes
GET    /app/notes  
POST   /app/notes  
GET    /app/notes/{id}  
PUT    /app/notes/{id}  
DELETE /app/notes/{id}  

### Versioning
GET  /app/notes/{id}/versions  
GET  /app/notes/{id}/versions/{version_number}  
POST /app/notes/{id}/versions/{version_number}/rollback  

Versioning behavior:
- Every note update creates a new version.
- Rollback restores an old version by creating a new version with the old content.
- History is immutable.

---

## ▶️ Running Locally

### Prerequisites
- Go installed
- PostgreSQL running

### Setup

1. Clone the repository:
git clone <your-repo-url>  
cd <repo-name>  

2. Set environment variables:
DB_URL=postgres://...  
TOKEN_SECRET=your_secret 

3. Run migrations (if applicable)

4. Start the server:
go run main.go  

## Notes
- Versioning is implemented using full snapshots for simplicity and correctness.
- Rollback creates a new version instead of mutating history.