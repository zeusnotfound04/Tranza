// User Types
export interface User {
  id: string;
  email: string;
  username: string;
  avatar?: string;
  provider: string;
  provider_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  access_token?: string;
  refresh_token?: string;
  expires_in?: number;
}

// Authentication Request Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface PreRegistrationRequest {
  username: string;
  email: string;
  password: string;
}

export interface EmailVerificationRequest {
  email: string;
  code: string;
  username: string;
  password: string;
}

export interface ResendVerificationRequest {
  email: string;
}

export interface OAuthCallbackRequest {
  code: string;
  provider: string;
  redirect_uri: string;
  state?: string;
}

// Wallet Types
export interface Wallet {
  id: string;
  user_id: string;
  balance: number;
  currency: string;
  status: string;
  ai_access_enabled: boolean;
  ai_daily_limit: number;
  ai_per_transaction_limit: number;
  created_at: string;
  updated_at: string;
}

export interface LoadMoneyRequest {
  amount: number;
}

export interface LoadMoneyResponse {
  order_id: string;
  amount: number;
  currency: string;
  transaction_id: string;
  razorpay_key_id: string;
}

export interface VerifyPaymentRequest {
  razorpay_payment_id: string;
  razorpay_order_id: string;
  razorpay_signature: string;
}

export interface UpdateWalletSettingsRequest {
  ai_access_enabled?: boolean;
  ai_daily_limit?: number;
  ai_per_transaction_limit?: number;
}

// Card Types
export interface LinkedCard {
  id: number;
  user_id: number;
  card_number: string;
  card_holder_name: string;
  expiry_month: number;
  expiry_year: number;
  cvv: string;
  card_type: string;
  daily_limit: number;
  monthly_limit: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CardLinkRequest {
  card_number: string;
  card_holder_name: string;
  expiry_month: number;
  expiry_year: number;
  cvv: string;
  card_type: string;
  daily_limit?: number;
  monthly_limit?: number;
}

export interface UpdateCardLimitRequest {
  daily_limit?: number;
  monthly_limit?: number;
}

// Transaction Types
export interface Transaction {
  id: string;
  wallet_id: string;
  user_id: string;
  type: string;
  category: string;
  amount: number;
  currency: string;
  description: string;
  status: string;
  reference_id?: string;
  payment_method: string;
  payment_gateway: string;
  gateway_transaction_id?: string;
  gateway_payment_id?: string;
  gateway_order_id?: string;
  fee_amount: number;
  tax_amount: number;
  net_amount: number;
  exchange_rate?: number;
  original_amount?: number;
  original_currency?: string;
  merchant_id?: string;
  merchant_name?: string;
  location?: string;
  ip_address?: string;
  user_agent?: string;
  ai_category?: string;
  ai_confidence?: number;
  ai_tags?: string[];
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  processed_at?: string;
}

export interface TransactionFilters {
  type?: string;
  category?: string;
  status?: string;
  payment_method?: string;
  start_date?: string;
  end_date?: string;
  min_amount?: number;
  max_amount?: number;
  merchant_name?: string;
  limit?: number;
  offset?: number;
}

export interface TransactionStats {
  total_transactions: number;
  total_amount: number;
  successful_transactions: number;
  failed_transactions: number;
  pending_transactions: number;
  average_transaction_amount: number;
  total_fees: number;
  by_category: Record<string, { count: number; amount: number }>;
  by_payment_method: Record<string, { count: number; amount: number }>;
}

export interface TransactionAnalytics {
  period: string;
  total_volume: number;
  transaction_count: number;
  average_amount: number;
  growth_rate: number;
  top_categories: Array<{ category: string; amount: number; count: number }>;
  top_merchants: Array<{ merchant: string; amount: number; count: number }>;
  success_rate: number;
  peak_hours: number[];
}

export interface MonthlyTransactionSummary {
  month: string;
  year: number;
  total_amount: number;
  total_count: number;
  by_category: Record<string, { amount: number; count: number }>;
  by_status: Record<string, number>;
}

export interface DailyTransactionSummary {
  date: string;
  total_amount: number;
  total_count: number;
  by_hour: Array<{ hour: number; amount: number; count: number }>;
}

export interface TransactionTrends {
  period: string;
  data_points: Array<{
    date: string;
    amount: number;
    count: number;
  }>;
  trend_direction: 'up' | 'down' | 'stable';
  growth_percentage: number;
}

// Payment Types
export interface PaymentOrder {
  id: string;
  amount: number;
  currency: string;
  receipt: string;
  status: string;
  created_at: string;
  notes?: Record<string, string>;
}

export interface CreateOrderRequest {
  amount: number;
  currency?: string;
  receipt?: string;
  notes?: Record<string, string>;
  description?: string;
}

export interface VerifyPaymentResponse {
  payment_id: string;
  order_id: string;
  status: string;
  amount: number;
  currency: string;
}

export interface WalletVerificationResponse {
  success: boolean;
  new_balance: number;
  transaction_id: string;
  message: string;
  amount: number;
}

// API Key Types
export interface APIKey {
  id: number;
  label: string;
  key_type: string;
  scopes: string[];
  usage_count: number;
  rate_limit: number;
  is_active: boolean;
  created_at: string;
  expires_at?: string;
  last_used_at?: string;
}

export interface ListAPIKeysResponse {
  keys: APIKey[];
  total: number;
}

export interface CreateAPIKeyRequest {
  label: string;
  ttl_hours: number;
}

export interface CreateAPIKeyResponse {
  api_key: string;
  expires_at: string;
}

export interface RevokeAPIKeyRequest {
  key_id: number;
}

// API Usage Types
export interface APIUsageLog {
  id: number;
  api_key_id: number;
  endpoint: string;
  method: string;
  status_code: number;
  request_size: number;
  response_size: number;
  response_time: number;
  ip_address: string;
  user_agent?: string;
  command?: string;
  amount_spent?: number;
  currency?: string;
  error_message?: string;
  timestamp: string;
}

export interface UsageStatsResponse {
  key_id: number;
  total_requests: number;
  total_amount_spent: number;
  spending_limit: number;
  remaining_limit: number;
  currency: string;
  avg_response_time: number;
  success_rate: number;
  last_used_at?: string;
  period_start: string;
  period_end: string;
  command_usage: CommandUsage[];
  time_series_data: TimeSeriesData[];
}

export interface CommandUsage {
  command: string;
  count: number;
  total_amount: number;
  avg_amount: number;
  success_rate: number;
  avg_response_time: number;
}

export interface TimeSeriesData {
  date: string;
  requests: number;
  amount_spent: number;
  avg_response_time: number;
  success_rate: number;
}

export interface DetailedUsageResponse {
  logs: APIUsageLog[];
  total: number;
  page: number;
  limit: number;
  has_more: boolean;
}

// Response Wrappers
export interface APIResponse<T> {
  message?: string;
  data?: T;
  user?: User;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

// Error Types
export interface APIErrorResponse {
  error: string;
  message?: string;
  code?: string;
  details?: any;
}

// OAuth Types
export interface OAuthURL {
  url: string;
  provider: string;
  state: string;
}

// Export all types
export type {
  // Re-export for convenience
  User as UserType,
  Wallet as WalletType,
  Transaction as TransactionType,
  LinkedCard as CardType,
  APIKey as APIKeyType,
};
