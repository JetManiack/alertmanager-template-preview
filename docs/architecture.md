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
   - Uses Prometheus Alertmanager libraries to parse and execute the template.
   - Returns the rendered result or detailed error messages.
2. **Web UI (React)**:
   - Input for the template (e.g., `{{ .GroupLabels.alertname }}`).
   - Input for the sample data (JSON or YAML).
   - Real-time or on-demand preview window.

## Backend Architecture
The backend will be structured into:
- `/cmd/server`: Application entry point.
- `/internal/template`: Core logic for Alertmanager template processing.
- `/internal/api`: HTTP handlers and routing.

### Key Libraries
- `github.com/prometheus/alertmanager/template`: Core rendering engine.
- `github.com/prometheus/common/model`: Alert and label models.
- `gin-gonic/gin` or standard `net/http`: HTTP framework.
