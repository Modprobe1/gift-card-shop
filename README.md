# 🚀 Crypto Exchange Platform

Современная платформа для обмена криптовалют, построенная на React + Go + MariaDB.

## 🛠 Технологический стек

- **Frontend:** React 18, TypeScript, Tailwind CSS
- **Backend:** Go (Gin framework), REST API
- **Database:** MariaDB
- **Deployment:** Docker & Docker Compose

## 📁 Структура проекта

```
crypto-exchange/
├── frontend/           # React приложение
│   ├── src/
│   │   ├── components/ # Переиспользуемые компоненты
│   │   ├── pages/      # Страницы приложения
│   │   ├── services/   # API сервисы
│   │   ├── utils/      # Утилиты
│   │   └── assets/     # Статические файлы
│   └── public/         # Публичные файлы
├── backend/            # Go API сервер
│   ├── cmd/           # Точки входа приложения
│   ├── internal/      # Внутренняя логика
│   │   ├── handlers/  # HTTP обработчики
│   │   ├── models/    # Модели данных
│   │   ├── services/  # Бизнес логика
│   │   ├── database/  # Работа с БД
│   │   └── config/    # Конфигурация
│   └── pkg/           # Общие пакеты
├── database/          # SQL схемы и миграции
├── docker/            # Docker конфигурации
└── docs/              # Документация
```

## 🔥 Основные функции

- ✅ Обмен криптовалют в реальном времени
- ✅ Актуальные курсы валют
- ✅ Создание заявок на обмен
- ✅ Система уведомлений
- ✅ Админ панель
- ✅ История транзакций
- ✅ Интеграция с внешними API курсов

## 🚀 Быстрый старт

### Предварительные требования
- Docker и Docker Compose
- Node.js 18+ (для разработки frontend)
- Go 1.21+ (для разработки backend)

### Запуск в Docker
```bash
# Клонирование и переход в директорию
cd crypto-exchange

# Запуск всего стека
docker-compose up -d

# Приложение будет доступно на:
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Adminer (DB): http://localhost:8081
```

### Разработка

#### Backend
```bash
cd backend
go mod tidy
go run cmd/main.go
```

#### Frontend
```bash
cd frontend
npm install
npm start
```

## 📊 API Документация

API документация доступна по адресу: `http://localhost:8080/docs`

### Основные эндпоинты:
- `GET /api/rates` - Получение курсов валют
- `POST /api/orders` - Создание заявки на обмен
- `GET /api/orders/:id` - Получение заявки
- `GET /api/currencies` - Список поддерживаемых валют

## 🔒 Безопасность

- Валидация всех входящих данных
- Rate limiting
- CORS настройки
- Защита от SQL инъекций
- Логирование всех операций

## 📝 Лицензия

MIT License - см. файл [LICENSE](LICENSE)