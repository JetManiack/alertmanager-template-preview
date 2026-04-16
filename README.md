# Alertmanager Template Preview

A web application for previewing Prometheus Alertmanager templates.

## Tech Stack
-   **Backend**: Go 1.26
-   **Frontend**: React JS

## Documentation
Index of project documentation:
-   [Architecture Overview](docs/architecture.md)
-   [Development Roadmap](docs/roadmap.md)
-   [Deployment Guide](docs/deployment.md)

## Getting Started

### Local Development
1. Build the UI and the server:
   ```bash
   make build
   ```
2. Run the server:
   ```bash
   ./bin/server
   ```

### Docker
1. Build the Docker image:
   ```bash
   docker build -t alertmanager-template-preview .
   ```
2. Run the container:
   ```bash
   docker run -p 8080:8080 alertmanager-template-preview
   ```
3. (Optional) Run with a real Prometheus backend:
   ```bash
   docker run -p 8080:8080 alertmanager-template-preview -p http://host.docker.internal:9090
   ```

## Security
- **XSS Protection**: HTML output is sanitized using DOMPurify.
- **DoS Protection**: API requests are limited to 1MB and rendering has a 5s timeout.
- **Auth**: This application does not include built-in authentication. It is recommended to run it behind a reverse proxy with Auth (e.g., OAuth2 Proxy, Nginx Basic Auth).
- **SSRF**: Prometheus proxying is limited to the URL provided at startup.

See the [Project Guidelines](.junie/guidelines.md) for contribution rules.
