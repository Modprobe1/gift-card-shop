#!/bin/bash

echo "📦 Создаю React компоненты..."

# frontend/src/services/api.ts
cat > frontend/src/services/api.ts << 'EOF'
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
EOF

# frontend/src/contexts/AuthContext.tsx
cat > frontend/src/contexts/AuthContext.tsx << 'EOF'
import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { User, authAPI } from '../services/api';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is logged in on mount
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const response = await authAPI.getProfile();
          setUser(response.data.user);
        } catch (error) {
          console.error('Auth check failed:', error);
          localStorage.removeItem('token');
        }
      }
      setLoading(false);
    };
    checkAuth();
  }, []);

  const login = async (email: string, password: string) => {
    try {
      const response = await authAPI.login({ email, password });
      localStorage.setItem('token', response.data.token);
      setUser(response.data.user);
    } catch (error: any) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error);
      }
      throw new Error('Login failed');
    }
  };

  const register = async (email: string, password: string, name: string) => {
    try {
      const response = await authAPI.register({ email, password, name });
      localStorage.setItem('token', response.data.token);
      setUser(response.data.user);
    } catch (error: any) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error);
      }
      throw new Error('Registration failed');
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  const value = {
    user,
    loading,
    login,
    register,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
EOF

# frontend/src/components/PrivateRoute.tsx
cat > frontend/src/components/PrivateRoute.tsx << 'EOF'
import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const PrivateRoute: React.FC = () => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh' 
      }}>
        <div>Loading...</div>
      </div>
    );
  }

  return user ? <Outlet /> : <Navigate to="/login" />;
};

export default PrivateRoute;
EOF

# frontend/src/pages/Login.tsx
cat > frontend/src/pages/Login.tsx << 'EOF'
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import './Auth.css';

const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(email, password);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2 className="auth-title">Welcome Back</h2>
        <p className="auth-subtitle">Sign in to your account</p>
        
        <form onSubmit={handleSubmit} className="auth-form">
          {error && <div className="error-message">{error}</div>}
          
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
            />
          </div>
          
          <button type="submit" className="auth-button" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
        
        <p className="auth-link">
          Don't have an account? <Link to="/register">Sign up</Link>
        </p>
      </div>
    </div>
  );
};

export default Login;
EOF

# frontend/src/pages/Register.tsx
cat > frontend/src/pages/Register.tsx << 'EOF'
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import './Auth.css';

const Register: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (password.length < 6) {
      setError('Password must be at least 6 characters');
      return;
    }

    setLoading(true);

    try {
      await register(email, password, name);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2 className="auth-title">Create Account</h2>
        <p className="auth-subtitle">Sign up to get started</p>
        
        <form onSubmit={handleSubmit} className="auth-form">
          {error && <div className="error-message">{error}</div>}
          
          <div className="form-group">
            <label htmlFor="name">Full Name</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter your full name"
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Create a password"
              required
              minLength={6}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Confirm your password"
              required
            />
          </div>
          
          <button type="submit" className="auth-button" disabled={loading}>
            {loading ? 'Creating account...' : 'Sign Up'}
          </button>
        </form>
        
        <p className="auth-link">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </div>
    </div>
  );
};

export default Register;
EOF

# frontend/src/pages/Dashboard.tsx
cat > frontend/src/pages/Dashboard.tsx << 'EOF'
import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import './Dashboard.css';

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Dashboard</h1>
        <button onClick={logout} className="logout-button">
          Logout
        </button>
      </div>
      
      <div className="dashboard-content">
        <div className="welcome-card">
          <h2>Welcome back, {user?.name}!</h2>
          <p>You're successfully logged into your account.</p>
        </div>
        
        <div className="profile-card">
          <h3>Profile Information</h3>
          <div className="profile-info">
            <div className="info-item">
              <span className="info-label">Name:</span>
              <span className="info-value">{user?.name}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Email:</span>
              <span className="info-value">{user?.email}</span>
            </div>
            <div className="info-item">
              <span className="info-label">User ID:</span>
              <span className="info-value">{user?.id}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Account Created:</span>
              <span className="info-value">{formatDate(user?.created_at)}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Last Updated:</span>
              <span className="info-value">{formatDate(user?.updated_at)}</span>
            </div>
          </div>
        </div>
        
        <div className="activity-card">
          <h3>Recent Activity</h3>
          <p>Your recent activity will appear here.</p>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
EOF

# frontend/src/pages/Auth.css
cat > frontend/src/pages/Auth.css << 'EOF'
/* Auth pages styles */
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.auth-card {
  background: white;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  padding: 40px;
  width: 100%;
  max-width: 400px;
  animation: slideUp 0.4s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.auth-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin: 0 0 10px 0;
  text-align: center;
}

.auth-subtitle {
  font-size: 16px;
  color: #666;
  margin: 0 0 30px 0;
  text-align: center;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 16px;
  transition: all 0.3s ease;
  outline: none;
}

.form-group input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.auth-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.auth-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
}

.auth-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 12px;
  border-radius: 6px;
  font-size: 14px;
  text-align: center;
  animation: shake 0.4s ease-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

.auth-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #666;
}

.auth-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.3s ease;
}

.auth-link a:hover {
  color: #764ba2;
  text-decoration: underline;
}

/* Responsive design */
@media (max-width: 480px) {
  .auth-card {
    padding: 30px 20px;
  }
  
  .auth-title {
    font-size: 24px;
  }
}
EOF

