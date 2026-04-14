import { useState, useEffect } from 'react'

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
    <div className="min-h-screen bg-gray-900 text-gray-100 p-8 font-sans">
      <header className="mb-8">
        <h1 className="text-3xl font-bold text-purple-400">Alertmanager Template Preview</h1>
        <p className="text-gray-400">Test your Prometheus Alertmanager templates in real-time.</p>
      </header>

      <main className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="space-y-6">
          <section>
            <h2 className="text-xl font-semibold mb-2 flex items-center">
              <span className="bg-purple-600 rounded-full w-6 h-6 flex items-center justify-center text-sm mr-2">1</span>
              Template
            </h2>
            <textarea
              className="w-full h-48 bg-gray-800 border border-gray-700 rounded-lg p-4 font-mono text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none transition-all"
              value={template}
              onChange={(e) => setTemplate(e.target.value)}
              placeholder="Enter your template here..."
            />
          </section>

          <section>
            <h2 className="text-xl font-semibold mb-2 flex items-center">
              <span className="bg-purple-600 rounded-full w-6 h-6 flex items-center justify-center text-sm mr-2">2</span>
              Alert Data (JSON)
            </h2>
            <textarea
              className="w-full h-64 bg-gray-800 border border-gray-700 rounded-lg p-4 font-mono text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none transition-all"
              value={data}
              onChange={(e) => setData(e.target.value)}
              placeholder="Enter alert data JSON here..."
            />
          </section>
        </div>

        <div className="flex flex-col h-full">
          <h2 className="text-xl font-semibold mb-2 flex items-center">
            <span className="bg-green-600 rounded-full w-6 h-6 flex items-center justify-center text-sm mr-2">3</span>
            Preview Result
            {loading && <span className="ml-4 text-sm text-gray-500 animate-pulse italic">Rendering...</span>}
          </h2>
          <div className="flex-grow min-h-[400px] bg-gray-800 border border-gray-700 rounded-lg p-6 overflow-auto shadow-inner relative">
            {error ? (
              <div className="bg-red-900/50 border border-red-700 text-red-200 p-4 rounded-md">
                <strong className="block font-bold mb-1">Rendering Error</strong>
                <pre className="whitespace-pre-wrap text-sm">{error}</pre>
              </div>
            ) : (
              <div className="prose prose-invert max-w-none h-full">
                <pre className="whitespace-pre-wrap font-mono text-sm bg-transparent p-0">{result || '(empty output)'}</pre>
              </div>
            )}
          </div>
        </div>
      </main>
      
      <footer className="mt-12 pt-8 border-t border-gray-800 text-center text-gray-500 text-sm">
        Built with Go 1.26, React, and urfave/cli v3.
      </footer>
    </div>
  )
}

export default App
