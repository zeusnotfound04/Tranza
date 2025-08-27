'use client';

import withAuth from '@/hooks/useAuth';
import { DashboardLayout } from '@/components/layout/Navigation';
import AIClothingAgent from '@/components/ai/AIClothingAgent';
import { aeonikPro } from '@/lib/fonts';

function AIShoppingPage() {
  return (
    <DashboardLayout>
      <div className={`h-full ${aeonikPro.className}`}>
        {/* Page Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200 bg-white">
          <div>
            <h1 className="text-3xl font-bold text-white">AI Shopping Assistant</h1>
            <p className="mt-2 text-gray-600">
              Use AI to find and purchase clothing with your Tranza wallet
            </p>
          </div>
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
            <span>AI Online</span>
          </div>
        </div>

        {/* AI Agent Container */}
        <div className="flex-1 h-[calc(100vh-180px)]">
          <AIClothingAgent />
        </div>
      </div>
    </DashboardLayout>
  );
}

export default withAuth(AIShoppingPage);
