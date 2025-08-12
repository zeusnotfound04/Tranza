'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { WalletService } from '@/lib/services';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Wallet as WalletIcon, BarChart3, TrendingUp } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import Link from 'next/link';
import { Wallet } from '@/types/api';

export default function WalletDashboard() {
  const { user } = useAuth();
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (user) {
      loadWalletData();
    }
  }, [user]);

  const loadWalletData = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await WalletService.getWallet();
      if (response.data) {
        setWallet(response.data);
      }
    } catch (err: any) {
      if (err.status === 404) {
        setError('Wallet not found. A wallet should have been created automatically during signup.');
      } else {
        setError(err.message || 'Failed to load wallet data');
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2 text-gray-600">Loading wallet...</span>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>
          {error}
          <Button 
            onClick={loadWalletData} 
            size="sm" 
            variant="outline" 
            className="ml-2"
          >
            Retry
          </Button>
        </AlertDescription>
      </Alert>
    );
  }

  if (!wallet) {
    return (
      <Card className="p-6 text-center">
        <h3 className="text-lg font-medium text-gray-900 mb-2">No Wallet Found</h3>
        <p className="text-gray-600 mb-4">
          Your wallet should have been created automatically. Please contact support.
        </p>
        <Button onClick={loadWalletData}>Retry Loading</Button>
      </Card>
    );
  }

  return (
    <div className={`space-y-6 ${aeonikPro.className}`}>
      {/* Main Wallet Balance Card */}
      <Card className="p-6 bg-gradient-to-r from-blue-500 to-purple-600 text-white">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-medium opacity-90">Wallet Balance</h2>
            <p className="text-3xl font-bold mt-2">
              ₹{wallet.balance.toLocaleString('en-IN', { 
                minimumFractionDigits: 2, 
                maximumFractionDigits: 2 
              })}
            </p>
            <p className="text-sm opacity-75 mt-1">
              Available Balance
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

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Link href="/wallet/load-money">
          <Card className="p-4 hover:shadow-md transition-shadow cursor-pointer border-2 hover:border-green-200">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                <WalletIcon className="text-green-600 w-6 h-6" />
              </div>
              <div>
                <h3 className="font-medium text-gray-900">Load Money</h3>
                <p className="text-sm text-gray-600">Add funds to your wallet</p>
              </div>
            </div>
          </Card>
        </Link>

        <Link href="/wallet/history">
          <Card className="p-4 hover:shadow-md transition-shadow cursor-pointer border-2 hover:border-blue-200">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                <BarChart3 className="text-blue-600 w-6 h-6" />
              </div>
              <div>
                <h3 className="font-medium text-gray-900">Transaction History</h3>
                <p className="text-sm text-gray-600">View all transactions</p>
              </div>
            </div>
          </Card>
        </Link>

        <Link href="/wallet/analytics">
          <Card className="p-4 hover:shadow-md transition-shadow cursor-pointer border-2 hover:border-purple-200">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
                <TrendingUp className="text-purple-600 w-6 h-6" />
              </div>
              <div>
                <h3 className="font-medium text-gray-900">Analytics</h3>
                <p className="text-sm text-gray-600">Spending insights & trends</p>
              </div>
            </div>
          </Card>
        </Link>

        <Link href="/wallet/settings">
          <Card className="p-4 hover:shadow-md transition-shadow cursor-pointer border-2 hover:border-yellow-200">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
                <span className="text-yellow-600 text-xl">⚙️</span>
              </div>
              <div>
                <h3 className="font-medium text-gray-900">Settings</h3>
                <p className="text-sm text-gray-600">AI limits & preferences</p>
              </div>
            </div>
          </Card>
        </Link>
      </div>

      {/* AI Spending Limits */}
      <Card className="p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">AI Agent Settings</h3>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900">AI Access</p>
              <p className="text-sm text-gray-600">Allow AI agents to spend from wallet</p>
            </div>
            <Badge variant={wallet.ai_access_enabled ? 'default' : 'secondary'}>
              {wallet.ai_access_enabled ? 'Enabled' : 'Disabled'}
            </Badge>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="font-medium text-gray-900">Daily Limit</p>
              <p className="text-sm text-gray-600">
                ₹{wallet.ai_daily_limit?.toLocaleString('en-IN') || '0'} per day
              </p>
            </div>
            <div>
              <p className="font-medium text-gray-900">Per Transaction Limit</p>
              <p className="text-sm text-gray-600">
                ₹{wallet.ai_per_transaction_limit?.toLocaleString('en-IN') || '0'} per transaction
              </p>
            </div>
          </div>

          <Link href="/wallet/settings">
            <Button variant="outline" size="sm">
              Update Limits
            </Button>
          </Link>
        </div>
      </Card>

      {/* Wallet Stats - Coming Soon */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <div className="text-center">
            <p className="text-2xl font-bold text-gray-900">
              ₹{wallet.balance.toLocaleString('en-IN')}
            </p>
            <p className="text-sm text-gray-600">Current Balance</p>
          </div>
        </Card>

        <Card className="p-4">
          <div className="text-center">
            <p className="text-2xl font-bold text-gray-900">
              ₹{wallet.ai_daily_limit?.toLocaleString('en-IN') || '0'}
            </p>
            <p className="text-sm text-gray-600">AI Daily Limit</p>
          </div>
        </Card>

        <Card className="p-4">
          <div className="text-center">
            <p className="text-2xl font-bold text-gray-900">
              ₹{wallet.ai_per_transaction_limit?.toLocaleString('en-IN') || '0'}
            </p>
            <p className="text-sm text-gray-600">AI Per Txn Limit</p>
          </div>
        </Card>
      </div>

      {/* Last Updated */}
      <div className="text-center">
        <p className="text-xs text-gray-500">
          Last updated: {new Date().toLocaleString()}
        </p>
        <Button 
          onClick={loadWalletData} 
          variant="outline" 
          size="sm" 
          className="mt-2"
        >
          Refresh
        </Button>
      </div>
    </div>
  );
}
