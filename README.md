# Go Books API

REST API for managing books, built with Go (Gin), GORM, and MySQL. Includes a full observability stack with OpenTelemetry, Prometheus, Grafana, Loki, and Promtail for metrics and log aggregation.

## Prerequisites

- Go 1.24+
- Docker & Docker Compose

## Quick Start

### 1. Create your environment file

```bash
cp environment_values.example environment_values
```

Edit `environment_values` if you need different credentials. The defaults work out of the box:

```
DB_USER=developer
DB_PASSWORD=admin
DB_NAME=books_db
DB_HOST=db
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=books_db
MYSQL_USER=developer
MYSQL_PASSWORD=admin
METRICS_PORT=2223
LOG_LEVEL=info
```

### 2. Start everything

```bash
docker compose up -d --build
```

This starts 6 services on a shared Docker network:

| Service | URL | Description |
|---|---|---|
| API | http://localhost:8080 | Books REST API |
| Metrics | http://localhost:2223/metrics | Raw Prometheus metrics (OpenTelemetry) |
| Prometheus | http://localhost:9090 | Metrics scraping & storage |
| Grafana | http://localhost:3000 | Dashboards & visualization |
| Loki | http://localhost:3100 | Log aggregation backend |
| Promtail | - | Log shipper (sends Docker logs to Loki) |

### 3. Verify everything works

```bash
# Create a book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Refactoring", "author": "Martin Fowler"}'

# List books
curl http://localhost:8080/books

# Check raw metrics
curl http://localhost:2223/metrics

# Open Grafana
open http://localhost:3000
```

### 4. Open the Dashboard

1. Go to http://localhost:3000 (credentials: `admin` / `admin`)
2. Navigate to **Dashboards** → **API Monitoring** → **Books API Monitoring**
3. Send some traffic to populate the panels, then set the time range to the last 5 minutes

### Stop everything

```bash
docker compose down
```

### Stop and reset database

```bash
docker compose down -v
```

## Local Development

### Option A: DB + monitoring in Docker, API locally

```bash
make run-local
```

Starts MySQL, Prometheus, Grafana, Loki, and Promtail in Docker, then runs the API on your machine.

**Important:** The Grafana dashboard will **not** show data in this mode on macOS. Docker Desktop containers cannot reach the host's metrics port (`host.docker.internal` routing is broken on Mac). The raw metrics endpoint still works at http://localhost:2223/metrics. To use the full dashboard, run everything in Docker with `docker compose up -d --build`.

### Option B: API locally with local DB

```bash
make run
```

Requires a MySQL instance running on `localhost:3306` with the schema from `scripts/init.sql`.

## Configuration

Environment variables are managed via the `environment_values` file (gitignored).

| Variable | Default (Docker) | Default (Local) | Description |
|---|---|---|---|
| `DB_USER` | developer | developer | MySQL user |
| `DB_PASSWORD` | admin | admin | MySQL password |
| `DB_HOST` | db | localhost | DB hostname |
| `DB_NAME` | books_db | books_db | Database name |
| `MYSQL_ROOT_PASSWORD` | rootpassword | - | MySQL root password |
| `MYSQL_DATABASE` | books_db | - | Auto-created database |
| `MYSQL_USER` | developer | - | Auto-created MySQL user |
| `MYSQL_PASSWORD` | admin | - | Auto-created MySQL password |
| `METRICS_PORT` | 2223 | 2223 | Prometheus metrics endpoint port |
| `LOG_LEVEL` | info | info | Log level (DEBUG, INFO, WARN, ERROR) |

`DB_HOST` must be `db` when running inside Docker, `localhost` when running the API on your host machine.

The database is automatically initialized on first startup via `scripts/init.sql` (creates the `books` table and inserts sample data).

## API Endpoints

### Create a Book

**POST** `/books`

```json
{
    "title": "The Refactoring",
    "author": "Martin Fowler"
}
```

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Refactoring", "author": "Martin Fowler"}'
```

### Get All Books

**GET** `/books`

```bash
curl http://localhost:8080/books
```

### Get One Book

**GET** `/books/:id`

```bash
curl http://localhost:8080/books/1
```

## Monitoring & Metrics

### Architecture

```
API (:8080) ──HTTP──> Prometheus (:9090) ──query──> Grafana (:3000)
     │                                                    ▲
     └──OTel SDK──> Prometheus exporter (:2223)           │
                                                          │
