'use client';

import { useState, useEffect } from 'react';
import { TransactionService } from '@/lib/services';
import { Transaction } from '@/types/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { BarChart3 } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';

interface PaymentHistoryProps {
  limit?: number;
  showHeader?: boolean;
}

export default function PaymentHistory({ limit, showHeader = true }: PaymentHistoryProps) {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const pageSize = limit || 10;

  useEffect(() => {
    loadTransactions(0);
  }, []);

  const loadTransactions = async (pageNum: number) => {
    try {
      if (pageNum === 0) {
        setLoading(true);
      }
      setError('');
      
      const response = await TransactionService.getTransactionHistory({
        limit: pageSize,
        offset: pageNum * pageSize
      });
      
      if (response.data) {
        if (pageNum === 0) {
          setTransactions(response.data);
        } else {
          setTransactions(prev => [...prev, ...response.data!]);
        }
        
        setHasMore(response.data.length === pageSize);
        setPage(pageNum);
      } else {
        throw new Error('Failed to load payment history');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load payment history');
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    loadTransactions(page + 1);
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'completed': return 'success';
      case 'failed': return 'destructive';
      case 'pending': return 'warning';
      default: return 'secondary';
    }
  };

  const getTransactionTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'load_money':
      case 'refund':
        return 'text-green-600';
      case 'send_money':
      case 'ai_agent':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getTransactionSign = (type: string) => {
    switch (type.toLowerCase()) {
      case 'load_money':
      case 'refund':
        return '+';
      case 'send_money':
      case 'ai_agent':
        return '-';
      default:
        return '';
    }
  };

  const getTransactionTitle = (type: string) => {
    switch (type.toLowerCase()) {
      case 'load_money':
        return 'Money Loaded';
      case 'send_money':
        return 'Money Sent';
      case 'ai_agent':
        return 'AI Transaction';
      case 'refund':
        return 'Refund Received';
      default:
        return 'Transaction';
    }
  };

  if (loading && page === 0) {
    return (
      <Card>
        {showHeader && (
          <CardHeader>
            <CardTitle>Payment History</CardTitle>
            <CardDescription>Your transaction history and payment records</CardDescription>
          </CardHeader>
        )}
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error && transactions.length === 0) {
    return (
      <Card>
        {showHeader && (
          <CardHeader>
            <CardTitle>Payment History</CardTitle>
            <CardDescription>Your transaction history and payment records</CardDescription>
          </CardHeader>
        )}
        <CardContent>
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
          <Button onClick={() => loadTransactions(0)} className="mt-4" variant="outline" size="sm">
            Retry
          </Button>
        </CardContent>
      </Card>
    );
  }

  const content = (
    <div className="space-y-4">
      {transactions.length === 0 ? (
        <div className="text-center py-8">
          <div className="flex justify-center mb-4">
            <BarChart3 className="w-12 h-12 text-gray-400" />
          </div>
          <p className="text-gray-500 mb-4">No payment history found</p>
          <p className="text-sm text-gray-400">Your transactions will appear here once you start using the wallet.</p>
        </div>
      ) : (
        <>
          {transactions.map((transaction, index) => (
            <div
              key={`${transaction.id}-${index}`}
              className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="flex-1">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold text-gray-900">
                    {getTransactionTitle(transaction.type)}
                  </h3>
                  <div className={`font-bold ${getTransactionTypeColor(transaction.type)}`}>
                    {getTransactionSign(transaction.type)}{formatCurrency(transaction.amount)}
                  </div>
                </div>
                
                <div className="flex items-center justify-between">
                  <div className="text-sm text-gray-600">
                    <p>{formatDate(transaction.created_at)}</p>
                    {transaction.reference_id && (
                      <p className="text-xs text-gray-400 mt-1">
                        Ref: {transaction.reference_id}
                      </p>
                    )}
                    {transaction.description && (
                      <p className="text-xs text-gray-500 mt-1">
                        {transaction.description}
                      </p>
                    )}
                  </div>
                  
                  <div className="text-right">
                    <Badge variant={getStatusBadgeVariant(transaction.status)} className="mb-1">
                      {transaction.status}
                    </Badge>
                    {transaction.gateway_payment_id && (
                      <p className="text-xs text-gray-400">
                        Payment ID: {transaction.gateway_payment_id.slice(-8)}
                      </p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}

          {hasMore && (
            <div className="text-center pt-4">
              <Button 
                onClick={loadMore} 
                variant="outline" 
                disabled={loading}
              >
                {loading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-600 mr-2"></div>
                    Loading...
                  </div>
                ) : (
                  'Load More'
                )}
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );

  if (!showHeader) {
    return content;
  }

  return (
    <Card className={aeonikPro.className}>
      <CardHeader>
        <CardTitle>Payment History</CardTitle>
        <CardDescription>Your transaction history and payment records</CardDescription>
      </CardHeader>
      <CardContent>
        {content}
      </CardContent>
    </Card>
  );
}
