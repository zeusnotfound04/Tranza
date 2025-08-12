import { apiClient } from '@/lib/api-client';
import {
  Transaction,
  TransactionFilters,
  TransactionStats,
  TransactionAnalytics,
  MonthlyTransactionSummary,
  DailyTransactionSummary,
  TransactionTrends,
  PaginatedResponse,
  APIResponse,
} from '@/types/api';

export class TransactionService {
  // Get transaction history with pagination and filters
  static async getTransactionHistory(filters: TransactionFilters = {}): Promise<PaginatedResponse<Transaction>> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    const queryString = params.toString();
    const endpoint = queryString ? `/api/v1/transactions?${queryString}` : '/api/v1/transactions';
    
    return apiClient.get(endpoint);
  }

  // Get specific transaction
  static async getTransaction(id: string): Promise<APIResponse<Transaction>> {
    return apiClient.get(`/api/v1/transactions/${id}`);
  }

  // Search transactions
  static async searchTransactions(filters: TransactionFilters & { query?: string }): Promise<PaginatedResponse<Transaction>> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    return apiClient.get(`/api/v1/transactions/search?${params.toString()}`);
  }

  // Get transactions by type
  static async getTransactionsByType(type: string, filters: TransactionFilters = {}): Promise<PaginatedResponse<Transaction>> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    const queryString = params.toString();
    const endpoint = queryString ? `/api/v1/transactions/type/${type}?${queryString}` : `/api/v1/transactions/type/${type}`;
    
    return apiClient.get(endpoint);
  }

  // Get transaction receipt
  static async getTransactionReceipt(id: string): Promise<APIResponse<any>> {
    return apiClient.get(`/api/v1/transactions/${id}/receipt`);
  }

  // Get transaction statistics
  static async getTransactionStats(filters: Omit<TransactionFilters, 'limit' | 'offset'> = {}): Promise<APIResponse<TransactionStats>> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    const queryString = params.toString();
    const endpoint = queryString ? `/api/v1/transactions/stats?${queryString}` : '/api/v1/transactions/stats';
    
    return apiClient.get(endpoint);
  }

  // Get transaction analytics
  static async getTransactionAnalytics(period: string = '30d'): Promise<APIResponse<TransactionAnalytics>> {
    return apiClient.get(`/api/v1/transactions/analytics?period=${period}`);
  }

  // Get monthly transaction summary
  static async getMonthlyTransactionSummary(year: number, month: number): Promise<APIResponse<MonthlyTransactionSummary>> {
    return apiClient.get(`/api/v1/transactions/summary/monthly?year=${year}&month=${month}`);
  }

  // Get daily transaction summary
  static async getDailyTransactionSummary(date: string): Promise<APIResponse<DailyTransactionSummary>> {
    return apiClient.get(`/api/v1/transactions/summary/daily?date=${date}`);
  }

  // Get transaction trends
  static async getTransactionTrends(period: string = '30d'): Promise<APIResponse<TransactionTrends>> {
    return apiClient.get(`/api/v1/transactions/trends?period=${period}`);
  }

  // Export transactions
  static async exportTransactions(filters: TransactionFilters & { format?: 'csv' | 'pdf' } = {}): Promise<Blob> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/transactions/export?${params.toString()}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Export failed');
    }

    return response.blob();
  }

  // Validate transaction (admin function)
  static async validateTransaction(id: string): Promise<APIResponse<any>> {
    return apiClient.post(`/api/v1/transactions/${id}/validate`);
  }

  // Retry failed transaction
  static async retryFailedTransaction(id: string): Promise<APIResponse<any>> {
    return apiClient.post(`/api/v1/transactions/${id}/retry`);
  }
}

export default TransactionService;
