#!/bin/bash

echo "🔍 Проверка открытых портов..."
echo ""

# Проверяем, какие порты слушаются
echo "📊 Порты, которые слушают контейнеры:"
netstat -tlnp 2>/dev/null | grep -E ':3000|:8080|:8081' || ss -tlnp | grep -E ':3000|:8080|:8081'

echo ""
echo "🐳 Статус контейнеров:"
docker-compose ps

echo ""
echo "🌐 Проверка доступности сервисов:"

# Backend API
echo -n "Backend API (8080): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ || echo "FAIL"

echo -n "Backend Health: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health || echo "FAIL"

echo -n "Backend Rates: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/rates || echo "FAIL"

# Frontend
echo -n "Frontend (3000): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ || echo "FAIL"

echo ""
echo "🔥 Firewall статус:"
sudo ufw status 2>/dev/null || echo "UFW не установлен"

echo ""
echo "💡 Если порты закрыты, выполните:"
echo "   sudo ufw allow 3000/tcp"
echo "   sudo ufw allow 8080/tcp"
echo "   sudo ufw allow 8081/tcp"