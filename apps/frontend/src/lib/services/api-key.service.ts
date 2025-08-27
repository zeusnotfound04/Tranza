import { apiClient } from '@/lib/api-client';
import {
  APIKey,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  RevokeAPIKeyRequest,
  APIResponse,
  UsageStatsResponse,
  DetailedUsageResponse,
  TimeSeriesData,
  CommandUsage,
} from '@/types/api';

export class APIKeyService {
  // Create new API key
  static async createAPIKey(data: CreateAPIKeyRequest): Promise<APIResponse<CreateAPIKeyResponse>> {
    return apiClient.post('/api/v1/api-keys', data);
  }

  // Revoke API key
  static async revokeAPIKey(data: RevokeAPIKeyRequest): Promise<APIResponse<any>> {
    return apiClient.delete('/api/v1/api-keys', data);
  }

  // Note: Getting API keys list is not implemented in backend
  // but would be useful for UI
  static async getAPIKeys(): Promise<APIResponse<APIKey[]>> {
    // This endpoint doesn't exist yet, but would be useful
    // return apiClient.get('/api/v1/api-keys');
    throw new Error('Get API keys endpoint not implemented yet');
  }

  // Get detailed usage statistics for an API key
  static async getUsageStats(keyId: number, days: number = 30): Promise<APIResponse<UsageStatsResponse>> {
    return apiClient.get(`/api/keys/${keyId}/usage/detailed?days=${days}`);
  }

  // Get paginated usage logs for an API key
  static async getUsageLogs(
    keyId: number, 
    page: number = 0, 
    limit: number = 50
  ): Promise<APIResponse<DetailedUsageResponse>> {
    return apiClient.get(`/api/keys/${keyId}/logs?page=${page}&limit=${limit}`);
  }

  // Get time series data for charts
  static async getTimeSeriesData(keyId: number, days: number = 30): Promise<APIResponse<TimeSeriesData[]>> {
    return apiClient.get(`/api/keys/${keyId}/usage/timeseries?days=${days}`);
  }

  // Get command usage data
  static async getCommandData(keyId: number, days: number = 30): Promise<APIResponse<CommandUsage[]>> {
    return apiClient.get(`/api/keys/${keyId}/usage/commands?days=${days}`);
  }
}

export default APIKeyService;
