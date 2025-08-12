'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import Link from 'next/link';
import { aeonikPro } from '@/lib/fonts';
import {
  Wallet,
  Send,
  BarChart3,
  CreditCard,
  Key,
  Bot,
  TrendingUp,
  Settings
} from 'lucide-react';

const quickActions = [
  {
    title: 'Load Money',
    description: 'Add funds to your wallet',
    icon: Wallet,
    href: '/wallet/load',
    primary: true,
  },
  {
    title: 'Send Money',
    description: 'Transfer funds to others',
    icon: Send,
    href: '/transactions/send',
  },
  {
    title: 'Transaction History',
    description: 'View all transactions',
    icon: BarChart3,
    href: '/transactions',
  },
  {
    title: 'Manage Cards',
    description: 'Link and manage cards',
    icon: CreditCard,
    href: '/cards',
  },
  {
    title: 'API Keys',
    description: 'Generate API keys for AI agents',
    icon: Key,
    href: '/api-keys',
  },
  {
    title: 'AI Agents',
    description: 'Monitor automated transactions',
    icon: Bot,
    href: '/ai-agents',
  },
  {
    title: 'Analytics',
    description: 'View spending insights',
    icon: TrendingUp,
    href: '/analytics',
  },
  {
    title: 'Settings',
    description: 'Account and wallet settings',
    icon: Settings,
    href: '/settings',
  },
];

export default function QuickActions() {
  return (
    <Card className={aeonikPro.className}>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
        <CardDescription>Frequently used features and tools</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {quickActions.map((action) => {
            const IconComponent = action.icon;
            return (
              <Link key={action.href} href={action.href}>
                <div className={`
                  p-4 rounded-lg border hover:shadow-md transition-all cursor-pointer h-full
                  ${action.primary 
                    ? 'bg-blue-50 border-blue-200 hover:bg-blue-100' 
                    : 'bg-white border-gray-200 hover:bg-gray-50'
                  }
                `}>
                  <div className="flex flex-col items-center text-center space-y-2">
                    <IconComponent className={`h-8 w-8 ${
                      action.primary ? 'text-blue-600' : 'text-gray-600'
                    }`} />
                    <div>
                      <h3 className={`font-semibold text-sm ${
                        action.primary ? 'text-blue-900' : 'text-gray-900'
                      }`}>
                        {action.title}
                      </h3>
                      <p className={`text-xs ${
                        action.primary ? 'text-blue-700' : 'text-gray-600'
                      }`}>
                        {action.description}
                      </p>
                    </div>
                  </div>
                </div>
              </Link>
            );
          })}
        </div>

        <div className="mt-6 p-4 bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg border border-blue-200">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-semibold text-blue-900">Need Help?</h3>
              <p className="text-sm text-blue-700">
                Check our documentation and tutorials
              </p>
            </div>
            <Link href="/help">
              <Button variant="outline" size="sm" className="border-blue-300 text-blue-700 hover:bg-blue-100">
                Get Help
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
