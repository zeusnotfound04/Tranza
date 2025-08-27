'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '@/hooks/useAuth';
import { TransactionService } from '@/lib/services';
import PaymentHistory from '@/components/wallet/PaymentHistory';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@/components/ui/Button';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Transaction, TransactionFilters } from '@/types/api';

function WalletHistoryPage() {
  const router = useRouter();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [totalTransactions, setTotalTransactions] = useState(0);
  const [filters, setFilters] = useState<TransactionFilters>({
    limit: 20,
    offset: 0
  });

  // Quick filter states
  const [quickFilter, setQuickFilter] = useState<'all' | 'today' | 'week' | 'month'>('all');
  const [typeFilter, setTypeFilter] = useState<string>('all');

  useEffect(() => {
    loadTransactions();
  }, [filters, quickFilter, typeFilter]);

  const loadTransactions = async () => {
    try {
      setLoading(true);
      setError('');

      // Build filters based on quick filter and type filter
      const currentFilters: TransactionFilters = { ...filters };

      // Apply date filters based on quick filter
      if (quickFilter !== 'all') {
        const now = new Date();
        let fromDate: Date;

        switch (quickFilter) {
          case 'today':
            fromDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            break;
          case 'week':
            fromDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
            break;
          case 'month':
            fromDate = new Date(now.getFullYear(), now.getMonth(), 1);
            break;
          default:
            fromDate = new Date(0);
        }

        currentFilters.start_date = fromDate.toISOString();
        currentFilters.end_date = now.toISOString();
      }

      // Apply type filter
      if (typeFilter !== 'all') {
        currentFilters.type = typeFilter;
      }

      const response = await TransactionService.getTransactionHistory(currentFilters);
      
      if (response.data) {
        setTransactions(response.data);
        setTotalTransactions(response.total || 0);
      }

    } catch (err: any) {
      setError(err.message || 'Failed to load transaction history');
    } finally {
      setLoading(false);
    }
  };

  const handleFiltersChange = (newFilters: TransactionFilters) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  };

  const handleQuickFilterChange = (filter: 'all' | 'today' | 'week' | 'month') => {
    setQuickFilter(filter);
    // Reset pagination when changing filters
    setFilters(prev => ({ ...prev, offset: 0 }));
  };

  const handleTypeFilterChange = (type: string) => {
    setTypeFilter(type);
    // Reset pagination when changing filters
    setFilters(prev => ({ ...prev, offset: 0 }));
  };

  const handlePagination = (newOffset: number) => {
    setFilters(prev => ({ ...prev, offset: newOffset }));
  };

  const handleExport = () => {
    if (!transactions.length) return;

    const csvData = transactions.map(txn => ({
      Date: new Date(txn.created_at).toLocaleDateString(),
      Time: new Date(txn.created_at).toLocaleTimeString(),
      Type: txn.type.charAt(0).toUpperCase() + txn.type.slice(1),
      Amount: txn.amount,
      Status: txn.status.charAt(0).toUpperCase() + txn.status.slice(1),
      Merchant: txn.merchant_name || 'N/A',
      Description: txn.description || 'N/A',
      'Transaction ID': txn.id
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
    link.download = `wallet-history-${quickFilter}-${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    window.URL.revokeObjectURL(url);
  };

  const getQuickStats = () => {
    if (!transactions.length) return { spent: 0, received: 0, net: 0 };

    const spent = transactions
      .filter(t => ['debit', 'payment'].includes(t.type) && t.status === 'completed')
      .reduce((sum, t) => sum + t.amount, 0);

    const received = transactions
      .filter(t => ['credit', 'deposit'].includes(t.type) && t.status === 'completed')
      .reduce((sum, t) => sum + t.amount, 0);

    return { spent, received, net: received - spent };
  };

  const stats = getQuickStats();

  if (loading && transactions.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading transaction history...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
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
                <h1 className="text-3xl font-bold text-gray-900">Transaction History</h1>
                <p className="text-gray-600">
                  {totalTransactions > 0 && `${totalTransactions} transaction${totalTransactions !== 1 ? 's' : ''} found`}
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-3">
              <Button
                onClick={() => router.push('/wallet/analytics')}
                variant="outline"
                className="flex items-center space-x-2"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                <span>Analytics</span>
              </Button>

              <Button
                onClick={handleExport}
                disabled={!transactions.length}
                variant="outline"
                className="flex items-center space-x-2"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span>Export</span>
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

        {/* Quick Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <Card className="p-4">
            <div className="flex items-center">
              <div className="p-2 bg-red-100 rounded-lg">
                <svg className="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-gray-600">Total Spent</p>
                <p className="text-xl font-semibold text-gray-900">
                  ₹{stats.spent.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
              </div>
            </div>
          </Card>

          <Card className="p-4">
            <div className="flex items-center">
              <div className="p-2 bg-green-100 rounded-lg">
                <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-gray-600">Total Received</p>
                <p className="text-xl font-semibold text-gray-900">
                  ₹{stats.received.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
              </div>
            </div>
          </Card>

          <Card className="p-4">
            <div className="flex items-center">
              <div className={`p-2 rounded-lg ${stats.net >= 0 ? 'bg-blue-100' : 'bg-orange-100'}`}>
                <svg className={`w-5 h-5 ${stats.net >= 0 ? 'text-blue-600' : 'text-orange-600'}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-gray-600">Net Change</p>
                <p className={`text-xl font-semibold ${stats.net >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {stats.net >= 0 ? '+' : ''}₹{Math.abs(stats.net).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
              </div>
            </div>
          </Card>
        </div>

        {/* Quick Filters */}
        <Card className="p-4 mb-6">
          <div className="flex flex-wrap items-center justify-between gap-4">
            {/* Date Filters */}
            <div className="flex items-center space-x-2">
              <span className="text-sm font-medium text-gray-700">Period:</span>
              <div className="flex space-x-1">
                {(['all', 'today', 'week', 'month'] as const).map((filter) => (
                  <button
                    key={filter}
                    onClick={() => handleQuickFilterChange(filter)}
                    className={`px-3 py-1 text-sm rounded-full transition-colors ${
                      quickFilter === filter
                        ? 'bg-blue-100 text-blue-800 font-medium'
                        : 'text-gray-600 hover:bg-gray-100'
                    }`}
                  >
                    {filter.charAt(0).toUpperCase() + filter.slice(1)}
                  </button>
                ))}
              </div>
            </div>

            {/* Type Filters */}
            <div className="flex items-center space-x-2">
              <span className="text-sm font-medium text-gray-700">Type:</span>
              <select
                value={typeFilter}
                onChange={(e) => handleTypeFilterChange(e.target.value)}
                className="px-3 py-1 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Types</option>
                <option value="credit">Credits</option>
                <option value="debit">Debits</option>
                <option value="payment">Payments</option>
                <option value="deposit">Deposits</option>
              </select>
            </div>
          </div>
        </Card>

        {/* Transaction List */}
        {!loading && transactions.length === 0 ? (
          <Card className="p-8 text-center">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No transactions found</h3>
            <p className="text-gray-600 mb-6">
              {quickFilter !== 'all' || typeFilter !== 'all'
                ? 'Try adjusting your filters to see more transactions.'
                : 'Start by adding money to your wallet to see transactions here.'
              }
            </p>
            <div className="flex justify-center space-x-3">
              {(quickFilter !== 'all' || typeFilter !== 'all') && (
                <Button
                  variant="outline"
                  onClick={() => {
                    setQuickFilter('all');
                    setTypeFilter('all');
                  }}
                >
                  Clear Filters
                </Button>
              )}
              <Button onClick={() => router.push('/wallet/load-money')}>
                Add Money
              </Button>
            </div>
          </Card>
        ) : (
          <div className="space-y-6">
            {/* Transaction Cards */}
            <div className="space-y-3">
              {transactions.map((transaction) => {
                const isDebit = ['debit', 'payment'].includes(transaction.type);
                const isCredit = ['credit', 'deposit'].includes(transaction.type);
                
                return (
                  <Card key={transaction.id} className="p-4">
                    <div className="flex items-center justify-between">
                      {/* Transaction Info */}
                      <div className="flex items-center space-x-4">
                        {/* Icon */}
                        <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                          isDebit ? 'bg-red-100' : 
                          isCredit ? 'bg-green-100' : 'bg-gray-100'
                        }`}>
                          {isDebit && (
                            <svg className="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                            </svg>
                          )}
                          {isCredit && (
                            <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
                            </svg>
                          )}
                          {!isDebit && !isCredit && (
                            <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
                            </svg>
                          )}
                        </div>

                        {/* Details */}
                        <div className="flex-1">
                          <div className="flex items-center space-x-2">
                            <h3 className="font-medium text-gray-900">
                              {transaction.merchant_name || 
                               (transaction.type === 'payment' ? 'Payment' :
                                transaction.type === 'deposit' ? 'Deposit' :
                                transaction.type === 'credit' ? 'Credit' :
                                transaction.type === 'debit' ? 'Debit' : 'Transaction')}
                            </h3>
                            <Badge
                              variant={
                                transaction.status === 'completed' ? 'default' :
                                transaction.status === 'pending' ? 'secondary' :
                                transaction.status === 'failed' ? 'destructive' : 'outline'
                              }
                            >
                              {transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1)}
                            </Badge>
                          </div>
                          <p className="text-sm text-gray-600 mt-1">
                            {transaction.description || `${transaction.type.charAt(0).toUpperCase() + transaction.type.slice(1)} transaction`}
                          </p>
                          <p className="text-xs text-gray-500 mt-1">
                            {new Date(transaction.created_at).toLocaleDateString('en-IN', {
                              year: 'numeric',
                              month: 'short',
                              day: 'numeric',
                              hour: '2-digit',
                              minute: '2-digit'
                            })}
                          </p>
                        </div>
                      </div>

                      {/* Amount */}
                      <div className="text-right">
                        <p className={`text-lg font-semibold ${
                          isDebit ? 'text-red-600' : 
                          isCredit ? 'text-green-600' : 'text-gray-900'
                        }`}>
                          {isDebit ? '-' : isCredit ? '+' : ''}₹{transaction.amount.toLocaleString('en-IN', {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2
                          })}
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                          ID: {transaction.id.slice(0, 8)}...
                        </p>
                      </div>
                    </div>
                  </Card>
                );
              })}
            </div>

            {/* Pagination */}
            {totalTransactions > filters.limit! && (
              <div className="flex items-center justify-between">
                <p className="text-sm text-gray-600">
                  Showing {filters.offset! + 1} to {Math.min(filters.offset! + filters.limit!, totalTransactions)} of {totalTransactions} transactions
                </p>
                
                <div className="flex items-center space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePagination(Math.max(0, filters.offset! - filters.limit!))}
                    disabled={filters.offset === 0}
                  >
                    Previous
                  </Button>
                  
                  <span className="text-sm text-gray-600">
                    Page {Math.floor(filters.offset! / filters.limit!) + 1} of {Math.ceil(totalTransactions / filters.limit!)}
                  </span>
                  
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePagination(filters.offset! + filters.limit!)}
                    disabled={filters.offset! + filters.limit! >= totalTransactions}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}

            {/* Load More for Mobile */}
            {totalTransactions > transactions.length && (
              <div className="text-center">
                <Button
                  variant="outline"
                  onClick={() => handlePagination(filters.offset! + filters.limit!)}
                  disabled={loading}
                  className="w-full sm:w-auto"
                >
                  {loading ? 'Loading...' : 'Load More Transactions'}
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default withAuth(WalletHistoryPage);
