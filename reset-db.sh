#!/bin/bash

echo "⚠️  ВНИМАНИЕ: Этот скрипт полностью сбросит базу данных!"
read -p "Вы уверены? (yes/no): " -n 3 -r
echo

if [[ $REPLY =~ ^yes$ ]]
then
    # Определяем docker-compose
    if command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE="docker-compose"
    else
        DOCKER_COMPOSE="docker compose"
    fi

    echo "🛑 Останавливаем контейнеры..."
    $DOCKER_COMPOSE down

    echo "🗑️  Удаляем volume с данными БД..."
    docker volume rm gift-card-shop_db_data 2>/dev/null || true

    echo "🚀 Запускаем заново..."
    $DOCKER_COMPOSE up -d

    echo "⏳ Ждем запуска БД..."
    sleep 10

    echo "✅ База данных сброшена и контейнеры запущены!"
    echo "📋 Проверяем статус:"
    $DOCKER_COMPOSE ps
else
    echo "❌ Отменено"
fi