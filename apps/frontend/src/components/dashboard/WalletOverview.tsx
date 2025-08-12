'use client';

import { useState, useEffect } from 'react';
import { WalletService } from '@/lib/services';
import { Wallet } from '@/types/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import Link from 'next/link';

interface WalletData {
  id: string;
  balance: number;
  currency: string;
  status: string;
  ai_daily_limit: number;
  ai_per_transaction_limit: number;
  // Extended properties for display (these might come from backend calculations)
  daily_spent?: number;
  monthly_spent?: number;
  ai_daily_spent?: number;
}

export default function WalletOverview() {
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadWalletData();
  }, []);

  const loadWalletData = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await WalletService.getWallet();
      if (response.data) {
        setWallet(response.data);
      } else {
        throw new Error(response.message || 'Failed to load wallet');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load wallet data');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Wallet Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Wallet Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
          <Button onClick={loadWalletData} className="mt-4" variant="outline">
            Retry
          </Button>
        </CardContent>
      </Card>
    );
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: wallet?.currency || 'INR',
    }).format(amount);
  };

  const getSpendingPercentage = (spent: number, limit: number) => {
    return limit > 0 ? Math.round((spent / limit) * 100) : 0;
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active': return 'success';
      case 'suspended': return 'destructive';
      case 'pending': return 'warning';
      default: return 'secondary';
    }
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Main Balance Card */}
      <Card className="md:col-span-2">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className='text-2xl'>Wallet Balance</CardTitle>
              <CardDescription>Your current available balance</CardDescription>
            </div>
            <Badge variant={getStatusBadgeVariant(wallet?.status || '')}>
              {wallet?.status || 'Unknown'}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold text-green-600 mb-4">
            {wallet ? formatCurrency(wallet.balance) : '₹0.00'}
          </div>
          <div className="flex space-x-2">
            <Link href="/wallet/load">
              <Button  size="sm">Load Money</Button>
            </Link>
            <Link href="/transactions/send">
              <Button variant="outline" size="sm">Send Money</Button>
            </Link>
          </div>
        </CardContent>
      </Card>

      {/* Wallet Info */}
      <Card>
        <CardHeader>
          <CardTitle>Wallet Information</CardTitle>
          <CardDescription>Account details and status</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600">Currency:</span>
              <span className="text-sm text-black font-semibold">{wallet?.currency || 'INR'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm text-gray-600">Status:</span>
              <Badge variant={getStatusBadgeVariant(wallet?.status || '')}>
                {wallet?.status || 'Unknown'}
              </Badge>
            </div>
            <div className="flex justify-between">
              <span className="text-sm text-gray-600">AI Access:</span>
              <Badge variant={wallet?.ai_access_enabled ? 'success' : 'secondary'}>
                {wallet?.ai_access_enabled ? 'Enabled' : 'Disabled'}
              </Badge>
            </div>
            <div className="pt-2">
              <Link href="/wallet/settings">
                <Button variant="outline" size="sm" className="w-full">
                  Wallet Settings
                </Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* AI Agent Limits */}
      <Card>
        <CardHeader>
          <CardTitle>AI Agent Limits</CardTitle>
          <CardDescription>Automated transaction controls</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Daily Limit:</span>
                <span className="font-bold text-black ">{formatCurrency(wallet?.ai_daily_limit || 0)}</span>
              </div>
            </div>
            
            <div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Per Transaction Limit:</span>
                <span className="font-bold text-black">{formatCurrency(wallet?.ai_per_transaction_limit || 0)}</span>
              </div>
            </div>
            
            <div className="pt-2 space-y-2">
              <Link href="/ai-agents">
                <Button variant="outline" size="sm" className="w-full">
                  Manage AI Agents
                </Button>
              </Link>
              <Link href="/api-keys">
                <Button variant="outline" size="sm" className="w-full">
                  API Keys
                </Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
