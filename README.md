<div align="center">

# Go Books API

**REST API for managing books** — built with Go, Gin, GORM & MySQL.

Full observability stack included: **OpenTelemetry**, **Prometheus**, **Grafana**, **Loki** & **Promtail**.

[![Tests](https://github.com/aguszunda/mvp_books/actions/workflows/pr-validation.yml/badge.svg)](https://github.com/aguszunda/mvp_books/actions/workflows/pr-validation.yml)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.11-00E0C4?logo=gin&logoColor=white)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-1.31-00ADD8?logo=database&logoColor=white)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [API Endpoints](#api-endpoints)
- [Bulk Insert](#bulk-insert)
- [Monitoring & Observability](#monitoring--observability)
- [Architecture](#architecture)
- [Testing](#testing)
- [Configuration](#configuration)
- [Make Commands](#make-commands)
- [Bruno API Client](#bruno-api-client)
- [Project Structure](#project-structure)

---

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Go Books API                                 │
│                                                                     │
│  ┌──────────┐    ┌──────────┐    ┌──────────────┐    ┌──────────┐  │
│  │  Gin      │───▶│ Service  │───▶│  Repository   │───▶│  MySQL   │  │
│  │ Handlers  │    │ (Logic)  │    │  (GORM/MySQL) │    │  :3306   │  │
│  └──────────┘    └──────────┘    └──────────────┘    └──────────┘  │
│       │                                                         │   │
│       │ OTel SDK ──▶ Prometheus Exporter (:2223)                │   │
│       │                                                         │   │
│       └── Docker logs ──▶ Promtail ──▶ Loki (:3100)             │   │
│                                                                 │   │
│  Prometheus (:9090) ◀── scrapes metrics                         │   │
│  Grafana    (:3000) ◀── dashboards + log search                 │   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Quick Start

### Prerequisites

- **Go** 1.24+
- **Docker** & **Docker Compose**

### 1. Setup environment

```bash
cp environment_values.example environment_values
```

The defaults work out of the box:

```bash
DB_USER=developer          DB_PASSWORD=admin
DB_HOST=db                 DB_NAME=books_db
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=books_db    MYSQL_USER=developer
MYSQL_PASSWORD=admin       METRICS_PORT=2223
LOG_LEVEL=info
```

### 2. Start everything

```bash
docker compose up -d --build
```

### 3. Verify

```bash
# Create a book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Refactoring", "author": "Martin Fowler"}'

# List books
curl http://localhost:8080/books

# Open Grafana
open http://localhost:3000
```

### Services

| Service | URL | Description |
|---------|-----|-------------|
| **API** | http://localhost:8080 | Books REST API |
| **Metrics** | http://localhost:2223/metrics | OpenTelemetry Prometheus exporter |
| **Prometheus** | http://localhost:9090 | Metrics scraping & storage |
| **Grafana** | http://localhost:3000 | Dashboards & log visualization |

> **Credentials:** Grafana `admin` / `admin`

### Stop

```bash
docker compose down          # Stop containers
docker compose down -v       # Stop + delete database
```

---

## API Endpoints

### Create a Book

```
POST /books
```

**Request:**

```json
{
    "title": "The Refactoring",
    "author": "Martin Fowler"
}
```

**Response** `201 Created`**:

```json
{
    "id": 1,
    "title": "The Refactoring",
    "author": "Martin Fowler"
}
```

### Get All Books

```
GET /books
```

### Get Book by ID

```
GET /books/:id
```

---

## Bulk Insert

A standalone Go script for inserting thousands of books via CSV. Uses the `POST /books` endpoint over HTTP — no direct DB access needed.

### Usage

```bash
# Direct
go run scripts/bulk_insert.go <csv_file> [api_url]

# Via Make
make bulk-insert CSV=scripts/books_sample.csv
```

### CSV Format

```csv
title,author
The Refactoring,Martin Fowler
Clean Code,Robert C. Martin
Domain-Driven Design,Eric Evans
```

### Example Output

```
=== Bulk Insert Started ===
CSV file: scripts/books_sample.csv
API URL:  http://localhost:8080
Total records to insert: 2500
----------------------------
[1/2500] OK - Created book (id=1): 'Golden Chronicles' by 'Sarah Clark'
[2/2500] OK - Created book (id=2): 'Lost Odyssey' by 'James Green'
[3/2500] FAIL (HTTP 400) - 'Bad' by '': {"error": "author is required"}
...
----------------------------
=== Bulk Insert Completed ===
Total:   2500 | Success: 2498 | Failed: 2
Duration: 45.23s | Avg per record: 18ms
```

> See [scripts/README.md](scripts/README.md) for full documentation.

---

## Monitoring & Observability

### Architecture

```
API (:8080) ──HTTP──▶ Prometheus (:9090) ──query──▶ Grafana (:3000)
     │                                                    ▲
     └──OTel SDK──▶ Prometheus exporter (:2223)           │
                                                           │
Docker logs ──Promtail──▶ Loki (:3100) ───────────────────┘
```

**Metrics pipeline:** The API instruments HTTP requests via [OpenTelemetry](https://opentelemetry.io/). A Prometheus exporter exposes metrics on `:2223` in Prometheus format. Prometheus scrapes every 10s. Grafana displays a pre-configured 7-panel dashboard.

**Logs pipeline:** [Promtail](https://grafana.com/docs/loki/latest/send-data/promtail/) ships Docker container logs to [Loki](https://grafana.com/oss/loki/), which stores and indexes them. Grafana connects to Loki as a data source for log search and correlation with metrics.

### Available Metrics

| OTel Metric | Prometheus Name | Type | Description |
|-------------|----------------|------|-------------|
| `http.server.request.count` | `http.server.request.count_total` | Counter | Total requests by method, route, status |
| `http.server.request.duration` | `http.server.request.duration_seconds` | Histogram | Request latency (5ms — 10s buckets) |
| `http.server.request.inflight` | `http.server.request.inflight` | Gauge | Concurrent requests in progress |

**Labels:** `http.method`, `http.route`, `http.response.status_code`

> Query in Prometheus/Grafana using `{__name__="http.server.request.count_total"}` (dots are preserved).

### Grafana Dashboard

Auto-provisioned at **Dashboards → API Monitoring → Books API Monitoring**.

| Panel | Type | Description |
|-------|------|-------------|
| Request Rate (req/s) | Timeseries | Requests per second by method and route |
| Request Latency (P50 / P95 / P99) | Timeseries | Latency percentiles over time |
| Error Rate (5xx) | Timeseries | Percentage of 5xx errors by route |
| In-Flight Requests | Timeseries | Currently active concurrent requests |
| Total Requests by Status Code | Piechart | Distribution of responses by status code |
| Avg Request Duration by Endpoint | Bargauge | Average response time per route |
| Request Duration Heatmap | Heatmap | Latency distribution across time |

---

## Architecture

```
cmd/api/main.go              Entrypoint: wires telemetry, middleware, routes, DB

internal/
  controller/                HTTP handlers (Gin)
  service/                   Business logic (use cases)
  repository/                Data persistence (GORM/MySQL)
  domain/                    Entities & repository interfaces
  middleware/
    logging.go               Structured request logging (slog)
    metrics.go               HTTP metrics (OTel: count, duration, in-flight)

platform/
  connection/                DB initialization (GORM)
  logger/                    Structured JSON logger (slog)
  metrics/                   OTel metric definitions
  telemetry.go               OTel SDK init + Prometheus exporter

config/
  prometheus/                Scrape config (targets app:2223)
  grafana/
    provisioning/            Datasource & dashboard provisioning
    dashboards/              Pre-configured monitoring dashboard
  loki/                      Loki server config
  promtail/                  Log shipper config

scripts/
  init.sql                   DB schema + seed data
  bulk_insert.go             Bulk insert script (Go)
  books_sample.csv           2500 example books for load testing
```

### Data Flow

```
Request → Gin Middleware (logging, metrics) → Handler → Service → Repository → MySQL
                                                                       │
                                                              OTel SDK ─┘──▶ Prometheus
```

---

## Testing

### Run all tests

```bash
make test
```

### Coverage report

```bash
make coverage
```

Coverage excludes infrastructure code (`platform/`, `internal/middleware/`, `cmd/`, `mocks/`, `internal/domain/`, `scripts/`). Business logic (`internal/service/`, `internal/repository/`) must maintain **≥85%** coverage — enforced in CI.

### Find untested functions

```bash
make uncovered
```

---

## Configuration

Environment variables are managed via the `environment_values` file (gitignored).

| Variable | Docker Default | Local Default | Description |
|----------|---------------|---------------|-------------|
| `DB_USER` | developer | developer | MySQL user |
| `DB_PASSWORD` | admin | admin | MySQL password |
| `DB_HOST` | db | localhost | DB hostname |
| `DB_NAME` | books_db | books_db | Database name |
| `MYSQL_ROOT_PASSWORD` | rootpassword | — | MySQL root password |
| `MYSQL_DATABASE` | books_db | — | Auto-created database |
| `MYSQL_USER` | developer | — | Auto-created MySQL user |
| `MYSQL_PASSWORD` | admin | — | Auto-created MySQL password |
| `METRICS_PORT` | 2223 | 2223 | Prometheus metrics endpoint port |
| `LOG_LEVEL` | info | info | Log level (DEBUG, INFO, WARN, ERROR) |

> `DB_HOST` must be `db` in Docker, `localhost` when running locally.

The database is auto-initialized on first startup via `scripts/init.sql`.

---

## Make Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Build and start everything in Docker |
| `make docker-down` | Stop all containers |
| `make run-local` | DB + monitoring in Docker, API runs locally |
| `make run` | API locally (requires local MySQL on :3306) |
| `make test` | Run all tests |
| `make coverage` | Run tests with filtered coverage report |
| `make uncovered` | Show functions with 0% coverage |
| `make bulk-insert CSV=<file>` | Bulk insert books from CSV file |

---

## Bruno API Client

A [Bruno](https://www.usebruno.com/) collection is included for testing the API endpoints.

### Requests

| Request | Method | Endpoint | Description |
|---------|--------|----------|-------------|
| Get Books | GET | `/books` | List all books |
| Insert Books | POST | `/books` | Create a new book |
| Get Book By Id | GET | `/books/:id` | Get a book by ID |

### Usage

1. Open Bruno and import the collection from `Bruno_api/collections_boocks/`
2. Make sure the API is running (`docker compose up -d --build` or `make run-local`)
3. Select a request and click **Send**

---

## Project Structure

```
mvp_books/
├── cmd/api/                 Application entrypoint
├── internal/
│   ├── controller/          HTTP handlers + tests
│   ├── service/             Business logic + tests + mocks
│   ├── repository/          Data persistence + tests
│   ├── domain/              Entities
│   └── middleware/          Logging & metrics middleware
├── platform/                Infrastructure (DB, logger, OTel)
├── config/                  Prometheus, Grafana, Loki, Promtail configs
├── scripts/                 DB init, bulk insert tool, sample data
├── Bruno_api/               API testing collection
├── .github/workflows/       CI pipeline (lint, test, coverage gate)
├── docker-compose.yaml      6 services orchestrated
├── Dockerfile               Multi-stage build (Go → Alpine)
├── Makefile                 Development shortcuts
└── go.mod                   Go 1.24 module definition
```

---

<div align="center">

**Built with Go + Gin + GORM + MySQL + OpenTelemetry + Prometheus + Grafana + Loki**

</div>