Docker logs ──Promtail──> Loki (:3100) ───────────────────┘
```

The API uses [OpenTelemetry](https://opentelemetry.io/) to instrument HTTP requests. A Prometheus exporter runs on a separate port (`:2223`) and exposes metrics in Prometheus format. Prometheus scrapes this endpoint every 10 seconds. Grafana connects to Prometheus as a data source and displays a pre-configured dashboard.

For log aggregation, [Promtail](https://grafana.com/docs/loki/latest/send-data/promtail/) ships Docker container logs to [Loki](https://grafana.com/oss/loki/), which stores and indexes them. Grafana connects to Loki as a second data source, enabling log search and correlation with metrics.

### Available Metrics

The OpenTelemetry Prometheus exporter converts OTel metric names to Prometheus format by replacing dots with underscores (e.g. `http.server.request.count` becomes `http.server.request.count_total`). Labels follow the same convention (`http.method`, `http.route`, etc.).

When querying in Prometheus/Grafana, use `{__name__="http.server.request.count_total"}` syntax because metric names contain dots.

| OTel Metric | Prometheus Name | Type | Description |
|---|---|---|---|
| `http.server.request.count` | `http.server.request.count_total` | Counter | Total requests by method, route, and status |
| `http.server.request.duration` | `http.server.request.duration_seconds` | Histogram | Request latency (bucket boundaries: 5ms to 10s) |
| `http.server.request.inflight` | `http.server.request.inflight` | Gauge | Concurrent requests in progress |

Labels:

| Label | Example Values |
|---|---|
| `http.method` | `GET`, `POST` |
| `http.route` | `/books`, `/books/:id` |
| `http.response.status_code` | `200`, `201`, `404` |

### Grafana Dashboards

Two dashboards are auto-provisioned:

#### API Monitoring

**Dashboards → API Monitoring → Books API Monitoring**

Connects to Prometheus. Includes 7 panels:

| Panel | Type | Description |
|---|---|---|
| Request Rate (req/s) | Timeseries | Requests per second by method and route |
| Request Latency (P50 / P95 / P99) | Timeseries | Latency percentiles over time |
| Error Rate (5xx) | Timeseries | Percentage of 5xx errors by route |
| In-Flight Requests | Timeseries | Currently active concurrent requests |
| Total Requests by Status Code | Piechart | Distribution of responses by status code |
| Avg Request Duration by Endpoint | Bargauge | Average response time per route |
| Request Duration Heatmap | Heatmap | Latency distribution across time |

The dashboard datasource connects to Prometheus via the Docker network (`http://prometheus:9090`). If you modify `config/grafana/dashboards/api-monitoring.json`, Grafana will reload it within 10 seconds.

#### Logs

**Dashboards → API Monitoring → Books Logs**

Connects to Loki. Includes 4 panels:

| Panel | Type | Description |
|---|---|---|
| All Logs | Logs | Raw log stream from all containers |
| API Errors | Logs | Filtered error-level log entries |
| Logs by Service | Table | Log count grouped by Docker service |
| Log Rate | Timeseries | Log ingestion rate over time |

### Raw Metrics

Access raw Prometheus metrics at http://localhost:2223/metrics

### Prometheus

Access Prometheus at http://localhost:9090 → **Status → Targets** to verify scraping is healthy. Both `app:2223` and `host.docker.internal:2223` should show as `UP`.

## Testing

```bash
make test
```

### Coverage

```bash
make coverage
```

Coverage excludes infrastructure code (`platform/`, `internal/middleware/`, `cmd/`, `mocks/`, `internal/domain/`) from the calculation. Business logic (`internal/service/`, `internal/repository/`) should maintain ≥85% coverage.

## Architecture

```
cmd/api/main.go              # Entrypoint: wires telemetry, middleware, routes, DB
internal/
  controller/                # HTTP handlers (Gin)
  service/                   # Business logic (use cases)
  repository/                # Data persistence (GORM/MySQL)
  domain/                    # Entities & repository interfaces
  middleware/
    logging.go               # Gin middleware: structured request logging (slog)
    metrics.go               # Gin middleware: records request count, duration, in-flight
platform/
  connection/                # DB initialization (GORM)
  logger/
    logger.go                # Structured JSON logger init (slog)
  metrics/
    metrics.go               # OTel metric definitions (counter, histogram, updowncounter)
  telemetry.go               # OTel SDK init + Prometheus exporter + metrics HTTP server
config/
  prometheus/
    prometheus.yml           # Scrape config targeting app:2223
  grafana/
    provisioning/
      datasources/
        datasource.yml       # Prometheus datasource (uid: prometheus)
        loki.yml             # Loki datasource (uid: loki)
      dashboards/
        dashboard.yml        # File-based dashboard provisioning
    dashboards/
      api-monitoring.json    # Pre-configured 7-panel metrics dashboard
      books-logs.json        # Pre-configured 4-panel logs dashboard
  loki/
    loki-config.yml          # Loki server configuration
  promtail/
    config.yml               # Promtail log shipper configuration
scripts/
  init.sql                   # DB schema + sample data (runs on first Docker startup)
Bruno_api/
  collections_boocks/        # Bruno API client collection for testing
```

## Make Commands

| Command | Description |
|---|---|
| `make run-local` | DB + monitoring in Docker, API runs locally |
| `make run` | API locally (requires local MySQL on :3306) |
| `make docker-up` | Build and start everything in Docker |
| `make docker-down` | Stop all containers |
| `make test` | Run all tests |
| `make coverage` | Run tests with filtered coverage report |
| `make uncovered` | Show functions with 0% coverage |

## Bruno API Client

A [Bruno](https://www.usebruno.com/) collection is included for testing the API endpoints.

### Requests

| Request | Method | Endpoint | Description |
|---|---|---|---|
| Get Books | GET | `/books` | List all books |
| Insert Books | POST | `/books` | Create a new book |
| Get Book By Id | GET | `/books/:id` | Get a book by ID |

### Usage

1. Open Bruno and import the collection from `Bruno_api/collections_boocks/`
2. Make sure the API is running (via `docker compose up -d --build` or `make run-local`)
3. Select a request and click **Send**
