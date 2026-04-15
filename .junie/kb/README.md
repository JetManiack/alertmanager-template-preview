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
- **Persistence**: Editor contents (`template` and `alertData`) and UI theme are persisted in `localStorage` to preserve user progress between page refreshes.
- **Error Highlighting**:
    - **YAML/JSON**: Real-time validation in the frontend using `js-yaml` with a visual indicator.
    - **Templates**: Parsing of Go template error strings (line:column) from the backend to show indicators in the UI.
- **Autocomplete**: Implemented code autocompletion for Alertmanager templates in CodeMirror 6.
    - **Static**: Suggests common template functions (`toUpper`, `toJson`, etc.) and top-level `.Data` fields (`.CommonLabels`, `.Alerts`, etc.).
    - **Dynamic**: Real-time extraction of keys from "Alert Data" (YAML/JSON) to provide specific suggestions for `.CommonLabels.<key>`, `.GroupLabels.<key>`, etc.
    - **Trigger**: Triggers on `.` for variables or when typing function names. Can be manually activated anywhere (e.g., via Ctrl+Space) to see all available functions and variables.
    - **Autocomplete Cursor Positioning**: In CodeMirror 6, to ensure the cursor stays after the inserted word, provide explicit `from` and `to` ranges in the `CompletionResult`. For variables following a dot (e.g., `.C`), set `from` to the position immediately after the dot, and omit the leading dot from the completion labels. This ensures the prefix is correctly replaced and the cursor position is correctly updated by the editor.
    - **Global Suggestions**: If the user hasn't typed a dot yet (e.g., they are typing a function name or at the start of a block), variables are suggested with a leading dot (e.g., `.CommonLabels`). This makes all template elements discoverable without needing to type a dot first.
    - **Multi-Mode Support**: The application now supports two separate modes: **Alertmanager** and **Prometheus**.
        - **Alertmanager Mode**: Uses `prometheus/alertmanager/template` for rendering. Context includes `.CommonLabels`, `.Alerts`, etc.
        - **Prometheus Mode**: Uses a custom renderer to mimic Prometheus alerting/recording rule templates. Context includes `.Labels`, `.Value`, `.ExternalURL`, etc.
        - **Functions**: Prometheus mode supports specific functions like `humanize`, `humanize1024`, `humanizePercentage`, `humanizeDuration`, and `humanizeTimestamp`.
        - **Tabs**: A tab switcher in the header allows users to switch between modes. Each mode's template and data are persisted independently in `localStorage`.
    - **YAML Support**: The backend uses `github.com/goccy/go-yaml` for unmarshaling alert data. This library is used because it correctly respects `json` struct tags (which are present in `prometheus/alertmanager/template.Data`), allowing both YAML and JSON input to be parsed into the same Go structures.
- **Automatic Rendering**: Debounced (500ms) automatic rendering on every change in Template or Alert Data fields. Manual "Run" button was removed to provide a more seamless experience.
- **CodeMirror Height**: To ensure the editor scroller takes up the full container height (preventing it from shrinking with short content), set `.cm-scroller { height: 100% !important; }` and `.cm-content, .cm-gutters { min-height: 100% !important; }`. Also ensure `.cm-editor` and its parent `.editor-container` are correctly expanded to 100% height. This keeps the horizontal scrollbar at the bottom of the editor pane regardless of content length.
- **Resizable Panels**: Implemented using `react-resizable-panels` (version 4+).
    - **Layout Persistence**: Proportions of panels are saved in `localStorage` using the `useDefaultLayout` hook with unique `id`s for horizontal and vertical groups. This ensures that the user's preferred workspace layout is preserved between sessions.
    - **Custom Styling**: The `Separator` component is styled with `resize-handle-horizontal/vertical` classes, using Bootstrap's border color by default and the primary color when hovered or active.
    - **Structure**: Uses a nested `Group` approach to separate the Template/Data editors from the Result preview.

### Known Issues & Solutions
- **404 on Assets in Production**:
    - **Symptom**: UI loads but assets (`/assets/index-...`) return 404.
    - **Cause**: Assets are linked relative to the root, but the server is set up to serve them under `/ui`.
    - **Fix**: Add `base: '/ui/'` to `ui/vite.config.js` and rebuild.
- **Prometheus Query Mocking (Deprecated)**:
    - **Issue**: Initial attempt to implement `query` with manual mocks proved inflexible.
    - **Action**: Rolled back mock implementation.
    - **Next Step**: Implement real integration with Prometheus API via backend proxy or snapshot data.

### Planned Features & Architecture Ideas
- **Prometheus Real Integration**:
    - **Option 1 (Backend Proxy)**: Add `--prometheus-url` CLI flag. Proxy `query` calls from templates to the real server.
    - **Option 2 (Direct Browser Access)**: Input Prometheus URL in UI. Frontend fetches data directly (requires CORS on Prometheus).
    - **Option 3 (Smart Snapshots)**: Allow pasting raw JSON output from Prometheus API `/api/v1/query` into the "Alert Data" field.

### Successful Patterns
-   (No entries yet)
