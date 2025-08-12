'use client';

import { useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';

interface AuthGuardProps {
  children: ReactNode;
  redirectTo?: string;
  requireAuth?: boolean;
  fallback?: ReactNode;
}

export default function AuthGuard({ 
  children, 
  redirectTo = '/auth/login',
  requireAuth = true,
  fallback = <div className="flex items-center justify-center min-h-screen">
    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
  </div>
}: AuthGuardProps) {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      if (requireAuth && !user) {
        router.push(redirectTo);
      } else if (!requireAuth && user) {
        router.push('/dashboard');
      }
    }
  }, [user, isLoading, requireAuth, redirectTo, router]);

  // Show loading state while checking authentication
  if (isLoading) {
    return <>{fallback}</>;
  }

  // If auth is required and user is not authenticated, don't render children
  if (requireAuth && !user) {
    return <>{fallback}</>;
  }

  // If auth is not required but user is authenticated (e.g., login page), redirect
  if (!requireAuth && user) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
