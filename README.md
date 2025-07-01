# Cryptocurrency Exchange Platform

Платформа для обмена криптовалют, аналогичная tsunami.cash.

## 🎯 Возможности

- 🔄 Обмен USDT ↔ BTC и USDT ↔ RUB  
- 📊 Реальные курсы с CoinGecko API
- 💰 2% профит встроен в курсы
- 🎨 Современный интерфейс в стиле tsunami.cash
- 📱 Адаптивный дизайн
- 🔗 Уникальные ссылки для каждой заявки

## 🛠 Технологии

- **Frontend**: React 18 + TypeScript
- **Backend**: Go + Gin Framework  
- **Database**: MariaDB 11
- **Deployment**: Docker + Docker Compose

## 🚀 Быстрый старт

### Продакшен (рекомендуемый)

```bash
# Клонируем репозиторий
git clone https://github.com/Modprobe1/gift-card-shop.git
cd gift-card-shop

# Запускаем продакшен версию
docker compose up -d
```

### Разработка (с live reload)

```bash
# Запускаем dev версию с автоперезагрузкой
docker compose -f docker-compose.dev.yml up -d
```

## 🌐 Доступ к сервисам

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080  
- **Database Admin**: http://localhost:8081
  - Сервер: `database`
  - Пользователь: `exchange_user`
  - Пароль: `exchange_password`
  - База: `crypto_exchange`

## 📚 API Endpoints

```
GET    /api/health              - Проверка работы
GET    /api/currencies          - Список валют
GET    /api/rates               - Курсы обмена
POST   /api/orders              - Создание заявки
GET    /api/orders/:number      - Информация о заявке
```

## 📁 Структура проекта

```
crypto-exchange/
├── backend/                 # Go API сервер
│   ├── cmd/main.go         # Точка входа
│   ├── internal/           # Внутренние пакеты
│   └── Dockerfile          # Docker образ
├── frontend/               # React приложение
│   ├── src/                # Исходный код
│   ├── public/             # Статические файлы
│   └── Dockerfile          # Docker образ
├── database/               # SQL схемы
│   └── init/               # Скрипты инициализации
├── docker-compose.yml      # Продакшен
├── docker-compose.dev.yml  # Разработка
└── README.md              # Документация
```

## ⚙️ Конфигурация

### Переменные окружения

**Backend**:
```env
DB_HOST=database
DB_PORT=3306
DB_USER=exchange_user
DB_PASSWORD=exchange_password
DB_NAME=crypto_exchange
GIN_MODE=release
PORT=8080
```

**Frontend**:
```env
REACT_APP_API_URL=http://localhost:8080
```

## 🚀 Развертывание на сервере

### Ubuntu 22.04 LTS

1. **Установка Docker**:
```bash
# Обновляем систему
sudo apt update && sudo apt upgrade -y

# Устанавливаем Docker
sudo apt install docker.io docker-compose-plugin -y
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

2. **Клонируем и запускаем**:
```bash
git clone https://github.com/Modprobe1/gift-card-shop.git
cd gift-card-shop
docker compose up -d
```

3. **Настройка Nginx** (опционально):
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📊 Мониторинг

```bash
# Просмотр логов
docker compose logs -f backend
docker compose logs -f frontend

# Статус контейнеров
docker compose ps

# Остановка
docker compose down

# Перезапуск
docker compose restart
```

## 💻 Разработка

### Локальная разработка

```bash
# Backend
cd backend
go mod download
go run cmd/main.go

# Frontend  
cd frontend
npm install
npm start
```

### Тестирование API

```bash
# Проверка здоровья
curl http://localhost:8080/api/health

# Получение курсов
curl http://localhost:8080/api/rates

# Создание заявки
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USDT_TRC20","to_currency":"BTC","from_amount":100,"client_name":"Test","client_phone":"+1234567890","client_email":"test@example.com","recipient_wallet":"bc1...","recipient_details":"Test wallet"}'
```

## 📝 Лицензия

MIT License

## 🆘 Поддержка

Если у вас возникли вопросы или проблемы, создайте Issue в репозитории.