import { apiClient } from '@/lib/api-client';
import {
  APIKey,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  RevokeAPIKeyRequest,
  APIResponse,
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
}

export default APIKeyService;
