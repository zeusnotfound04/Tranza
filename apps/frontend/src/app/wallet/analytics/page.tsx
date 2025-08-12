'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '@/hooks/useAuth';
import { TransactionService, WalletService } from '@/lib/services';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Transaction, Wallet } from '@/types/api';

interface AnalyticsData {
  totalTransactions: number;
  totalSpent: number;
  totalReceived: number;
  avgTransactionAmount: number;
  monthlyData: {
    month: string;
    spent: number;
    received: number;
    transactions: number;
  }[];
  categoryBreakdown: {
    category: string;
    amount: number;
    count: number;
    percentage: number;
  }[];
  recentTrends: {
    period: string;
    change: number;
    type: 'increase' | 'decrease' | 'stable';
  };
}

function TransactionAnalyticsPage() {
  const router = useRouter();
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedPeriod, setSelectedPeriod] = useState<'30' | '90' | '365' | 'all'>('90');

  useEffect(() => {
    loadData();
  }, [selectedPeriod]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');

      // Load wallet info
      const walletResponse = await WalletService.getWallet();
      if (walletResponse.data) {
        setWallet(walletResponse.data);
      }

      // Load transactions based on selected period
      let fromDate: string | undefined;
      const toDate = new Date().toISOString();
      
      if (selectedPeriod !== 'all') {
        const days = parseInt(selectedPeriod);
        const date = new Date();
        date.setDate(date.getDate() - days);
        fromDate = date.toISOString();
      }

      const transactionsResponse = await TransactionService.getTransactionHistory({
        limit: 1000,
        ...(fromDate && { from_date: fromDate, to_date: toDate })
      });

      if (transactionsResponse.data) {
        const txns = transactionsResponse.data;
        setTransactions(txns);
        generateAnalytics(txns);
      }

    } catch (err: any) {
      setError(err.message || 'Failed to load analytics data');
    } finally {
      setLoading(false);
    }
  };

  const generateAnalytics = (txns: Transaction[]) => {
    if (txns.length === 0) {
      setAnalytics(null);
      return;
    }

    const totalTransactions = txns.length;
    const completedTxns = txns.filter(t => t.status === 'completed');
    
    const totalSpent = completedTxns
      .filter(t => ['debit', 'payment'].includes(t.type))
      .reduce((sum, t) => sum + t.amount, 0);
    
    const totalReceived = completedTxns
      .filter(t => ['credit', 'deposit'].includes(t.type))
      .reduce((sum, t) => sum + t.amount, 0);

    const avgTransactionAmount = completedTxns.length > 0 
      ? completedTxns.reduce((sum, t) => sum + t.amount, 0) / completedTxns.length 
      : 0;

    // Generate monthly data
    const monthlyMap = new Map<string, { spent: number; received: number; transactions: number }>();
    
    completedTxns.forEach(txn => {
      const date = new Date(txn.created_at);
      const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
      
      if (!monthlyMap.has(monthKey)) {
        monthlyMap.set(monthKey, { spent: 0, received: 0, transactions: 0 });
      }
      
      const monthData = monthlyMap.get(monthKey)!;
      monthData.transactions++;
      
      if (['debit', 'payment'].includes(txn.type)) {
        monthData.spent += txn.amount;
      } else if (['credit', 'deposit'].includes(txn.type)) {
        monthData.received += txn.amount;
      }
    });

    const monthlyData = Array.from(monthlyMap.entries())
      .map(([month, data]) => ({ month, ...data }))
      .sort((a, b) => a.month.localeCompare(b.month))
      .slice(-6); // Last 6 months

    // Generate category breakdown (based on merchant or transaction type)
    const categoryMap = new Map<string, { amount: number; count: number }>();
    
    completedTxns.forEach(txn => {
      const category = txn.merchant_name || 
                     (txn.type === 'payment' ? 'Payments' :
                      txn.type === 'deposit' ? 'Deposits' :
                      txn.type === 'credit' ? 'Credits' :
                      txn.type === 'debit' ? 'Debits' : 'Other');
      
      if (!categoryMap.has(category)) {
        categoryMap.set(category, { amount: 0, count: 0 });
      }
      
      const catData = categoryMap.get(category)!;
      catData.amount += txn.amount;
      catData.count++;
    });

    const totalCategoryAmount = Array.from(categoryMap.values()).reduce((sum, cat) => sum + cat.amount, 0);
    
    const categoryBreakdown = Array.from(categoryMap.entries())
      .map(([category, data]) => ({
        category,
        amount: data.amount,
        count: data.count,
        percentage: totalCategoryAmount > 0 ? (data.amount / totalCategoryAmount) * 100 : 0
      }))
      .sort((a, b) => b.amount - a.amount)
      .slice(0, 5); // Top 5 categories

    // Calculate recent trends (comparing last 30 days vs previous 30 days)
    const now = new Date();
    const thirtyDaysAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
    const sixtyDaysAgo = new Date(now.getTime() - 60 * 24 * 60 * 60 * 1000);

    const recentTxns = completedTxns.filter(t => new Date(t.created_at) >= thirtyDaysAgo);
    const previousTxns = completedTxns.filter(t => {
      const date = new Date(t.created_at);
      return date >= sixtyDaysAgo && date < thirtyDaysAgo;
    });

    const recentTotal = recentTxns.reduce((sum, t) => sum + t.amount, 0);
    const previousTotal = previousTxns.reduce((sum, t) => sum + t.amount, 0);

    let change = 0;
    let type: 'increase' | 'decrease' | 'stable' = 'stable';

    if (previousTotal > 0) {
      change = ((recentTotal - previousTotal) / previousTotal) * 100;
      type = change > 5 ? 'increase' : change < -5 ? 'decrease' : 'stable';
    } else if (recentTotal > 0) {
      change = 100;
      type = 'increase';
    }

    const recentTrends = {
      period: 'Last 30 days vs Previous 30 days',
      change: Math.abs(change),
      type
    };

    setAnalytics({
      totalTransactions,
      totalSpent,
      totalReceived,
      avgTransactionAmount,
      monthlyData,
      categoryBreakdown,
      recentTrends
    });
  };

  const exportData = () => {
    if (!analytics || !transactions.length) return;

    const csvData = transactions.map(txn => ({
      Date: new Date(txn.created_at).toLocaleDateString(),
      Type: txn.type,
      Amount: txn.amount,
      Status: txn.status,
      Merchant: txn.merchant_name || 'N/A',
      Description: txn.description || 'N/A',
      ID: txn.id
    }));

    const headers = Object.keys(csvData[0]);
    const csvContent = [
      headers.join(','),
      ...csvData.map(row => headers.map(h => `"${row[h as keyof typeof row] || ''}"`).join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `transactions-${selectedPeriod}-days-${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    window.URL.revokeObjectURL(url);
  };

  const formatCurrency = (amount: number) => {
    return `₹${amount.toLocaleString('en-IN', { 
      minimumFractionDigits: 2, 
      maximumFractionDigits: 2 
    })}`;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading analytics...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
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
                <h1 className="text-3xl font-bold text-gray-900">Transaction Analytics</h1>
                <p className="text-gray-600">Insights and trends from your transaction history</p>
              </div>
            </div>

            <div className="flex items-center space-x-4">
              {/* Period Selector */}
              <select
                value={selectedPeriod}
                onChange={(e) => setSelectedPeriod(e.target.value as any)}
                className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="30">Last 30 days</option>
                <option value="90">Last 90 days</option>
                <option value="365">Last year</option>
                <option value="all">All time</option>
              </select>

              <Button
                onClick={exportData}
                disabled={!analytics || !transactions.length}
                variant="outline"
              >
                Export CSV
              </Button>
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {!analytics && !loading && (
          <Card className="p-8 text-center">
            <p className="text-gray-600">No transaction data available for the selected period.</p>
            <Button 
              onClick={() => router.push('/wallet/load-money')}
              className="mt-4"
            >
              Add Money to Wallet
            </Button>
          </Card>
        )}

        {analytics && (
          <div className="space-y-8">
            {/* Overview Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                  </div>
                  <div className="ml-4">
                    <h3 className="text-sm font-medium text-gray-600">Total Transactions</h3>
                    <p className="text-2xl font-bold text-gray-900">{analytics.totalTransactions}</p>
                  </div>
                </div>
              </Card>

              <Card className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-red-100 rounded-lg">
                    <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                    </svg>
                  </div>
                  <div className="ml-4">
                    <h3 className="text-sm font-medium text-gray-600">Total Spent</h3>
                    <p className="text-2xl font-bold text-gray-900">{formatCurrency(analytics.totalSpent)}</p>
                  </div>
                </div>
              </Card>

              <Card className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-green-100 rounded-lg">
                    <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
                    </svg>
                  </div>
                  <div className="ml-4">
                    <h3 className="text-sm font-medium text-gray-600">Total Received</h3>
                    <p className="text-2xl font-bold text-gray-900">{formatCurrency(analytics.totalReceived)}</p>
                  </div>
                </div>
              </Card>

              <Card className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-purple-100 rounded-lg">
                    <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                    </svg>
                  </div>
                  <div className="ml-4">
                    <h3 className="text-sm font-medium text-gray-600">Avg Transaction</h3>
                    <p className="text-2xl font-bold text-gray-900">{formatCurrency(analytics.avgTransactionAmount)}</p>
                  </div>
                </div>
              </Card>
            </div>

            {/* Trends and Monthly Data */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Recent Trends */}
              <Card className="p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Trends</h3>
                <div className="flex items-center space-x-4">
                  <div className={`p-3 rounded-lg ${
                    analytics.recentTrends.type === 'increase' ? 'bg-green-100' :
                    analytics.recentTrends.type === 'decrease' ? 'bg-red-100' :
                    'bg-gray-100'
                  }`}>
                    {analytics.recentTrends.type === 'increase' && (
                      <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
                      </svg>
                    )}
                    {analytics.recentTrends.type === 'decrease' && (
                      <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                      </svg>
                    )}
                    {analytics.recentTrends.type === 'stable' && (
                      <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
                      </svg>
                    )}
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">{analytics.recentTrends.period}</p>
                    <p className={`text-lg font-semibold ${
                      analytics.recentTrends.type === 'increase' ? 'text-green-600' :
                      analytics.recentTrends.type === 'decrease' ? 'text-red-600' :
                      'text-gray-600'
                    }`}>
                      {analytics.recentTrends.type === 'stable' ? 'Stable' :
                       `${analytics.recentTrends.change.toFixed(1)}% ${analytics.recentTrends.type}`}
                    </p>
                  </div>
                </div>
              </Card>

              {/* Current Balance */}
              {wallet && (
                <Card className="p-6 bg-gradient-to-r from-blue-500 to-purple-600 text-white">
                  <h3 className="text-lg font-medium opacity-90 mb-2">Current Balance</h3>
                  <p className="text-3xl font-bold mb-4">
                    {formatCurrency(wallet.balance)}
                  </p>
                  <div className="flex items-center justify-between text-sm opacity-75">
                    <span>Status: {wallet.status.charAt(0).toUpperCase() + wallet.status.slice(1)}</span>
                    <span>AI Access: {wallet.ai_access_enabled ? 'Enabled' : 'Disabled'}</span>
                  </div>
                </Card>
              )}
            </div>

            {/* Category Breakdown */}
            <Card className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Categories</h3>
              <div className="space-y-4">
                {analytics.categoryBreakdown.map((category, index) => (
                  <div key={category.category} className="flex items-center space-x-4">
                    <div className="w-8 text-center">
                      <Badge variant="outline">#{index + 1}</Badge>
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-medium text-gray-900">{category.category}</span>
                        <span className="text-gray-600">{formatCurrency(category.amount)}</span>
                      </div>
                      <div className="flex items-center justify-between text-sm text-gray-500">
                        <span>{category.count} transactions</span>
                        <span>{category.percentage.toFixed(1)}% of total</span>
                      </div>
                      <div className="mt-2 bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${category.percentage}%` }}
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </Card>

            {/* Monthly Data */}
            {analytics.monthlyData.length > 0 && (
              <Card className="p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Monthly Overview</h3>
                <div className="overflow-x-auto">
                  <table className="min-w-full">
                    <thead>
                      <tr className="border-b border-gray-200">
                        <th className="text-left py-2 font-medium text-gray-600">Month</th>
                        <th className="text-right py-2 font-medium text-gray-600">Transactions</th>
                        <th className="text-right py-2 font-medium text-gray-600">Spent</th>
                        <th className="text-right py-2 font-medium text-gray-600">Received</th>
                        <th className="text-right py-2 font-medium text-gray-600">Net</th>
                      </tr>
                    </thead>
                    <tbody>
                      {analytics.monthlyData.map((month) => {
                        const net = month.received - month.spent;
                        return (
                          <tr key={month.month} className="border-b border-gray-100">
                            <td className="py-3 text-gray-900">
                              {new Date(`${month.month}-01`).toLocaleDateString('en-US', { 
                                year: 'numeric', 
                                month: 'long' 
                              })}
                            </td>
                            <td className="py-3 text-right text-gray-900">{month.transactions}</td>
                            <td className="py-3 text-right text-red-600">{formatCurrency(month.spent)}</td>
                            <td className="py-3 text-right text-green-600">{formatCurrency(month.received)}</td>
                            <td className={`py-3 text-right font-medium ${net >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                              {formatCurrency(Math.abs(net))}
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>
              </Card>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default withAuth(TransactionAnalyticsPage);
