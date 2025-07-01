# 🚀 Crypto Exchange Backend

REST API сервер для криптообменника, написанный на Go с использованием Gin framework.

## ✨ Возможности

- 🔄 **Реальные курсы валют** через CoinGecko API
- 💱 **Калькулятор обмена** с комиссией
- 📋 **Система заявок** с отслеживанием статуса
- 🗄️ **База данных** MariaDB с GORM ORM
- 📊 **Админ панель** для управления заявками
- 🔍 **Логирование** всех операций
- ⚡ **CORS** для интеграции с frontend

## 🛠 Технологии

- **Go 1.21+**
- **Gin** - HTTP framework
- **GORM** - ORM для работы с БД
- **MariaDB** - база данных
- **Logrus** - структурированное логирование
- **CoinGecko API** - актуальные курсы

## 📁 Структура проекта

```
backend/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── database/            # Подключение к БД
│   ├── handlers/            # HTTP обработчики
│   ├── models/              # GORM модели
│   └── services/            # Бизнес логика
├── go.mod                   # Go модули
├── go.sum                   # Зависимости
├── Dockerfile              # Docker образ
└── .env.example            # Пример переменных окружения
```

## 🚀 Быстрый старт

### 1. Установка зависимостей
```bash
go mod tidy
```

### 2. Настройка переменных окружения
```bash
cp .env.example .env
# Отредактируйте .env под ваши настройки
```

### 3. Запуск базы данных
```bash
# Из корня проекта
docker-compose up database -d
```

### 4. Запуск сервера
```bash
go run cmd/main.go
```

Сервер будет доступен на `http://localhost:8080`

## 📊 API Эндпоинты

### Основные
- `GET /health` - Проверка состояния
- `GET /api/v1/currencies` - Список валют
- `GET /api/v1/rates` - Курсы обмена
- `POST /api/v1/calculate` - Калькулятор
- `POST /api/v1/orders` - Создание заявки
- `GET /api/v1/orders/:number` - Получение заявки

### Админские
- `GET /api/v1/admin/orders` - Список заявок
- `PUT /api/v1/admin/orders/:id/status` - Изменение статуса
- `GET /api/v1/admin/statistics` - Статистика

📚 **Полная документация:** [API Documentation](../docs/api-documentation.md)

## 🧪 Тестирование

### Автоматическое тестирование
```bash
# Из корня проекта
./docs/test-api.sh
```

### Ручное тестирование с curl
```bash
# Получить валюты
curl http://localhost:8080/api/v1/currencies

# Рассчитать обмен
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USDT_TRC20","to_currency":"RUB_TBANK","from_amount":100}'

# Создать заявку
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USDT_TRC20",
    "to_currency": "RUB_TBANK",
    "from_amount": 100,
    "client_phone": "+79123456789",
    "client_name": "Тест Тестов",
    "client_email": "test@example.com",
    "recipient_wallet": "TRX7n4VkyYF7PgAwVUGdNdJg5dgrEpqhPH"
  }'
```

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|-------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `3306` |
| `DB_USER` | Пользователь БД | `exchange_user` |
| `DB_PASSWORD` | Пароль БД | `exchange_password` |
| `DB_NAME` | Имя базы данных | `crypto_exchange` |
| `PORT` | Порт сервера | `8080` |
| `GIN_MODE` | Режим Gin | `debug` |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `API_UPDATE_INTERVAL` | Интервал обновления курсов (сек) | `60` |

## 🔄 Обновление курсов

Backend автоматически обновляет курсы валют каждые 60 секунд (настраивается) из CoinGecko API:

- **BTC** ↔ **USDT/RUB**
- **ETH** ↔ **USDT/RUB** 
- **USDT** ↔ **RUB**

## 📋 Модели данных

### Currency (Валюта)
```go
type Currency struct {
    ID        uint    `json:"id"`
    Code      string  `json:"code"`     // USDT_TRC20, RUB_TBANK
    Name      string  `json:"name"`     // Tether USDT
    Symbol    string  `json:"symbol"`   // USDT
    Network   string  `json:"network"`  // TRC20
    MinAmount float64 `json:"min_amount"`
    // ...
}
```

### ExchangeOrder (Заявка)
```go
type ExchangeOrder struct {
    ID             uint      `json:"id"`
    OrderNumber    string    `json:"order_number"`
    FromCurrencyID uint      `json:"from_currency_id"`
    ToCurrencyID   uint      `json:"to_currency_id"`
    FromAmount     float64   `json:"from_amount"`
    ToAmount       float64   `json:"to_amount"`
    StatusID       uint      `json:"status_id"`
    ClientPhone    string    `json:"client_phone"`
    ClientEmail    string    `json:"client_email"`
    // ...
}
```

## 🔍 Логирование

Все операции логируются в структурированном JSON формате:

```json
{
  "level": "info",
  "msg": "Order created",
  "order_id": 1,
  "order_number": "ORD17059876541234",
  "from": "USDT_TRC20",
  "to": "RUB_TBANK",
  "amount": 100,
  "time": "2024-01-15T10:30:00Z"
}
```

## 🐳 Docker

### Сборка образа
```bash
docker build -t crypto-exchange-backend .
```

### Запуск контейнера
```bash
docker run -p 8080:8080 \
  -e DB_HOST=database \
  -e DB_PASSWORD=yourpassword \
  crypto-exchange-backend
```

## 🔧 Разработка

### Добавление новой валюты
1. Добавьте запись в таблицу `currencies`
2. Обновите метод `UpdateRatesFromAPI` в `rate_service.go`
3. Добавьте иконку валюты

### Добавление нового статуса заявки
1. Добавьте запись в таблицу `order_statuses`
2. Обновите логику в `order_service.go`

## 📝 TODO

- [ ] Аутентификация для админ панели
- [ ] Rate limiting
- [ ] Кэширование курсов валют
- [ ] WebSocket для реального времени
- [ ] Telegram бот интеграция
- [ ] Email уведомления