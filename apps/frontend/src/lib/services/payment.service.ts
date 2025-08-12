import { apiClient } from '@/lib/api-client';
import {
  PaymentOrder,
  CreateOrderRequest,
  VerifyPaymentRequest,
  VerifyPaymentResponse,
  APIResponse,
} from '@/types/api';

export class PaymentService {
  // Create Razorpay order
  static async createOrder(data: CreateOrderRequest): Promise<APIResponse<PaymentOrder>> {
    return apiClient.post('/api/v1/payments/orders', data);
  }

  // Verify payment
  static async verifyPayment(data: VerifyPaymentRequest): Promise<APIResponse<VerifyPaymentResponse>> {
    return apiClient.post('/api/v1/payments/verify', data);
  }

  // Get order details
  static async getOrder(orderId: string): Promise<APIResponse<PaymentOrder>> {
    return apiClient.get(`/api/v1/payments/orders/${orderId}`);
  }

  // Get payment details
  static async getPayment(paymentId: string): Promise<APIResponse<any>> {
    return apiClient.get(`/api/v1/payments/payments/${paymentId}`);
  }
}

export default PaymentService;
