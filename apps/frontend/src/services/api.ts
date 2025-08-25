// Base API configuration and client
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface TransferValidationRequest {
  amount: string;
  recipient_type: 'upi' | 'phone';
  recipient_value: string;
}

export interface TransferValidationResponse {
  valid: boolean;
  errors: string[];
  warnings: string[];
  transfer_fee: string;
  total_amount: string;
  estimated_time: string;
}

export interface CreateTransferRequest {
  amount: string;
  recipient_type: 'upi' | 'phone';
  recipient_value: string;
  recipient_name?: string;
  description?: string;
}

export interface CreateTransferResponse {
  transfer_id: string;
  reference_id: string;
  amount: string;
  transfer_fee: string;
  total_amount: string;
  status: string;
  recipient: string;
  estimated_time: string;
}

export interface TransferStatusResponse {
  transfer_id: string;
  reference_id: string;
  status: string;
  amount: string;
  recipient: string;
  estimated_time: string;
  created_at: string;
  updated_at: string;
}

export interface WalletBalanceResponse {
  balance: number;
  currency: string;
  last_updated: string;
}

export interface TransferHistoryResponse {
  transfers: TransferStatusResponse[];
  total_count: number;
  page: number;
  limit: number;
}

// API Client class for making authenticated requests
export class TranzaAPIClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || API_BASE_URL;
  }

  setAuthToken(token: string): void {
    this.token = token;
  }

  private async makeRequest<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    // Get the current access token
    let accessToken = this.token;
    if (!accessToken) {
      try {
        const tokenResponse = await fetch("/api/token", {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (tokenResponse.ok) {
          const tokenData = await tokenResponse.json();
          accessToken = tokenData.data.token;
          console.log("✅ Retrieved access token from /api/token:", tokenData);
          console.log("🔑 Access token value:", accessToken);
        } else {
          console.error("❌ Failed to retrieve access token, status:", tokenResponse.status);
          const errorText = await tokenResponse.text();
          console.error("❌ Token response error:", errorText);
        }
      } catch (error) {
        console.error("Error fetching access token:", error);
      }
    }

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {}),
    };
    
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
      console.log("🚀 Making request to:", url);
      console.log("📋 Request headers:", headers);
      console.log("🔐 Authorization header:", headers['Authorization']);
    } else {
      console.warn("⚠️ No access token available for request to:", url);
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      // Check if response has content before trying to parse JSON
      const contentType = response.headers.get('content-type');
      let data: any = {};

      if (contentType && contentType.includes('application/json')) {
        try {
          const text = await response.text();
          if (text.trim()) {
            data = JSON.parse(text);
          }
        } catch (jsonError) {
          console.error('JSON Parse Error:', jsonError);
          return {
            success: false,
            error: 'Invalid JSON response from server',
          };
        }
      } else {
        // Non-JSON response - get text content
        const text = await response.text();
        if (text.trim()) {
          console.warn('Non-JSON response:', text);
          data = { message: text };
        }
      }

      if (!response.ok) {
        return {
          success: false,
          error: data.error || data.message || `HTTP ${response.status}: ${response.statusText}`,
          data: data,
        };
      }

      // Handle different response structures from the backend
      if (data.success !== undefined) {
        // Backend returns { success: boolean, data: any, message?: string }
        return {
          success: data.success,
          data: data.data,
          message: data.message,
          error: data.success ? undefined : data.error || data.message,
        };
      } else {
        // Direct data response
        return {
          success: true,
          data: data,
        };
      }
    } catch (error) {
      console.error('API Request Error:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Network error',
      };
    }
  }

  // Authentication methods
  async login(email: string, password: string): Promise<APIResponse<{ token: string; user: any }>> {
    return this.makeRequest('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async register(userData: {
    first_name: string;
    last_name: string;
    email: string;
    password: string;
  }): Promise<APIResponse<{ token: string; user: any }>> {
    return this.makeRequest('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async getUserProfile(): Promise<APIResponse<any>> {
    return this.makeRequest('/api/v1/profile');
  }

  async verifyToken(): Promise<APIResponse<{ valid: boolean; user?: any }>> {
    return this.makeRequest('/api/auth/verify');
  }

  // Wallet methods
  async getWalletBalance(): Promise<APIResponse<WalletBalanceResponse>> {
    return this.makeRequest('/api/v1/wallet/balance');
  }

  async loadMoney(amount: number, method: string): Promise<APIResponse<any>> {
    return this.makeRequest('/api/v1/wallet/load', {
      method: 'POST',
      body: JSON.stringify({ amount, method }),
    });
  }

  // External transfer methods
  async validateTransfer(request: TransferValidationRequest): Promise<APIResponse<TransferValidationResponse>> {
    return this.makeRequest('/api/external/transfers/validate', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async createTransfer(request: CreateTransferRequest): Promise<APIResponse<CreateTransferResponse>> {
    return this.makeRequest('/api/external/transfers', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async getTransferStatus(transferId: string): Promise<APIResponse<TransferStatusResponse>> {
    return this.makeRequest(`/api/external/transfers/${transferId}/status`);
  }

  async getTransferHistory(page = 1, limit = 10): Promise<APIResponse<TransferHistoryResponse>> {
    return this.makeRequest(`/api/external/transfers/history?page=${page}&limit=${limit}`);
  }

  async getTransferFees(amount: string, recipientType: 'upi' | 'phone'): Promise<APIResponse<{ fee: string; total: string }>> {
    return this.makeRequest(`/api/external/transfers/fees?amount=${amount}&type=${recipientType}`);
  }

  // API Key management (for bot users)
  async generateAPIKey(label: string, ttlHours: number = 8760): Promise<APIResponse<{ api_key: string; key_id: string }>> {
    return this.makeRequest('/api/v1/keys', {
      method: 'POST',
      body: JSON.stringify({ label, ttl_hours: ttlHours }),
    });
  }

  async generateBotAPIKey(label: string, workspaceId: string, botUserId: string, ttlHours: number = 8760): Promise<APIResponse<{ api_key: string; key_id: string }>> {
    return this.makeRequest('/api/v1/keys/bot', {
      method: 'POST',
      body: JSON.stringify({ 
        label, 
        workspace_id: workspaceId, 
        bot_user_id: botUserId, 
        ttl_hours: ttlHours 
      }),
    });
  }

  async getAPIKeys(): Promise<APIResponse<any[]>> {
    return this.makeRequest('/api/v1/keys');
  }

  async revokeAPIKey(keyId: string): Promise<APIResponse<any>> {
    return this.makeRequest(`/api/v1/keys/${keyId}`, {
      method: 'DELETE',
    });
  }
}

// Export singleton instance
export const apiClient = new TranzaAPIClient();

// Helper functions for common operations
export const validateAmount = (amount: string): boolean => {
  const num = parseFloat(amount);
  return !isNaN(num) && num > 0 && num <= 100000; // Max ₹1 lakh per transfer
};

export const validateUPI = (upiId: string): boolean => {
  const upiRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+$/;
  return upiRegex.test(upiId);
};

export const validatePhone = (phone: string): boolean => {
  const phoneRegex = /^[6-9]\d{9}$/; // Indian mobile number
  return phoneRegex.test(phone);
};

export const formatCurrency = (amount: number | string): string => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
  }).format(num);
};

export const formatDateTime = (dateString: string): string => {
  return new Intl.DateTimeFormat('en-IN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateString));
};
