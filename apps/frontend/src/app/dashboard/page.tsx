'use client';

import withAuth from '@/hooks/useAuth';
import { DashboardLayout } from '@/components/layout/Navigation';
import WalletOverview from '@/components/dashboard/WalletOverview';
import RecentTransactions from '@/components/dashboard/RecentTransactions';
import QuickActions from '@/components/dashboard/QuickActions';
import AIAgentStatus from '@/components/dashboard/AIAgentStatus';
import APIKeyManagement from '@/components/APIKeyManagement';
import { aeonikPro } from '@/lib/fonts';
import { useState } from 'react';
import { Key, CreditCard, BarChart3, Bot } from 'lucide-react';
// import DynamicScrollIslandTocDemo from '@/components/ui/dynamic-scroll-island-toc/demo';

function Dashboard() {
  const [activeTab, setActiveTab] = useState('overview');

  const tabs = [
    { id: 'overview', name: 'Overview', icon: BarChart3 },
    { id: 'wallet', name: 'Wallet', icon: CreditCard },
    { id: 'api-keys', name: 'API Keys', icon: Key },
    { id: 'bot', name: 'Bot Setup', icon: Bot },
  ];

  return (
    <DashboardLayout>
      <div className={`space-y-6 p-6 ${aeonikPro.className}`}>
        {/* Page Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-black">Dashboard</h1>
            <p className="mt-2 text-black">
              Welcome back! Here's an overview of your financial activity.
            </p>
          </div>
          <div className="text-sm text-black">
            Last updated: {new Date().toLocaleString()}
          </div>
        </div>

        {/* Tab Navigation */}
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`group inline-flex items-center py-2 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <Icon
                    className={`mr-2 h-5 w-5 ${
                      activeTab === tab.id ? 'text-blue-500' : 'text-gray-400 group-hover:text-gray-500'
                    }`}
                  />
                  {tab.name}
                </button>
              );
            })}
          </nav>
        </div>

        {/* Tab Content */}
        {activeTab === 'overview' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left Column - Main Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Wallet Overview */}
              <WalletOverview />
              
              {/* Recent Transactions */}
              <RecentTransactions />
            </div>

            {/* Right Column - Sidebar Content */}
            <div className="space-y-6">
              {/* AI Agent Status */}
              <AIAgentStatus />
              
              {/* Quick Actions */}
              <QuickActions />
            </div>
          </div>
        )}

        {activeTab === 'wallet' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-3">
              <WalletOverview />
              <div className="mt-6">
                <RecentTransactions />
              </div>
            </div>
          </div>
        )}

        {activeTab === 'api-keys' && (
          <div>
            <APIKeyManagement />
          </div>
        )}

        {activeTab === 'bot' && (
          <div className="bg-white rounded-lg border p-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Slack Bot Setup</h2>
            <div className="space-y-4">
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="font-semibold text-blue-800 mb-2">🤖 Setup Instructions:</h3>
                <ol className="list-decimal list-inside space-y-2 text-blue-700">
                  <li>Go to the <strong>API Keys</strong> tab and create a new API key</li>
                  <li>Copy the generated API key</li>
                  <li>In Slack, use the command: <code className="bg-blue-100 px-2 py-1 rounded">/auth your-api-key</code></li>
                  <li>Start using bot commands like <code className="bg-blue-100 px-2 py-1 rounded">/fetch-balance</code></li>
                </ol>
              </div>
              
              <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
                <h3 className="font-semibold text-gray-800 mb-2">📱 Available Commands:</h3>
                <ul className="space-y-1 text-gray-700">
                  <li><code className="bg-gray-100 px-2 py-1 rounded">/auth &lt;api-key&gt;</code> - Authenticate with your API key</li>
                  <li><code className="bg-gray-100 px-2 py-1 rounded">/fetch-balance</code> - Get your wallet balance</li>
                  <li><code className="bg-gray-100 px-2 py-1 rounded">/send-money &lt;amount&gt; &lt;upi|phone&gt; &lt;recipient&gt;</code> - Send money</li>
                  <li><code className="bg-gray-100 px-2 py-1 rounded">/logout</code> - Clear your session</li>
                  <li><code className="bg-gray-100 px-2 py-1 rounded">/help</code> - Show help</li>
                </ul>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

export default Dashboard