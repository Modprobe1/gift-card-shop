import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Типы данных
export interface Currency {
  id: number;
  code: string;
  name: string;
  symbol: string;
  network: string;
  icon_url: string;
  is_active: boolean;
  min_amount: number;
  max_amount?: number;
  decimals: number;
}

export interface ExchangeRate {
  from_currency: Currency;
  to_currency: Currency;
  rate: number;
  reverse_rate: number;
  min_amount: number;
  max_amount?: number;
  updated_at: string;
}

export interface CalculateRequest {
  from_currency: string;
  to_currency: string;
  from_amount: number;
}

export interface CalculateResponse {
  from_amount: number;
  to_amount: number;
  rate: number;
  commission: number;
  commission_rate: number;
  exchange_rate: ExchangeRate;
}

export interface CreateOrderRequest {
  from_currency: string;
  to_currency: string;
  from_amount: number;
  client_phone: string;
  client_name: string;
  client_email: string;
  client_telegram?: string;
  recipient_wallet: string;
  recipient_details?: string;
}

export interface OrderStatus {
  id: number;
  name: string;
  description: string;
}

export interface ExchangeOrder {
  id: number;
  order_number: string;
  from_currency: Currency;
  to_currency: Currency;
  from_amount: number;
  to_amount: number;
  rate: number;
  status: OrderStatus;
  client_phone: string;
  client_name: string;
  client_email: string;
  client_telegram?: string;
  recipient_wallet: string;
  recipient_details?: string;
  expires_at: string;
  created_at: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
}

// API методы
export const exchangeApi = {
  // Получить все валюты
  getCurrencies: async (): Promise<Currency[]> => {
    const response = await api.get<ApiResponse<Currency[]>>('/currencies');
    return response.data.data;
  },

  // Получить валюту по коду
  getCurrency: async (code: string): Promise<Currency> => {
    const response = await api.get<ApiResponse<Currency>>(`/currencies/${code}`);
    return response.data.data;
  },

  // Получить все курсы обмена
  getExchangeRates: async (): Promise<ExchangeRate[]> => {
    const response = await api.get<ApiResponse<ExchangeRate[]>>('/rates');
    return response.data.data;
  },

  // Получить курс между валютами
  getExchangeRate: async (from: string, to: string): Promise<ExchangeRate> => {
    const response = await api.get<ApiResponse<ExchangeRate>>(`/rates/${from}/${to}`);
    return response.data.data;
  },

  // Рассчитать обмен
  calculate: async (request: CalculateRequest): Promise<CalculateResponse> => {
    const response = await api.post<ApiResponse<CalculateResponse>>('/calculate', request);
    return response.data.data;
  },

  // Создать заявку
  createOrder: async (request: CreateOrderRequest): Promise<ExchangeOrder> => {
    const response = await api.post<ApiResponse<ExchangeOrder>>('/orders', request);
    return response.data.data;
  },

  // Получить заявку по номеру
  getOrder: async (orderNumber: string): Promise<ExchangeOrder> => {
    const response = await api.get<ApiResponse<ExchangeOrder>>(`/orders/${orderNumber}`);
    return response.data.data;
  },

  // Проверка здоровья API
  healthCheck: async (): Promise<any> => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api;