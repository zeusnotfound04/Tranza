import { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useAuth } from '../hooks/useAuth';
import { apiClient } from '../services/api';
import { Key, Plus, Trash2, Copy, Eye, EyeOff, Clock, CheckCircle } from 'lucide-react';
import APIKeyCard from './APIKeyCard';

interface APIKey {
  id: string;
  label: string; // Changed from 'name' to 'label'
  key_preview?: string;
  key_type: string;
  scopes: string[];
  created_at: string;
  expires_at: string;
  last_used_at: string;
  usage_count: number;
  is_active: boolean;
  rate_limit: number;
}

// Modal component that renders as a portal
const CreateAPIKeyModal = ({ 
  isOpen, 
  onClose, 
  onSubmit, 
  formData, 
  setFormData, 
  creating 
}: {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (e: React.FormEvent) => void;
  formData: { name: string; password: string };
  setFormData: (data: { name: string; password: string }) => void;
  creating: boolean;
}) => {
  if (!isOpen || typeof document === 'undefined') return null;

  return createPortal(
    <div 
      className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
      style={{ 
        zIndex: 999999,
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0
      }}
      onClick={(e) => {
        if (e.target === e.currentTarget) {
          onClose();
        }
      }}
    >
      <div 
        className="border border-gray-700 rounded-xl shadow-2xl w-full max-w-md transform transition-all duration-200 scale-100 relative"
        style={{ maxHeight: '80vh', backgroundColor: '#1f1f1f' }}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-white">Create Universal API Key</h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-200 transition-colors p-1 hover:bg-gray-800 rounded-full"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <form onSubmit={onSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                API Key Name <span className="text-red-400">*</span>
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-4 py-3 border border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors text-white bg-gray-800"
                placeholder="e.g., Slack Bot Integration, Mobile App"
                required
                autoFocus
              />
              <p className="text-xs text-gray-400 mt-1">Choose a descriptive name to identify this key</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Password <span className="text-red-400">*</span>
              </label>
              <input
                type="password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                className="w-full px-4 py-3 border border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors text-white bg-gray-800"
                placeholder="Enter a secure password"
                required
                minLength={6}
              />
              <p className="text-xs text-gray-400 mt-1">This password will be required to view your API key later (minimum 6 characters)</p>
            </div>

            <div className="p-4 bg-gradient-to-r from-green-900/30 to-blue-900/30 border border-green-700/50 rounded-lg">
              <div className="flex items-start space-x-2">
                <CheckCircle className="w-5 h-5 text-green-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-green-300 mb-1">Universal Access Key</p>
                  <p className="text-xs text-green-200">
                    This API key provides full access to all Tranza features including wallet operations, transfers, Slack bot integration, and future services.
                  </p>
                </div>
              </div>
            </div>

            <div className="flex gap-3 pt-2">
              <button
                type="button"
                onClick={onClose}
                className="flex-1 bg-gray-700 hover:bg-gray-600 text-gray-200 py-3 px-4 rounded-lg font-medium transition-colors border border-gray-600"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={creating || !formData.name.trim() || !formData.password.trim()}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white py-3 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center gap-2"
              >
                {creating ? (
                  <>
                    <div className="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></div>
                    Creating...
                  </>
                ) : (
                  <>
                    <Plus className="w-4 h-4" />
                    Create Key
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>,
    document.body
  );
};

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
    password: string;
  }>({
    name: '',
    password: '',
  });

  useEffect(() => {
    fetchAPIKeys();
  }, []);

  // Prevent body scroll when modal is open
  useEffect(() => {
    if (showCreateForm) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    
    // Cleanup function to restore scroll when component unmounts
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [showCreateForm]);

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
    if (!createForm.name.trim() || !createForm.password.trim()) return;

    setCreating(true);
    setError(null);

    try {
      const response = await apiClient.generateAPIKey(createForm.name, createForm.password, 8760); // Default to 1 year TTL

      if (response.success && response.data) {
        setNewApiKey(response.data.api_key);
        setShowCreateForm(false);
        setCreateForm({ name: '', password: '' });
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

  const handleViewAPIKey = async (keyId: string, password: string): Promise<string | null> => {
    try {
      const response = await apiClient.viewAPIKey(keyId, password);
      if (response.success && response.data) {
        return response.data.api_key;
      } else {
        throw new Error(response.error || 'Failed to retrieve API key');
      }
    } catch (err: any) {
      throw new Error(err.message || 'Network error');
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
        <h1 className="text-3xl font-bold text-white flex items-center gap-3">
          <Key className="w-8 h-8 text-blue-400" />
          API Key Management
        </h1>
        <p className="text-gray-400">Manage universal API keys for all integrations including Slack bot, mobile apps, and external services</p>
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
        <div className="mb-6 p-6 bg-green-900/30 border border-green-700/50 rounded-lg">
          <div className="flex items-center gap-2 text-green-300 mb-3">
            <CheckCircle className="w-5 h-5" />
            <span className="font-medium">API Key Created Successfully!</span>
          </div>
          <p className="text-green-200 text-sm mb-3">
            Please copy and save this API key. You won't be able to see it again.
          </p>
          <div className="flex items-center gap-2 bg-gray-800 p-3 rounded border border-gray-600">
            <code className="flex-1 font-mono text-sm text-white">{newApiKey}</code>
            <button
              onClick={() => copyToClipboard(newApiKey)}
              className="text-green-400 hover:text-green-300 p-1"
              title="Copy API Key"
            >
              <Copy className="w-4 h-4" />
            </button>
          </div>
          <button
            onClick={() => setNewApiKey(null)}
            className="mt-3 text-green-400 hover:text-green-300 text-sm font-medium"
          >
            I've saved the key, close this
          </button>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-6 p-4 bg-red-900/30 border border-red-700/50 rounded-lg">
          <p className="text-red-200">{error}</p>
        </div>
      )}

      {/* API Keys List */}
      <div className="border border-gray-800 rounded-lg min-h-[400px]" style={{ backgroundColor: '#1f1f1f' }}>
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin w-8 h-8 border-4 border-blue-400 border-t-transparent rounded-full mx-auto mb-4"></div>
            <p className="text-gray-400">Loading API keys...</p>
          </div>
        ) : apiKeys.length === 0 ? (
          <div className="p-8 text-center">
            <Key className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-white mb-2">No API Keys</h3>
            <p className="text-gray-400 mb-4">Create your first API key to get started with integrations.</p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
            >
              Create API Key
            </button>
          </div>
        ) : (
          <div className="space-y-4  bg-[#121212]">
            {apiKeys.map((apiKey) => (
              <APIKeyCard
                key={apiKey.id}
                apiKey={apiKey}
                onRevoke={() => handleRevokeAPIKey(apiKey.id)}
                onCopy={copyToClipboard}
                onViewKey={handleViewAPIKey}
                formatDate={formatDate}
              />
            ))}
          </div>
        )}
      </div>

      {/* Portal Modal */}
      <CreateAPIKeyModal
        isOpen={showCreateForm}
        onClose={() => setShowCreateForm(false)}
        onSubmit={handleCreateAPIKey}
        formData={createForm}
        setFormData={setCreateForm}
        creating={creating}
      />
    </div>
  );
}
