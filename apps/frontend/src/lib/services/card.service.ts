import { apiClient } from '@/lib/api-client';
import {
  LinkedCard,
  CardLinkRequest,
  UpdateCardLimitRequest,
  APIResponse,
} from '@/types/api';

export class CardService {
  // Link a new card
  static async linkCard(data: CardLinkRequest): Promise<APIResponse<any>> {
    return apiClient.post('/api/v1/cards', data);
  }

  static async getCards(): Promise<APIResponse<LinkedCard[]>> {
    return apiClient.get('/api/v1/cards');
  }

  static async deleteCard(cardId: number): Promise<APIResponse<any>> {
    return apiClient.delete(`/api/v1/cards/${cardId}`);
  }

  static async updateLimit(cardId: number, data: UpdateCardLimitRequest): Promise<APIResponse<any>> {
    return apiClient.put(`/api/v1/cards/${cardId}/limit`, data);
  }
}

export default CardService;
