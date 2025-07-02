#!/bin/bash

echo "🚀 Starting Auth App..."

# Проверка MongoDB
if ! command -v mongod &> /dev/null; then
    echo "⚠️  MongoDB не установлен. Пожалуйста, установите MongoDB или используйте Docker Compose."
    exit 1
fi

# Запуск MongoDB (если не запущен)
if ! pgrep -x "mongod" > /dev/null; then
    echo "📦 Запуск MongoDB..."
    mongod --fork --logpath /tmp/mongodb.log --dbpath /tmp/mongodb
fi

# Backend
echo "🔧 Запуск Backend..."
cd backend
if [ ! -f .env ]; then
    cp .env.example .env
    echo "📝 Создан файл .env для backend"
fi
go run main.go &
BACKEND_PID=$!

# Frontend
echo "⚛️  Запуск Frontend..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

echo "✅ Приложение запущено!"
echo "   Frontend: http://localhost:3000"
echo "   Backend: http://localhost:8080"
echo ""
echo "Нажмите Ctrl+C для остановки..."

# Ожидание и обработка завершения
trap "echo '🛑 Остановка сервисов...'; kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait