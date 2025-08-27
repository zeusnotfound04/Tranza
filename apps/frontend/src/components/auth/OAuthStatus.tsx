'use client';

import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/ui/Button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Card } from '@tranza/ui/components/ui/card-ui';
import Link from 'next/link';

interface OAuthStatusProps {
  showDetails?: boolean;
}

export default function OAuthStatus({ showDetails = false }: OAuthStatusProps) {
  const { user, isLoading, isAuthenticated, logout } = useAuth();

  if (isLoading) {
    return (
      <div className="flex items-center space-x-2">
        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
        <span className="text-sm text-gray-600">Checking authentication...</span>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="flex items-center space-x-2">
        <Badge variant="secondary">Not Authenticated</Badge>
        {showDetails && (
          <div className="flex space-x-2">
            <Link href="/auth/login">
              <Button size="sm" variant="outline">Sign In</Button>
            </Link>
            <Link href="/auth/register">
              <Button size="sm">Sign Up</Button>
            </Link>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center space-x-2">
        <Badge variant="default" className="bg-green-100 text-green-800">
          Authenticated via {user?.provider || 'Unknown'}
        </Badge>
        {user?.is_active && (
          <Badge variant="default" className="bg-blue-100 text-blue-800">
            Active
          </Badge>
        )}
      </div>

      {showDetails && user && (
        <Card className="p-4 bg-gray-50">
          <div className="space-y-2">
            <div className="flex items-center space-x-2">
              {user.avatar && (
                <img
                  src={user.avatar}
                  alt="Profile"
                  className="w-8 h-8 rounded-full"
                />
              )}
              <div>
                <p className="font-semibold text-sm">{user.username}</p>
                <p className="text-xs text-gray-600">{user.email}</p>
              </div>
            </div>

            <div className="text-xs text-gray-500 space-y-1">
              <p><strong>Provider:</strong> {user.provider}</p>
              <p><strong>User ID:</strong> {user.id}</p>
              <p><strong>Created:</strong> {new Date(user.created_at).toLocaleDateString()}</p>
              <p><strong>Status:</strong> {user.is_active ? 'Active' : 'Inactive'}</p>
              {user.provider_id && (
                <p><strong>Provider ID:</strong> {user.provider_id}</p>
              )}
            </div>

            <div className="flex space-x-2 pt-2">
              <Link href="/dashboard">
                <Button size="sm" variant="outline">Dashboard</Button>
              </Link>
              <Button size="sm" variant="destructive" onClick={logout}>
                Sign Out
              </Button>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
}
