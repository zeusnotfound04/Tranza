import { apiClient } from '@/lib/api-client';
import {
  User,
  APIResponse,
} from '@/types/api';

export interface UpdateProfileRequest {
  username?: string;
  email?: string;
  avatar?: string;
}

export class ProfileService {
  // Get current user profile
  static async getProfile(): Promise<APIResponse<User>> {
    return apiClient.get('/api/v1/profile');
  }

  // Update user profile
  static async updateProfile(data: UpdateProfileRequest): Promise<APIResponse<User>> {
    return apiClient.put('/api/v1/profile', data);
  }
}

export default ProfileService;
