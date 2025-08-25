import { NextApiRequest, NextApiResponse } from 'next';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// AI Clothing Types
export interface AIClothingOrderRequest {
  prompt: string;
  address_id?: string;
  budget?: number;
  category?: string;
}

export interface AIClothingOrderResponse {
  order_created: boolean;
  required_info?: string[];
  suggested_products?: SuggestedProduct[];
  message: string;
  requires_confirmation: boolean;
  order_analysis?: ClothingOrderAnalysis;
  selected_address?: Address;
}

export interface SuggestedProduct {
  id: string;
  name: string;
  brand: string;
  price: number;
  currency: string;
  image_url: string;
  rating: number;
  description: string;
  url: string;
  website: string;
}

export interface ClothingOrderAnalysis {
  category: string;
  size: string;
  color: string;
  brand: string;
  occasion: string;
  min_price: number;
  max_price: number;
}

export interface SelectedProductForOrder {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
  size?: string;
  color?: string;
  url: string;
  website: string;
}

export interface ConfirmAIClothingOrderRequest {
  selected_products: SelectedProductForOrder[];
  address_id: string;
}

export interface ConfirmAIClothingOrderResponse {
  success: boolean;
  order_id?: string;
  message: string;
  total_amount?: number;
  remaining_balance?: number;
  required_amount?: number;
  current_balance?: number;
  selected_products?: SelectedProductForOrder[];
}

// Address Types
export interface Address {
  id: string;
  name: string;
  phone: string;
  address_line: string;
  city: string;
  state: string;
  pin_code: string;
  country: string;
  landmark?: string;
  is_default: boolean;
  address_type: string;
  created_at: string;
}

export interface AddressCreateRequest {
  name: string;
  phone: string;
  address_line: string;
  city: string;
  state: string;
  pin_code: string;
  country?: string;
  landmark?: string;
  is_default?: boolean;
  address_type?: string;
}

// API Response wrapper
export interface APIResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

// AI Clothing Service
export class AIClothingService {
  private static getAuthHeaders(): Record<string, string> {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('access_token');
      return {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      };
    }
    return { 'Content-Type': 'application/json' };
  }

  static async processAIClothingOrder(request: AIClothingOrderRequest): Promise<AIClothingOrderResponse> {
    const response = await fetch(`${BASE_URL}/api/v1/ai/clothing/order`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<AIClothingOrderResponse> = await response.json();
    return result.data;
  }

  static async confirmAIClothingOrder(request: ConfirmAIClothingOrderRequest): Promise<ConfirmAIClothingOrderResponse> {
    const response = await fetch(`${BASE_URL}/api/v1/ai/clothing/confirm`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<ConfirmAIClothingOrderResponse> = await response.json();
    return result.data;
  }
}

// Address Service
export class AddressService {
  private static getAuthHeaders(): Record<string, string> {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('access_token');
      return {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      };
    }
    return { 'Content-Type': 'application/json' };
  }

  static async getAddresses(): Promise<Address[]> {
    const response = await fetch(`${BASE_URL}/api/v1/addresses`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<Address[]> = await response.json();
    return result.data;
  }

  static async getDefaultAddress(): Promise<Address> {
    const response = await fetch(`${BASE_URL}/api/v1/addresses/default`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<Address> = await response.json();
    return result.data;
  }

  static async createAddress(address: AddressCreateRequest): Promise<Address> {
    const response = await fetch(`${BASE_URL}/api/v1/addresses`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(address),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<Address> = await response.json();
    return result.data;
  }

  static async setDefaultAddress(addressId: string): Promise<void> {
    const response = await fetch(`${BASE_URL}/api/v1/addresses/${addressId}/default`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  }
}
