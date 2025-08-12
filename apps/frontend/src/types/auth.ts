// Authentication types matching Go backend models
export interface User {
  id: number;
  email: string;
  username: string;
  avatar?: string;
  provider: string;
  provider_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Legacy AuthResponse for localStorage (deprecated)
export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  expires_in: number;
}

// New response type for HttpOnly cookie authentication
export interface CookieAuthResponse {
  user: User;
  expires_in: number;
  message?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

// Email verification types
export interface PreRegistrationRequest {
  email: string;
  username: string;
  password: string;
}

export interface EmailVerificationRequest {
  email: string;
  code: string;
}

export interface ResendVerificationRequest {
  email: string;
}

export interface PreRegistrationResponse {
  message: string;
  email: string;
  expires_at: string;
}

export interface EmailVerificationResponse {
  message: string;
  user: User;
}

export interface AuthState {
  user: User | null;
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;
}

export interface AuthContextType extends AuthState {
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshAuth: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}