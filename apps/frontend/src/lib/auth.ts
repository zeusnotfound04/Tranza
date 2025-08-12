// Authentication API service layer using HttpOnly cookies
import { 
  User, 
  CookieAuthResponse, 
  LoginRequest, 
  RegisterRequest,
  PreRegistrationRequest,
  EmailVerificationRequest,
  ResendVerificationRequest,
  PreRegistrationResponse,
  EmailVerificationResponse
} from '@/types/auth';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class AuthError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'AuthError';
  }
}

class AuthService {
  private getHeaders(): HeadersInit {
    return {
      'Content-Type': 'application/json',
    };
  }

  private getRequestOptions(method: string = 'GET', body?: any): RequestInit {
    return {
      method,
      headers: this.getHeaders(),
      credentials: 'include', // Essential for HttpOnly cookies
      body: body ? JSON.stringify(body) : undefined,
    };
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new AuthError(
        errorData.error || `HTTP error! status: ${response.status}`,
        response.status
      );
    }

    return response.json();
  }

  // Legacy register method (deprecated - use preRegister + verifyEmail instead)
  async register(data: RegisterRequest): Promise<CookieAuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/signup`, 
      this.getRequestOptions('POST', data)
    );

    const result = await this.handleResponse<{ message: string; data: CookieAuthResponse }>(response);
    return result.data;
  }

  // New two-step registration process
  async preRegister(data: PreRegistrationRequest): Promise<PreRegistrationResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/pre-register`, 
      this.getRequestOptions('POST', data)
    );

    const result = await this.handleResponse<{ message: string; data: PreRegistrationResponse }>(response);
    return result.data;
  }

  async verifyEmail(data: EmailVerificationRequest): Promise<EmailVerificationResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/verify-email`, 
      this.getRequestOptions('POST', data)
    );

    const result = await this.handleResponse<{ message: string; data: EmailVerificationResponse }>(response);
    return result.data;
  }

  async resendVerificationCode(data: ResendVerificationRequest): Promise<PreRegistrationResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/resend-verification`, 
      this.getRequestOptions('POST', data)
    );

    const result = await this.handleResponse<{ message: string; data: PreRegistrationResponse }>(response);
    return result.data;
  }

  async login(data: LoginRequest): Promise<CookieAuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, 
      this.getRequestOptions('POST', data)
    );

    const result = await this.handleResponse<{ message: string; data: CookieAuthResponse }>(response);
    return result.data;
  }

  async refreshToken(): Promise<CookieAuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, 
      this.getRequestOptions('POST')
    );

    const result = await this.handleResponse<{ message: string; data: CookieAuthResponse }>(response);
    return result.data;
  }

  async validateToken(): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/auth/validate`, 
      this.getRequestOptions('GET')
    );

    const result = await this.handleResponse<{ message: string; user: User }>(response);
    return result.user;
  }

  async logout(): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, 
      this.getRequestOptions('POST')
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new AuthError(
        errorData.error || `HTTP error! status: ${response.status}`,
        response.status
      );
    }
  }

  async getOAuthUrl(provider: 'google' | 'github', state: string): Promise<string> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/${provider}?state=${state}`, 
      this.getRequestOptions('GET')
    );

    const result = await this.handleResponse<{ url: string }>(response);
    return result.url;
  }

  async handleOAuthCallback(
    provider: string,
    code: string,
    state?: string,
    redirectUri?: string
  ): Promise<CookieAuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/callback`, 
      this.getRequestOptions('POST', {
        provider,
        code,
        state,
        redirect_uri: redirectUri,
      })
    );

    const result = await this.handleResponse<{ message: string; data: CookieAuthResponse }>(response);
    return result.data;
  }
}

export const authService = new AuthService();
export { AuthError };