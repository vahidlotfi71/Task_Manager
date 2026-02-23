# Task Manager REST API

A production-ready Task Management API built with Go (Gin), PostgreSQL.

## Features
- **RESTful CRUD** for Tasks with Soft Delete & Restore.
- **Advanced Filtering & Pagination** (by Status, Assignee).
- **PostgreSQL** Database with GORM.
- **Dockerized** multi-stage environment.
- **Test-Driven** (Unit & Integration tests with Mocking).

## 🛠 Tech Stack
Go 1.24, Gin, GORM, PostgreSQL, Redis, Prometheus, Docker.

## ⚙️ How to Run (Docker Compose)
1. Clone the repository.
2. Run the services:
```bash
   docker-compose up -d --build
   
