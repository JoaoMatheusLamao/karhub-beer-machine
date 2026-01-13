# Karhub Beer Machine

Backend service to manage beer styles and always serve the **best beer for a given temperature**, following Clean Architecture principles.

---

## 📌 Overview

This project implements a backend service that:

* Manages beer styles and their ideal temperature ranges (CRUD)
* Selects the best beer style for a given temperature using a **domain rule**
* Retrieves a playlist related to the selected beer style (Spotify integration stub)
* Exposes a RESTful API
* Provides a CLI for administrative tasks (seed)
* Is fully tested (domain, application, infrastructure, HTTP)

---

## 🧠 Business Rule (Core Logic)

Each beer style has:

* Minimum temperature
* Maximum temperature

### Selection rule

Given an input temperature **T**:

1. Calculate the **average temperature** of each style

   ```
   avg = (min + max) / 2
   ```

2. Select the style whose **average temperature is closest to T**
3. If there is a tie:

   * Choose by **alphabetical order (lexicographical)**
4. The temperature **does not need to be inside the range**

This rule lives in the **domain layer** and is fully unit tested.

---

## 🏗 Architecture

The project follows **Clean Architecture** with clear separation of responsibilities:

```
cmd/
  api/        -> HTTP entrypoint
  cli/        -> CLI entrypoint (Cobra)
internal/
  domain/     -> Entities, business rules, repository interfaces
  application -> Use cases (CRUD + best beer selection)
  interfaces/ -> HTTP handlers and routes
  infrastructure/
    persistence/
      memory/ -> Ristretto-based in-memory repository
```

### Dependency rule

```
interfaces → application → domain
infrastructure → application → domain
```

The domain is completely independent of frameworks, HTTP, database, or external APIs.

---

## 🔐 Environment Variables

The application uses environment variables for configuration.

### Required

```env
HTTP_PORT=8080
SPOTIFY_CLIENT_ID=your_spotify_client_id 
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
```

## 🚀 How to Run

### Prerequisites

* Go **1.21+**
* Git

---

### Run the HTTP API

```bash
docker compose up --build -d
```

The API will be available at:

```
http://localhost:8080
```

Health check:

```bash
GET /health
```

---

## 🧰 CLI (Administrative)

The project provides a CLI built with **Cobra**.

### Seed beer styles

```bash
docker compose run --rm \
  -e API_BASE_URL=http://api:8080 \
  --entrypoint /app/cli \
  api seed
```

This populates the repository with predefined beer styles.

---

## 🌐 HTTP API

### Create beer style

```http
POST /beer-styles
```

```json
{
  "id": "1",
  "name": "IPA",
  "minTemp": -7,
  "maxTemp": 10
}
```

---

### Update beer style

```http
PUT /beer-styles/{id}
```

```json
{
  "name": "Imperial IPA",
  "minTemp": -8,
  "maxTemp": 12
}
```

---

### Delete beer style

```http
DELETE /beer-styles/{id}
```

---

### List beer styles

```http
GET /beer-styles
```

---

### Find best beer for a temperature (core endpoint)

```http
POST /beer-styles/best
```

```json
{
  "temperature": -7
}
```

Response:

```json
{
  "beerStyle": "Dunkel",
  "playlist": {
    "name": "Dunkel Playlist",
    "tracks": []
  }
}
```

---

## 🧪 Tests

The project contains **unit and integration-style tests** covering:

* Domain rules
* Use cases
* In-memory repository (Ristretto v2)
* HTTP handlers and routing

### Run all tests

```bash
go test ./...
```

### Verbose output

```bash
go test -v ./...
```

### Disable test cache

```bash
go test -count=1 ./...
```

---

## 📦 In-Memory Cache

The repository uses **Ristretto v2 (generic)** as an in-memory cache:

* Thread-safe
* High-performance
* Explicit key tracking for iteration
* Suitable for local development and tests

The repository is injected via interface and can be replaced by Postgres without changing the core logic.

---

## 🎧 Spotify Integration

The project includes a full integration with the Spotify Web API using the
**Client Credentials Flow**, following the official Spotify documentation and
the reference implementation from the `spotify/v2` Go SDK.

### Current limitation

At the time of development, Spotify is temporarily blocking the creation of new
applications, displaying the message:

> “New integrations are currently on hold”

Because of this limitation, it may not be possible to generate new
`Client ID` and `Client Secret` credentials.

### Fallback strategy (important)

To avoid blocking the execution of the application, the system implements a
**fallback mechanism**:

* When valid Spotify credentials are available, the real Spotify integration is used
* When credentials are missing or Spotify is unavailable, the application
  automatically falls back to a **stub (mock) implementation**

This ensures that:

* The API is always functional
* Business rules can be fully validated
* The architecture remains production-ready

No structural changes are required to switch between stub and real integration.

---

## ⚡ Caching Strategy

To reduce calls to external services and improve performance, the application
implements an in-memory cache using **Ristretto v2**.

### What is cached

* Spotify playlists, indexed by beer style name

### Cache details

* Cache key format: `spotify:playlist:<beer_style>`
* Time-to-live (TTL): **10 minutes**
* Cache is completely transparent to the application layer

The cache is implemented using the **Decorator Pattern**, wrapping the
`SpotifyGateway` without introducing coupling between layers.

---

## 🧩 Design Decisions

* Clean Architecture with explicit use cases
* Table-driven tests
* Manual mocks (no mock frameworks)
* No business logic in HTTP handlers
* CLI reuses application use cases
* Composition root in `main.go`

---

## 🐳 Running with Docker

The project includes a multi-stage Dockerfile and a docker-compose setup.

### Start the API

```bash
docker compose up --build
docker compose run --rm api /app/cli seed
```

---

## ℹ️ Notes

The Spotify integration code is production-ready and fully implemented.
Due to a temporary limitation on Spotify's side, new application credentials
may not be creatable at this moment.

For this reason, a stub implementation is used as a fallback to guarantee
application availability and testability.
