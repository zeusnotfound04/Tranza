import { apiClient } from '@/lib/api-client';
import {
  Wallet,
  LoadMoneyRequest,
  LoadMoneyResponse,
  VerifyPaymentRequest,
  VerifyPaymentResponse,
  WalletVerificationResponse,
  UpdateWalletSettingsRequest,
  APIResponse,
} from '@/types/api';

export class WalletService {
  // Get wallet details
  static async getWallet(): Promise<APIResponse<Wallet>> {
    return apiClient.get('/api/v1/wallet');
  }

  // Update wallet settings
  static async updateSettings(data: UpdateWalletSettingsRequest): Promise<APIResponse<any>> {
    return apiClient.put('/api/v1/wallet/settings', data);
  }

  // Create load money order
  static async createLoadMoneyOrder(data: LoadMoneyRequest): Promise<APIResponse<LoadMoneyResponse>> {
    return apiClient.post('/api/v1/wallet/load', data);
  }

  // Verify payment and credit wallet
  static async verifyPayment(data: VerifyPaymentRequest): Promise<APIResponse<WalletVerificationResponse>> {
    return apiClient.post('/api/v1/wallet/verify-payment', data);
  }
}

export default WalletService;
