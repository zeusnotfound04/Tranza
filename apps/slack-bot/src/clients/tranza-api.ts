import axios, { AxiosInstance, AxiosResponse } from 'axios';

export interface TranzaAPIConfig {
  baseURL: string;
  apiKey: string;
  timeout?: number;
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
}

export interface WalletBalanceResponse {
  user_id: string;
  message: string;
  // TODO: Add actual balance fields when wallet service is integrated
}

export interface APIErrorResponse {
  error: string;
  message?: string;
  details?: any;
}

// Create and configure the axios client
export const createAPIClient = (config: TranzaAPIConfig): AxiosInstance => {
  const client = axios.create({
    baseURL: config.baseURL,
    timeout: config.timeout || 30000,
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': config.apiKey,
    },
  });

  // Add request interceptor for logging
  client.interceptors.request.use(
    (config) => {
      console.log(`🔗 API Request: ${config.method?.toUpperCase()} ${config.url}`);
      return config;
    },
    (error) => {
      console.error('🚨 API Request Error:', error);
      return Promise.reject(error);
    }
  );

  // Add response interceptor for error handling
  client.interceptors.response.use(
    (response) => {
      console.log(`✅ API Response: ${response.status} ${response.config.url}`);
      return response;
    },
    (error) => {
      console.error('🚨 API Response Error:', {
        status: error.response?.status,
        url: error.config?.url,
        data: error.response?.data,
      });
      return Promise.reject(handleAPIError(error));
    }
  );

  return client;
};

// Handle API errors consistently
const handleAPIError = (error: any): Error => {
  if (error.response) {
    // Server responded with error status
    const { status, data } = error.response;
    
    if (status === 401) {
      return new Error('Invalid or expired API key');
    } else if (status === 403) {
      return new Error('Insufficient permissions for this operation');
    } else if (status === 429) {
      return new Error('Rate limit exceeded. Please try again later');
    } else if (status >= 500) {
      return new Error('Backend service temporarily unavailable');
    } else {
      const errorMessage = data?.error || data?.message || `API Error: ${status}`;
      return new Error(errorMessage);
    }
  } else if (error.request) {
    // Network error
    return new Error('Unable to connect to Tranza backend service');
  } else {
    // Other error
    return new Error(`Request failed: ${error.message}`);
  }
};

/**
 * Validate a transfer before creating it
 */
export const validateTransfer = async (
  client: AxiosInstance,
  request: TransferValidationRequest
): Promise<TransferValidationResponse> => {
  try {
    const response: AxiosResponse<{ data: TransferValidationResponse }> = await client.post(
      '/api/bot/transfers/validate',
      request
    );
    return response.data.data;
  } catch (error) {
    console.error('❌ Transfer validation failed:', error);
    throw error;
  }
};

/**
 * Create a new transfer
 */
export const createTransfer = async (
  client: AxiosInstance,
  request: CreateTransferRequest
): Promise<CreateTransferResponse> => {
  try {
    const response: AxiosResponse<{ data: CreateTransferResponse }> = await client.post(
      '/api/bot/transfers',
      request
    );
    return response.data.data;
  } catch (error) {
    console.error('❌ Transfer creation failed:', error);
    throw error;
  }
};

/**
 * Get transfer status by ID
 */
export const getTransferStatus = async (
  client: AxiosInstance,
  transferId: string
): Promise<TransferStatusResponse> => {
  try {
    const response: AxiosResponse<{ data: TransferStatusResponse }> = await client.get(
      `/api/bot/transfers/${transferId}/status`
    );
    return response.data.data;
  } catch (error) {
    console.error('❌ Get transfer status failed:', error);
    throw error;
  }
};

/**
 * Get wallet balance
 */
export const getWalletBalance = async (
  client: AxiosInstance
): Promise<WalletBalanceResponse> => {
  try {
    const response: AxiosResponse<{ data: WalletBalanceResponse }> = await client.get(
      '/api/bot/wallet/balance'
    );
    return response.data.data;
  } catch (error) {
    console.error('❌ Get wallet balance failed:', error);
    throw error;
  }
};

/**
 * Test the API connection and authentication
 */
export const testConnection = async (client: AxiosInstance): Promise<boolean> => {
  try {
    await getWalletBalance(client);
    return true;
  } catch (error) {
    console.error('❌ API connection test failed:', error);
    return false;
  }
};

/**
 * Update the API key for a client
 */
export const updateAPIKey = (client: AxiosInstance, newApiKey: string): void => {
  client.defaults.headers['X-API-Key'] = newApiKey;
  console.log('🔄 API key updated successfully');
};
