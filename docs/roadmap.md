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
- [ ] Add syntax highlighting for templates and JSON (CodeMirror/Monaco).
- [ ] Implement real-time rendering.
- [ ] Add "Shareable Links" (saving state to backend or URL encoding).
- [ ] Pre-load common Alertmanager objects (labels, annotations).

## Phase 4: Production Readiness
- [ ] Embed frontend build into the Go binary.
- [ ] Add Dockerfile for easy deployment.
- [ ] Documentation for usage and deployment.
