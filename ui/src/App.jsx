import { useState, useEffect, useCallback, useMemo } from 'react'
import 'bootstrap/dist/css/bootstrap.min.css';
import { Nav } from 'react-bootstrap';
import { Group, Panel, Separator, useDefaultLayout } from 'react-resizable-panels';
import CodeMirror from '@uiw/react-codemirror';
import { autocompletion } from '@codemirror/autocomplete';
import { yaml } from '@codemirror/lang-yaml';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';
import { vscodeDark, vscodeLight } from '@uiw/codemirror-theme-vscode';
import { SunFill, MoonStarsFill, ExclamationTriangleFill } from 'react-bootstrap-icons';
import jsYaml from 'js-yaml';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { createTemplateCompletionSource } from './completions';

const DEFAULT_AM_DATA = JSON.stringify({
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
}, null, 2);

const DEFAULT_PROM_DATA = JSON.stringify({
  labels: { alertname: "HighCPU", instance: "localhost:9090" },
  externalLabels: { region: "us-east-1" },
  externalURL: "http://prometheus.example.com",
  value: 95.5
}, null, 2);

function App() {
  const [theme, setTheme] = useState(localStorage.getItem('theme') || 'light');
  const [mode, setMode] = useState(localStorage.getItem('activeMode') || 'alertmanager');
  
  const [amTemplate, setAmTemplate] = useState(localStorage.getItem('amTemplate') || '{{ .CommonLabels.alertname }}');
  const [amData, setAmData] = useState(localStorage.getItem('amData') || DEFAULT_AM_DATA);
  
  const [promTemplate, setPromTemplate] = useState(localStorage.getItem('promTemplate') || 'Alert {{ .Labels.alertname }} value is {{ .Value | humanize }}');
  const [promData, setPromData] = useState(localStorage.getItem('promData') || DEFAULT_PROM_DATA);
  const [previewMode, setPreviewMode] = useState(localStorage.getItem('previewMode') || 'text');
  
  const currentTemplate = mode === 'alertmanager' ? amTemplate : promTemplate;
  const currentData = mode === 'alertmanager' ? amData : promData;

  const [result, setResult] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [dataError, setDataError] = useState(null);

  // Parse template error to get location
  const templateError = useMemo(() => {
    if (!error) return null;
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

  // Save mode, template and data to localStorage
  useEffect(() => {
    localStorage.setItem('activeMode', mode);
  }, [mode]);

  useEffect(() => {
    localStorage.setItem('amTemplate', amTemplate);
    localStorage.setItem('amData', amData);
  }, [amTemplate, amData]);

  useEffect(() => {
    localStorage.setItem('promTemplate', promTemplate);
    localStorage.setItem('promData', promData);
  }, [promTemplate, promData]);

  useEffect(() => {
    localStorage.setItem('previewMode', previewMode);
  }, [previewMode]);

  // YAML/JSON Validation
  useEffect(() => {
    try {
      if (currentData.trim() === '') {
        setDataError('Data cannot be empty');
        return;
      }
      jsYaml.load(currentData);
      setDataError(null);
    } catch (err) {
      setDataError(err.message);
    }
  }, [currentData]);

  const parsedData = useMemo(() => {
    try {
      return jsYaml.load(currentData);
    } catch {
      return null;
    }
  }, [currentData]);

  const templateExtensions = useMemo(() => {
    return [
      StreamLanguage.define(go),
      autocompletion({ override: [createTemplateCompletionSource(parsedData, mode)] })
    ];
  }, [parsedData, mode]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  const handleRender = useCallback(async () => {
    if (dataError) {
      setError('Cannot render: Invalid YAML/JSON in Data field');
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
        body: JSON.stringify({ template: currentTemplate, data: currentData, mode }),
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
  }, [currentTemplate, currentData, mode, dataError]);

  useEffect(() => {
    const timer = setTimeout(() => {
      handleRender();
    }, 500);
    return () => clearTimeout(timer);
  }, [handleRender]);

  const cmTheme = theme === 'dark' ? vscodeDark : vscodeLight;

  const handleTemplateChange = (val) => {
    if (mode === 'alertmanager') setAmTemplate(val);
    else setPromTemplate(val);
  };

  const handleDataChange = (val) => {
    if (mode === 'alertmanager') setAmData(val);
    else setPromData(val);
  };

  const { defaultLayout: horizontalLayout, onLayoutChanged: onHorizontalLayoutChanged } = useDefaultLayout({ id: "horizontal-layout" });
  const { defaultLayout: verticalLayout, onLayoutChanged: onVerticalLayoutChanged } = useDefaultLayout({ id: "vertical-layout" });

  return (
    <div className="vh-100 d-flex flex-column">
      <header className="header">
        <div className="container-fluid d-flex align-items-center justify-content-between">
          <div className="d-flex align-items-center">
            <h6 className="mb-0 header-title me-4">Template Preview</h6>
            <Nav variant="underline" activeKey={mode} onSelect={(k) => setMode(k)}>
              <Nav.Item>
                <Nav.Link eventKey="alertmanager" className="py-1">Alertmanager</Nav.Link>
              </Nav.Item>
              <Nav.Item>
                <Nav.Link eventKey="prometheus" className="py-1">Prometheus</Nav.Link>
              </Nav.Item>
            </Nav>
          </div>
          <button className="theme-toggle" onClick={toggleTheme} title="Toggle Dark/Light Mode">
            {theme === 'light' ? <MoonStarsFill size={18} /> : <SunFill size={18} />}
          </button>
        </div>
      </header>

      <main className="main-container">
        <Group orientation="horizontal" defaultLayout={horizontalLayout} onLayoutChanged={onHorizontalLayoutChanged}>
          <Panel defaultSize={50} minSize={20}>
            <Group orientation="vertical" defaultLayout={verticalLayout} onLayoutChanged={onVerticalLayoutChanged}>
              <Panel defaultSize={50} minSize={20}>
                <div className="editor-pane h-100">
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
                      value={currentTemplate}
                      height="100%"
                      theme={cmTheme}
                      extensions={templateExtensions}
                      onChange={handleTemplateChange}
                      basicSetup={{
                        lineNumbers: true,
                        foldGutter: true,
                        highlightActiveLine: true,
                      }}
                    />
                  </div>
                </div>
              </Panel>
              <Separator className="resize-handle-vertical" />
              <Panel defaultSize={50} minSize={20}>
                <div className="editor-pane h-100">
                  <div className="editor-label">
                    <span>{mode === 'alertmanager' ? 'Alert Data (YAML/JSON)' : 'Rule Data (YAML/JSON)'}</span>
                    {dataError && (
                      <span className="badge-error" title={dataError}>
                        <ExclamationTriangleFill className="me-1" />
                        Invalid YAML/JSON
                      </span>
                    )}
                  </div>
                  <div className="editor-container">
                    <CodeMirror
                      value={currentData}
                      height="100%"
                      theme={cmTheme}
                      extensions={[yaml()]}
                      onChange={handleDataChange}
                      basicSetup={{
                        lineNumbers: true,
                        foldGutter: true,
                        highlightActiveLine: true,
                      }}
                    />
                  </div>
                </div>
              </Panel>
            </Group>
          </Panel>
          <Separator className="resize-handle-horizontal" />
          <Panel defaultSize={50} minSize={20}>
            <div className="editor-pane h-100">
              <div className="editor-label d-flex align-items-center justify-content-between">
                <div>
                  Result
                  {loading && <small className="text-success ms-2 italic">Rendering...</small>}
                </div>
                <Nav variant="pills" activeKey={previewMode} onSelect={(k) => setPreviewMode(k)} className="preview-nav">
                  <Nav.Item>
                    <Nav.Link eventKey="text">Text</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="html">HTML</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="markdown">Markdown</Nav.Link>
                  </Nav.Item>
                </Nav>
              </div>
              <div className="preview-content">
                {error ? (
                  <div className="alert alert-danger mb-0 rounded-0 border-0">
                    <strong className="d-block mb-1">Rendering Error</strong>
                    <pre className="mb-0 text-break" style={{whiteSpace: 'pre-wrap'}}>{error}</pre>
                  </div>
                ) : (
                  <div className={`preview-result ${previewMode}-mode`}>
                    {previewMode === 'text' && (
                      <pre className="mb-0">{result || '(empty output)'}</pre>
                    )}
                    {previewMode === 'html' && (
                      <div className="html-preview" dangerouslySetInnerHTML={{ __html: result || '<i>(empty output)</i>' }} />
                    )}
                    {previewMode === 'markdown' && (
                      <div className="markdown-preview">
                        <ReactMarkdown remarkPlugins={[remarkGfm]}>{result || '*(empty output)*'}</ReactMarkdown>
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>
          </Panel>
        </Group>
      </main>
    </div>
  )
}

export default App
