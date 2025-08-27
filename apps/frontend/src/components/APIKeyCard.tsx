import React, { useState } from 'react';
import { 
  Clock, 
  Eye, 
  EyeOff, 
  Key, 
  Trash2, 
  Copy, 
  Calendar,
  Activity,
  Shield,
  AlertCircle,
  Lock
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface APIKey {
  id: string;
  label: string;
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

interface APIKeyCardProps {
  apiKey: APIKey;
  onRevoke: () => void;
  onCopy: (text: string) => void;
  onViewKey: (keyId: string, password: string) => Promise<string | null>;
  formatDate: (date: string) => string;
}

const APIKeyCard: React.FC<APIKeyCardProps> = ({ 
  apiKey, 
  onRevoke, 
  onCopy,
  onViewKey,
  formatDate 
}) => {
  const [isViewingKey, setIsViewingKey] = useState(false);
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [actualKey, setActualKey] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isExpired = new Date(apiKey.expires_at) < new Date();
  const isNeverUsed = apiKey.last_used_at === "0001-01-01T05:30:00+05:30" || !apiKey.last_used_at;
  
  const getStatusColor = () => {
    if (!apiKey.is_active) return "text-red-300 bg-red-900/30 border-red-700";
    if (isExpired) return "text-orange-300 bg-orange-900/30 border-orange-700";
    return "text-green-300 bg-green-900/30 border-green-700";
  };

  const getStatusText = () => {
    if (!apiKey.is_active) return "Inactive";
    if (isExpired) return "Expired";
    return "Active";
  };

  const getKeyTypeDisplay = () => {
    return apiKey.key_type.charAt(0).toUpperCase() + apiKey.key_type.slice(1);
  };

  const handleViewKey = async () => {
    if (!password.trim()) {
      setError('Password is required');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const key = await onViewKey(apiKey.id, password);
      if (key) {
        setActualKey(key);
        setIsViewingKey(false); // Hide password form, show key
      }
    } catch (err: any) {
      setError(err.message || 'Failed to retrieve API key');
    } finally {
      setIsLoading(false);
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const resetView = () => {
    setIsViewingKey(false);
    setPassword('');
    setActualKey(null);
    setError(null);
  };

  return (
    <div className={cn(
      "group relative border border-gray-700 dark:border-gray-600 rounded-xl p-6 transition-all duration-200",
      "hover:border-gray-600 dark:hover:border-gray-500 hover:shadow-xl hover:shadow-gray-900/25 hover:-translate-y-0.5"
    )}
    style={{ backgroundColor: '#1f1f1f' }}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center space-x-3">
          <div className="p-2.5 rounded-lg border border-gray-600" style={{ backgroundColor: '#2a2a2a' }}>
            <Key className="w-5 h-5 text-blue-400 dark:text-blue-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-white dark:text-white group-hover:text-blue-400 dark:group-hover:text-blue-400 transition-colors">
              {apiKey.label}
            </h3>
            <div className="flex items-center space-x-2 mt-1">
              <span className="text-sm font-medium text-gray-300 dark:text-gray-300">
                {getKeyTypeDisplay()} Key
              </span>
              <span className="text-gray-400 dark:text-gray-400">•</span>
              <span className={cn(
                "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border",
                getStatusColor()
              )}>
                <Shield className="w-3 h-3 mr-1" />
                {getStatusText()}
              </span>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={() => setIsViewingKey(true)}
            className="p-2 text-gray-400 hover:text-green-400 rounded-lg transition-colors"
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#2a2a2a'}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            title="View API key"
          >
            <Eye className="w-4 h-4" />
          </button>
          <button
            onClick={onRevoke}
            className="p-2 text-gray-400 hover:text-red-400 rounded-lg transition-colors"
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#2a2a2a'}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            title="Revoke API key"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* API Key View Section */}
      {isViewingKey && !actualKey && (
        <div className="mb-4 p-4 border border-gray-600 rounded-lg" style={{ backgroundColor: '#2a2a2a' }}>
          <div className="flex items-center mb-3">
            <Lock className="w-4 h-4 text-yellow-400 mr-2" />
            <h4 className="text-sm font-medium text-white">Enter Password to View API Key</h4>
          </div>
          <div className="space-y-3">
            <div>
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                className="w-full px-3 py-2 border border-gray-600 rounded-lg text-sm text-white bg-gray-800 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                onKeyPress={(e) => e.key === 'Enter' && handleViewKey()}
              />
            </div>
            {error && (
              <p className="text-red-400 text-xs">{error}</p>
            )}
            <div className="flex items-center justify-between">
              <button
                onClick={togglePasswordVisibility}
                className="text-xs text-gray-400 hover:text-gray-200 flex items-center"
              >
                {showPassword ? <EyeOff className="w-3 h-3 mr-1" /> : <Eye className="w-3 h-3 mr-1" />}
                {showPassword ? 'Hide' : 'Show'} password
              </button>
              <div className="flex space-x-2">
                <button
                  onClick={resetView}
                  className="px-3 py-1 text-xs text-gray-400 hover:text-gray-200 border border-gray-600 rounded"
                >
                  Cancel
                </button>
                <button
                  onClick={handleViewKey}
                  disabled={isLoading || !password.trim()}
                  className="px-3 py-1 text-xs bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white rounded"
                >
                  {isLoading ? 'Loading...' : 'View Key'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Actual API Key Display */}
      {actualKey && (
        <div className="mb-4">
          <div className="flex items-center justify-between mb-2">
            <label className="text-sm font-medium text-gray-300">API Key</label>
            <button
              onClick={resetView}
              className="text-xs text-gray-400 hover:text-gray-200"
            >
              Hide Key
            </button>
          </div>
          <div className="relative">
            <code className="block w-full px-3 py-2.5 border border-green-600 rounded-lg text-sm font-mono text-green-400 select-all" style={{ backgroundColor: '#1a2e1a' }}>
              {actualKey}
            </code>
            <button
              onClick={() => onCopy(actualKey)}
              className="absolute right-2 top-1/2 transform -translate-y-1/2 p-1.5 text-green-400 hover:text-green-300 rounded transition-colors"
              title="Copy to clipboard"
            >
              <Copy className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}

      {/* Scopes */}
      <div className="mb-4">
        <label className="text-sm font-medium text-gray-300 mb-2 block">Permissions</label>
        <div className="flex flex-wrap gap-1.5">
          {apiKey.scopes.map((scope, index) => (
            <span
              key={index}
              className="inline-flex items-center px-2.5 py-1 text-blue-300 text-xs font-medium rounded-full border border-gray-600"
              style={{ backgroundColor: '#2a2a2a' }}
            >
              {scope === '*' ? 'Universal Access' : scope}
            </span>
          ))}
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div className="border border-gray-600 rounded-lg p-3" style={{ backgroundColor: '#2a2a2a' }}>
          <div className="flex items-center space-x-2 mb-1">
            <Activity className="w-4 h-4 text-gray-400" />
            <span className="text-sm font-medium text-gray-300">Usage</span>
          </div>
          <div className="text-lg font-bold text-white">
            {apiKey.usage_count.toLocaleString()}
          </div>
          <div className="text-xs text-gray-400">requests</div>
        </div>

        <div className="border border-gray-600 rounded-lg p-3" style={{ backgroundColor: '#2a2a2a' }}>
          <div className="flex items-center space-x-2 mb-1">
            <Clock className="w-4 h-4 text-gray-400" />
            <span className="text-sm font-medium text-gray-300">Rate Limit</span>
          </div>
          <div className="text-lg font-bold text-white">
            {apiKey.rate_limit || 'Unlimited'}
          </div>
          <div className="text-xs text-gray-400">per hour</div>
        </div>
      </div>

      {/* Dates */}
      <div className="space-y-3 pt-4 border-t border-gray-700">
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-2 text-gray-400">
            <Calendar className="w-4 h-4" />
            <span className="font-medium">Created</span>
          </div>
          <span className="text-white font-medium">
            {formatDate(apiKey.created_at)}
          </span>
        </div>

        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-2 text-gray-400">
            <Clock className="w-4 h-4" />
            <span className="font-medium">Expires</span>
          </div>
          <span className={cn(
            "font-medium",
            isExpired ? "text-red-400" : "text-white"
          )}>
            {formatDate(apiKey.expires_at)}
          </span>
        </div>

        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-2 text-gray-400">
            <Activity className="w-4 h-4" />
            <span className="font-medium">Last Used</span>
          </div>
          <span className="text-white font-medium">
            {isNeverUsed ? (
              <span className="inline-flex items-center text-gray-400">
                <AlertCircle className="w-3 h-3 mr-1" />
                Never used
              </span>
            ) : (
              formatDate(apiKey.last_used_at)
            )}
          </span>
        </div>
      </div>

      {/* Warning for expired/inactive keys */}
      {(!apiKey.is_active || isExpired) && (
        <div className={cn(
          "mt-4 p-3 rounded-lg border-l-4",
          !apiKey.is_active 
            ? "bg-red-900/30 border-red-500 text-red-300" 
            : "bg-orange-900/30 border-orange-500 text-orange-300"
        )}>
          <div className="flex items-center space-x-2">
            <AlertCircle className="w-4 h-4" />
            <span className="text-sm font-medium">
              {!apiKey.is_active 
                ? "This API key has been deactivated and cannot be used."
                : "This API key has expired and needs to be renewed."
              }
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default APIKeyCard;
