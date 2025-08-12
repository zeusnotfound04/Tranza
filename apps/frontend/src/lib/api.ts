// Authenticated API utility for making requests with HttpOnly cookies
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'ApiError';
  }
}

class ApiClient {
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
      throw new ApiError(
        errorData.error || `HTTP error! status: ${response.status}`,
        response.status
      );
    }

    // Handle empty responses (like 204 No Content)
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      return {} as T;
    }

    return response.json();
  }

  // Generic GET request
  async get<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, 
      this.getRequestOptions('GET')
    );
    return this.handleResponse<T>(response);
  }

  // Generic POST request
  async post<T>(endpoint: string, data?: any): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, 
      this.getRequestOptions('POST', data)
    );
    return this.handleResponse<T>(response);
  }

  // Generic PUT request
  async put<T>(endpoint: string, data?: any): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, 
      this.getRequestOptions('PUT', data)
    );
    return this.handleResponse<T>(response);
  }

  // Generic DELETE request
  async delete<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, 
      this.getRequestOptions('DELETE')
    );
    return this.handleResponse<T>(response);
  }

  // Generic PATCH request
  async patch<T>(endpoint: string, data?: any): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, 
      this.getRequestOptions('PATCH', data)
    );
    return this.handleResponse<T>(response);
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Convenience functions for common patterns
export const api = {
  // Profile/User endpoints
  getProfile: () => apiClient.get('/api/profile'),
  updateProfile: (data: any) => apiClient.put('/api/profile', data),
  
  // Transactions endpoints (example)
  getTransactions: () => apiClient.get('/api/transactions'),
  createTransaction: (data: any) => apiClient.post('/api/transactions', data),
  getTransaction: (id: string) => apiClient.get(`/api/transactions/${id}`),
  
  // API Keys endpoints (example)
  getApiKeys: () => apiClient.get('/api/keys'),
  createApiKey: (data: any) => apiClient.post('/api/keys', data),
  deleteApiKey: (id: string) => apiClient.delete(`/api/keys/${id}`),
  
  // Generic methods
  get: <T>(endpoint: string) => apiClient.get<T>(endpoint),
  post: <T>(endpoint: string, data?: any) => apiClient.post<T>(endpoint, data),
  put: <T>(endpoint: string, data?: any) => apiClient.put<T>(endpoint, data),
  delete: <T>(endpoint: string) => apiClient.delete<T>(endpoint),
  patch: <T>(endpoint: string, data?: any) => apiClient.patch<T>(endpoint, data),
};