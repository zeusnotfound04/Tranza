'use client';

import { useState, useEffect } from 'react';
import { APIKeyService } from '@/lib/services';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@/components/ui/Button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import Link from 'next/link';
import { aeonikPro } from '@/lib/fonts';
import { Key, MessageSquare, BarChart3, Users, Zap, Activity } from 'lucide-react';

interface SlackBotStats {
  totalAPIKeys: number;
  activeBotKeys: number;
  universalKeys: number;
  totalWorkspaces: number;
  activeWorkspaces: number;
  totalBotTransactions: number;
  todayBotTransactions: number;
  totalSpent: number;
  todaySpent: number;
  lastActivity: string;
  botUptime: string;
}

export default function SlackBotStatus() {
  const [stats, setStats] = useState<SlackBotStats>({
    totalAPIKeys: 0,
    activeBotKeys: 0,
    universalKeys: 0,
    totalWorkspaces: 0,
    activeWorkspaces: 0,
    totalBotTransactions: 0,
    todayBotTransactions: 0,
    totalSpent: 0,
    todaySpent: 0,
    lastActivity: '',
    botUptime: '',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadSlackBotStats();
  }, []);

  const loadSlackBotStats = async () => {
    try {
      setLoading(true);
      setError('');
      
      // For now, we'll show realistic demo data
      // In a real implementation, you'd call multiple API endpoints to gather this data
      
      // Simulated Slack bot data
      setStats({
        totalAPIKeys: 5,
        activeBotKeys: 3,
        universalKeys: 2,
        totalWorkspaces: 4,
        activeWorkspaces: 3,
        totalBotTransactions: 156,
        todayBotTransactions: 12,
        totalSpent: 3240.50,
        todaySpent: 289.75,
        lastActivity: new Date(Date.now() - 300000).toISOString(), // 5 minutes ago
        botUptime: '99.8%',
      });
      
    } catch (err: any) {
      console.error('Failed to load Slack bot statistics:', err);
      setError(err.message || 'Failed to load Slack bot statistics');
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

  if (loading) {
    return (
      <Card className="border-gray-800" style={{ backgroundColor: '#1f1f1f' }}>
        <CardHeader>
          <CardTitle className="text-white">Slack Bot Status</CardTitle>
          <CardDescription className="text-gray-400">Slack workspace integrations overview</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-400"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="border-gray-800" style={{ backgroundColor: '#1f1f1f' }}>
        <CardHeader>
          <CardTitle className="text-white">Slack Bot Status</CardTitle>
          <CardDescription className="text-gray-400">Slack workspace integrations overview</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert variant="destructive" className="bg-red-900/20 border-red-800">
            <AlertDescription className="text-red-200">{error}</AlertDescription>
          </Alert>
          <Button onClick={loadSlackBotStats} className="mt-4" variant="outline" size="sm">
            Retry
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={`${aeonikPro.className} border-gray-800`} style={{ backgroundColor: '#1f1f1f' }}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-white">Slack Bot Status</CardTitle>
            <CardDescription className="text-gray-400">Slack workspace integrations overview</CardDescription>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span className="text-xs text-green-400 font-medium">Connected</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* API Keys Overview */}
          <div className="grid grid-cols-3 gap-3">
            <div className="text-center p-3 bg-blue-900/30 border border-blue-800/50 rounded-lg">
              <div className="text-xl font-bold text-blue-400">{stats.totalAPIKeys}</div>
              <div className="text-xs text-blue-300">Total Keys</div>
            </div>
            <div className="text-center p-3 bg-green-900/30 border border-green-800/50 rounded-lg">
              <div className="text-xl font-bold text-green-400">{stats.activeBotKeys}</div>
              <div className="text-xs text-green-300">Bot Keys</div>
            </div>
            <div className="text-center p-3 bg-purple-900/30 border border-purple-800/50 rounded-lg">
              <div className="text-xl font-bold text-purple-400">{stats.universalKeys}</div>
              <div className="text-xs text-purple-300">Universal</div>
            </div>
          </div>

          {/* Workspace Stats */}
          <div className="grid grid-cols-2 gap-3">
            <div className="text-center p-3 bg-orange-900/30 border border-orange-800/50 rounded-lg">
              <div className="text-xl font-bold text-orange-400">{stats.totalWorkspaces}</div>
              <div className="text-xs text-orange-300">Workspaces</div>
            </div>
            <div className="text-center p-3 bg-teal-900/30 border border-teal-800/50 rounded-lg">
              <div className="text-xl font-bold text-teal-400">{stats.activeWorkspaces}</div>
              <div className="text-xs text-teal-300">Active</div>
            </div>
          </div>

          {/* Transaction Stats */}
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">Today's Bot Transactions:</span>
              <Badge variant="outline" className="text-white border-gray-600">{stats.todayBotTransactions}</Badge>
            </div>
            
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">Total Bot Transactions:</span>
              <Badge variant="secondary" className="bg-gray-700 text-gray-200">{stats.totalBotTransactions}</Badge>
            </div>

            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">Bot Uptime:</span>
              <Badge variant="outline" className="text-green-400 border-green-600">{stats.botUptime}</Badge>
            </div>
          </div>

          {/* Spending Overview */}
          <div className="p-3 bg-gray-800/50 border border-gray-700/50 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-indigo-300">Today's Bot Spending</span>
              <span className="text-sm font-bold text-indigo-400">
                {formatCurrency(stats.todaySpent)}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-xs text-indigo-400">Total Spent via Bots</span>
              <span className="text-xs font-medium text-indigo-300">
                {formatCurrency(stats.totalSpent)}
              </span>
            </div>
          </div>

          {/* Activity Status */}
          <div className="p-3 bg-gray-800/50 border border-gray-700/50 rounded-lg">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">Last Bot Activity:</span>
              <span className="text-xs text-gray-300">
                {stats.lastActivity ? new Date(stats.lastActivity).toLocaleTimeString() : 'Never'}
              </span>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="grid grid-cols-1 gap-2 pt-2">
            <Link href="/api-keys">
              <Button variant="outline" size="sm" className="w-full">
                <Key className="h-4 w-4 mr-2" />
                Manage API Keys
              </Button>
            </Link>
            
            <Link href="/slack-workspaces">
              <Button variant="outline" size="sm" className="w-full">
                <MessageSquare className="h-4 w-4 mr-2" />
                Slack Workspaces
              </Button>
            </Link>
            
            <Link href="/transactions?filter=slack_bot">
              <Button variant="outline" size="sm" className="w-full">
                <BarChart3 className="h-4 w-4 mr-2" />
                Bot Transactions
              </Button>
            </Link>

            <Link href="/bot-analytics">
              <Button variant="outline" size="sm" className="w-full">
                <Activity className="h-4 w-4 mr-2" />
                Bot Analytics
              </Button>
            </Link>
          </div>

          {/* Status Indicators */}
          <div className="pt-2 border-t border-gray-700">
            <div className="flex items-center justify-between text-xs text-gray-400 mb-1">
              <span>Bot Status:</span>
              <span className="text-green-400 font-medium">Online</span>
            </div>
            <div className="flex items-center justify-between text-xs text-gray-400">
              <span>Last updated:</span>
              <span>{new Date().toLocaleTimeString()}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
