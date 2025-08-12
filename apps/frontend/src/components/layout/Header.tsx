'use client';

import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { Button } from '@tranza/ui/components/ui/button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Wallet } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import Logo from '@/components/ui/Logo';
import { UserAvatar } from '@/components/ui/UserAvatar';

export default function Header() {
  const { user, logout, isLoading } = useAuth();
  const [showUserMenu, setShowUserMenu] = useState(false);
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await logout();
      router.push('/');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Logo href={user ? '/dashboard' : '/'} />
          </div>

          {/* Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            {user ? (
              <>
                <Link href="/dashboard" className="text-gray-700 hover:text-blue-600 font-medium">
                  Dashboard
                </Link>
                <Link href="/wallet" className="text-gray-700 hover:text-blue-600 font-medium">
                  Wallet
                </Link>
                <Link href="/transactions" className="text-gray-700 hover:text-blue-600 font-medium">
                  Transactions
                </Link>
                <Link href="/ai-agents" className="text-gray-700 hover:text-blue-600 font-medium">
                  AI Agents
                </Link>
                <Link href="/api-keys" className="text-gray-700 hover:text-blue-600 font-medium">
                  API Keys
                </Link>
              </>
            ) : (
              <>
                <Link href="/features" className="text-gray-700 hover:text-blue-600 font-medium">
                  Features
                </Link>
                <Link href="/pricing" className="text-gray-700 hover:text-blue-600 font-medium">
                  Pricing
                </Link>
                <Link href="/docs" className="text-gray-700 hover:text-blue-600 font-medium">
                  Docs
                </Link>
              </>
            )}
          </nav>

         
          {/* <div className="flex items-center space-x-4">
            {isLoading ? (
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            ) : user ? (
              <>
                <div className="hidden lg:flex items-center space-x-2">
                  <Link href="/wallet/load">
                    <Button size="sm" variant="outline" className="text-xs">
                      <Wallet className="w-3 h-3 mr-1" />
                      Load Money
                    </Button>
                  </Link>
                </div>

                <div className="relative">
                  <button
                    onClick={() => setShowUserMenu(!showUserMenu)}
                    className="flex items-center space-x-2 p-2 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <UserAvatar user={user} size="md" />
                    <div className="hidden md:block text-left">
                      <div className="text-sm font-medium text-gray-900">{user.username}</div>
                      <div className="flex items-center space-x-1">
                        <Badge variant="outline" className="text-xs">
                          {user.provider}
                        </Badge>
                      </div>
                    </div>
                    <svg 
                      className="w-4 h-4 text-gray-500" 
                      fill="none" 
                      stroke="currentColor" 
                      viewBox="0 0 24 24"
                    >
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  </button>

                  {showUserMenu && (
                    <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border z-50">
                      <div className="py-1">
                        <div className="px-4 py-2 border-b">
                          <div className="text-sm font-medium text-gray-900">{user.username}</div>
                          <div className="text-xs text-gray-500">{user.email}</div>
                        </div>
                        
                        <Link 
                          href="/profile" 
                          className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                          onClick={() => setShowUserMenu(false)}
                        >
                          Profile Settings
                        </Link>
                        
                        <Link 
                          href="/wallet/settings" 
                          className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                          onClick={() => setShowUserMenu(false)}
                        >
                          Wallet Settings
                        </Link>
                        
                        <Link 
                          href="/security" 
                          className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                          onClick={() => setShowUserMenu(false)}
                        >
                          Security
                        </Link>
                        
                        <Link 
                          href="/help" 
                          className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                          onClick={() => setShowUserMenu(false)}
                        >
                          Help & Support
                        </Link>
                        
                        <div className="border-t">
                          <button
                            onClick={() => {
                              setShowUserMenu(false);
                              handleLogout();
                            }}
                            className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                          >
                            Sign Out
                          </button>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </>
            ) : (
              <div className="flex items-center space-x-3">
                <Link href="/auth/login">
                  <Button variant="outline" size="sm">
                    Sign In
                  </Button>
                </Link>
                <Link href="/auth/register">
                  <Button size="sm">
                    Get Started
                  </Button>
                </Link>
              </div>
            )}
          </div> */}
          {/*  */}
        </div>
      </div>

      {/* Close dropdown when clicking outside */}
      {showUserMenu && (
        <div 
          className="fixed inset-0 z-40"
          onClick={() => setShowUserMenu(false)}
        ></div>
      )}
    </header>
  );
}
