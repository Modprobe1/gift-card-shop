# Auth App - React + Golang + MongoDB

Полноценное приложение с авторизацией и личным кабинетом.

## Технологии

- **Frontend**: React + Vite, React Router, Axios
- **Backend**: Golang + Gin Framework
- **Database**: MongoDB
- **Auth**: JWT токены

## Функциональность

- ✅ Регистрация новых пользователей
- ✅ Вход в систему
- ✅ Защищенный личный кабинет
- ✅ Редактирование профиля
- ✅ JWT авторизация
- ✅ Валидация данных
- ✅ Красивый UI

## Быстрый старт

### Вариант 1: Docker Compose (рекомендуется)

```bash
# Клонируйте репозиторий и перейдите в папку
cd auth-app

# Запустите все сервисы
docker-compose up -d

# Приложение будет доступно:
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# MongoDB: localhost:27017
```

### Вариант 2: Локальный запуск

#### Требования
- Go 1.21+
- Node.js 18+
- MongoDB

#### Backend
```bash
cd backend

# Создайте .env файл
cp .env.example .env

# Установите зависимости
go mod download

# Запустите сервер
go run main.go
```

#### Frontend
```bash
cd frontend

# Установите зависимости
npm install

# Запустите dev сервер
npm run dev
```

## Структура проекта

```
auth-app/
├── backend/              # Golang backend
│   ├── controllers/      # HTTP контроллеры
│   ├── models/          # Модели данных
│   ├── middleware/      # Middleware (auth и др.)
│   ├── config/          # Конфигурация БД
│   ├── utils/           # Утилиты (JWT, хеширование)
│   └── main.go          # Точка входа
│
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React компоненты
│   │   ├── pages/       # Страницы приложения
│   │   ├── context/     # React Context (AuthContext)
│   │   └── App.jsx      # Главный компонент
│   └── package.json
│
└── docker-compose.yml   # Docker конфигурация
```

## API Endpoints

### Публичные
- `POST /api/register` - Регистрация
- `POST /api/login` - Вход

### Защищенные (требуют JWT)
- `GET /api/profile` - Получить профиль
- `PUT /api/profile` - Обновить профиль

## Переменные окружения

### Backend (.env)
```
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

## Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены для авторизации
- CORS настроен для frontend
- Валидация входных данных

## Разработка

### Добавление новых роутов

1. Добавьте новый метод в `controllers/`
2. Зарегистрируйте роут в `main.go`
3. Для защищенных роутов используйте `middleware.AuthMiddleware()`

### Изменение моделей

1. Обновите структуру в `models/user.go`
2. Обновите валидацию и поля в контроллерах

## Лицензия

MIT