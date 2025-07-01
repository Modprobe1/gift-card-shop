# ⚡ Быстрый запуск криптообменника

## 🚀 Локальная разработка

### 1. Запуск backend
```bash
cd crypto-exchange

# Запуск базы данных
docker-compose up database -d

# Запуск Go сервера
cd backend
go run cmd/main.go
```
Backend будет доступен на: http://localhost:8080

### 2. Запуск frontend
```bash
# В новом терминале
cd crypto-exchange/frontend
npm start
```
Frontend будет доступен на: http://localhost:3000

## 🐳 Запуск через Docker

### Весь стек одной командой
```bash
cd crypto-exchange
docker-compose up -d
```

**Сервисы:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Adminer (DB): http://localhost:8081

## 🧪 Тестирование

### API тесты
```bash
cd crypto-exchange
./docs/test-api.sh
```

### Ручная проверка
1. Откройте http://localhost:3000
2. Выберите валюты (USDT TRC20 → RUB)
3. Введите сумму (например, 100)
4. Заполните контактные данные
5. Нажмите "Обменять сейчас"
6. Проверьте созданную заявку

## 📋 Функциональность

✅ **Готово:**
- Форма обмена с реальными курсами
- Калькулятор комиссии
- Создание заявок в БД
- Переход на страницу заявки
- API для всех операций

✅ **URL структура:**
- Главная: `/`
- Заявка: `/order/{orderNumber}/{orderId}`

✅ **Интеграция:**
- Frontend ↔ Backend API
- Реальные курсы CoinGecko
- Сохранение в MariaDB

## 🔧 Настройки

### Переменные окружения
```bash
# Backend (.env)
DB_HOST=localhost
DB_PASSWORD=your_password
API_UPDATE_INTERVAL=60

# Frontend
REACT_APP_API_URL=http://localhost:8080/api/v1
```

### Порты
- Frontend: 3000
- Backend: 8080
- Database: 3306
- Adminer: 8081

## 📊 Структура проекта

```
crypto-exchange/
├── backend/           # Go API (✅ готов)
├── frontend/          # React App (✅ готов)
├── database/          # SQL схемы (✅ готов)
├── docs/             # Документация (✅ готов)
└── docker-compose.yml # Оркестрация (✅ готов)
```

## 🎯 Следующие шаги

1. **Локальное тестирование** - запустите и протестируйте
2. **Деплой на сервер** - следуйте [ubuntu-deploy.md](ubuntu-deploy.md)
3. **Настройка домена** - настройте DNS и SSL
4. **Мониторинг** - добавьте мониторинг продакшена

---

**Поздравляем! 🎉 Ваш криптообменник готов к работе!**