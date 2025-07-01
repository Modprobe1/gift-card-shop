# 📚 API Документация

## Базовый URL
```
http://localhost:8080/api/v1
```

## 🔍 Эндпоинты

### 💰 Валюты

#### Получить все валюты
```http
GET /currencies
```

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "USDT_TRC20",
      "name": "Tether USDT",
      "symbol": "USDT",
      "network": "TRC20",
      "icon_url": "",
      "is_active": true,
      "min_amount": 10.0,
      "max_amount": null,
      "decimals": 6
    }
  ]
}
```

#### Получить валюту по коду
```http
GET /currencies/{code}
```

### 📈 Курсы обмена

#### Получить все курсы
```http
GET /rates
```

#### Получить курс между валютами
```http
GET /rates/{from}/{to}
```
Пример: `GET /rates/USDT_TRC20/RUB_TBANK`

### 🧮 Калькулятор

#### Рассчитать обмен
```http
POST /calculate
```

**Тело запроса:**
```json
{
  "from_currency": "USDT_TRC20",
  "to_currency": "RUB_TBANK", 
  "from_amount": 100.0
}
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "from_amount": 100.0,
    "to_amount": 7859.0,
    "rate": 78.597,
    "commission": 117.885,
    "commission_rate": 1.5,
    "exchange_rate": {
      "from_currency": {...},
      "to_currency": {...},
      "rate": 78.597,
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### 📋 Заявки

#### Создать заявку
```http
POST /orders
```

**Тело запроса:**
```json
{
  "from_currency": "USDT_TRC20",
  "to_currency": "RUB_TBANK",
  "from_amount": 100.0,
  "client_phone": "+79123456789",
  "client_name": "Иван Иванов",
  "client_email": "ivan@example.com",
  "client_telegram": "@ivan_telegram",
  "recipient_wallet": "TRX7n4VkyYF7PgAwVUGdNdJg5dgrEpqhPH",
  "recipient_details": "Дополнительная информация"
}
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "order_number": "ORD17059876541234",
    "from_currency": {...},
    "to_currency": {...},
    "from_amount": 100.0,
    "to_amount": 7741.115,
    "rate": 78.597,
    "status": {
      "id": 1,
      "name": "pending",
      "description": "Ожидает обработки"
    },
    "expires_at": "2024-01-15T11:00:00Z",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Получить заявку по номеру
```http
GET /orders/{order_number}
```

### 🔐 Админские эндпоинты

#### Получить список заявок
```http
GET /admin/orders?page=1&limit=20&status=pending
```

#### Обновить статус заявки
```http
PUT /admin/orders/{id}/status
```

**Тело запроса:**
```json
{
  "status": "confirmed"
}
```

#### Получить статистику
```http
GET /admin/statistics
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "total_orders": 150,
    "today_orders": 12,
    "week_orders": 85,
    "status_stats": [
      {"status_name": "pending", "count": 5},
      {"status_name": "completed", "count": 140},
      {"status_name": "cancelled", "count": 5}
    ]
  }
}
```

## 📋 Статусы заявок

| Статус | Описание |
|--------|----------|
| `pending` | Ожидает обработки |
| `confirmed` | Подтверждена |
| `processing` | В обработке |
| `completed` | Завершена |
| `cancelled` | Отменена |
| `expired` | Истекла |

## ❌ Коды ошибок

| Код | Описание |
|-----|----------|
| `400` | Неверный запрос |
| `404` | Не найдено |
| `500` | Внутренняя ошибка сервера |

## 🔧 Примеры использования

### JavaScript (Fetch API)
```javascript
// Получить курсы
const rates = await fetch('http://localhost:8080/api/v1/rates')
  .then(r => r.json());

// Создать заявку
const order = await fetch('http://localhost:8080/api/v1/orders', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    from_currency: 'USDT_TRC20',
    to_currency: 'RUB_TBANK',
    from_amount: 100,
    client_phone: '+79123456789',
    client_name: 'Иван Иванов',
    client_email: 'ivan@example.com',
    recipient_wallet: 'TRX7n4VkyYF7PgAwVUGdNdJg5dgrEpqhPH'
  })
}).then(r => r.json());
```

### cURL
```bash
# Получить валюты
curl http://localhost:8080/api/v1/currencies

# Рассчитать обмен
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USDT_TRC20","to_currency":"RUB_TBANK","from_amount":100}'
```