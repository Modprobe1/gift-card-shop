#!/bin/bash

echo "🔍 Диагностика сетевых проблем..."
echo ""

# Проверяем IP адрес сервера
echo "🌐 IP адрес сервера:"
curl -s ifconfig.me || curl -s ipinfo.io/ip || echo "Не удалось получить внешний IP"
echo ""

# Проверяем статус контейнеров
echo "🐳 Статус контейнеров:"
docker-compose ps
echo ""

# Проверяем какие порты слушаются
echo "📊 Открытые порты:"
netstat -tlnp 2>/dev/null | grep -E ':3000|:8080|:8081' || ss -tlnp | grep -E ':3000|:8080|:8081'
echo ""

# Проверяем firewall
echo "🔥 Статус firewall:"
sudo ufw status || echo "UFW не установлен"
echo ""

# Проверяем локальную доступность сервисов
echo "🧪 Локальные тесты:"

echo -n "Backend (/) : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ 2>/dev/null || echo "FAIL"

echo -n "Backend (/api/health) : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health 2>/dev/null || echo "FAIL"

echo -n "Backend (/api/v1/rates) : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/rates 2>/dev/null || echo "FAIL"

echo -n "Frontend : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ 2>/dev/null || echo "FAIL"

echo -n "Adminer : "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/ 2>/dev/null || echo "FAIL"

echo ""
echo "🔧 Исправления:"

# Открываем порты в firewall
echo "📂 Открываем порты в firewall..."
sudo ufw allow 3000/tcp 2>/dev/null || echo "Не удалось открыть порт 3000"
sudo ufw allow 8080/tcp 2>/dev/null || echo "Не удалось открыть порт 8080"  
sudo ufw allow 8081/tcp 2>/dev/null || echo "Не удалось открыть порт 8081"

# Проверяем Docker сеть
echo ""
echo "🌐 Docker сети:"
docker network ls

echo ""
echo "🔍 Проверяем Docker bridge:"
docker network inspect bridge | grep -A 10 "Containers" || echo "Проблемы с Docker bridge"

echo ""
echo "💡 Рекомендации:"
echo "1. Проверьте что сервер доступен извне:"
echo "   ping 45.12.111.208"
echo ""
echo "2. Если проблемы с портами, попробуйте:"
echo "   sudo ufw disable"
echo "   sudo systemctl stop iptables"
echo ""
echo "3. Перезапустите Docker если нужно:"
echo "   sudo systemctl restart docker"
echo "   docker-compose down && docker-compose up -d"
echo ""
echo "4. Проверьте логи:"
echo "   docker-compose logs frontend"
echo "   docker-compose logs backend"