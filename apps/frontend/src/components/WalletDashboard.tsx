import { useState, useEffect } from 'react';
import { useAuth, useWallet } from '../hooks/useAuth';
import { formatCurrency } from '../services/api';
import { Wallet, Send, History, Plus, Eye, EyeOff } from 'lucide-react';

export default function WalletDashboard() {
  const { user } = useAuth();
  const { balance, loading, error, fetchBalance } = useWallet();
  const [showBalance, setShowBalance] = useState(true);

  useEffect(() => {
    fetchBalance();
  }, []);

  return (
    <div className="max-w-4xl mx-auto p-6">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Welcome back, {user?.first_name}!
        </h1>
        <p className="text-gray-600">Manage your wallet and send money securely</p>
      </div>

      {/* Wallet Balance Card */}
      <div className="bg-gradient-to-br from-blue-600 to-blue-800 text-white rounded-xl p-6 mb-6 shadow-lg">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <Wallet className="w-8 h-8" />
            <h2 className="text-xl font-semibold">Wallet Balance</h2>
          </div>
          <button
            onClick={() => setShowBalance(!showBalance)}
            className="p-2 hover:bg-white/20 rounded-lg transition-colors"
          >
            {showBalance ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
          </button>
        </div>
        
        <div className="text-3xl font-bold mb-2">
          {loading ? (
            <div className="animate-pulse bg-white/20 h-8 w-32 rounded"></div>
          ) : showBalance ? (
            formatCurrency(balance)
          ) : (
            '••••••'
          )}
        </div>
        
        {error && (
          <p className="text-red-200 text-sm">{error}</p>
        )}
        
        <button
          onClick={fetchBalance}
          disabled={loading}
          className="text-sm bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors disabled:opacity-50"
        >
          {loading ? 'Refreshing...' : 'Refresh Balance'}
        </button>
      </div>

      {/* Quick Actions */}
      <div className="grid md:grid-cols-3 gap-4 mb-8">
        <QuickActionCard
          icon={<Plus className="w-6 h-6" />}
          title="Add Money"
          description="Load money to your wallet"
          href="/wallet/load"
          color="green"
        />
        <QuickActionCard
          icon={<Send className="w-6 h-6" />}
          title="Send Money"
          description="Transfer to UPI or phone"
          href="/transfer"
          color="blue"
        />
        <QuickActionCard
          icon={<History className="w-6 h-6" />}
          title="Transaction History"
          description="View all transactions"
          href="/transactions"
          color="purple"
        />
      </div>

      {/* Recent Activity Preview */}
      <RecentActivityPreview />
    </div>
  );
}

function QuickActionCard({ 
  icon, 
  title, 
  description, 
  href, 
  color 
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
  href: string;
  color: 'green' | 'blue' | 'purple';
}) {
  const colorClasses = {
    green: 'from-green-500 to-green-600 hover:from-green-600 hover:to-green-700',
    blue: 'from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700',
    purple: 'from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700',
  };

  return (
    <a
      href={href}
      className={`block bg-gradient-to-br ${colorClasses[color]} text-white rounded-lg p-6 hover:shadow-lg transition-all duration-200 transform hover:scale-105`}
    >
      <div className="flex items-center gap-3 mb-3">
        {icon}
        <h3 className="font-semibold">{title}</h3>
      </div>
      <p className="text-sm opacity-90">{description}</p>
    </a>
  );
}

function RecentActivityPreview() {
  const [recentTransfers, setRecentTransfers] = useState([]);
  const [loading, setLoading] = useState(false);

  // This would typically fetch recent transactions
  // For now, showing placeholder

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Recent Activity</h3>
        <a
          href="/transactions"
          className="text-blue-600 hover:text-blue-700 text-sm font-medium"
        >
          View All
        </a>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="animate-pulse flex items-center gap-4">
              <div className="w-10 h-10 bg-gray-200 rounded-full"></div>
              <div className="flex-1 space-y-2">
                <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
              <div className="h-4 bg-gray-200 rounded w-20"></div>
            </div>
          ))}
        </div>
      ) : recentTransfers.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <History className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No recent transactions</p>
          <p className="text-sm">Your transaction history will appear here</p>
        </div>
      ) : (
        <div className="space-y-3">
          {/* Recent transfers would be mapped here */}
        </div>
      )}
    </div>
  );
}
