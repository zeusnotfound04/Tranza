"use client"
import { useState, useEffect, useContext, createContext, ReactNode } from 'react';
import { apiClient, APIResponse } from '../services/api';
import { User } from '@/types/api';
import Cookies from 'js-cookie';

interface RegisterData {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  loading: boolean;
  isLoading: boolean; // Add this for compatibility
  isAuthenticated: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<APIResponse<any>>;
  register: (userData: RegisterData) => Promise<APIResponse<any>>;
  logout: () => void;
  refreshUser: () => Promise<void>;
  getToken: () => Promise<string | null>;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const isAuthenticated = !!user && !!token;

  // Initialize auth state from cookies
  useEffect(() => {
    const initAuth = async () => {
      try {
        const savedToken = Cookies.get('tranza_token');
        if (savedToken) {
          setToken(savedToken);
          apiClient.setAuthToken(savedToken);
          
          // Verify token and get user data
          const response = await apiClient.verifyToken();
          if (response.success && response.data?.valid && response.data?.user) {
            setUser(response.data.user);
          } else {
            // Try to get user profile directly
            const profileResponse = await apiClient.getUserProfile();
            if (profileResponse.success && profileResponse.data) {
              setUser(profileResponse.data);
            } else {
              // Token is invalid, clear it
              Cookies.remove('tranza_token');
              setToken(null);
              apiClient.setAuthToken('');
            }
          }
        }
      } catch (error) {
        console.error('Auth initialization error:', error);
        Cookies.remove('tranza_token');
        setToken(null);
        apiClient.setAuthToken('');
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  const login = async (email: string, password: string): Promise<APIResponse<any>> => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.login(email, password);
      
      if (response.success && response.data) {
        const { token: newToken, user: userData } = response.data;
        
        // Save token and user data
        setToken(newToken);
        setUser(userData);
        apiClient.setAuthToken(newToken);
        
        // Save token to cookies (expires in 30 days)
        Cookies.set('tranza_token', newToken, { expires: 30, secure: true, sameSite: 'strict' });
      } else {
        setError(response.error || 'Login failed');
      }
      
      return response;
    } catch (error) {
      console.error('Login error:', error);
      const errorMessage = 'Login failed';
      setError(errorMessage);
      return {
        success: false,
        error: errorMessage,
      };
    } finally {
      setLoading(false);
    }
  };

  const register = async (userData: RegisterData): Promise<APIResponse<any>> => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.register(userData);
      
      if (response.success && response.data) {
        const { token: newToken, user: newUser } = response.data;
        
        // Save token and user data
        setToken(newToken);
        setUser(newUser);
        apiClient.setAuthToken(newToken);
        
        // Save token to cookies
        Cookies.set('tranza_token', newToken, { expires: 30, secure: true, sameSite: 'strict' });
      } else {
        setError(response.error || 'Registration failed');
      }
      
      return response;
    } catch (error) {
      console.error('Register error:', error);
      const errorMessage = 'Registration failed';
      setError(errorMessage);
      return {
        success: false,
        error: errorMessage,
      };
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    setError(null);
    Cookies.remove('tranza_token');
    apiClient.setAuthToken('');
  };

  const clearError = () => {
    setError(null);
  };

  const refreshUser = async () => {
    if (!token) return;
    
    try {
      const response = await apiClient.getUserProfile();
      if (response.success && response.data) {
        setUser(response.data);
      }
    } catch (error) {
      console.error('Refresh user error:', error);
    }
  };

  const getToken = async (): Promise<string | null> => {
    try {
      // First, try to get token from current state
      if (token) {
        return token;
      }

      // If no token in state, try to get from API route
      const response = await fetch('/api/token', {
        method: 'GET',
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        if (data.success && data.data?.token) {
          // Update the auth state with the retrieved token
          const retrievedToken = data.data.token;
          console.log('Retrieved token from /api/token:', retrievedToken);
          setToken(retrievedToken);
          apiClient.setAuthToken(retrievedToken);
          
          // Try to get user data with this token
          const userResponse = await apiClient.getUserProfile();
          if (userResponse.success && userResponse.data) {
            setUser(userResponse.data);
          }
          
          return retrievedToken;
        }
      }

      // If API route fails, try to get from cookies directly
      const cookieToken = Cookies.get('tranza_token');
      if (cookieToken) {
        setToken(cookieToken);
        apiClient.setAuthToken(cookieToken);
        return cookieToken;
      }

      return null;
    } catch (error) {
      console.error('Get token error:', error);
      return null;
    }
  };

  const value: AuthContextType = {
    user,
    token,
    loading,
    isLoading: loading, // Add this alias for compatibility
    isAuthenticated,
    error,
    login,
    register,
    logout,
    refreshUser,
    getToken,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Hook for wallet-specific operations
export const useWallet = () => {
  const { isAuthenticated } = useAuth();
  const [balance, setBalance] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBalance = async () => {
    if (!isAuthenticated) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.getWalletBalance();
      if (response.success && response.data) {
        setBalance(response.data.balance);
      } else {
        setError(response.error || 'Failed to fetch balance');
      }
    } catch (err) {
      setError('Network error');
    } finally {
      setLoading(false);
    }
  };

  const loadMoney = async (amount: number, method: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.loadMoney(amount, method);
      if (response.success) {
        await fetchBalance(); // Refresh balance
      } else {
        setError(response.error || 'Failed to load money');
      }
      return response;
    } catch (err) {
      setError('Network error');
      return { success: false, error: 'Network error' };
    } finally {
      setLoading(false);
    }
  };

  return {
    balance,
    loading,
    error,
    fetchBalance,
    loadMoney,
  };
};

// Higher-order component for authentication
export const withAuth = <P extends object>(Component: React.ComponentType<P>) => {
  const AuthenticatedComponent = (props: P) => {
    const { user, loading } = useAuth();
    
    if (loading) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <h2 className="mt-4 text-lg font-medium text-gray-900">Loading...</h2>
          </div>
        </div>
      );
    }
    
    if (!user) {
      if (typeof window !== 'undefined') {
        window.location.href = '/auth/login';
      }
      return null;
    }
    
    return <Component {...props} />;
  };
  
  AuthenticatedComponent.displayName = `withAuth(${Component.displayName || Component.name})`;
  return AuthenticatedComponent;
};

export default withAuth;
