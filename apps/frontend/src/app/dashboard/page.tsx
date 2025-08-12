'use client';

import withAuth from '@/hooks/useAuth';
import { DashboardLayout } from '@/components/layout/Navigation';
import WalletOverview from '@/components/dashboard/WalletOverview';
import RecentTransactions from '@/components/dashboard/RecentTransactions';
import QuickActions from '@/components/dashboard/QuickActions';
import AIAgentStatus from '@/components/dashboard/AIAgentStatus';
import { aeonikPro } from '@/lib/fonts';
import DynamicScrollIslandTocDemo from '@/components/ui/dynamic-scroll-island-toc/demo';

function Dashboard() {
  return (
    <DashboardLayout >
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

        {/* Main Dashboard Grid */}
        <div className="grid  grid-cols-1 lg:grid-cols-3 gap-6">
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
      </div>
    </DashboardLayout>
  );
}

export default withAuth(Dashboard);