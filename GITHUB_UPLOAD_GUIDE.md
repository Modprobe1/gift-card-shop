# 📋 Простая инструкция для загрузки через GitHub веб-интерфейс

## Вариант 1: Используйте GitHub Desktop (Самый простой)

1. Скачайте GitHub Desktop: https://desktop.github.com/
2. Клонируйте репозиторий: https://github.com/Modprobe1/gift-card-shop
3. Скопируйте все файлы из этого списка в папку репозитория
4. Commit и Push через GitHub Desktop

## Вариант 2: Загрузка через веб-интерфейс GitHub

### Шаг 1: Очистите репозиторий
1. Зайдите на https://github.com/Modprobe1/gift-card-shop
2. Удалите все старые файлы

### Шаг 2: Создайте базовые файлы
Создавайте файлы через "Add file" → "Create new file":

#### 1. `.gitignore`
```
backend/.env
backend/auth-backend
frontend/node_modules/
frontend/build/
.DS_Store
```

#### 2. `docker-compose.yml`
```yaml
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
```

#### 3. `backend/go.mod`
```go
module auth-backend

go 1.20

require (
    github.com/gin-gonic/gin v1.10.1
    github.com/golang-jwt/jwt/v4 v4.5.2
    github.com/joho/godotenv v1.5.1
    go.mongodb.org/mongo-driver v1.17.4
    golang.org/x/crypto v0.39.0
)
```

### Шаг 3: Frontend файлы

#### 1. `frontend/package.json`
```json
{
  "name": "frontend",
  "version": "0.1.0",
  "private": true,
  "dependencies": {
    "@types/node": "^16.18.11",
    "@types/react": "^18.0.27",
    "@types/react-dom": "^18.0.10",
    "@types/react-router-dom": "^5.3.3",
    "axios": "^1.7.2",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.23.1",
    "react-scripts": "5.0.1",
    "typescript": "^4.9.4"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test",
    "eject": "react-scripts eject"
  },
  "eslintConfig": {
    "extends": [
      "react-app"
    ]
  },
  "browserslist": {
    "production": [
      ">0.2%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  }
}
```

#### 2. `frontend/public/index.html`
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%PUBLIC_URL%/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta name="description" content="Auth application" />
    <title>Auth App</title>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
  </body>
</html>
```

## Вариант 3: Скачайте архив и загрузите

Я создал архивы проекта:
- `auth-app.tar.gz` - полный архив
- `auth-app.zip` - ZIP архив без node_modules

К сожалению, в веб-версии Cursor нет прямой возможности скачать файлы.

## 💡 Альтернативное решение

Если у вас есть доступ к любому компьютеру с Git:

1. Создайте временную папку
2. Скопируйте туда все файлы из моих инструкций
3. Выполните:
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/Modprobe1/gift-card-shop.git
git push -f origin main
```

## Нужна помощь?

Если нужно, я могу:
1. Показать содержимое любого конкретного файла
2. Создать Gist со всеми файлами
3. Разбить файлы на несколько частей для удобного копирования

Что предпочитаете?