# frontend/src/pages/Dashboard.css
cat > frontend/src/pages/Dashboard.css << 'EOF'
/* Dashboard styles */
.dashboard-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.dashboard-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.dashboard-header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
}

.logout-button {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid white;
  color: white;
  padding: 8px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.logout-button:hover {
  background: white;
  color: #667eea;
}

.dashboard-content {
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr;
  gap: 30px;
}

.welcome-card,
.profile-card,
.activity-card {
  background: white;
  border-radius: 10px;
  padding: 30px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-align: center;
}

.welcome-card h2 {
  margin: 0 0 10px 0;
  font-size: 32px;
  font-weight: 700;
}

.welcome-card p {
  margin: 0;
  font-size: 18px;
  opacity: 0.9;
}

.profile-card h3,
.activity-card h3 {
  margin: 0 0 25px 0;
  font-size: 20px;
  font-weight: 700;
  color: #333;
  padding-bottom: 15px;
  border-bottom: 2px solid #f0f0f0;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 600;
  color: #666;
  font-size: 14px;
}

.info-value {
  color: #333;
  font-size: 14px;
  font-weight: 500;
  word-break: break-all;
}

.activity-card {
  min-height: 200px;
}

.activity-card p {
  color: #666;
  font-size: 14px;
  text-align: center;
  margin-top: 40px;
}

/* Responsive design */
@media (min-width: 768px) {
  .dashboard-content {
    grid-template-columns: 1fr 1fr;
  }
  
  .welcome-card {
    grid-column: 1 / -1;
  }
}

@media (max-width: 640px) {
  .dashboard-header {
    padding: 15px 20px;
    flex-direction: column;
    gap: 15px;
  }
  
  .dashboard-header h1 {
    font-size: 24px;
  }
  
  .dashboard-content {
    padding: 20px;
    gap: 20px;
  }
  
  .welcome-card,
  .profile-card,
  .activity-card {
    padding: 20px;
  }
  
  .welcome-card h2 {
    font-size: 24px;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
EOF

# Create .gitignore
cat > .gitignore << 'EOF'
# Backend
backend/auth-backend
backend/*.exe
backend/*.dll
backend/*.so
backend/*.dylib
backend/.env

# Frontend
frontend/node_modules/
frontend/.pnp
frontend/.pnp.js
frontend/coverage/
frontend/build/
frontend/.DS_Store
frontend/.env.local
frontend/.env.development.local
frontend/.env.test.local
frontend/.env.production.local
frontend/npm-debug.log*
frontend/yarn-debug.log*
frontend/yarn-error.log*

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# MongoDB data (if running locally without Docker)
data/
EOF

# Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:latest
    container_name: auth-app-mongodb
    restart: always
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password123
      MONGO_INITDB_DATABASE: authapp
    volumes:
      - mongodb_data:/data/db

volumes:
  mongodb_data:
EOF

# Create README.md
cat > README.md << 'EOF'
# Auth App - React + Golang + MongoDB

Полноценное приложение с авторизацией и регистрацией пользователей.

## Технологии

- **Frontend**: React, TypeScript, React Router, Axios
- **Backend**: Golang, Gin Framework, JWT
- **Database**: MongoDB

## Функциональность

- Регистрация новых пользователей
- Авторизация с JWT токенами
- Защищенный личный кабинет
- Красивый современный UI

## Запуск проекта

### 1. Запуск MongoDB

```bash
# С помощью Docker Compose
docker-compose up -d

# Или установите MongoDB локально
```

### 2. Запуск Backend

```bash
cd backend

# Установка зависимостей (если еще не установлены)
go mod download

# Запуск сервера
go run main.go
```

Backend будет доступен на http://localhost:8080

### 3. Запуск Frontend

```bash
cd frontend

# Установка зависимостей (если еще не установлены)
npm install

# Запуск в режиме разработки
npm start
```

Frontend будет доступен на http://localhost:3000

## API Endpoints

- `POST /api/register` - Регистрация нового пользователя
- `POST /api/login` - Авторизация пользователя
- `GET /api/profile` - Получение профиля (требует авторизации)

## Переменные окружения

### Backend (.env)

```
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your-super-secret-jwt-key-change-this
PORT=8080
```

### Frontend

API URL настроен в `src/services/api.ts`

## Структура проекта

```
auth-app/
├── backend/
│   ├── config/         # Конфигурация БД
│   ├── controllers/    # Контроллеры
│   ├── middleware/     # Middleware
│   ├── models/         # Модели данных
│   ├── utils/          # Утилиты (JWT)
│   └── main.go         # Точка входа
├── frontend/
│   ├── src/
│   │   ├── components/ # React компоненты
│   │   ├── contexts/   # React Context
│   │   ├── pages/      # Страницы
│   │   └── services/   # API сервисы
│   └── package.json
└── docker-compose.yml
```

## Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены для авторизации
- CORS настроен для защиты API
- Проверка валидации на клиенте и сервере
EOF

# Create type declaration file for React
cat > frontend/src/react-app-env.d.ts << 'EOF'
/// <reference types="react-scripts" />
EOF

echo "✅ Все React компоненты созданы!"
echo "✅ Стили созданы!"
echo "✅ Все файлы готовы!"
echo ""
echo "🎉 ПРОЕКТ ПОЛНОСТЬЮ СОЗДАН!"
echo ""
echo "Теперь выполните эти команды:"
echo "1. git add ."
echo "2. git commit -m 'Initial commit: Full-stack auth application'"
echo "3. git push origin main"