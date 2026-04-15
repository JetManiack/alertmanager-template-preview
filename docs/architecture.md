# Architecture Overview

## Project Goals
The project aims to provide a web interface for developers and SREs to preview and test Prometheus Alertmanager templates without deploying to a live Alertmanager instance.

## Technical Stack
- **Backend**: Go 1.26.1, using `urfave/cli/v3` for CLI.
- **Frontend**: React JS (using Vite), embedded in the Go binary with `go:embed`.
- **Communication**: REST API (JSON).

## High-Level Components
1. **API Server (Go)**:
   - Receives template text and sample alert data.
   - Routes requests to either **Alertmanager** or **Prometheus** rendering engines.
   - Proxies `query` calls to a real Prometheus server if configured via `--prometheus-url`.
   - Returns the rendered result, error messages, and health/metric data.
2. **Web UI (React)**:
   - Three-pane IDE-like layout (Template, Data, Result).
   - Mode switcher for Alertmanager/Prometheus context.
   - CodeMirror 6 editors with syntax highlighting and autocompletion.
   - Real-time preview with automatic debounced rendering.
   - Preview modes: Text, HTML, and Markdown.

## Backend Architecture
The backend is structured as follows:
- `/cmd/server`: Application entry point, CLI flag parsing, and HTTP server lifecycle.
- `/internal/template`: Core logic for template processing.
  - `alertmanager.go`: Alertmanager-specific rendering logic.
  - `prometheus.go`: Prometheus-specific logic and API proxying for the `query` function.
  - `functions.go`: Shared template functions (e.g., `humanize`, `toTime`, `round`) used by both engines.
- `/internal/api`: HTTP handlers, routing (Gin), and middleware.

### Key Libraries
- `github.com/prometheus/alertmanager/template`: Core rendering engine for Alertmanager.
- `github.com/prometheus/client_golang/prometheus`: Metrics collection.
- `github.com/urfave/cli/v3`: Command-line interface.
- `github.com/gin-gonic/gin`: HTTP web framework.
- `github.com/goccy/go-yaml`: Flexible YAML/JSON unmarshaling.

## Observability & Lifecycle
- **Healthchecks**: `/healthz` endpoint for liveness/readiness probes.
- **Metrics**: `/metrics` endpoint exposing standard Go and application-specific metrics in Prometheus format.
- **Graceful Shutdown**: The server listens for `SIGINT` and `SIGTERM` signals and allows up to 5 seconds for pending requests to complete before exiting.
