import { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { 
  runAllTests, 
  TestResults, 
  TestDetail,
  TestConfig 
} from '../lib/testing';
import { 
  Play, 
  CheckCircle, 
  XCircle, 
  Clock, 
  AlertTriangle,
  Settings,
  Download,
  RefreshCw
} from 'lucide-react';

interface TestSuite {
  name: string;
  results: TestResults | null;
  running: boolean;
}

export default function IntegrationTesting() {
  const { token } = useAuth();
  const [testConfig, setTestConfig] = useState<TestConfig>({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
    testToken: token || undefined,
    timeout: 30000,
  });

  const [allResults, setAllResults] = useState<any>(null);
  const [running, setRunning] = useState(false);
  const [showConfig, setShowConfig] = useState(false);
  const [lastRun, setLastRun] = useState<Date | null>(null);

  useEffect(() => {
    if (token) {
      setTestConfig(prev => ({ ...prev, testToken: token }));
    }
  }, [token]);

  const runTests = async () => {
    setRunning(true);
    setAllResults(null);

    try {
      const results = await runAllTests(testConfig);
      setAllResults(results);
      setLastRun(new Date());
    } catch (error) {
      console.error('Test execution failed:', error);
    } finally {
      setRunning(false);
    }
  };

  const exportResults = () => {
    if (!allResults) return;

    const data = {
      timestamp: new Date().toISOString(),
      config: testConfig,
      results: allResults,
    };

    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `tranza-test-results-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Play className="w-8 h-8 text-green-600" />
          Integration Testing
        </h1>
        <p className="text-gray-600">Test API endpoints, validation, and integration flows</p>
      </div>

      {/* Controls */}
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={runTests}
            disabled={running}
            className="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white px-6 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            {running ? (
              <>
                <RefreshCw className="w-5 h-5 animate-spin" />
                Running Tests...
              </>
            ) : (
              <>
                <Play className="w-5 h-5" />
                Run All Tests
              </>
            )}
          </button>

          <button
            onClick={() => setShowConfig(!showConfig)}
            className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <Settings className="w-5 h-5" />
            Config
          </button>

          {allResults && (
            <button
              onClick={exportResults}
              className="bg-blue-100 hover:bg-blue-200 text-blue-700 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
            >
              <Download className="w-5 h-5" />
              Export Results
            </button>
          )}
        </div>

        {lastRun && (
          <div className="text-sm text-gray-600">
            Last run: {lastRun.toLocaleString()}
          </div>
        )}
      </div>

      {/* Configuration Panel */}
      {showConfig && (
        <div className="mb-6 bg-white border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Test Configuration</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                API Base URL
              </label>
              <input
                type="url"
                value={testConfig.baseURL}
                onChange={(e) => setTestConfig({ ...testConfig, baseURL: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                placeholder="http://localhost:8080"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Timeout (ms)
              </label>
              <input
                type="number"
                value={testConfig.timeout}
                onChange={(e) => setTestConfig({ ...testConfig, timeout: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                min="1000"
                max="60000"
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Test Token (Optional)
              </label>
              <input
                type="password"
                value={testConfig.testToken || ''}
                onChange={(e) => setTestConfig({ ...testConfig, testToken: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                placeholder="Leave empty to use current session token"
              />
            </div>
          </div>
        </div>
      )}

      {/* Overall Results */}
      {allResults && (
        <div className="mb-6 bg-white border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Overall Results</h3>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">{allResults.overall.totalPassed}</div>
              <div className="text-sm text-gray-600">Passed</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-red-600">{allResults.overall.totalFailed}</div>
              <div className="text-sm text-gray-600">Failed</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{allResults.overall.totalTests}</div>
              <div className="text-sm text-gray-600">Total</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600">{allResults.overall.successRate.toFixed(1)}%</div>
              <div className="text-sm text-gray-600">Success Rate</div>
            </div>
          </div>

          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-green-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${allResults.overall.successRate}%` }}
            ></div>
          </div>
        </div>
      )}

      {/* Test Suite Results */}
      {allResults && (
        <div className="space-y-6">
          <TestSuiteCard title="API Tests" results={allResults.api} />
          <TestSuiteCard title="Form Validation Tests" results={allResults.validation} />
          <TestSuiteCard title="Transfer Flow Tests" results={allResults.transfer} />
          <TestSuiteCard title="Performance Tests" results={allResults.performance} />
        </div>
      )}

      {/* Running State */}
      {running && !allResults && (
        <div className="bg-white border border-gray-200 rounded-lg p-8 text-center">
          <RefreshCw className="w-12 h-12 text-blue-600 animate-spin mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Running Tests</h3>
          <p className="text-gray-600">Please wait while we test your integrations...</p>
        </div>
      )}

      {/* Empty State */}
      {!running && !allResults && (
        <div className="bg-white border border-gray-200 rounded-lg p-8 text-center">
          <Play className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Ready to Test</h3>
          <p className="text-gray-600 mb-4">
            Run comprehensive tests to verify your API endpoints, form validation, and integration flows.
          </p>
          <button
            onClick={runTests}
            className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
          >
            Start Testing
          </button>
        </div>
      )}
    </div>
  );
}

function TestSuiteCard({ title, results }: { title: string; results: TestResults }) {
  const [expanded, setExpanded] = useState(false);

  const getStatusIcon = (status: TestDetail['status']) => {
    switch (status) {
      case 'pass':
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case 'fail':
        return <XCircle className="w-5 h-5 text-red-600" />;
      case 'skip':
        return <Clock className="w-5 h-5 text-yellow-600" />;
      default:
        return <AlertTriangle className="w-5 h-5 text-gray-600" />;
    }
  };

  const successRate = results.total > 0 ? (results.passed / results.total) * 100 : 0;

  return (
    <div className="bg-white border border-gray-200 rounded-lg">
      <div 
        className="p-6 cursor-pointer hover:bg-gray-50 transition-colors"
        onClick={() => setExpanded(!expanded)}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
            <div className="flex items-center gap-2">
              <span className="text-sm text-green-600 font-medium">{results.passed} passed</span>
              <span className="text-sm text-red-600 font-medium">{results.failed} failed</span>
              <span className="text-sm text-gray-600">({successRate.toFixed(1)}%)</span>
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <div className="w-24 bg-gray-200 rounded-full h-2">
              <div 
                className="bg-green-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${successRate}%` }}
              ></div>
            </div>
            <button className="text-gray-400 hover:text-gray-600">
              {expanded ? '−' : '+'}
            </button>
          </div>
        </div>
      </div>

      {expanded && (
        <div className="border-t border-gray-200 p-6">
          <div className="space-y-3">
            {results.details.map((test, index) => (
              <div 
                key={index}
                className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  {getStatusIcon(test.status)}
                  <div>
                    <div className="font-medium text-gray-900">{test.name}</div>
                    <div className="text-sm text-gray-600">{test.message}</div>
                    {test.error && (
                      <div className="text-sm text-red-600 mt-1">Error: {test.error}</div>
                    )}
                  </div>
                </div>
                <div className="text-sm text-gray-500">
                  {test.duration}ms
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
