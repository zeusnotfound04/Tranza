'use client';

import React, { useState, useEffect, useContext, createContext, ReactNode } from 'react';
import { AuthService, APIError } from '@/lib/services';
import { User, LoginRequest, PreRegistrationRequest, EmailVerificationRequest } from '@/types/api';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  preRegister: (data: PreRegistrationRequest) => Promise<{ message: string }>;
  verifyEmail: (data: EmailVerificationRequest) => Promise<void>;
  resendVerification: (email: string) => Promise<{ message: string }>;
  refreshToken: () => Promise<void>;
  redirectToLogin: (returnUrl?: string) => void;
  error: string | null;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const isAuthenticated = !!user;

  const clearError = () => setError(null);

  // Check if user is authenticated on mount
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      setIsLoading(true);
      const response = await AuthService.getCurrentUser();
      if (response.user) {
        setUser(response.user);
      }
    } catch (error) {
      if (error instanceof APIError && error.status === 401) {
        // User is not authenticated, which is fine
        setUser(null);
      } else {
        console.error('Auth check failed:', error);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (credentials: LoginRequest) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await AuthService.login(credentials);
      if (response.user) {
        setUser(response.user);
      } else {
        throw new Error('Login response missing user data');
      }
    } catch (error) {
      const message = error instanceof APIError ? error.message : 'Login failed';
      setError(message);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      setIsLoading(true);
      await AuthService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      setIsLoading(false);
    }
  };

  const preRegister = async (data: PreRegistrationRequest) => {
    try {
      setError(null);
      const response = await AuthService.preRegister(data);
      return { message: response.message || 'Verification code sent successfully' };
    } catch (error) {
      const message = error instanceof APIError ? error.message : 'Pre-registration failed';
      setError(message);
      throw error;
    }
  };

  const verifyEmail = async (data: EmailVerificationRequest) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await AuthService.verifyEmail(data);
      if (response.user) {
        setUser(response.user);
      } else {
        throw new Error('Email verification response missing user data');
      }
    } catch (error) {
      const message = error instanceof APIError ? error.message : 'Email verification failed';
      setError(message);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const resendVerification = async (email: string) => {
    try {
      setError(null);
      const response = await AuthService.resendVerification({ email });
      return { message: response.message || 'Verification code resent successfully' };
    } catch (error) {
      const message = error instanceof APIError ? error.message : 'Failed to resend verification code';
      setError(message);
      throw error;
    }
  };

  const refreshToken = async () => {
    try {
      const response = await AuthService.refreshToken();
      if (response.data?.user) {
        setUser(response.data.user);
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
      setUser(null);
      // Redirect to login if refresh fails
      if (typeof window !== 'undefined') {
        window.location.href = '/auth/login';
      }
      throw error;
    }
  };

  // Auto-refresh token when it's about to expire
  useEffect(() => {
    if (!user) return;

    // Set up token refresh interval (refresh every 50 minutes, token expires in 1 hour)
    const refreshInterval = setInterval(async () => {
      try {
        await refreshToken();
      } catch (error) {
        console.error('Auto refresh failed:', error);
        // Will redirect to login in refreshToken function
      }
    }, 50 * 60 * 1000); // 50 minutes

    return () => clearInterval(refreshInterval);
  }, [user]);

  // Handle auth redirects
  const redirectToLogin = (returnUrl?: string) => {
    if (typeof window !== 'undefined') {
      const params = returnUrl ? `?returnUrl=${encodeURIComponent(returnUrl)}` : '';
      window.location.href = `/auth/login${params}`;
    }
  };

  const handleLogout = async () => {
    try {
      setIsLoading(true);
      await AuthService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      setIsLoading(false);
      // Clear any stored OAuth state
      if (typeof window !== 'undefined') {
        sessionStorage.removeItem('oauth_state');
        sessionStorage.removeItem('oauth_provider');
      }
      redirectToLogin();
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout: logout,
    preRegister,
    verifyEmail,
    resendVerification,
    refreshToken,
    redirectToLogin,
    error,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

// Higher-order component for protecting routes
export function withAuth<P extends object>(Component: React.ComponentType<P>) {
  const AuthenticatedComponent = (props: P) => {
    const { isAuthenticated, isLoading, redirectToLogin } = useAuth();

    useEffect(() => {
      if (!isLoading && !isAuthenticated) {
        // Get current path for return URL
        const currentPath = typeof window !== 'undefined' ? window.location.pathname + window.location.search : '';
        redirectToLogin(currentPath);
      }
    }, [isAuthenticated, isLoading, redirectToLogin]);

    if (isLoading) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
            <p className="mt-2 text-sm text-gray-600">Loading...</p>
          </div>
        </div>
      );
    }

    if (!isAuthenticated) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
            <p className="mt-2 text-sm text-gray-600">Redirecting to login...</p>
          </div>
        </div>
      );
    }

    return <Component {...props} />;
  };

  AuthenticatedComponent.displayName = `withAuth(${Component.displayName || Component.name})`;
  return AuthenticatedComponent;
}

export default withAuth;
