# Web Service Template For Go

A **robust Go microservice template** with:

- Multiple HTTP servers (API, Metrics, Admin, optional pprof)
- Graceful goroutine management with `tomb`
- Dynamic configuration and hot reload with Viper
- IPv4/IPv6/Dual-stack support
- Dynamic logging with `slog`
- Optional pprof endpoints protected by Basic Auth
- Structured logging in JSON

---

## Table of Contents

- [Web Service Template For Go](#web-service-template-for-go)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
  - [Running the Service](#running-the-service)
  - [Admin Endpoints](#admin-endpoints)
  - [Hot Reload](#hot-reload)
  - [pprof](#pprof)
  - [Logging](#logging)
  - [Testing](#testing)
    - [API](#api)
    - [Metrics](#metrics)
    - [Admin reload](#admin-reload)
    - [pprof](#pprof-1)
  - [Signal Handling](#signal-handling)
  - [Notes](#notes)

---

## Prerequisites

- Go 1.21+  
- Optional: Docker (for containerized testing)  

---

## Installation

```bash
# Clone repository
git clone https://github.com/alex-cos/svctemplate.git
cd svctemplate

# Build the binary
make build
```

## Configuration

YAML (config.yaml)

```yaml
LogLevel: "info"

API:
  Host: "127.0.0.1"
  HTTPPort: 8086
  ShutdownGrace: 5s
  EnableTLS: false
  TLSCertFile: "certs/server.crt"
  TLSKeyFile: "certs/server.key"
  TLSClientAuth: "none" # "none", "require", "verify"
  TLSCAFile: "certs/ca.crt"

Admin:
  Host: "127.0.0.1"
  HTTPPort: 8070
  ShutdownGrace: 5s
  EnableTLS: false
  TLSCertFile: "certs/server.crt"
  TLSKeyFile: "certs/server.key"
  TLSClientAuth: "none" # "none", "require", "verify"
  TLSCAFile: "certs/ca.crt"

Metrics:
  Host: "127.0.0.1"
  HTTPPort: 8090
  ShutdownGrace: 5s
  EnableTLS: false
  TLSCertFile: "certs/server.crt"
  TLSKeyFile: "certs/server.key"
  TLSClientAuth: "none" # "none", "require", "verify"
  TLSCAFile: "certs/ca.crt"

Pprof:
  Enable: true
  Port: 8060
  ShutdownGrace: 5s
  User: "admin"
  Pass: "secret123"
```

## Environment Variables

Prefix: MYAPP_

Examples:

```sh
export MYAPP_APICONFIG_HTTPPORT=8086
export MYAPP_PPROFCONFIG_ENABLE=true
export MYAPP_LOGLEVEL=debug
```

Viper automatically merges environment variables and YAML values.

## Running the Service

```bash
# Using the built binary
./webserv[.exe]

# Or directly
go run . server
```

## Admin Endpoints

| Endpoint      | Method  | Description                  |
| :------------ | :------ | :--------------------------- |
| /healthz      | GET     | Health check                 |
| /readyz       | GET     | Readiness probe              |
| /admin/reload | POST    | Trigger configuration reload |

## Hot Reload

- Any change to config.yaml triggers automatic reload.
- Log level changes apply immediately.
- Structural changes (ports, pprof enable, credentials, or network mode) cause a graceful restart of all servers.
- Manual reload can be triggered via:

```bash
curl -X POST http://localhost:8070/admin/reload
```

## pprof

- Enabled with EnablePprof: true
- Runs on PprofPort
- Optional Basic Auth via PprofUser and PprofPass
- Access example:

```bash
curl -u admin:secret123 http://localhost:8060/debug/pprof/
```

If disabled or not authenticated, the endpoint returns 401 or is unavailable.

## Logging

- Structured JSON logs using Go's slog
- Levels: debug, info, warn, error
- Runtime log level can be adjusted via config file hot reload

Example log entry:

```json
{
  "time":"2025-10-09T16:25:21Z",
  "level":"INFO",
  "msg":"request processed",
  "method":"GET",
  "path":"/api/hello",
  "status":200,
  "duration":0.000321
}
```

## Testing

### API

```bash
curl http://localhost:8086/api/hello
```

### Metrics

```bash
curl http://localhost:8090/metrics
```

### Admin reload

```bash
curl -X POST http://localhost:8070/admin/reload
```

### pprof

```bash
curl -u admin:secret123 http://localhost:8060/debug/pprof/
```

## Signal Handling

- SIGINT and SIGTERM trigger graceful shutdown.
- Shutdown grace period controlled via ShutdownGrace (default: 5s).
- All running servers are stopped concurrently through tomb.

Example log when stopping:

```json
{
  "level":"WARN",
  "msg":"os signal received, initiating shutdown",
  "signal":"SIGTERM"
}
```

## Notes

- Designed as a production-grade microservice template for Go.
- Demonstrates clean concurrency and dynamic configuration management.
- Compatible with both IPv4 and IPv6 environments.
- Easily containerized with Docker or deployed to Kubernetes.
- Can serve as a starting point for microservice observability stacks (Prometheus, Grafana, OpenTelemetry).
