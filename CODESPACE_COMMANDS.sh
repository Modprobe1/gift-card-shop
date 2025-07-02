#!/bin/bash

# Команды для выполнения в GitHub Codespaces

# 1. Создаем структуру директорий
mkdir -p backend/{config,controllers,middleware,models,utils}
mkdir -p frontend/src/{components,contexts,pages,services}
mkdir -p frontend/public

# 2. Создаем backend файлы
cat > backend/go.mod << 'EOF'
module auth-backend

go 1.20

require (
    github.com/gin-gonic/gin v1.10.1
    github.com/golang-jwt/jwt/v4 v4.5.2
    github.com/joho/godotenv v1.5.1
    go.mongodb.org/mongo-driver v1.17.4
    golang.org/x/crypto v0.39.0
)
EOF

# 3. Создаем .gitignore
cat > .gitignore << 'EOF'
# Backend
backend/auth-backend
backend/*.exe
backend/.env

# Frontend
frontend/node_modules/
frontend/build/

# OS
.DS_Store
EOF

# 4. Создаем docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:latest
    container_name: auth-app-mongodb
    restart: always
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password123
      MONGO_INITDB_DATABASE: authapp
    volumes:
      - mongodb_data:/data/db

volumes:
  mongodb_data:
EOF

echo "Базовая структура создана! Теперь скопируйте остальные файлы."