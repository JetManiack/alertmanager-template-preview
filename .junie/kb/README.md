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
- **Code Editor**: Replaced standard `textarea` with `CodeMirror 6` (via `@uiw/react-codemirror`) to provide syntax highlighting, line numbers, and better user experience.
- **Syntax Highlighting**:
    - **YAML/JSON**: Uses `@codemirror/lang-yaml` for "Alert Data" field (YAML is a superset of JSON).
    - **Go Templates**: Uses `@codemirror/legacy-modes/mode/go` as a close approximation for highlighting.
- **Theme Switching**: Implemented using Bootstrap 5's `data-bs-theme` attribute and CodeMirror's theme system (`@uiw/codemirror-theme-vscode`). Theme preference is persisted in `localStorage`.
- **Error Highlighting**:
    - **YAML/JSON**: Real-time validation in the frontend using `js-yaml` with a visual indicator.
    - **Templates**: Parsing of Go template error strings (line:column) from the backend to show indicators in the UI.
- **YAML Support**: The backend uses `github.com/goccy/go-yaml` for unmarshaling alert data. This library is used because it correctly respects `json` struct tags (which are present in `prometheus/alertmanager/template.Data`), allowing both YAML and JSON input to be parsed into the same Go structures.
- **Automatic Rendering**: Debounced (500ms) automatic rendering on every change in Template or Alert Data fields. Manual "Run" button was removed to provide a more seamless experience.

### Known Issues & Solutions
- **404 on Assets in Production**:
    - **Symptom**: UI loads but assets (`/assets/index-...`) return 404.
    - **Cause**: Assets are linked relative to the root, but the server is set up to serve them under `/ui`.
    - **Fix**: Add `base: '/ui/'` to `ui/vite.config.js` and rebuild.

### Successful Patterns
-   (No entries yet)
