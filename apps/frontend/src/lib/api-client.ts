// API Configuration and Base Service
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class APIError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: any
  ) {
    super(message);
    this.name = 'APIError';
  }
}

// Token management
class TokenManager {
  private static readonly ACCESS_TOKEN_KEY = 'access_token';
  private static readonly REFRESH_TOKEN_KEY = 'refresh_token';

  static setTokens(accessToken: string, refreshToken?: string) {
    if (typeof window !== 'undefined') {
      localStorage.setItem(this.ACCESS_TOKEN_KEY, accessToken);
      if (refreshToken) {
        localStorage.setItem(this.REFRESH_TOKEN_KEY, refreshToken);
      }
    }
  }

  static async getAccessToken(): Promise<string | null> {
        console.log("🔍 [TokenManager] Requesting token from /api/token");
    
    const tokenResponse = await fetch("/api/token", {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (tokenResponse.ok) {
      const tokenData = await tokenResponse.json();
      console.log("✅ [TokenManager] Token response:", tokenData);
      console.log("🔑 [TokenManager] Access token:", tokenData.data.token);
      return tokenData.data.token;
    }

    console.error("❌ [TokenManager] Failed to retrieve access token, status:", tokenResponse.status);
    const errorText = await tokenResponse.text();
    console.error("❌ [TokenManager] Error response:", errorText);
    return null;
  }

  static getRefreshToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(this.REFRESH_TOKEN_KEY);
    }
    return null;
  }

  static clearTokens() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(this.ACCESS_TOKEN_KEY);
      localStorage.removeItem(this.REFRESH_TOKEN_KEY);
    }
  }
}

// Base API client with both cookie and token support
class APIClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    // Add Bearer token if available
    const accessToken = await TokenManager.getAccessToken();
    console.log("🔍 DEBUG API Client: Raw access token:", accessToken);
    console.log("🔍 DEBUG API Client: Token type:", typeof accessToken);
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>,
    };

    if (accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
      console.log("✅ DEBUG API Client: Authorization header set:", headers.Authorization.substring(0, 20) + "...");
    } else {
      console.log("❌ DEBUG API Client: No access token available!");
    }

    const config: RequestInit = {
      credentials: 'include', // Still include for backward compatibility
      headers,
      ...options,
    };

    // Debug logs
    console.log('DEBUG API Client: Making request to:', url);
    console.log('DEBUG API Client: Request headers:', headers);
    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        let errorData;
        try {
          errorData = await response.json();
        } catch {
          errorData = { message: response.statusText };
        }
        
        throw new APIError(
          errorData.error || errorData.message || 'Request failed',
          response.status,
          errorData
        );
      }

      // Handle responses that might not have JSON body
      const contentType = response.headers.get('content-type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }
      
      return {} as T;
    } catch (error) {
      if (error instanceof APIError) {
        throw error;
      }
      throw new APIError('Network error occurred', 0, error);
    }
  }

  async get<T>(endpoint: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET', headers });
  }

  async post<T>(
    endpoint: string,
    data?: any,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
      headers,
    });
  }

  async put<T>(
    endpoint: string,
    data?: any,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
      headers,
    });
  }

  async delete<T>(
    endpoint: string,
    data?: any,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(endpoint, { 
      method: 'DELETE', 
      body: data ? JSON.stringify(data) : undefined,
      headers 
    });
  }

  async patch<T>(
    endpoint: string,
    data?: any,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
      headers,
    });
  }
}

export const apiClient = new APIClient();
export { APIError, TokenManager };
