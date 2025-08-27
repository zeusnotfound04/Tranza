'use client';

import { useState } from 'react';
import Header from './Header';
import Sidebar from './Sidebar';
import Logo from '@/components/ui/Logo';

interface NavigationProps {
  children: React.ReactNode;
  showSidebar?: boolean;
}

export default function Navigation({ children, showSidebar = true }: NavigationProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="min-h-screen bg-[#121212]">
      
      
      <div className="flex">
        {/* Sidebar */}
        {showSidebar && (
          <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
        )}

        {/* Main Content */}
        <main className={`flex-1 ${showSidebar ? 'lg:ml-64' : ''}`}>
          {/* Mobile menu button */}
          {showSidebar && (
            <div className="lg:hidden">
              <div className="flex items-center justify-between h-12 px-4 bg-[#121212] border-b border-gray-200">
                <button
                  type="button"
                  className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
                  onClick={() => setSidebarOpen(true)}
                >
                  <span className="sr-only">Open main menu</span>
                  <svg
                    className="h-6 w-6"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                  </svg>
                </button>
                <Logo href="/dashboard" size="sm" />
                <div className="w-10"></div> {/* Spacer for centering */}
              </div>
            </div>
          )}

          {/* Content */}
          <div className="p-4  sm:p-6 lg:p-8">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}

// Layout wrapper for pages that need sidebar
export function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <Navigation  showSidebar={true}>
      {children}
    </Navigation>
  );
}

// Layout wrapper for pages that don't need sidebar
export function SimpleLayout({ children }: { children: React.ReactNode }) {
  return (
    <Navigation showSidebar={false}>
      {children}
    </Navigation>
  );
}
