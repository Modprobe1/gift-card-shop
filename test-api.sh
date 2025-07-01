#!/bin/bash

echo "🧪 Тестирование API..."
echo ""

# Тестируем корневой endpoint
echo "📍 Тест корневого endpoint:"
curl -v http://localhost:8080/ 2>&1 | head -20
echo ""

# Тестируем health endpoint
echo "💓 Тест health endpoint:"
curl -v http://localhost:8080/api/health 2>&1 | head -20
echo ""

# Тестируем rates endpoint
echo "💱 Тест rates endpoint:"
curl -v http://localhost:8080/api/v1/rates 2>&1 | head -20
echo ""

# Проверяем что слушается на портах
echo "👂 Что слушается на портах:"
netstat -tlnp | grep -E ':8080|:3000'
echo ""

# Проверяем процессы в контейнерах
echo "🐳 Процессы в контейнерах:"
docker-compose exec backend ps aux || echo "Backend недоступен"
echo ""
docker-compose exec frontend ps aux || echo "Frontend недоступен"