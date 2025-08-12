'use client';

import { useState, useEffect } from 'react';
import { APIKeyService } from '@/lib/services';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import Link from 'next/link';
import { aeonikPro } from '@/lib/fonts';
import { Key, Bot, BarChart3 } from 'lucide-react';

interface AIAgentStats {
  totalKeys: number;
  activeKeys: number;
  totalTransactions: number;
  todayTransactions: number;
  totalSpent: number;
  todaySpent: number;
}

export default function AIAgentStatus() {
  const [stats, setStats] = useState<AIAgentStats>({
    totalKeys: 0,
    activeKeys: 0,
    totalTransactions: 0,
    todayTransactions: 0,
    totalSpent: 0,
    todaySpent: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadAIAgentStats();
  }, []);

  const loadAIAgentStats = async () => {
    try {
      setLoading(true);
      setError('');
      
      // For now, we'll show placeholder data since we don't have all the endpoints
      // In a real implementation, you'd call multiple API endpoints to gather this data
      
      // Simulated data - replace with actual API calls when available
      setStats({
        totalKeys: 3,
        activeKeys: 2,
        totalTransactions: 47,
        todayTransactions: 5,
        totalSpent: 2456.78,
        todaySpent: 150.00,
      });
      
    } catch (err: any) {
      setError(err.message || 'Failed to load AI agent statistics');
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
      <Card>
        <CardHeader>
          <CardTitle>AI Agent Status</CardTitle>
          <CardDescription>Automated transaction overview</CardDescription>
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
          <CardTitle>AI Agent Status</CardTitle>
          <CardDescription>Automated transaction overview</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
          <Button onClick={loadAIAgentStats} className="mt-4" variant="outline" size="sm">
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
            <CardTitle>AI Agent Status</CardTitle>
            <CardDescription>Automated transaction overview</CardDescription>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span className="text-xs text-green-600 font-medium">Active</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* API Keys Overview */}
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-3 bg-blue-50 rounded-lg">
              <div className="text-2xl font-bold text-blue-600">{stats.totalKeys}</div>
              <div className="text-xs text-blue-700">Total Keys</div>
            </div>
            <div className="text-center p-3 bg-green-50 rounded-lg">
              <div className="text-2xl font-bold text-green-600">{stats.activeKeys}</div>
              <div className="text-xs text-green-700">Active Keys</div>
            </div>
          </div>

          {/* Transaction Stats */}
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Today's Transactions:</span>
              <Badge variant="outline">{stats.todayTransactions}</Badge>
            </div>
            
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Total Transactions:</span>
              <Badge variant="secondary">{stats.totalTransactions}</Badge>
            </div>
          </div>

          {/* Spending Overview */}
          <div className="p-3 bg-purple-50 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-purple-800">Today's Spending</span>
              <span className="text-sm font-bold text-purple-600">
                {formatCurrency(stats.todaySpent)}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-xs text-purple-700">Total Spent</span>
              <span className="text-xs font-medium text-purple-600">
                {formatCurrency(stats.totalSpent)}
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
            
            <Link href="/ai-agents">
              <Button variant="outline" size="sm" className="w-full">
                <Bot className="h-4 w-4 mr-2" />
                View AI Agents
              </Button>
            </Link>
            
            <Link href="/transactions?filter=ai_agent">
              <Button variant="outline" size="sm" className="w-full">
                <BarChart3 className="h-4 w-4 mr-2" />
                AI Transactions
              </Button>
            </Link>
          </div>

          {/* Status Indicators */}
          <div className="pt-2 border-t">
            <div className="flex items-center justify-between text-xs text-gray-500">
              <span>Last updated:</span>
              <span>{new Date().toLocaleTimeString()}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
