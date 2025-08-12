'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';

export default function LogoutPage() {
  const router = useRouter();
  const { logout, user, isLoading } = useAuth();
  const [loggingOut, setLoggingOut] = useState(false);

  useEffect(() => {
    const performLogout = async () => {
      if (!user && !isLoading) {
        // Already logged out, redirect to home
        router.push('/');
        return;
      }

      if (user && !loggingOut) {
        setLoggingOut(true);
        try {
          await logout();
          // Redirect will be handled by the auth context
        } catch (error) {
          console.error('Logout failed:', error);
          // Still redirect even if logout fails
          router.push('/');
        }
      }
    };

    performLogout();
  }, [user, isLoading, logout, router, loggingOut]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <h2 className="mt-4 text-lg font-medium text-gray-900">
          Signing you out...
        </h2>
        <p className="mt-2 text-sm text-gray-600">
          Please wait while we securely log you out
        </p>
      </div>
    </div>
  );
}
