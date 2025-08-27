'use client';

import Link from 'next/link';
import { useAuth } from '@/hooks/useAuth';
import { useEffect, useState } from 'react';
import { WalletService, TransactionService, APIKeyService } from '@/lib/services';
import Navbar from '@/components/ui/Navbar';
import OAuthStatus from '@/components/auth/OAuthStatus';
import OAuthButtons from '@/components/auth/OAuthButtons';
import { 
  FiCreditCard, 
  FiDollarSign, 
  FiKey, 
  FiTrendingUp, 
  FiShield, 
  FiZap, 
  FiUsers, 
  FiArrowRight,
  FiCheck,
  FiSend,
  FiPlus,
  FiSettings,
  FiBarChart
} from 'react-icons/fi';

interface DashboardStats {
  walletBalance: number;
  transactionCount: number;
  apiKeyCount: number;
}

export default function Home() {
  const { user, isLoading } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    walletBalance: 0,
    transactionCount: 0,
    apiKeyCount: 0,
  });
  const [statsLoading, setStatsLoading] = useState(false);

  // Load dashboard stats when user is authenticated
  useEffect(() => {
    if (user && !isLoading) {
      loadDashboardStats();
    }
  }, [user, isLoading]);

  const loadDashboardStats = async () => {
    try {
      setStatsLoading(true);
      
      // Load wallet balance
      try {
        const walletResponse = await WalletService.getWallet();
        if (walletResponse.data) {
          setStats(prev => ({ ...prev, walletBalance: walletResponse.data!.balance }));
        }
      } catch (error) {
        console.log('Wallet not found - user may need to create one');
        setStats(prev => ({ ...prev, walletBalance: 0 }));
      }

      // Load transaction count
      try {
        const transactionResponse = await TransactionService.getTransactionHistory({ limit: 1, offset: 0 });
        setStats(prev => ({ ...prev, transactionCount: transactionResponse.total || 0 }));
      } catch (error) {
        console.log('No transactions found');
        setStats(prev => ({ ...prev, transactionCount: 0 }));
      }

      // API keys count - placeholder since backend doesn't have get endpoint yet
      setStats(prev => ({ ...prev, apiKeyCount: 0 }));
      
    } catch (error) {
      console.error('Failed to load dashboard stats:', error);
    } finally {
      setStatsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-white via-gray-50 to-blue-50 dark:from-gray-900 dark:via-gray-800 dark:to-blue-900">
      {/* Background decoration */}
      <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8ZGVmcz4KICAgIDxwYXR0ZXJuIGlkPSJncmlkIiB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHBhdHRlcm5Vbml0cz0idXNlclNwYWNlT25Vc2UiPgogICAgICA8cGF0aCBkPSJNIDQwIDAgTCAwIDAgMCA0MCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJyZ2JhKDAsIDAsIDAsIDAuMDUpIiBzdHJva2Utd2lkdGg9IjEiLz4KICAgIDwvcGF0dGVybj4KICA8L2RlZnM+CiAgPHJlY3Qgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIgZmlsbD0idXJsKCNncmlkKSIvPgo8L3N2Zz4=')] opacity-20"></div>
      
      <Navbar />
      
      <main className="relative max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Hero Section */}
          <div className="text-center mb-16">
            <div className="mb-8">
              <div className="inline-flex items-center justify-center w-20 h-20 bg-white rounded-3xl mb-6 shadow-2xl border border-gray-200">
                <img 
                  src="/logo.png" 
                  alt="Tranza Logo" 
                  className="w-16 h-16 object-contain"
                />
              </div>
            </div>
            <h1 className="text-6xl font-bold bg-gradient-to-r from-gray-900 via-gray-700 to-blue-600 dark:from-gray-100 dark:via-gray-300 dark:to-blue-400 bg-clip-text text-transparent mb-6">
              Welcome to Tranza
            </h1>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-12 max-w-2xl mx-auto leading-relaxed">
              Experience the future of financial transactions with our secure, 
              AI-powered platform designed for modern digital payments.
            </p>

            {isLoading ? (
              <div className="flex justify-center">
                <div className="relative">
                  <div className="animate-spin rounded-full h-12 w-12 border-4 border-blue-500 border-t-transparent"></div>
                  <div className="absolute inset-0 rounded-full bg-gradient-to-r from-blue-500 to-blue-600 opacity-20 animate-pulse"></div>
                </div>
              </div>
            ) : user ? (
              <div className="max-w-4xl mx-auto">
                {/* User Welcome Card */}
                <div className="bg-white/80 dark:bg-gray-800/80 backdrop-blur-lg border border-gray-200 dark:border-gray-700 rounded-3xl p-8 mb-12 shadow-2xl">
                  <div className="text-center">
                    <div className="mb-6">
                      {user.avatar ? (
                        <img
                          src={user.avatar}
                          alt="Profile"
                          className="h-20 w-20 rounded-full mx-auto ring-4 ring-blue-500/50 shadow-lg"
                        />
                      ) : (
                        <div className="h-20 w-20 rounded-full bg-gradient-to-r from-blue-600 to-blue-500 flex items-center justify-center mx-auto ring-4 ring-blue-500/50 shadow-lg">
                          <span className="text-white text-2xl font-bold">
                            {user.username.charAt(0).toUpperCase()}
                          </span>
                        </div>
                      )}
                    </div>
                    <h2 className="text-3xl font-bold text-gray-800 dark:text-gray-100 mb-3">
                      Welcome back, {user.username}!
                    </h2>
                    <p className="text-gray-600 dark:text-gray-300 mb-6 text-lg">{user.email}</p>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-600 dark:text-gray-400 mb-8">
                      <div className="flex items-center justify-center space-x-2">
                        <FiUsers className="text-blue-400" />
                        <span>Member since {new Date(user.created_at).toLocaleDateString()}</span>
                      </div>
                      <div className="flex items-center justify-center space-x-2">
                        <FiShield className="text-green-400" />
                        <span>Provider: {user.provider}</span>
                      </div>
                      <div className="flex items-center justify-center space-x-2">
                        <FiCheck className="text-green-400" />
                        <span>Status: {user.is_active ? 'Active' : 'Inactive'}</span>
                      </div>
                    </div>
                    <Link
                      href="/dashboard"
                      className="inline-flex items-center px-8 py-4 bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-700 hover:to-blue-600 text-white font-semibold rounded-2xl transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl"
                    >
                      <span>Access Dashboard</span>
                      <FiArrowRight className="ml-2" />
                    </Link>
                  </div>
                </div>
              </div>
            ) : (
              <div className="max-w-lg mx-auto">
                {/* Welcome Card */}
                <div className="bg-white/80 backdrop-blur-lg border border-gray-200 rounded-3xl p-8 mb-8 shadow-2xl">
                  <h2 className="text-2xl font-bold text-gray-800 mb-4 text-center">
                    Start Your Financial Journey
                  </h2>
                  <p className="text-gray-600 mb-8 text-center leading-relaxed">
                    Join thousands of users who trust Tranza for secure, 
                    fast, and intelligent financial transactions.
                  </p>
                  
                  {/* Features List */}
                  <div className="space-y-4 mb-8">
                    <div className="flex items-center text-gray-700">
                      <div className="w-8 h-8 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                        <FiShield className="text-white text-sm" />
                      </div>
                      <span>Enterprise-grade security with OAuth 2.0</span>
                    </div>
                    <div className="flex items-center text-gray-700">
                      <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-full flex items-center justify-center mr-3">
                        <FiZap className="text-white text-sm" />
                      </div>
                      <span>Instant wallet creation and management</span>
                    </div>
                    <div className="flex items-center text-gray-700">
                      <div className="w-8 h-8 bg-gradient-to-r from-gray-600 to-gray-500 rounded-full flex items-center justify-center mr-3">
                        <FiTrendingUp className="text-white text-sm" />
                      </div>
                      <span>Real-time transaction tracking</span>
                    </div>
                    <div className="flex items-center text-gray-700">
                      <div className="w-8 h-8 bg-gradient-to-r from-slate-600 to-slate-500 rounded-full flex items-center justify-center mr-3">
                        <FiKey className="text-white text-sm" />
                      </div>
                      <span>Advanced API management for developers</span>
                    </div>
                  </div>

                  {/* OAuth Buttons */}
                  <div className="mb-8">
                    <OAuthButtons mode="register" />
                  </div>

            

                 
                </div>
              </div>
            )}
          </div>

          {user && (
            <div className="space-y-12">
              {/* Stats Cards */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white/80 backdrop-blur-lg border border-gray-200 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-xl flex items-center justify-center shadow-lg">
                        <FiCreditCard className="text-white text-xl" />
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-600 truncate">
                          Total Transactions
                        </dt>
                        <dd className="text-2xl font-bold text-gray-800">
                          {statsLoading ? (
                            <div className="animate-pulse bg-gray-200 h-8 w-16 rounded-lg"></div>
                          ) : (
                            <Link href="/transactions" className="hover:text-cyan-600 transition-colors">
                              {stats.transactionCount.toLocaleString()}
                            </Link>
                          )}
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>

                <div className="bg-white/80 backdrop-blur-lg border border-gray-200 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-500 rounded-xl flex items-center justify-center shadow-lg">
                        <FiDollarSign className="text-white text-xl" />
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-600 truncate">
                          Wallet Balance
                        </dt>
                        <dd className="text-2xl font-bold text-gray-800">
                          {statsLoading ? (
                            <div className="animate-pulse bg-gray-200 h-8 w-20 rounded-lg"></div>
                          ) : (
                            <Link href="/wallet" className="hover:text-emerald-600 transition-colors">
                              ₹{stats.walletBalance.toLocaleString()}
                            </Link>
                          )}
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>

                <div className="bg-white/80 backdrop-blur-lg border border-gray-200 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-12 h-12 bg-gradient-to-r from-gray-600 to-gray-500 rounded-xl flex items-center justify-center shadow-lg">
                        <FiKey className="text-white text-xl" />
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-600 truncate">
                          API Keys
                        </dt>
                        <dd className="text-2xl font-bold text-gray-800">
                          {statsLoading ? (
                            <div className="animate-pulse bg-gray-200 h-8 w-16 rounded-lg"></div>
                          ) : (
                            <Link href="/api-keys" className="hover:text-gray-600 transition-colors">
                              {stats.apiKeyCount}
                            </Link>
                          )}
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>

              {/* Quick Actions */}
              <div>
                <h3 className="text-2xl font-bold text-white mb-8 text-center">Quick Actions</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                  <Link
                    href="/wallet/load"
                    className="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center text-center">
                      <div className="w-16 h-16 bg-gradient-to-r from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg group-hover:shadow-green-500/25 transition-all duration-300">
                        <FiPlus className="text-white text-2xl" />
                      </div>
                      <h4 className="text-lg font-semibold text-white mb-2">Load Money</h4>
                      <p className="text-sm text-gray-300">Add funds to your wallet instantly</p>
                    </div>
                  </Link>

                  <Link
                    href="/transactions/new"
                    className="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center text-center">
                      <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg group-hover:shadow-blue-500/25 transition-all duration-300">
                        <FiSend className="text-white text-2xl" />
                      </div>
                      <h4 className="text-lg font-semibold text-white mb-2">Send Money</h4>
                      <p className="text-sm text-gray-300">Transfer funds securely</p>
                    </div>
                  </Link>

                  <Link
                    href="/cards"
                    className="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center text-center">
                      <div className="w-16 h-16 bg-gradient-to-r from-gray-600 to-gray-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg group-hover:shadow-gray-500/25 transition-all duration-300">
                        <FiCreditCard className="text-white text-2xl" />
                      </div>
                      <h4 className="text-lg font-semibold text-white mb-2">Manage Cards</h4>
                      <p className="text-sm text-gray-300">Link and manage your cards</p>
                    </div>
                  </Link>

                  <Link
                    href="/analytics"
                    className="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center text-center">
                      <div className="w-16 h-16 bg-gradient-to-r from-slate-600 to-slate-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg group-hover:shadow-slate-500/25 transition-all duration-300">
                        <FiBarChart className="text-white text-2xl" />
                      </div>
                      <h4 className="text-lg font-semibold text-white mb-2">Analytics</h4>
                      <p className="text-sm text-gray-300">View detailed insights</p>
                    </div>
                  </Link>
                </div>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
