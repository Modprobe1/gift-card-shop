import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if it exists
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export interface LoginData {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  created_at?: string;
  updated_at?: string;
}

export interface AuthResponse {
  message: string;
  token: string;
  user: User;
}

export const authAPI = {
  register: (data: RegisterData) => 
    api.post<AuthResponse>('/register', data),
  
  login: (data: LoginData) => 
    api.post<AuthResponse>('/login', data),
  
  getProfile: () => 
    api.get<{ user: User }>('/profile'),
};

export default api;