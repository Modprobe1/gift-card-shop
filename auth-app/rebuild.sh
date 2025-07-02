#!/bin/bash

echo "🧹 Очистка и пересборка Auth App..."

# Остановка контейнеров
echo "🛑 Остановка контейнеров..."
docker-compose down

# Удаление старых образов
echo "🗑️  Удаление старых образов..."
docker rmi auth-app-backend auth-app-frontend 2>/dev/null || true

# Очистка кеша сборки
echo "🧹 Очистка кеша Docker..."
docker builder prune -f

# Пересборка с нуля
echo "🔨 Пересборка приложения..."
docker-compose build --no-cache

# Запуск
echo "🚀 Запуск приложения..."
docker-compose up -d

echo "✅ Готово! Приложение запущено."
echo "   Frontend: http://localhost:3000"
echo "   Backend: http://localhost:8080"