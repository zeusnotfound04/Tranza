'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';
import { AuthService } from '@/lib/services';
import Cookies from 'js-cookie';
import { apiClient } from '@/services/api';

export default function AuthCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, loading } = useAuth();
  const [processing, setProcessing] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const error = searchParams.get('error');
        const errorDescription = searchParams.get('error_description');
        
        // Get stored OAuth data
        const storedState = sessionStorage.getItem('oauth_state');
        const storedProvider = sessionStorage.getItem('oauth_provider') || 'google';
        
        // Clean up stored data
        sessionStorage.removeItem('oauth_state');
        sessionStorage.removeItem('oauth_provider');

        // Check for OAuth errors
        if (error) {
          setError(`Authentication failed: ${errorDescription || error}`);
          setProcessing(false);
          return;
        }

        // Check for required parameters
        if (!code) {
          setError('Authorization code not received');
          setProcessing(false);
          return;
        }

        // Validate state to prevent CSRF attacks
        if (state && storedState && state !== storedState) {
          setError('Invalid state parameter - possible security issue');
          setProcessing(false);
          return;
        }

        // Exchange code for tokens
        const authResponse = await AuthService.handleOAuthCallback({
          provider: storedProvider,
          code,
          state: state || undefined,
          redirect_uri: `${window.location.origin}/auth/callback`
        });

        // Check if the OAuth callback was successful
        if (authResponse.data) {
          const { access_token, user } = authResponse.data;
          
          if (access_token && user) {
            // Set the token using the API endpoint with GET method
            const tokenResponse = await fetch("/api/token", {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              },
            });
            
            console.log("Response from /api/token ", tokenResponse);
            
            if (tokenResponse.ok) {
              // Set the token in the API client
              apiClient.setAuthToken(access_token);
              
              // Force a page reload to ensure auth state is properly initialized
              window.location.href = '/dashboard';
              return;
            } else {
              setError('Failed to retrieve authentication token');
              setProcessing(false);
              return;
            }
          }
        }

        // Fallback: redirect to dashboard and let the auth context handle it
        router.push('/dashboard');
        
      } catch (err: any) {
        console.error('OAuth callback error:', err);
        setError(err.message || 'Authentication failed');
        setProcessing(false);
      }
    };

    // Only process if we're not already authenticated
    if (!loading && !user) {
      handleOAuthCallback();
    } else if (!loading && user) {
      // Already authenticated, redirect to dashboard
      router.push('/dashboard');
    }
  }, [searchParams, router, user, loading]);

  // Loading state
  if (processing || loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <h2 className="mt-4 text-lg font-medium text-gray-900">
            Completing authentication...
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Please wait while we set up your account
          </p>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6">
          <div className="flex items-center mb-4">
            <div className="flex-shrink-0">
              <svg className="h-8 w-8 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.728-.833-2.498 0L3.316 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-lg font-medium text-gray-900">
                Authentication Failed
              </h3>
            </div>
          </div>
          
          <p className="text-sm text-gray-600 mb-4">
            {error}
          </p>
          
          <div className="flex space-x-3">
            <button
              onClick={() => router.push('/auth/login')}
              className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              Try Again
            </button>
            <button
              onClick={() => router.push('/')}
              className="flex-1 bg-gray-200 text-gray-900 px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
            >
              Go Home
            </button>
          </div>
        </div>
      </div>
    );
  }

  return null;
}
