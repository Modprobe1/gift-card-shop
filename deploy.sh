#!/bin/bash

echo "🚀 Начинаем развертывание Crypto Exchange..."

# Проверяем, какой docker-compose доступен
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    echo "❌ Docker Compose не найден!"
    exit 1
fi

echo "📦 Используем: $DOCKER_COMPOSE"

# Останавливаем старые контейнеры
echo "🛑 Останавливаем старые контейнеры..."
$DOCKER_COMPOSE down

# Очищаем неиспользуемые образы
echo "🧹 Очищаем старые образы..."
docker system prune -f

# Собираем новые образы
echo "🔨 Собираем новые образы..."
$DOCKER_COMPOSE build --no-cache

# Запускаем контейнеры
echo "▶️  Запускаем контейнеры..."
$DOCKER_COMPOSE up -d

# Ждем пока все запустится
echo "⏳ Ждем запуска сервисов..."
sleep 15

# Проверяем статус
echo "✅ Проверяем статус контейнеров:"
$DOCKER_COMPOSE ps

echo ""
echo "📋 Логи backend:"
$DOCKER_COMPOSE logs --tail=20 backend

echo ""
echo "📊 Проверяем работу API:"
for i in {1..5}; do
    if curl -s http://localhost:8080/api/health | jq '.' 2>/dev/null; then
        echo "✅ API работает!"
        break
    else
        echo "⏳ Попытка $i/5: API еще не готов..."
        sleep 5
    fi
done

echo ""
echo "🎉 Развертывание завершено!"
echo ""
echo "📍 Доступ к сервисам:"
echo "   - Frontend: http://$(curl -s ifconfig.me 2>/dev/null || echo "YOUR_IP"):3000"
echo "   - Backend API: http://$(curl -s ifconfig.me 2>/dev/null || echo "YOUR_IP"):8080"
echo "   - Adminer: http://$(curl -s ifconfig.me 2>/dev/null || echo "YOUR_IP"):8081"
echo ""
echo "📝 Команды для диагностики:"
echo "   - Логи всех сервисов: $DOCKER_COMPOSE logs -f"
echo "   - Логи backend: $DOCKER_COMPOSE logs -f backend"
echo "   - Перезапуск backend: $DOCKER_COMPOSE restart backend"