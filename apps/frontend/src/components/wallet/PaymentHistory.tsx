'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { TransactionService } from '@/lib/services';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Input } from '@tranza/ui/components/ui/input';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Wallet, CreditCard } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import { Transaction } from '@/types/api';

interface PaymentHistoryProps {
  limit?: number;
  showFilters?: boolean;
}

export default function PaymentHistory({ limit, showFilters = true }: PaymentHistoryProps) {
  const { user } = useAuth();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  
  // Filters
  const [typeFilter, setTypeFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  
  const itemsPerPage = limit || 10;

  useEffect(() => {
    if (user) {
      loadTransactions();
    }
  }, [user, currentPage, typeFilter, statusFilter]);

  const loadTransactions = async () => {
    try {
      setLoading(true);
      setError('');
      
      const offset = (currentPage - 1) * itemsPerPage;
      
      const response = await TransactionService.getTransactionHistory({
        limit: itemsPerPage,
        offset,
        type: typeFilter || undefined,
        status: statusFilter || undefined,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
        merchant_name: searchTerm || undefined
      });
      
      if (response.data) {
        setTransactions(response.data);
        setTotalCount(response.total || 0);
        setTotalPages(Math.ceil((response.total || 0) / itemsPerPage));
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load transactions');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setCurrentPage(1);
    loadTransactions();
  };

  const clearFilters = () => {
    setTypeFilter('');
    setStatusFilter('');
    setSearchTerm('');
    setStartDate('');
    setEndDate('');
    setCurrentPage(1);
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'success':
        return 'default';
      case 'pending':
        return 'secondary';
      case 'failed':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const getTransactionIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'load_money':
        return <Wallet className="w-4 h-4 text-green-600" />;
      case 'ai_agent':
        return <CreditCard className="w-4 h-4 text-purple-600" />;
      case 'transfer':
        return <CreditCard className="w-4 h-4 text-blue-600" />;
      case 'refund':
        return <Wallet className="w-4 h-4 text-orange-600" />;
      default:
        return <CreditCard className="w-4 h-4 text-gray-600" />;
    }
  };

  const formatTransactionType = (type: string) => {
    return type.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
  };

  const formatAmount = (amount: number, type: string) => {
    const isInflow = ['load_money', 'refund'].includes(type.toLowerCase());
    const sign = isInflow ? '+' : '-';
    const color = isInflow ? 'text-green-600' : 'text-red-600';
    
    return (
      <span className={`font-semibold ${color}`}>
        {sign}₹{amount.toLocaleString('en-IN', { 
          minimumFractionDigits: 2, 
          maximumFractionDigits: 2 
        })}
      </span>
    );
  };

  if (loading && transactions.length === 0) {
    return (
      <div className="flex justify-center items-center h-32">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2 text-gray-600">Loading transactions...</span>
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${aeonikPro.className}`}>
      {/* Filters */}
      {showFilters && (
        <Card className="p-4">
          <div className="space-y-4">
            <h3 className="font-medium text-gray-900">Filters</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Type
                </label>
                <select
                  value={typeFilter}
                  onChange={(e) => setTypeFilter(e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">All Types</option>
                  <option value="load_money">Load Money</option>
                  <option value="ai_agent">AI Agent</option>
                  <option value="transfer">Transfer</option>
                  <option value="refund">Refund</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Status
                </label>
                <select
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">All Status</option>
                  <option value="success">Success</option>
                  <option value="pending">Pending</option>
                  <option value="failed">Failed</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Start Date
                </label>
                <input
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  End Date
                </label>
                <input
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>
            
            <div className="flex space-x-2">
              <Input
                placeholder="Search by description, reference ID, or merchant..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                className="flex-1"
              />
              <Button onClick={handleSearch} disabled={loading}>
                Search
              </Button>
              <Button variant="outline" onClick={clearFilters}>
                Clear
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Error */}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>
            {error}
            <Button onClick={loadTransactions} size="sm" variant="outline" className="ml-2">
              Retry
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Transactions List */}
      <Card className="overflow-hidden">
        {/* Header */}
        <div className="px-6 py-3 bg-gray-50 border-b">
          <div className="flex items-center justify-between">
            <h3 className="font-medium text-gray-900">
              Payment History
              {totalCount > 0 && (
                <span className="ml-2 text-sm text-gray-500">
                  ({totalCount} total)
                </span>
              )}
            </h3>
            {loading && (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            )}
          </div>
        </div>

        {/* Transactions */}
        <div className="divide-y divide-gray-200">
          {transactions.length === 0 ? (
            <div className="px-6 py-8 text-center">
              <p className="text-gray-500">No transactions found</p>
              {(typeFilter || statusFilter || searchTerm) && (
                <Button variant="outline" size="sm" onClick={clearFilters} className="mt-2">
                  Clear filters to see all transactions
                </Button>
              )}
            </div>
          ) : (
            transactions.map((transaction) => (
              <div key={transaction.id} className="px-6 py-4 hover:bg-gray-50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="flex items-center justify-center w-8 h-8">
                      {getTransactionIcon(transaction.type)}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">
                        {formatTransactionType(transaction.type)}
                      </p>
                      <p className="text-sm text-gray-600">
                        {transaction.description || 'No description'}
                      </p>
                      <p className="text-xs text-gray-500">
                        {new Date(transaction.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>
                  
                  <div className="text-right">
                    <div className="mb-1">
                      {formatAmount(transaction.amount, transaction.type)}
                    </div>
                    <Badge variant={getStatusBadgeVariant(transaction.status)}>
                      {transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1)}
                    </Badge>
                    {transaction.reference_id && (
                      <p className="text-xs text-gray-500 mt-1">
                        Ref: {transaction.reference_id.slice(-8)}
                      </p>
                    )}
                  </div>
                </div>
                
                {/* Additional Details */}
                {transaction.payment_method && (
                  <div className="mt-2 text-xs text-gray-500">
                    Payment Method: {transaction.payment_method.toUpperCase()}
                  </div>
                )}
                
                {transaction.status === 'failed' && (
                  <div className="mt-2 text-xs text-red-600">
                    Transaction failed - Please try again
                  </div>
                )}
              </div>
            ))
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="px-6 py-3 bg-gray-50 border-t flex items-center justify-between">
            <p className="text-sm text-gray-700">
              Showing {(currentPage - 1) * itemsPerPage + 1} to{' '}
              {Math.min(currentPage * itemsPerPage, totalCount)} of {totalCount} results
            </p>
            <div className="flex space-x-1">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1 || loading}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                disabled={currentPage === totalPages || loading}
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}
