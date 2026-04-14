# Knowledge Base

This is a structured Knowledge Base for the Alertmanager Template Preview project. 

## Key Learnings
### Nuances and Features
- **Go Version**: The project uses Go 1.26.1 (as per local environment and user request). This is slightly ahead of common stable releases (1.23-1.24).
- **Project Layout**: Standard Go layout is used with `/cmd` for entry points, `/internal` for core logic, and `/docs` for project documentation.
- **Go Embed**: Static assets for the Web-UI are stored in `assets/ui/dist` and embedded into the binary using a dedicated `assets` package. This avoids path issues with `go:embed` not supporting `..`.
- **UI Framework**: Transitioned from Tailwind v4 to Bootstrap 5 for better layout control and a more "IDE-like" appearance.
- **IDE Layout (jsfiddle-style)**: Implemented a 3-pane layout for Template, Alert Data (JSON), and Result.
- **Gin Static Serving**: When using `r.StaticFS("/ui", ...)` and a redirect from `/` to `/ui/`, ensure no conflicting routes like `r.GET("/ui/", ...)` are manually registered, as `StaticFS` already handles directory root requests.
- **Vite Base Path**: When serving the UI under a prefix (e.g., `/ui`), the `base` configuration in `vite.config.js` must match that prefix (e.g., `base: '/ui/'`). Otherwise, the browser will attempt to load assets from the root path, leading to 404 errors.

### Known Issues & Solutions
- **404 on Assets in Production**:
    - **Symptom**: UI loads but assets (`/assets/index-...`) return 404.
    - **Cause**: Assets are linked relative to the root, but the server is set up to serve them under `/ui`.
    - **Fix**: Add `base: '/ui/'` to `ui/vite.config.js` and rebuild.

### Successful Patterns
-   (No entries yet)
