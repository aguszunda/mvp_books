# Bulk Insert Script

Script en Go para insertar libros masivamente en la API mediante un archivo CSV, utilizando el endpoint `POST /books`.

## Prerrequisitos

- Go 1.24+
- La API debe estar corriendo (ver `README.md` principal)

## Formato del CSV

El archivo CSV debe tener el header `title,author`:

```csv
title,author
The Refactoring,Martin Fowler
Clean Code,Robert C. Martin
Domain-Driven Design,Eric Evans
```

> Se incluye `scripts/books_sample.csv` con 2500 libros de ejemplo para pruebas.

## Uso

### Con go run

```bash
go run scripts/bulk_insert.go <csv_file> [api_url]
```

- `csv_file` (requerido): ruta al archivo CSV
- `api_url` (opcional): URL base de la API (default: `http://localhost:8080`)

**Ejemplos:**

```bash
# Usar la URL por defecto (localhost:8080)
go run scripts/bulk_insert.go scripts/books_sample.csv

# Especificar una URL custom
go run scripts/bulk_insert.go scripts/books_sample.csv http://192.168.1.100:8080
```

### Con Make

```bash
make bulk-insert CSV=scripts/books_sample.csv
make bulk-insert CSV=scripts/books_sample.csv API_URL=http://localhost:9090
```

## Salida

El script muestra logs detallados del progreso:

```
=== Bulk Insert Started ===
CSV file: scripts/books_sample.csv
API URL:  http://localhost:8080
Total records to insert: 2500
----------------------------
[1/2500] OK - Created book (id=1): 'Golden Chronicles' by 'Sarah Clark'
[2/2500] OK - Created book (id=2): 'Lost Odyssey' by 'James Green'
[3/2500] FAIL (HTTP 400) - 'Bad Title' by '': {"error": 'author' is required"}
...
----------------------------
=== Bulk Insert Completed ===
Total:   2500
Success: 2498
Failed:  2
Duration: 45.23s
Avg per record: 18ms
```

| Log | Significado |
|---|---|
| `[n/total] OK` | Registro creado exitosamente (HTTP 201) |
| `[n/total] FAIL` | Error de la API (HTTP 4xx/5xx) con body de respuesta |
| `[n/total] ERROR` | Error de conexión o marshaling |
| `[n/total] SKIP` | Fila del CSV con menos de 2 columnas |

## Exit Codes

| Código | Significado |
|---|---|
| `0` | Todos los registros insertados correctamente |
| `1` | Al menos un registro falló |

## Ejecución de pruebas rápidas

```bash
# Levantar la API
docker compose up -d --build

# Esperar que arranque
sleep 5

# Insertar 2500 libros
go run scripts/bulk_insert.go scripts/books_sample.csv

# Verificar
curl http://localhost:8080/books | python3 -m json.tool | head
```
