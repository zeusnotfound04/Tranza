'use client';

import withAuth, { useAuth } from '@/hooks/useAuth';
import { DashboardLayout } from '@/components/layout/Navigation';
import WalletOverview from '@/components/dashboard/WalletOverview';
import RecentTransactions from '@/components/dashboard/RecentTransactions';
import QuickActions from '@/components/dashboard/QuickActions';
import SlackBotStatus from '@/components/dashboard/SlackBotStatus';
import APIKeyManagement from '@/components/APIKeyManagement';
import { Button } from '@/components/ui/Button';
import { aeonikPro } from '@/lib/fonts';
import { useState } from 'react';
import { Key, CreditCard, BarChart3, Bot } from 'lucide-react';
import { constants } from 'buffer';
// import DynamicScrollIslandTocDemo from '@/components/ui/dynamic-scroll-island-toc/demo';

function Dashboard() {
  const [activeTab, setActiveTab] = useState('overview');
  const { getToken, user } = useAuth();
  console.log("User Details:", user);
  const tabs = [
    { id: 'overview', name: 'Overview', icon: BarChart3 },
    { id: 'wallet', name: 'Wallet', icon: CreditCard },
    { id: 'api-keys', name: 'API Keys', icon: Key },
    { id: 'bot', name: 'Bot Setup', icon: Bot },
  ];
  console.log("Active Tab:", getToken);
  return (
    <DashboardLayout>
      <div className={`space-y-6 bg-[#121212] p-6 ${aeonikPro.className}`}>
        {/* Page Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">Dashboard</h1>
            <p className="mt-2 text-gray-400">
              Welcome back! Here's an overview of your financial activity.
            </p>
          </div>
          <div className="text-sm text-gray-400">
            Last updated: {new Date().toLocaleString()}
          </div>
        </div>

        {/* Enhanced Tab Navigation */}
        <div className="relative">
          {/* Tab Navigation Container */}
          <div className="relative bg-[#1f1f1f] backdrop-blur-sm rounded-xl border border-gray-700/50 p-2">
            <div className="flex space-x-1 relative">
              {/* Animated Background Indicator */}
              <div 
                className="absolute top-0 left-0 h-full bg-slate-300  rounded-lg transition-all duration-300 ease-out shadow-lg"
                style={{
                  width: `${100 / tabs.length}%`,
                  transform: `translateX(${tabs.findIndex(tab => tab.id === activeTab) * 100}%)`,
                }}
              />
              
              {/* Tab Buttons */}
              {tabs.map((tab, index) => {
                const Icon = tab.icon;
                const isActive = activeTab === tab.id;
                
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`
                      relative z-10 flex-1 flex items-center justify-center py-3 px-4 
                      font-medium text-sm rounded-lg transition-all duration-300 ease-out
                      group hover:scale-105 active:scale-95
                      ${isActive 
                        ? 'text-black font-semibold shadow-sm' 
                        : 'text-gray-400 hover:text-white'
                      }
                    `}
                  >
                    <Icon
                      className={`
                        mr-2 h-5 w-5 transition-all duration-300 ease-out
                        ${isActive 
                          ? 'text-black transform rotate-3' 
                          : 'text-gray-500 group-hover:text-white group-hover:scale-110'
                        }
                      `}
                    />
                    <span className="relative overflow-hidden">
                      <span 
                        className={`
                          block transition-transform duration-300 ease-out
                          ${isActive ? 'transform translate-y-0' : 'transform translate-y-0 group-hover:-translate-y-1'}
                        `}
                      >
                        {tab.name}
                      </span>
                    </span>
                    
                    {/* Hover Effect */}
                    <div 
                      className={`
                        absolute inset-0 rounded-lg transition-all duration-300 ease-out
                        ${!isActive ? 'bg-white/0 group-hover:bg-white/5' : ''}
                      `}
                    />
                  </button>
                );
              })}
            </div>
          </div>
          
          {/* Subtle glow effect */}
          <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent rounded-xl blur-xl opacity-50 -z-10" />
        </div>

        {/* Tab Content with Fade Animation */}
        <div className="relative">
          {activeTab === 'overview' && (
            <div className="animate-in fade-in duration-500 grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Left Column - Main Content */}
              <div className="lg:col-span-2 space-y-6">
                {/* Wallet Overview */}
                <div className="transform transition-all duration-300 hover:scale-[1.02]">
                  <WalletOverview />
                </div>
                
                {/* Recent Transactions */}
                <div className="transform transition-all duration-300 hover:scale-[1.02]">
                  <RecentTransactions />
                </div>
              </div>

              {/* Right Column - Sidebar Content */}
              <div className="space-y-6">
                {/* Slack Bot Status */}
                <div className="transform transition-all duration-300 hover:scale-[1.02]">
                  <SlackBotStatus />
                </div>
                
                {/* Quick Actions */}
                {/* <div className="transform transition-all duration-300 hover:scale-[1.02]">
                  <QuickActions />
                </div> */}
              </div>
            </div>
          )}

          {activeTab === 'wallet' && (
            <div className="animate-in fade-in duration-500 grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-3">
                <div className="transform transition-all duration-300 hover:scale-[1.01]">
                  <WalletOverview />
                </div>
                <div className="mt-6 transform transition-all duration-300 hover:scale-[1.01]">
                  <RecentTransactions />
                </div>
              </div>
            </div>
          )}

          {activeTab === 'api-keys' && (
            <div className="animate-in fade-in duration-500">
              <div className="transform transition-all duration-300 hover:scale-[1.005]">
                <div className="bg-[#1f1f1f] border border-gray-700 rounded-xl p-6">
                  <div className="flex items-center justify-between mb-6">
                    <h2 className="text-2xl font-bold text-white">API Keys Management</h2>
                    <Button
                      onClick={() => window.location.href = '/dashboard/api-keys'}
                      className="bg-blue-600 hover:bg-blue-700 text-white"
                    >
                      <Key className="w-4 h-4 mr-2" />
                      Manage API Keys
                    </Button>
                  </div>
                  <p className="text-gray-400 mb-4">
                    Create and manage API keys for programmatic access to Tranza services. Monitor usage, set limits, and track performance.
                  </p>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="bg-[#2a2a2a] border border-gray-600 rounded-lg p-4">
                      <h3 className="font-semibold text-white mb-2">🔑 Create Keys</h3>
                      <p className="text-sm text-gray-400">Generate secure API keys with custom labels and expiration times</p>
                    </div>
                    <div className="bg-[#2a2a2a] border border-gray-600 rounded-lg p-4">
                      <h3 className="font-semibold text-white mb-2">📊 Monitor Usage</h3>
                      <p className="text-sm text-gray-400">Track API requests, spending, and performance metrics in real-time</p>
                    </div>
                    <div className="bg-[#2a2a2a] border border-gray-600 rounded-lg p-4">
                      <h3 className="font-semibold text-white mb-2">🔒 Security</h3>
                      <p className="text-sm text-gray-400">Set spending limits, rotate keys, and monitor for suspicious activity</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'bot' && (
            <div className="animate-in fade-in duration-500">
              <div className="rounded-xl border border-gray-800 p-8 transform transition-all duration-300 hover:scale-[1.02] hover:shadow-2xl" style={{ backgroundColor: '#1f1f1f' }}>
                <h2 className="text-3xl font-bold text-white mb-6 flex items-center">
                  <Bot className="mr-3 h-8 w-8 text-gray-300" />
                  Slack Bot Setup
                </h2>
                
                <div className="space-y-6">
                  {/* Setup Instructions */}
                  <div className="bg-gradient-to-r from-gray-800/50 to-gray-900/50 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
                    <h3 className="font-semibold text-white mb-4 flex items-center text-lg">
                      🤖 Setup Instructions:
                    </h3>
                    <ol className="list-decimal list-inside space-y-3 text-gray-300">
                      <li className="hover:text-white transition-colors duration-200">
                        Go to the <strong className="text-white">API Keys</strong> tab and create a new API key
                      </li>
                      <li className="hover:text-white transition-colors duration-200">
                        Copy the generated API key
                      </li>
                      <li className="hover:text-white transition-colors duration-200">
                        In Slack, use the command: 
                        <code className="bg-black/50 text-green-400 px-3 py-1 rounded-md ml-2 font-mono">
                          /auth your-api-key
                        </code>
                      </li>
                      <li className="hover:text-white transition-colors duration-200">
                        Start using bot commands like 
                        <code className="bg-black/50 text-green-400 px-3 py-1 rounded-md ml-2 font-mono">
                          /fetch-balance
                        </code>
                      </li>
                    </ol>
                  </div>
                  
                  {/* Available Commands */}
                  <div className="bg-gradient-to-r from-gray-800/50 to-gray-900/50 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
                    <h3 className="font-semibold text-white mb-4 flex items-center text-lg">
                      📱 Available Commands:
                    </h3>
                    <ul className="space-y-3 text-gray-300">
                      {[
                        { cmd: '/auth <api-key>', desc: 'Authenticate with your API key' },
                        { cmd: '/fetch-balance', desc: 'Get your wallet balance' },
                        { cmd: '/send-money <amount> <upi|phone> <recipient>', desc: 'Send money' },
                        { cmd: '/logout', desc: 'Clear your session' },
                        { cmd: '/help', desc: 'Show help' },
                      ].map((command, index) => (
                        <li key={index} className="flex items-start space-x-3 hover:text-white transition-colors duration-200 group">
                          <code className="bg-black/50 text-blue-400 px-3 py-1 rounded-md font-mono text-sm group-hover:bg-black/70 transition-colors duration-200 flex-shrink-0">
                            {command.cmd}
                          </code>
                          <span className="text-gray-400 group-hover:text-gray-300">
                            {command.desc}
                          </span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}

export default Dashboard