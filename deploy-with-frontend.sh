#!/bin/bash

echo "🚀 Развертывание Crypto Exchange с интегрированным фронтендом..."

# Определяем docker-compose
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    DOCKER_COMPOSE="docker compose"
fi

echo "📦 Используем: $DOCKER_COMPOSE"

# Останавливаем контейнеры
echo "🛑 Останавливаем контейнеры..."
$DOCKER_COMPOSE down

# Собираем фронтенд локально
echo "🔨 Собираем фронтенд..."
cd frontend

# Создаем .env с правильным API URL
echo "REACT_APP_API_URL=http://45.12.111.208:8080/api/v1" > .env

# Устанавливаем зависимости и собираем
npm install
npm run build

cd ..

# Собираем backend
echo "🔨 Собираем backend..."
$DOCKER_COMPOSE build --no-cache backend

# Запускаем контейнеры
echo "▶️  Запускаем контейнеры..."
$DOCKER_COMPOSE up -d

# Ждем пока backend запустится
echo "⏳ Ждем запуска backend..."
sleep 15

# Копируем фронтенд файлы в backend контейнер
echo "📂 Копируем фронтенд в backend..."
docker cp frontend/build/. crypto_exchange_backend:/app/web/

# Перезапускаем backend чтобы подхватить статические файлы
echo "🔄 Перезапускаем backend..."
$DOCKER_COMPOSE restart backend

# Ждем еще немного
echo "⏳ Финальная настройка..."
sleep 10

# Проверяем статус
echo "✅ Проверяем статус:"
$DOCKER_COMPOSE ps

echo ""
echo "🧪 Тестируем сервисы:"

# Тест API
echo -n "API (/) : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ 2>/dev/null || echo "FAIL"

echo -n "API (/api/health) : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health 2>/dev/null || echo "FAIL"

echo -n "Frontend через backend : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/app 2>/dev/null || echo "FAIL"

echo ""
echo "🎉 Развертывание завершено!"
echo ""
echo "📍 Доступ к сервисам:"
echo "   - 🌐 Фронтенд (форма): http://45.12.111.208:8080/app"
echo "   - 🔧 Backend API: http://45.12.111.208:8080/"
echo "   - 💾 Adminer: http://45.12.111.208:8081/"
echo ""
echo "🎯 ГЛАВНОЕ: Откройте http://45.12.111.208:8080/app чтобы увидеть форму!"