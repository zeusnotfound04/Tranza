'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '@/hooks/useAuth';
import { WalletService } from '@/lib/services';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@/components/ui/Button';
import { Input } from '@tranza/ui/components/ui/input';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Wallet, UpdateWalletSettingsRequest } from '@/types/api';

function WalletSettingsPage() {
  const router = useRouter();
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Form state
  const [aiAccessEnabled, setAiAccessEnabled] = useState(false);
  const [aiDailyLimit, setAiDailyLimit] = useState('');
  const [aiPerTransactionLimit, setAiPerTransactionLimit] = useState('');

  useEffect(() => {
    loadWallet();
  }, []);

  const loadWallet = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await WalletService.getWallet();
      if (response.data) {
        const walletData = response.data;
        setWallet(walletData);
        
        // Set form values
        setAiAccessEnabled(walletData.ai_access_enabled);
        setAiDailyLimit(walletData.ai_daily_limit?.toString() || '1000');
        setAiPerTransactionLimit(walletData.ai_per_transaction_limit?.toString() || '100');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load wallet settings');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveSettings = async () => {
    try {
      setSaving(true);
      setError('');
      setSuccess('');

      // Validate inputs
      const dailyLimit = parseFloat(aiDailyLimit);
      const transactionLimit = parseFloat(aiPerTransactionLimit);

      if (isNaN(dailyLimit) || dailyLimit < 0 || dailyLimit > 50000) {
        setError('Daily limit must be between ₹0 and ₹50,000');
        return;
      }

      if (isNaN(transactionLimit) || transactionLimit < 0 || transactionLimit > 10000) {
        setError('Per transaction limit must be between ₹0 and ₹10,000');
        return;
      }

      if (transactionLimit > dailyLimit) {
        setError('Per transaction limit cannot exceed daily limit');
        return;
      }

      const updateData: UpdateWalletSettingsRequest = {
        ai_access_enabled: aiAccessEnabled,
        ai_daily_limit: dailyLimit,
        ai_per_transaction_limit: transactionLimit
      };

      await WalletService.updateSettings(updateData);
      
      setSuccess('Settings updated successfully!');
      
      // Reload wallet data
      await loadWallet();
      
      // Clear success message after 3 seconds
      setTimeout(() => setSuccess(''), 3000);

    } catch (err: any) {
      setError(err.message || 'Failed to update settings');
    } finally {
      setSaving(false);
    }
  };

  const handleReset = () => {
    if (wallet) {
      setAiAccessEnabled(wallet.ai_access_enabled);
      setAiDailyLimit(wallet.ai_daily_limit?.toString() || '1000');
      setAiPerTransactionLimit(wallet.ai_per_transaction_limit?.toString() || '100');
      setError('');
      setSuccess('');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading wallet settings...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center space-x-4">
            <Button
              variant="outline"
              onClick={() => router.back()}
              className="flex items-center space-x-2"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              <span>Back</span>
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Wallet Settings</h1>
              <p className="text-gray-600">Manage your wallet preferences and AI agent limits</p>
            </div>
          </div>
        </div>

        {/* Success/Error Messages */}
        {success && (
          <Alert className="mb-6 bg-green-50 border-green-200">
            <AlertDescription className="text-green-800">
              ✅ {success}
            </AlertDescription>
          </Alert>
        )}

        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Current Wallet Info */}
        {wallet && (
          <Card className="p-6 mb-8 bg-gradient-to-r from-blue-500 to-purple-600 text-white">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-medium opacity-90">Current Balance</h2>
                <p className="text-3xl font-bold mt-1">
                  ₹{wallet.balance.toLocaleString('en-IN', { 
                    minimumFractionDigits: 2, 
                    maximumFractionDigits: 2 
                  })}
                </p>
              </div>
              <div className="text-right">
                <Badge 
                  variant={wallet.status === 'active' ? 'default' : 'secondary'}
                  className="bg-white/20 text-white border-white/30"
                >
                  {wallet.status.charAt(0).toUpperCase() + wallet.status.slice(1)}
                </Badge>
                <p className="text-xs opacity-75 mt-2">
                  Wallet ID: {wallet.id.slice(0, 8)}...
                </p>
              </div>
            </div>
          </Card>
        )}

        {/* AI Agent Settings */}
        <Card className="p-6">
          <div className="space-y-6">
            <div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">AI Agent Access</h3>
              <p className="text-gray-600 mb-4">
                Control how AI agents can access and spend from your wallet
              </p>
            </div>

            {/* AI Access Toggle */}
            <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
              <div>
                <h4 className="font-medium text-gray-900">Enable AI Agent Access</h4>
                <p className="text-sm text-gray-600">
                  Allow AI agents to make transactions using your wallet balance
                </p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={aiAccessEnabled}
                  onChange={(e) => setAiAccessEnabled(e.target.checked)}
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            {/* Spending Limits */}
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Daily Spending Limit for AI Agents
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">₹</span>
                  <Input
                    type="number"
                    value={aiDailyLimit}
                    onChange={(e) => setAiDailyLimit(e.target.value)}
                    placeholder="1000"
                    min="0"
                    max="50000"
                    className="pl-8"
                    disabled={!aiAccessEnabled}
                  />
                </div>
                <p className="text-xs text-gray-500 mt-1">
                  Maximum amount AI agents can spend per day (₹0 - ₹50,000)
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Per Transaction Limit for AI Agents
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">₹</span>
                  <Input
                    type="number"
                    value={aiPerTransactionLimit}
                    onChange={(e) => setAiPerTransactionLimit(e.target.value)}
                    placeholder="100"
                    min="0"
                    max="10000"
                    className="pl-8"
                    disabled={!aiAccessEnabled}
                  />
                </div>
                <p className="text-xs text-gray-500 mt-1">
                  Maximum amount per single AI transaction (₹0 - ₹10,000)
                </p>
              </div>
            </div>

            {/* Current Settings Preview */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900 mb-2">Settings Summary</h4>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                <div>
                  <span className="text-gray-600">AI Access:</span>
                  <p className="font-medium">
                    {aiAccessEnabled ? (
                      <span className="text-green-600">Enabled</span>
                    ) : (
                      <span className="text-red-600">Disabled</span>
                    )}
                  </p>
                </div>
                <div>
                  <span className="text-gray-600">Daily Limit:</span>
                  <p className="font-medium">
                    ₹{parseFloat(aiDailyLimit || '0').toLocaleString('en-IN')}
                  </p>
                </div>
                <div>
                  <span className="text-gray-600">Per Transaction:</span>
                  <p className="font-medium">
                    ₹{parseFloat(aiPerTransactionLimit || '0').toLocaleString('en-IN')}
                  </p>
                </div>
              </div>
            </div>

            {/* Security Notice */}
            <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
              <div className="flex items-start space-x-2">
                <svg className="w-5 h-5 text-yellow-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.728-.833-2.498 0L3.316 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
                <div>
                  <h5 className="font-medium text-yellow-800">Security Notice</h5>
                  <p className="text-sm text-yellow-700 mt-1">
                    Only enable AI agent access if you trust the AI services you're using. 
                    Set conservative limits to minimize risk. You can disable or adjust these 
                    settings at any time.
                  </p>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex space-x-3 pt-4">
              <Button
                onClick={handleSaveSettings}
                disabled={saving || loading}
                className="flex-1"
              >
                {saving ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Saving...
                  </div>
                ) : (
                  'Save Settings'
                )}
              </Button>
              
              <Button
                variant="outline"
                onClick={handleReset}
                disabled={saving || loading}
              >
                Reset
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}

export default withAuth(WalletSettingsPage);
