'use client';

import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { Badge } from '@tranza/ui/components/ui/badge';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { aeonikPro } from '@/lib/fonts';
import { UserAvatar } from '@/components/ui/UserAvatar';
import Logo from '@/components/ui/Logo';
import {
  BarChart3,
  TrendingUp,
  Wallet,
  ArrowUp,
  CreditCard,
  FileText,
  Send,
  CreditCard as PaymentIcon,
  Bot,
  Key,
  Plug,
  User,
  Shield,
  Bell,
  ChevronRight,
  ShoppingBag
} from 'lucide-react';

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const sidebarItems = [
  {
    title: 'Overview',
    items: [
      { name: 'Dashboard', href: '/dashboard', icon: BarChart3 },
      { name: 'Analytics', href: '/analytics', icon: TrendingUp },
    ]
  },
  {
    title: 'Wallet',
    items: [
      { name: 'Balance', href: '/wallet', icon: Wallet },
      { name: 'Load Money', href: '/wallet/load', icon: ArrowUp },
      { name: 'Cards', href: '/cards', icon: CreditCard },
    ]
  },
  {
    title: 'Transactions',
    items: [
      { name: 'History', href: '/transactions', icon: FileText },
      { name: 'Send Money', href: '/transactions/send', icon: Send },
      { name: 'Payments', href: '/payments', icon: PaymentIcon },
    ]
  },
  {
    title: 'AI & Automation',
    items: [
      { name: 'AI Shopping', href: '/ai-shopping', icon: ShoppingBag },
      { name: 'AI Agents', href: '/ai-agents', icon: Bot },
      { name: 'API Keys', href: '/dashboard/api-keys', icon: Key },
      { name: 'Integrations', href: '/integrations', icon: Plug },
    ]
  },
  {
    title: 'Settings',
    items: [
      { name: 'Profile', href: '/profile', icon: User },
      { name: 'Security', href: '/security', icon: Shield },
      { name: 'Notifications', href: '/notifications', icon: Bell },
    ]
  }
];

