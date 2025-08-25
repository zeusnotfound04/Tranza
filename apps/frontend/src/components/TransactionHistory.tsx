import { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { apiClient, formatCurrency, formatDateTime } from '../services/api';
import { 
  History, 
  ArrowUpRight, 
  ArrowDownLeft, 
  Clock, 
  CheckCircle, 
  XCircle,
  AlertCircle,
  Filter,
  Download,
  Search,
  RefreshCw
} from 'lucide-react';

interface Transaction {
  id: string;
  type: 'transfer_out' | 'transfer_in' | 'wallet_load';
  amount: string;
  fee?: string;
  total_amount: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  recipient?: string;
  sender?: string;
  description?: string;
  reference_id: string;
  created_at: string;
  updated_at: string;
}

interface FilterOptions {
  type: string;
  status: string;
  dateFrom: string;
  dateTo: string;
  search: string;
}

export default function TransactionHistory() {
  const { user } = useAuth();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [showFilters, setShowFilters] = useState(false);
  
  const [filters, setFilters] = useState<FilterOptions>({
    type: 'all',
    status: 'all',
    dateFrom: '',
    dateTo: '',
    search: '',
  });

  const limit = 10;

  useEffect(() => {
    fetchTransactions();
  }, [page, filters]);

  const fetchTransactions = async () => {
    setLoading(true);
    setError(null);

    try {
      // This would be replaced with actual API call
      // For now, showing mock data structure
      
      const response = await apiClient.getTransferHistory(page, limit);
      
      if (response.success && response.data) {
        // Map API response to local Transaction interface
        const mappedTransactions: Transaction[] = (response.data.transfers || []).map((transfer: any) => ({
          id: transfer.transfer_id || transfer.id,
          type: 'transfer_out' as const,
          amount: transfer.amount,
          fee: '0', // Default fee if not provided
          total_amount: transfer.amount,
          status: transfer.status as Transaction['status'],
          recipient: transfer.recipient,
          description: 'Transfer',
          reference_id: transfer.reference_id || transfer.transfer_id,
          created_at: transfer.created_at,
          updated_at: transfer.updated_at,
        }));
        
        setTransactions(mappedTransactions);
        setTotalPages(Math.ceil((response.data.total_count || 0) / limit));
      } else {
        // Mock data for demonstration
        const mockTransactions: Transaction[] = [
          {
            id: '1',
            type: 'transfer_out',
            amount: '500.00',
            fee: '5.00',
            total_amount: '505.00',
            status: 'completed',
            recipient: 'john@paytm',
            description: 'Payment for lunch',
            reference_id: 'TRX001',
            created_at: new Date(Date.now() - 86400000).toISOString(),
            updated_at: new Date(Date.now() - 86400000).toISOString(),
          },
          {
            id: '2',
            type: 'wallet_load',
            amount: '1000.00',
            total_amount: '1000.00',
            status: 'completed',
            description: 'Wallet top-up',
            reference_id: 'TRX002',
            created_at: new Date(Date.now() - 172800000).toISOString(),
            updated_at: new Date(Date.now() - 172800000).toISOString(),
          },
          {
            id: '3',
            type: 'transfer_out',
            amount: '250.00',
            fee: '2.50',
            total_amount: '252.50',
            status: 'pending',
            recipient: '9876543210',
            description: 'Payment to friend',
            reference_id: 'TRX003',
            created_at: new Date(Date.now() - 3600000).toISOString(),
            updated_at: new Date(Date.now() - 3600000).toISOString(),
          },
          {
            id: '4',
            type: 'transfer_out',
            amount: '100.00',
            fee: '1.00',
            total_amount: '101.00',
            status: 'failed',
            recipient: 'invalid@upi',
            description: 'Failed transfer',
            reference_id: 'TRX004',
            created_at: new Date(Date.now() - 7200000).toISOString(),
            updated_at: new Date(Date.now() - 7200000).toISOString(),
          },
        ];
        
        setTransactions(mockTransactions);
        setTotalPages(1);
      }
    } catch (err) {
      setError('Failed to fetch transaction history');
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: Transaction['status']) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case 'pending':
        return <Clock className="w-5 h-5 text-yellow-600" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-600" />;
      case 'cancelled':
        return <AlertCircle className="w-5 h-5 text-gray-600" />;
      default:
        return <Clock className="w-5 h-5 text-gray-600" />;
    }
  };

  const getTransactionIcon = (type: Transaction['type']) => {
    switch (type) {
      case 'transfer_out':
        return <ArrowUpRight className="w-5 h-5 text-red-600" />;
      case 'transfer_in':
        return <ArrowDownLeft className="w-5 h-5 text-green-600" />;
      case 'wallet_load':
        return <ArrowDownLeft className="w-5 h-5 text-blue-600" />;
      default:
        return <History className="w-5 h-5 text-gray-600" />;
    }
  };

  const getStatusColor = (status: Transaction['status']) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      case 'cancelled':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const exportTransactions = () => {
    const data = {
      exported_at: new Date().toISOString(),
      user_id: user?.id,
      transactions: transactions,
      filters: filters,
    };

    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `transactions-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <History className="w-8 h-8 text-blue-600" />
          Transaction History
        </h1>
        <p className="text-gray-600">View and manage your transaction history</p>
      </div>

      {/* Controls */}
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <Filter className="w-5 h-5" />
            Filters
          </button>

          <button
            onClick={fetchTransactions}
            disabled={loading}
            className="bg-blue-100 hover:bg-blue-200 text-blue-700 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </button>

          <button
            onClick={exportTransactions}
            className="bg-green-100 hover:bg-green-200 text-green-700 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <Download className="w-5 h-5" />
            Export
          </button>
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="mb-6 bg-white border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Filter Transactions</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Type</label>
              <select
                value={filters.type}
                onChange={(e) => setFilters({ ...filters, type: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Types</option>
                <option value="transfer_out">Transfers Out</option>
                <option value="transfer_in">Transfers In</option>
                <option value="wallet_load">Wallet Loads</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
              <select
                value={filters.status}
                onChange={(e) => setFilters({ ...filters, status: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Status</option>
                <option value="completed">Completed</option>
                <option value="pending">Pending</option>
                <option value="failed">Failed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">From Date</label>
              <input
                type="date"
                value={filters.dateFrom}
                onChange={(e) => setFilters({ ...filters, dateFrom: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">To Date</label>
              <input
                type="date"
                value={filters.dateTo}
                onChange={(e) => setFilters({ ...filters, dateTo: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Search</label>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <input
                  type="text"
                  value={filters.search}
                  onChange={(e) => setFilters({ ...filters, search: e.target.value })}
                  className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Search..."
                />
              </div>
            </div>
          </div>

          <div className="mt-4 flex gap-2">
            <button
              onClick={() => setFilters({
                type: 'all',
                status: 'all',
                dateFrom: '',
                dateTo: '',
                search: '',
              })}
              className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg text-sm transition-colors"
            >
              Clear Filters
            </button>
          </div>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-700">{error}</p>
        </div>
      )}

      {/* Transactions List */}
      <div className="bg-white border border-gray-200 rounded-lg">
        {loading ? (
          <div className="p-8 text-center">
            <RefreshCw className="w-8 h-8 text-blue-600 animate-spin mx-auto mb-4" />
            <p className="text-gray-600">Loading transactions...</p>
          </div>
        ) : transactions.length === 0 ? (
          <div className="p-8 text-center">
            <History className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No Transactions</h3>
            <p className="text-gray-600">You haven't made any transactions yet.</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {transactions.map((transaction) => (
              <TransactionCard key={transaction.id} transaction={transaction} />
            ))}
          </div>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-between">
          <div className="text-sm text-gray-600">
            Page {page} of {totalPages}
          </div>
          
          <div className="flex gap-2">
            <button
              onClick={() => setPage(page - 1)}
              disabled={page === 1}
              className="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:text-gray-400 text-gray-700 px-4 py-2 rounded-lg transition-colors"
            >
              Previous
            </button>
            <button
              onClick={() => setPage(page + 1)}
              disabled={page === totalPages}
              className="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:text-gray-400 text-gray-700 px-4 py-2 rounded-lg transition-colors"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );

  function TransactionCard({ transaction }: { transaction: Transaction }) {
    return (
      <div className="p-6 hover:bg-gray-50 transition-colors">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              {getTransactionIcon(transaction.type)}
              {getStatusIcon(transaction.status)}
            </div>
            
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-1">
                <h3 className="font-medium text-gray-900">
                  {transaction.type === 'transfer_out' && 'Money Sent'}
                  {transaction.type === 'transfer_in' && 'Money Received'}
                  {transaction.type === 'wallet_load' && 'Wallet Loaded'}
                </h3>
                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(transaction.status)}`}>
                  {transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1)}
                </span>
              </div>
              
              <div className="text-sm text-gray-600 space-y-1">
                {transaction.recipient && (
                  <div>To: {transaction.recipient}</div>
                )}
                {transaction.sender && (
                  <div>From: {transaction.sender}</div>
                )}
                {transaction.description && (
                  <div>{transaction.description}</div>
                )}
                <div>Reference: {transaction.reference_id}</div>
                <div>{formatDateTime(transaction.created_at)}</div>
              </div>
            </div>
          </div>

          <div className="text-right">
            <div className="font-semibold text-lg">
              {transaction.type === 'transfer_out' ? '-' : '+'}
              {formatCurrency(transaction.amount)}
            </div>
            {transaction.fee && parseFloat(transaction.fee) > 0 && (
              <div className="text-sm text-gray-600">
                Fee: {formatCurrency(transaction.fee)}
              </div>
            )}
            <div className="text-sm text-gray-600">
              Total: {formatCurrency(transaction.total_amount)}
            </div>
          </div>
        </div>
      </div>
    );
  }
}
