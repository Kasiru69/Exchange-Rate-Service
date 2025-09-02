# Exchange Rate Service

A high-performance Currency Exchange Rate Service built with **Go**. It provides real-time exchange rates, currency conversion, and historical data retrieval. The service includes intelligent caching, graceful error handling, fallback mechanisms, and is fully containerized with Docker.

---

## 🧩 Features

* Real-time exchange rates (latest)
* Currency conversion (including historical conversions by date)
* Historical rates for date ranges
* Supported currencies endpoint
* Intelligent in-memory caching with configurable TTL
* Graceful error handling and sensible HTTP responses
* Docker-ready (single-container deployment)
* Simple, well-documented HTTP API for easy integration

---

## 🛠️ Prerequisites

Before running this project, install:

* **Go 1.21+** — [https://go.dev](https://go.dev)
* **Docker** (optional, for containerized deployment) — [https://www.docker.com](https://www.docker.com)
* **Git** — [https://git-scm.com](https://git-scm.com)

> Make sure `GOPATH` / `GOROOT` are configured if you rely on a non-default Go setup.

---

## 🚀 Quick Start

### Build and run with Docker (recommended)

```bash
# Build Docker image
docker build -t exchange-rate-service .

# Run container (exposes port 8080)
docker run -d --name exchange-service -p 8080:8080 exchange-rate-service

# Stop container
docker stop exchange-service

# Start container again
docker start exchange-service

# View container logs
docker logs -f exchange-service
```

### Run locally (without Docker)

```bash
# Or directly run
go run cmd/server/main.go 
```

The server listens on port **8080** by default. Adjustable via configuration or environment variables.

---

## 📡 API Endpoints

> All example `curl` requests assume the server runs on `http://localhost:8080`.

### Health check

```bash
curl http://localhost:8080/health
```

**Example response**

```json
{
  "status": "ok",
  "service": "exchange-rate-service",
  "timestamp": "2025-09-02T20:25:00+05:30",
  "version": "1.0.0"
}
```

---

### Convert currency

Basic conversion (USD → INR, amount 100):

```bash
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100"
```

Default amount (1):

```bash
curl "http://localhost:8080/convert?from=USD&to=INR"
```

Historical conversion (specific date):

```bash
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100&date=2025-08-01"
```

Different pair (EUR → GBP):

```bash
curl "http://localhost:8080/convert?from=EUR&to=GBP&amount=50"
```

**Example response**

```json
{
    "amount": 43.50805,
    "from_currency": "EUR",
    "to_currency": "GBP",
    "rate": 0.870161,
    "date": "2025-09-02",
    "timestamp": "2025-09-02T21:41:05+05:30"
}
```

---

### Latest exchange rates

Default base currency is `USD`.

```bash
curl "http://localhost:8080/api/v1/latest"
```

Custom base currency:

```bash
curl "http://localhost:8080/api/v1/latest?base=EUR"
curl "http://localhost:8080/api/v1/latest?base=GBP"
```

**Example response**

```json
{
  "base_currency": "USD",
  "rates": {
    "EUR": 0.85,
    "GBP": 0.73,
    "INR": 83.25,
    "JPY": 110.50
  },
  "timestamp": "2025-09-02T20:25:00+05:30",
  "date": "2025-09-02"
}
```

---

### Historical exchange rates (date range)

Get rates over a date range for a currency pair:

```bash
curl "http://localhost:8080/api/v1/historical?from=USD&to=INR&start_date=2025-08-28&end_date=2025-08-30"
```

Different pair example:

```bash
curl "http://localhost:8080/api/v1/historical?from=EUR&to=USD&start_date=2025-08-01&end_date=2025-08-05"
```

**Example response**

```json
{
  "from_currency": "USD",
  "to_currency": "INR",
  "rates": {
    "2025-08-28": {
      "from_currency": "USD",
      "to_currency": "INR",
      "rate": 83.25,
      "timestamp": "2025-09-02T20:25:00+05:30",
      "date": "2025-08-28"
    },
    "2025-08-29": {
      "from_currency": "USD",
      "to_currency": "INR",
      "rate": 83.30,
      "timestamp": "2025-09-02T20:25:00+05:30",
      "date": "2025-08-29"
    }
  },
  "start_date": "2025-08-28",
  "end_date": "2025-08-30"
}
```

---

### Supported currencies

```bash
curl http://localhost:8080/api/v1/currencies
```

**Example response**

```json
{
  "currencies": ["USD", "INR", "EUR", "JPY", "GBP"],
  "count": 5
}
```

---

## ❗ Error testing examples

Try these to validate error handling:

```bash
# Invalid currency
curl "http://localhost:8080/convert?from=USD&to=INVALID"

# Missing parameters
curl "http://localhost:8080/convert?from=USD"

# Invalid date format
curl "http://localhost:8080/convert?from=USD&to=INR&date=invalid-date"

# Future date
curl "http://localhost:8080/convert?from=USD&to=INR&date=2025-12-31"
```

The service returns descriptive HTTP status codes and JSON error payloads for invalid input — e.g., `400 Bad Request` for malformed queries and `422 Unprocessable Entity` for semantic errors such as unsupported currencies.

---

