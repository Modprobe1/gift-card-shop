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