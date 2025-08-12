'use client';

import { useState, useEffect, useContext, createContext, ReactNode } from 'react';
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
      if (response.user) {
        setUser(response.user);
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
      setUser(null);
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout,
    preRegister,
    verifyEmail,
    resendVerification,
    refreshToken,
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

export default useAuth;
