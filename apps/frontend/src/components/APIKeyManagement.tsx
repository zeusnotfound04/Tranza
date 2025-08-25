import { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { apiClient } from '../services/api';
import { Key, Plus, Trash2, Copy, Eye, EyeOff, Bot, Clock, CheckCircle } from 'lucide-react';

interface APIKey {
  id: string;
  name: string;
  key_preview: string;
  scopes: string[];
  created_at: string;
  last_used: string;
  usage_count: number;
  is_bot: boolean;
  active: boolean;
}

export default function APIKeyManagement() {
  const { user } = useAuth();
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newApiKey, setNewApiKey] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<{
    name: string;
  }>({
    name: '',
  });

  useEffect(() => {
    fetchAPIKeys();
  }, []);

  const fetchAPIKeys = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.getAPIKeys();
      if (response.success && response.data) {
        console.log('Fetched API Keys:', response.data);  
        // Fix: Extract the keys array from the response data
        const keysData = response.data.keys || response.data;
        setApiKeys(Array.isArray(keysData) ? keysData : []);
      } else {
        setError(response.error || 'Failed to fetch API keys');
      }
    } catch (err) {
      setError('Network error');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateAPIKey = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!createForm.name.trim()) return;

    setCreating(true);
    setError(null);

    try {
      const response = await apiClient.generateAPIKey(createForm.name, 8760); // Default to 1 year TTL

      if (response.success && response.data) {
        setNewApiKey(response.data.api_key);
        setShowCreateForm(false);
        setCreateForm({ name: '' });
        await fetchAPIKeys();
      } else {
        setError(response.error || 'Failed to create API key');
      }
    } catch (err) {
      setError('Network error');
    } finally {
      setCreating(false);
    }
  };

  const handleRevokeAPIKey = async (keyId: string) => {
    if (!confirm('Are you sure you want to revoke this API key? This action cannot be undone.')) {
      return;
    }

    try {
      const response = await apiClient.revokeAPIKey(keyId);
      if (response.success) {
        await fetchAPIKeys();
      } else {
        setError(response.error || 'Failed to revoke API key');
      }
    } catch (err) {
      setError('Network error');
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // You could add a toast notification here
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Key className="w-8 h-8 text-blue-600" />
          API Key Management
        </h1>
        <p className="text-gray-600">Manage universal API keys for all integrations including Slack bot, mobile apps, and external services</p>
      </div>

      {/* Create API Key Button */}
      <div className="mb-6">
        <button
          onClick={() => setShowCreateForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Create Universal API Key
        </button>
      </div>

      {/* New API Key Display */}
      {newApiKey && (
        <div className="mb-6 p-6 bg-green-50 border border-green-200 rounded-lg">
          <div className="flex items-center gap-2 text-green-700 mb-3">
            <CheckCircle className="w-5 h-5" />
            <span className="font-medium">API Key Created Successfully!</span>
          </div>
          <p className="text-green-600 text-sm mb-3">
            Please copy and save this API key. You won't be able to see it again.
          </p>
          <div className="flex items-center gap-2 bg-white p-3 rounded border">
            <code className="flex-1 font-mono text-sm text-gray-800">{newApiKey}</code>
            <button
              onClick={() => copyToClipboard(newApiKey)}
              className="text-green-600 hover:text-green-700 p-1"
              title="Copy API Key"
            >
              <Copy className="w-4 h-4" />
            </button>
          </div>
          <button
            onClick={() => setNewApiKey(null)}
            className="mt-3 text-green-600 hover:text-green-700 text-sm font-medium"
          >
            I've saved the key, close this
          </button>
        </div>
      )}

      {/* Create Form Modal */}
      {showCreateForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Create Universal API Key</h2>
            
            <form onSubmit={handleCreateAPIKey}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Name
                </label>
                <input
                  type="text"
                  value={createForm.name}
                  onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g., Slack Bot, Mobile App"
                  required
                />
              </div>

              <div className="mb-6">
                <div className="p-3 bg-green-50 border border-green-200 rounded-lg">
                  <p className="text-sm text-green-700">
                    <strong>Universal Access:</strong> This API key will work with all features including wallet operations, transfers, Slack bot integration, and future features.
                  </p>
                </div>
              </div>

              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 py-2 px-4 rounded-lg font-medium transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={creating || !createForm.name.trim()}
                  className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg font-medium transition-colors"
                >
                  {creating ? 'Creating...' : 'Create Key'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-700">{error}</p>
        </div>
      )}

      {/* API Keys List */}
      <div className="bg-white rounded-lg border border-gray-200">
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto mb-4"></div>
            <p className="text-gray-600">Loading API keys...</p>
          </div>
        ) : apiKeys.length === 0 ? (
          <div className="p-8 text-center">
            <Key className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No API Keys</h3>
            <p className="text-gray-600 mb-4">Create your first API key to get started with integrations.</p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
            >
              Create API Key
            </button>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {apiKeys.map((apiKey) => (
              <APIKeyCard
                key={apiKey.id}
                apiKey={apiKey}
                onRevoke={() => handleRevokeAPIKey(apiKey.id)}
                formatDate={formatDate}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function APIKeyCard({ 
  apiKey, 
  onRevoke, 
  formatDate 
}: { 
  apiKey: APIKey; 
  onRevoke: () => void;
  formatDate: (date: string) => string;
}) {
  return (
    <div className="p-6">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-2">
            <h3 className="text-lg font-semibold text-gray-900">{apiKey.name}</h3>
            {apiKey.is_bot && (
              <span className="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-700 text-xs font-medium rounded-full">
                <Bot className="w-3 h-3" />
                Bot
              </span>
            )}
            <span className={`inline-flex items-center px-2 py-1 text-xs font-medium rounded-full ${
              apiKey.active 
                ? 'bg-green-100 text-green-700' 
                : 'bg-red-100 text-red-700'
            }`}>
              {apiKey.active ? 'Active' : 'Revoked'}
            </span>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-600 mb-3">
            <div>
              <span className="font-medium">Key Preview:</span>
              <div className="font-mono">{apiKey.key_preview}</div>
            </div>
            <div>
              <span className="font-medium">Created:</span>
              <div>{formatDate(apiKey.created_at)}</div>
            </div>
            <div>
              <span className="font-medium">Last Used:</span>
              <div>{apiKey.last_used ? formatDate(apiKey.last_used) : 'Never'}</div>
            </div>
          </div>

          <div className="mb-3">
            <span className="text-sm font-medium text-gray-700">Scopes:</span>
            <div className="flex flex-wrap gap-1 mt-1">
              {apiKey.scopes.map((scope) => (
                <span
                  key={scope}
                  className="inline-block px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded"
                >
                  {scope}
                </span>
              ))}
            </div>
          </div>

          <div className="text-sm text-gray-600">
            <span className="font-medium">Usage Count:</span> {apiKey.usage_count.toLocaleString()} requests
          </div>
        </div>

        <div className="flex items-center gap-2 ml-4">
          {apiKey.active && (
            <button
              onClick={onRevoke}
              className="text-red-600 hover:text-red-700 p-2 hover:bg-red-50 rounded-lg transition-colors"
              title="Revoke API Key"
            >
              <Trash2 className="w-4 h-4" />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