export default function Sidebar({ isOpen, onClose }: SidebarProps) {
  const { user } = useAuth();
  const pathname = usePathname();
  const [expandedSections, setExpandedSections] = useState<string[]>(['Overview', 'Wallet']);

  const toggleSection = (title: string) => {
    setExpandedSections(prev =>
      prev.includes(title)
        ? prev.filter(t => t !== title)
        : [...prev, title]
    );
  };

  const isActiveLink = (href: string) => {
    return pathname === href || pathname.startsWith(href + '/');
  };

  if (!user) return null;

  return (
    <>
      {/* Mobile backdrop */}
      {isOpen && (
        <div 
          className="fixed inset-0 z-40 bg-black bg-opacity-50 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <div className={`
        fixed inset-y-0 left-0 z-50 w-72 bg-gradient-to-b from-slate-50 to-white shadow-2xl transform transition-all duration-500 ease-out ${aeonikPro.className}
        ${isOpen ? 'translate-x-0' : '-translate-x-full'}
        lg:translate-x-0 lg:rounded-r-3xl border-r border-slate-200/60 backdrop-blur-xl
      `}>
        
        {/* Logo Section */}
        <div className="px-6 py-2 border-b border-slate-200/60 bg-white/50 backdrop-blur-sm lg:rounded-tr-3xl">
          <div className="flex items-center justify-center">
            <Logo size="lg" />
          </div>
        </div>


        {/* User Info */}
        <div className="p-6 border-b border-slate-200/60 bg-gradient-to-r from-white to-slate-50/50">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <UserAvatar user={user} size="lg" />
              <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-emerald-400 border-2 border-white rounded-full animate-pulse"></div>
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-xl font-bold text-slate-900 truncate mb-1">
                {user.username}
              </div>
              <div className="flex items-center space-x-2">
                <Badge 
                  variant={user.provider === 'google' ? 'default' : user.provider === 'github' ? 'secondary' : 'outline'}
                  className={`text-xs px-3 py-1 flex items-center space-x-1.5 rounded-full transition-all duration-300 hover:scale-105 ${
                    user.provider === 'google' 
                      ? 'bg-gradient-to-r from-red-50 to-red-100 text-red-700 border-red-200 hover:from-red-100 hover:to-red-200 shadow-sm' 
                      : user.provider === 'github' 
                      ? 'bg-gradient-to-r from-slate-800 to-slate-900 text-white border-slate-800 hover:from-slate-700 hover:to-slate-800 shadow-sm'
                      : 'bg-gradient-to-r from-blue-50 to-blue-100 text-blue-700 border-blue-200 hover:from-blue-100 hover:to-blue-200 shadow-sm'
                  }`}
                >
                  {user.provider === 'google' && (
                    <svg className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                    </svg>
                  )}
                  {user.provider === 'github' && (
                    <svg className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                  )}
                  {user.provider === 'email' && (
                    <svg className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"/>
                    </svg>
                  )}
                  <span className="capitalize font-medium">{user.provider}</span>
                </Badge>
              </div>
            </div>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-6 py-6 space-y-3 overflow-y-auto scrollbar-thin scrollbar-thumb-slate-300 scrollbar-track-transparent">
          {sidebarItems.map((section) => (
            <div key={section.title} className="space-y-2">
              <button
                onClick={() => toggleSection(section.title)}
                className="flex items-center justify-between w-full px-4 py-3 text-sm font-semibold text-slate-700 rounded-xl hover:bg-slate-100/80 transition-all duration-300 hover:scale-[1.02] group"
              >
                <span className="text-slate-600 uppercase tracking-wider text-xs font-bold">{section.title}</span>
                <ChevronRight 
                  className={`w-4 h-4 transform transition-all duration-300 text-slate-400 group-hover:text-slate-600 ${
                    expandedSections.includes(section.title) ? 'rotate-90' : ''
                  }`}
                />
              </button>
              
              {expandedSections.includes(section.title) && (
                <div className="ml-2 space-y-1 animate-in slide-in-from-top-2 duration-300">
                  {section.items.map((item) => {
                    const IconComponent = item.icon;
                    const isActive = isActiveLink(item.href);
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={onClose}
                        className={`
                          flex items-center px-4 py-3 text-sm rounded-xl transition-all duration-300 group relative overflow-hidden
                          ${isActive
                            ? 'bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-lg shadow-blue-500/25 scale-[1.02]'
                            : 'text-slate-700 hover:bg-slate-100/80 hover:scale-[1.02] hover:translate-x-1'
                          }
                        `}
                      >
                        {isActive && (
                          <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-700 opacity-10 animate-pulse"></div>
                        )}
                        <IconComponent className={`mr-3 h-5 w-5 transition-all duration-300 ${
                          isActive ? 'text-white' : 'text-slate-500 group-hover:text-slate-700'
                        }`} />
                        <span className="font-medium relative z-10">{item.name}</span>
                        {isActive && (
                          <div className="ml-auto w-2 h-2 bg-white rounded-full animate-pulse"></div>
                        )}
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          ))}
        </nav>

        {/* Footer */}
        <div className="p-6 border-t border-slate-200/60 bg-gradient-to-r from-white to-slate-50/30 lg:rounded-br-3xl">
          <div className="text-xs text-slate-500 space-y-2">
            <div className="flex items-center justify-between">
              <p className="font-semibold text-slate-600">Tranza</p>
              <span className="px-2 py-1 bg-emerald-100 text-emerald-700 rounded-full text-xs font-medium">v1.0.0</span>
            </div>
            <div className="flex items-center space-x-4">
              <Link 
                href="/help" 
                className="text-blue-600 hover:text-blue-700 transition-colors duration-200 flex items-center space-x-1 group"
              >
                <svg className="w-3 h-3 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span>Help</span>
              </Link>
              <Link 
                href="/support" 
                className="text-slate-500 hover:text-slate-700 transition-colors duration-200 flex items-center space-x-1 group"
              >
                <svg className="w-3 h-3 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
                <span>Support</span>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
