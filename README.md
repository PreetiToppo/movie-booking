# High-Concurrency Cinema Booking System

A production-ready backend for a movie ticket booking system built in **Go**, designed to solve the classic "double-booking" concurrency problem.

## 🏗️ Architecture & Tech Stack
- **Language:** Go (Standard Library `net/http`, `sync`)
- **Database/Cache:** Redis (for distributed locking and TTL-based seat holds)
- **Pattern:** Repository Pattern with Dependency Injection (allows swapping Redis for In-Memory/Mutex stores)
- **Frontend:** Vanilla JS/HTML/CSS with polling for near real-time updates

## 🚀 Key Features
- **Pessimistic Locking:** Utilizes Redis `SET NX` atomic operations to ensure exactly one user can hold a specific seat at a time.
- **Checkout Sessions (TTL):** Seats are held for 2 minutes. If the user doesn't confirm, Redis automatically releases the seat.
- **Concurrency Tested:** Validated using `go test -race` with 100,000 simulated concurrent users attempting to book the same seat.

## 🛠️ How to Run Locally
1. Start Redis:
   ```bash
   docker-compose up -d

2. Run the Go server:
    ```bash
    go run cmd/main.go

3. Open your browser and go to http://localhost:8080 to view the UI.