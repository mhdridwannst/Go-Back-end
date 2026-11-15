# Go Backend Engineering Log

A structured repository documenting my progression in backend engineering with **Go (Golang)**. Rather than relying immediately on high-level abstractions, this repository focuses on building systems from the ground up—starting from raw `net/http` handlers and direct SQL queries to clean architecture, authentication, and containerized deployments.

![Go Version](https://shields.io)
![PostgreSQL](https://shields.io)
![License](https://shields.io)

---

## Repository Structure

```text
.
├── Backend/
│   └── api-baju/          # Standard library REST API with PostgreSQL persistence
└── README.md
```

---

## Engineering Milestones

### Phase 1: Core Fundamentals & Persistence
- [x] **HTTP Server & Routing:** Handling raw requests, dynamic path parsing, and HTTP methods via `net/http`.
- [x] **JSON Serialization:** Struct tags, request body decoding, and standardized response headers.
- [x] **Relational Database Integration:** Connecting Go's `database/sql` interface to PostgreSQL using `lib/pq` with parameterized queries to prevent SQL injection.
- [ ] **Configuration Management:** Externalizing database credentials via `.env` and environment variables.

### Phase 2: Architecture & Security
- [ ] **Layered Architecture:** Decoupling transport (`handler`), business logic (`service`), and data access (`repository`).
- [ ] **Middleware Pipeline:** Implementing request logging, panic recovery, and CORS handling.
- [ ] **Authentication & Authorization:** Password hashing with `bcrypt` and stateless session management via JWT.

### Phase 3: Reliability & Deployment
- [ ] **Automated Testing:** Table-driven unit tests and HTTP handler integration tests (`httptest`).
- [ ] **Containerization:** Multi-stage `Dockerfile` and `docker-compose` setup for reproducible environments.

---

## Modules

### `Backend/api-baju`
A dependency-light RESTful service managing apparel inventory, built purely with Go's standard library and PostgreSQL.

**Endpoints**

| Method | Path | Description | Status Codes |
| :--- | :--- | :--- | :--- |
| `GET` | `/shirts` | Retrieve all items ordered by ID | `200`, `500` |
| `GET` | `/shirts/{id}` | Retrieve a single item by ID | `200`, `400`, `404`, `500` |
| `POST` | `/shirts` | Create a new item record | `201`, `400`, `500` |
| `PUT` | `/shirts/{id}` | Update an existing item | `200`, `400`, `404`, `500` |
| `DELETE` | `/shirts/{id}` | Remove an item record | `200`, `400`, `404`, `500` |

**Local Setup**

1. **Initialize PostgreSQL Schema**
   ```sql
   CREATE DATABASE db_baju;

   CREATE TABLE shirts (
       id SERIAL PRIMARY KEY,
       brand VARCHAR(100) NOT NULL,
       size VARCHAR(10) NOT NULL
   );
   ```

2. **Install Dependencies & Run**
   ```bash
   cd Backend/api-baju
   go mod tidy
   go run main.go
   ```
   *(Ensure `connStr` in `main.go` matches your local PostgreSQL port and credentials).*

---

## Verification

```bash
# Fetch collection
curl -i http://localhost:8081/shirts

# Insert record
curl -i -X POST http://localhost:8081/shirts \
  -H "Content-Type: application/json" \
  -d '{"brand": "Fashion", "size": "XL"}'

# Update record
curl -i -X PUT http://localhost:8081/shirts/5 \
  -H "Content-Type: application/json" \
  -d '{"brand": "Fashion", "size": "L"}'

# Delete record
curl -i -X DELETE http://localhost:8081/shirts/5
```
