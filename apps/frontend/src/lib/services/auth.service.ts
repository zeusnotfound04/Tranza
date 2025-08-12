import { apiClient, TokenManager } from '@/lib/api-client';
import {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  PreRegistrationRequest,
  EmailVerificationRequest,
  ResendVerificationRequest,
  OAuthCallbackRequest,
  OAuthURL,
  User,
  APIResponse,
} from '@/types/api';

export class AuthService {
  // Email Verification Flow
  static async preRegister(data: PreRegistrationRequest): Promise<APIResponse<any>> {
    return apiClient.post('/auth/pre-register', data);
  }

  static async verifyEmail(data: EmailVerificationRequest): Promise<APIResponse<AuthResponse>> {
    return apiClient.post('/auth/verify-email', data);
  }

  static async resendVerification(data: ResendVerificationRequest): Promise<APIResponse<any>> {
    return apiClient.post('/auth/resend-verification', data);
  }

  // Standard Authentication
  static async login(data: LoginRequest): Promise<APIResponse<AuthResponse>> {
    const response = await apiClient.post<APIResponse<AuthResponse>>('/auth/login', data);
    
    // Store tokens if login successful
    if (response.data?.access_token) {
      TokenManager.setTokens(response.data.access_token, response.data.refresh_token);
    }
    
    return response;
  }

  static async logout(): Promise<APIResponse<any>> {
    const response = await apiClient.post<APIResponse<any>>('/auth/logout');
    
    // Clear tokens on logout
    TokenManager.clearTokens();
    
    return response;
  }

  static async refreshToken(): Promise<APIResponse<AuthResponse>> {
    return apiClient.post('/auth/refresh');
  }

  static async validateToken(): Promise<APIResponse<User>> {
    return apiClient.get('/auth/validate');
  }

  static async getCurrentUser(): Promise<APIResponse<User>> {
    return apiClient.get('/auth/me');
  }

  // OAuth Authentication
  static async getOAuthURL(provider: string, state?: string): Promise<OAuthURL> {
    const stateParam = state || this.generateState();
    const params = new URLSearchParams({ state: stateParam });
    
    return apiClient.get(`/auth/oauth/${provider}?${params.toString()}`);
  }

  static async handleOAuthCallback(data: OAuthCallbackRequest): Promise<APIResponse<AuthResponse>> {
    return apiClient.post('/auth/oauth/callback', data);
  }

  // Utility methods for OAuth flow
  private static generateState(): string {
    return Math.random().toString(36).substring(7) + Date.now().toString();
  }

  static storeOAuthState(state: string): void {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('oauth_state', state);
    }
  }

  static validateOAuthState(state: string): boolean {
    if (typeof window !== 'undefined') {
      const storedState = sessionStorage.getItem('oauth_state');
      sessionStorage.removeItem('oauth_state'); // Clean up
      return storedState === state;
    }
    return false;
  }

  // Legacy endpoints (for backward compatibility)
  static async register(data: RegisterRequest): Promise<APIResponse<AuthResponse>> {
    return apiClient.post('/auth/register', data);
  }
}

export default AuthService;
