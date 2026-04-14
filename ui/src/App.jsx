import { useState, useEffect } from 'react'
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
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

  const handleRender = async () => {
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

      const json = await response.json();
      if (response.ok) {
        setResult(json.result);
      } else {
        setError(json.error || 'Failed to render template');
      }
    } catch (err) {
      setError('Connection error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      handleRender();
    }, 500);
    return () => clearTimeout(timer);
  }, [template, data]);

  return (
    <div className="vh-100 d-flex flex-column">
      <header className="header">
        <h6 className="mb-0 fw-bold me-4">Alertmanager Template Preview</h6>
        <button className="btn-run" onClick={handleRender} disabled={loading}>
          {loading ? 'Running...' : 'Run'}
        </button>
      </header>

      <main className="main-container">
        <div className="left-panel">
          <div className="top-pane">
            <div className="editor-pane">
              <div className="editor-label">Template</div>
              <textarea
                className="editor-textarea"
                value={template}
                onChange={(e) => setTemplate(e.target.value)}
              />
            </div>
          </div>
          <div className="bottom-pane">
            <div className="editor-pane border-top-0">
              <div className="editor-label">Alert Data (JSON)</div>
              <textarea
                className="editor-textarea"
                value={data}
                onChange={(e) => setData(e.target.value)}
              />
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
