#!/bin/bash

echo "🚀 Начинаем развертывание Crypto Exchange..."

# Останавливаем старые контейнеры
echo "🛑 Останавливаем старые контейнеры..."
docker compose down

# Очищаем неиспользуемые образы
echo "🧹 Очищаем старые образы..."
docker system prune -f

# Собираем новые образы
echo "🔨 Собираем новые образы..."
docker compose build --no-cache

# Запускаем контейнеры
echo "▶️  Запускаем контейнеры..."
docker compose up -d

# Ждем пока все запустится
echo "⏳ Ждем запуска сервисов..."
sleep 10

# Проверяем статус
echo "✅ Проверяем статус контейнеров:"
docker compose ps

echo ""
echo "📊 Проверяем работу API:"
curl -s http://localhost:8080/api/health | jq '.' || echo "❌ API не отвечает"

echo ""
echo "🎉 Развертывание завершено!"
echo ""
echo "📍 Доступ к сервисам:"
echo "   - Frontend: http://$(curl -s ifconfig.me):3000"
echo "   - Backend API: http://$(curl -s ifconfig.me):8080"
echo "   - Adminer: http://$(curl -s ifconfig.me):8081"
echo ""
echo "📝 Логи можно посмотреть командой: docker compose logs -f"