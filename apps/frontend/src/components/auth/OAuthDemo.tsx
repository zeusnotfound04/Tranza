'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { Button } from '@tranza/ui/components/ui/button';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import OAuthButtons from '@/components/auth/OAuthButtons';
import OAuthStatus from '@/components/auth/OAuthStatus';

export default function OAuthDemo() {
  const { user, isLoading, logout, refreshToken, isAuthenticated } = useAuth();
  const [tokenRefreshStatus, setTokenRefreshStatus] = useState<string>('');
  const [sessionInfo, setSessionInfo] = useState<any>(null);

  useEffect(() => {
    // Capture session info from URL params (for demo purposes)
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search);
      if (urlParams.get('oauth_success')) {
        setSessionInfo({
          provider: sessionStorage.getItem('oauth_provider'),
          timestamp: new Date().toISOString(),
          success: true
        });
      }
    }
  }, []);

  const handleTokenRefresh = async () => {
    try {
      setTokenRefreshStatus('Refreshing...');
      await refreshToken();
      setTokenRefreshStatus('Token refreshed successfully!');
      setTimeout(() => setTokenRefreshStatus(''), 3000);
    } catch (error: any) {
      setTokenRefreshStatus(`Refresh failed: ${error.message}`);
      setTimeout(() => setTokenRefreshStatus(''), 5000);
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      setSessionInfo(null);
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          OAuth Authentication Demo
        </h1>
        <p className="text-gray-600 mb-4">
          Complete Phase 6 Implementation: OAuth Login Flow
        </p>
        <div className="flex justify-center mb-6">
          <OAuthStatus showDetails={false} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Authentication Panel */}
        <Card className="p-6">
          <h2 className="text-xl font-semibold mb-4">Authentication</h2>
          
          {!isAuthenticated ? (
            <div className="space-y-4">
              <p className="text-sm text-gray-600 mb-4">
                Sign in with Google or GitHub to test the OAuth flow:
              </p>
              <OAuthButtons mode="login" />
              
              <div className="text-xs text-gray-500 mt-4 p-3 bg-blue-50 rounded">
                <p><strong>OAuth Flow:</strong></p>
                <ol className="list-decimal list-inside mt-1 space-y-1">
                  <li>Click OAuth provider button</li>
                  <li>Redirected to provider (Google/GitHub)</li>
                  <li>Authorize application</li>
                  <li>Redirected back with authorization code</li>
                  <li>Code exchanged for user data + JWT tokens</li>
                  <li>Wallet created automatically</li>
                  <li>User logged in with secure cookies</li>
                </ol>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <OAuthStatus showDetails={true} />
              
              <div className="flex space-x-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={handleTokenRefresh}
                  disabled={!!tokenRefreshStatus}
                >
                  Test Token Refresh
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={handleLogout}
                >
                  Logout
                </Button>
              </div>

              {tokenRefreshStatus && (
                <Alert>
                  <AlertDescription>{tokenRefreshStatus}</AlertDescription>
                </Alert>
              )}
            </div>
          )}
        </Card>

        {/* Session Information */}
        <Card className="p-6">
          <h2 className="text-xl font-semibold mb-4">Session Info</h2>
          
          <div className="space-y-4">
            <div className="text-sm">
              <p><strong>Authentication State:</strong></p>
              <div className="flex items-center space-x-2 mt-1">
                <Badge variant={isAuthenticated ? "default" : "secondary"}>
                  {isAuthenticated ? "Authenticated" : "Not Authenticated"}
                </Badge>
                {isLoading && (
                  <Badge variant="secondary">Loading...</Badge>
                )}
              </div>
            </div>

            {user && (
              <div className="text-sm">
                <p><strong>User Details:</strong></p>
                <div className="bg-gray-50 p-3 rounded mt-1 text-xs font-mono">
                  <pre>{JSON.stringify({
                    id: user.id,
                    username: user.username,
                    email: user.email,
                    provider: user.provider,
                    provider_id: user.provider_id,
                    is_active: user.is_active,
                    created_at: user.created_at
                  }, null, 2)}</pre>
                </div>
              </div>
            )}

            {sessionInfo && (
              <div className="text-sm">
                <p><strong>OAuth Session:</strong></p>
                <div className="bg-green-50 p-3 rounded mt-1 text-xs">
                  <p>✅ OAuth login successful!</p>
                  <p>Provider: {sessionInfo.provider}</p>
                  <p>Time: {new Date(sessionInfo.timestamp).toLocaleString()}</p>
                </div>
              </div>
            )}

            <div className="text-sm">
              <p><strong>Security Features:</strong></p>
              <ul className="text-xs text-gray-600 mt-1 space-y-1">
                <li>✅ CSRF protection with state parameter</li>
                <li>✅ HttpOnly cookies for token storage</li>
                <li>✅ Automatic token refresh</li>
                <li>✅ Secure session management</li>
                <li>✅ Automatic wallet creation</li>
                <li>✅ Return URL support</li>
              </ul>
            </div>
          </div>
        </Card>
      </div>

      {/* Phase 6 Requirements Checklist */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold mb-4">Phase 6 Implementation Status</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h3 className="font-medium text-green-800 mb-2">✅ Phase 6.1: OAuth Login Flow</h3>
            <ul className="text-sm space-y-1">
              <li>✅ OAuth provider integration (Google/GitHub)</li>
              <li>✅ Authorization code flow</li>
              <li>✅ Secure token exchange</li>
              <li>✅ Automatic wallet creation</li>
              <li>✅ JWT token management</li>
              <li>✅ Secure cookie storage</li>
            </ul>
          </div>
          <div>
            <h3 className="font-medium text-green-800 mb-2">✅ Phase 6.2: Authentication State Management</h3>
            <ul className="text-sm space-y-1">
              <li>✅ Auth context implementation</li>
              <li>✅ Token refresh mechanism</li>
              <li>✅ Auth redirect handling</li>
              <li>✅ Logout flow</li>
              <li>✅ Return URL support</li>
              <li>✅ Session persistence</li>
            </ul>
          </div>
        </div>
      </Card>
    </div>
  );
}
