'use client';

import { useState, useEffect } from 'react';
import { TransactionService } from '@/lib/services';
import { Transaction } from '@/types/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { TransactionBadge } from '@/components/ui/TransactionBadge';
import { Wallet, ArrowUpRight, ArrowDownLeft, BarChart3, ExternalLink } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import Link from 'next/link';

export default function RecentTransactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadRecentTransactions();
  }, []);

  const loadRecentTransactions = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await TransactionService.getTransactionHistory({
        limit: 5,
        offset: 0
      });
      
      if (response.data) {
        setTransactions(response.data);
      } else {
        throw new Error('Failed to load transactions');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load recent transactions');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
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

  const getTransactionIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'load_money':
        return <Wallet className="w-5 h-5 text-green-600" />;
      case 'send_money':
        return <ArrowUpRight className="w-5 h-5 text-blue-600" />;
      case 'ai_agent':
        return <ArrowDownLeft className="w-5 h-5 text-purple-600" />;
      case 'refund':
        return <ArrowDownLeft className="w-5 h-5 text-orange-600" />;
      default:
        return <Wallet className="w-5 h-5 text-gray-600" />;
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
        return 'Refund';
      default:
        return 'Transaction';
    }
  };

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent Transactions</CardTitle>
          <CardDescription>Your latest transaction activity</CardDescription>
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
          <CardTitle>Recent Transactions</CardTitle>
          <CardDescription>Your latest transaction activity</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
          <Button onClick={loadRecentTransactions} className="mt-4" variant="outline" size="sm">
            Retry
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={aeonikPro.className}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Recent Transactions</CardTitle>
            <CardDescription>Your latest transaction activity</CardDescription>
          </div>
          <Link href="/transactions">
            <Button variant="outline" size="sm">View All</Button>
          </Link>
        </div>
      </CardHeader>
      <CardContent>
        {transactions.length === 0 ? (
          <div className="text-center py-8">
            <div className="flex justify-center mb-4">
              <BarChart3 className="w-12 h-12 text-gray-400" />
            </div>
            <p className="text-gray-500">No transactions yet</p>
            <Link href="/wallet/load" className="mt-4 inline-block">
              <Button size="sm">Load Money to Get Started</Button>
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {transactions.map((transaction) => (
              <div
                key={transaction.id}
                className="flex items-center justify-between p-4 bg-gradient-to-r from-white to-slate-50/30 rounded-xl hover:from-slate-50 hover:to-slate-100/50 transition-all duration-300 border border-slate-200/60 hover:border-slate-300/60 hover:shadow-md"
              >
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <TransactionBadge 
                      type={transaction.type} 
                      showStatus={false} 
                      size="md" 
                      variant="minimal"
                    />
                  </div>
                  <div>
                    <p className="font-semibold text-slate-900">
                      {getTransactionTitle(transaction.type)}
                    </p>
                    <div className="flex items-center space-x-2">
                      <p className="text-sm text-slate-600">
                        {formatDate(transaction.created_at)}
                      </p>
                      {transaction.reference_id && (
                        <span className="text-xs text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full">
                          #{transaction.reference_id.slice(-6)}
                        </span>
                      )}
                    </div>
                    {transaction.description && (
                      <p className="text-xs text-slate-500 mt-1">
                        {transaction.description}
                      </p>
                    )}
                  </div>
                </div>
                <div className="text-right space-y-1">
                  <p className={`font-bold text-lg ${
                    transaction.type === 'load_money' || transaction.type === 'refund' 
                      ? 'text-emerald-600' 
                      : 'text-red-600'
                  }`}>
                    {transaction.type === 'load_money' || transaction.type === 'refund' ? '+' : '-'}
                    {formatCurrency(transaction.amount)}
                  </p>
                  <TransactionBadge 
                    type={transaction.type} 
                    status={transaction.status} 
                    showIcon={false} 
                    size="sm"
                  />
                </div>
              </div>
            ))}
            
            <div className="pt-2">
              <Link href="/transactions" className="block">
                <Button variant="outline" className="w-full">
                  View All Transactions
                </Button>
              </Link>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
