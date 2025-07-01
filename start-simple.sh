#!/bin/bash

echo "🚀 Простой запуск Crypto Exchange с Caddy"
echo ""

# Определяем docker-compose
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    DOCKER_COMPOSE="docker compose"
fi

echo "📦 Используем: $DOCKER_COMPOSE"

# Останавливаем все
echo "🛑 Останавливаем старые контейнеры..."
$DOCKER_COMPOSE down

# Создаем .env для фронтенда
echo "📝 Создаем .env..."
echo "REACT_APP_API_URL=http://45.12.111.208:8080/api/v1" > frontend/.env

# Собираем все заново
echo "🔨 Собираем контейнеры..."
$DOCKER_COMPOSE build --no-cache

# Запускаем
echo "▶️  Запускаем все сервисы..."
$DOCKER_COMPOSE up -d

# Ждем
echo "⏳ Ждем запуска (60 секунд)..."
sleep 60

# Проверяем статус
echo "✅ Статус контейнеров:"
$DOCKER_COMPOSE ps

echo ""
echo "🧪 Тестируем:"

# Backend
echo -n "Backend API: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ 2>/dev/null || echo "FAIL"

echo -n "Backend Health: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health 2>/dev/null || echo "FAIL"

echo -n "Backend Rates: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/rates 2>/dev/null || echo "FAIL"

# Frontend
echo -n "Frontend (Caddy): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ 2>/dev/null || echo "FAIL"

echo ""
echo "🎉 Готово!"
echo ""
echo "🎯 ОТКРОЙТЕ В БРАУЗЕРЕ:"
echo "   💥 ФОРМА: http://45.12.111.208:3000/"
echo "   🔧 API: http://45.12.111.208:8080/"
echo ""
echo "📋 Если что-то не работает:"
echo "   $DOCKER_COMPOSE logs frontend"
echo "   $DOCKER_COMPOSE logs backend"