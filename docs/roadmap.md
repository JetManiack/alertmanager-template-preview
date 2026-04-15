# Development Roadmap

## Phase 1: Core Backend
- [x] Initialize Go modules and basic project structure.
- [x] Implement core template rendering logic using `alertmanager/template`. (TDD)
- [x] Create REST API endpoint: `POST /api/render` with template and JSON data.
- [x] Add basic validation for templates and alert data.
- [x] Implement CLI using `urfave/cli/v3`.
- [x] Set up `go:embed` for Web-UI.

## Phase 2: Frontend MVP
- [x] Initialize React frontend project (Vite).
- [x] Build form with template and alert data inputs.
- [x] Add real-time preview (500ms debounce).
- [x] Integrate with backend API.
- [x] Modern styling with Tailwind CSS v4.
- [x] Embedded UI into Go binary.

## Phase 3: Advanced Features & UX
- [x] Add syntax highlighting for templates and JSON/YAML (CodeMirror 6).
- [x] Implement real-time rendering with debounce.
- [x] Add "Shareable Links" (URL encoding with compression).
- [x] Context-aware autocompletion for templates.
- [x] Support for both Alertmanager and Prometheus modes.

## Phase 4: Production Readiness
- [x] Embed frontend build into the Go binary.
- [x] Add Dockerfile for easy deployment.
- [x] GitHub Actions for CI/CD and Docker publishing.
- [x] Documentation for usage and deployment.
