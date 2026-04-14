import { useState, useEffect, useCallback, useMemo } from 'react'
import 'bootstrap/dist/css/bootstrap.min.css';
import CodeMirror from '@uiw/react-codemirror';
import { json } from '@codemirror/lang-json';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';
import { vscodeDark, vscodeLight } from '@uiw/codemirror-theme-vscode';
import { SunFill, MoonStarsFill, ExclamationTriangleFill } from 'react-bootstrap-icons';

function App() {
  const [theme, setTheme] = useState(localStorage.getItem('theme') || 'light');
  const [template, setTemplate] = useState('{{ .CommonLabels.alertname }}');
  const [data, setData] = useState(JSON.stringify({
    receiver: "webhook",
    status: "firing",
    alerts: [
      {
        status: "firing",
        labels: { alertname: "HighCPU", severity: "critical" },
        annotations: { summary: "CPU is high" },
        startsAt: "2023-01-01T00:00:00Z"
      }
    ],
    commonLabels: { alertname: "HighCPU" },
    commonAnnotations: { summary: "CPU is high" },
    groupLabels: { alertname: "HighCPU" },
    externalURL: "http://prometheus.example.com"
  }, null, 2));
  
  const [result, setResult] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [jsonError, setJsonError] = useState(null);

  // Parse template error to get location
  const templateError = useMemo(() => {
    if (!error) return null;
    // Go template errors look like: template: :1:12: ...
    const match = error.match(/template: :(\d+):(\d+):/);
    if (match) {
      return {
        line: parseInt(match[1]),
        column: parseInt(match[2]),
        message: error
      };
    }
    return null;
  }, [error]);

  // Apply theme to document
  useEffect(() => {
    document.documentElement.setAttribute('data-bs-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  // JSON Validation
  useEffect(() => {
    try {
      if (data.trim() === '') {
        setJsonError('Data cannot be empty');
        return;
      }
      JSON.parse(data);
      setJsonError(null);
    } catch (err) {
      setJsonError(err.message);
    }
  }, [data]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  const handleRender = useCallback(async () => {
    if (jsonError) {
      setError('Cannot render: Invalid JSON in Alert Data');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const response = await fetch('/api/render', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ template, data }),
      });

      const jsonResponse = await response.json();
      if (response.ok) {
        setResult(jsonResponse.result);
      } else {
        setError(jsonResponse.error || 'Failed to render template');
      }
    } catch (err) {
      setError('Connection error: ' + err.message);
    } finally {
      setLoading(false);
    }
  }, [template, data, jsonError]);

  useEffect(() => {
    const timer = setTimeout(() => {
      handleRender();
    }, 500);
    return () => clearTimeout(timer);
  }, [handleRender]);

  const cmTheme = theme === 'dark' ? vscodeDark : vscodeLight;

  return (
    <div className="vh-100 d-flex flex-column">
      <header className="header">
        <div className="container-fluid d-flex align-items-center justify-content-between">
          <div className="d-flex align-items-center">
            <h6 className="mb-0 header-title me-4">Alertmanager Template Preview</h6>
          </div>
          <button className="theme-toggle" onClick={toggleTheme} title="Toggle Dark/Light Mode">
            {theme === 'light' ? <MoonStarsFill size={18} /> : <SunFill size={18} />}
          </button>
        </div>
      </header>

      <main className="main-container">
        <div className="left-panel">
          <div className="top-pane">
            <div className="editor-pane">
              <div className="editor-label">
                <span>Template</span>
                {templateError && (
                  <span className="badge-error" title={templateError.message}>
                    <ExclamationTriangleFill className="me-1" />
                    Syntax Error: Line {templateError.line}
                  </span>
                )}
              </div>
              <div className="editor-container">
                <CodeMirror
                  value={template}
                  height="100%"
                  theme={cmTheme}
                  extensions={[StreamLanguage.define(go)]}
                  onChange={(value) => setTemplate(value)}
                  basicSetup={{
                    lineNumbers: true,
                    foldGutter: true,
                    highlightActiveLine: true,
                  }}
                />
              </div>
            </div>
          </div>
          <div className="bottom-pane">
            <div className="editor-pane border-top-0">
              <div className="editor-label">
                <span>Alert Data (JSON)</span>
                {jsonError && (
                  <span className="badge-error" title={jsonError}>
                    <ExclamationTriangleFill className="me-1" />
                    Invalid JSON
                  </span>
                )}
              </div>
              <div className="editor-container">
                <CodeMirror
                  value={data}
                  height="100%"
                  theme={cmTheme}
                  extensions={[json()]}
                  onChange={(value) => setData(value)}
                  basicSetup={{
                    lineNumbers: true,
                    foldGutter: true,
                    highlightActiveLine: true,
                  }}
                />
              </div>
            </div>
          </div>
        </div>
        <div className="right-panel">
          <div className="editor-pane border-start-0">
            <div className="editor-label">
              Result
              {loading && <small className="text-success ms-2 italic">Rendering...</small>}
            </div>
            <div className="preview-content">
              {error ? (
                <div className="alert alert-danger mb-0 rounded-0 border-0">
                  <strong className="d-block mb-1">Rendering Error</strong>
                  <pre className="mb-0 text-break" style={{whiteSpace: 'pre-wrap'}}>{error}</pre>
                </div>
              ) : (
                <div className="preview-result">
                  <pre className="mb-0">{result || '(empty output)'}</pre>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
