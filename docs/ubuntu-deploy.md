# 🚀 Деплой криптообменника на Ubuntu 22.04

Подробная инструкция по развертыванию проекта на сервере Ubuntu 22.04 с использованием Docker.

## 📋 Предварительные требования

- Ubuntu 22.04 LTS
- Минимум 2GB RAM, 20GB диска
- Права sudo
- Интернет соединение

## 🛠 1. Подготовка сервера

### Обновление системы
```bash
sudo apt update && sudo apt upgrade -y
```

### Установка основных пакетов
```bash
sudo apt install -y curl wget git htop nano ufw
```

### Настройка файрвола
```bash
# Разрешаем SSH
sudo ufw allow ssh

# Разрешаем порты для приложения
sudo ufw allow 80      # HTTP
sudo ufw allow 443     # HTTPS
sudo ufw allow 3000    # Frontend
sudo ufw allow 8080    # Backend API

# Включаем файрвол
sudo ufw enable
```

## 🐳 2. Установка Docker

### Удаление старых версий (если есть)
```bash
sudo apt remove docker docker-engine docker.io containerd runc
```

### Установка Docker
```bash
# Добавляем официальный GPG ключ Docker
sudo apt update
sudo apt install ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Добавляем репозиторий Docker
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Устанавливаем Docker Engine
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### Настройка Docker для текущего пользователя
```bash
# Добавляем пользователя в группу docker
sudo usermod -aG docker $USER

# Перезагружаемся или выполняем
newgrp docker

# Проверяем установку
docker --version
docker compose version
```

## 📥 3. Клонирование проекта

### Клонирование с GitHub (если проект в репозитории)
```bash
cd /home/$USER
git clone https://github.com/your-username/crypto-exchange.git
cd crypto-exchange
```

### Или создание проекта вручную
```bash
mkdir -p /home/$USER/crypto-exchange
cd /home/$USER/crypto-exchange

# Загружаем архив проекта
wget your-project-archive.zip
unzip your-project-archive.zip
```

## ⚙️ 4. Настройка переменных окружения

### Создание .env файла для production
```bash
cd /home/$USER/crypto-exchange

# Создаем .env файл
cat > .env << 'EOF'
# Database Configuration
DB_HOST=database
DB_PORT=3306
DB_USER=exchange_user
DB_PASSWORD=your_secure_password_here
DB_NAME=crypto_exchange

# Server Configuration
PORT=8080
GIN_MODE=release
LOG_LEVEL=info

# External API Configuration
COINGECKO_URL=https://api.coingecko.com/api/v3
API_UPDATE_INTERVAL=60
REQUEST_TIMEOUT=10

# Frontend Configuration
REACT_APP_API_URL=http://your-server-ip:8080/api/v1
EOF
```

### Замените переменные на ваши значения:
```bash
# Замените your_secure_password_here на сильный пароль
nano .env

# Замените your-server-ip на IP вашего сервера
sed -i 's/your-server-ip/ВАША_IP_АДРЕС/g' .env
```

## 🔧 5. Обновление docker-compose.yml для production

```bash
# Создаем production версию docker-compose
cat > docker-compose.prod.yml << 'EOF'
version: '3.8'

services:
  # MariaDB Database
  database:
    image: mariadb:11
    container_name: crypto_exchange_db
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - db_data:/var/lib/mysql
      - ./database/init:/docker-entrypoint-initdb.d
    networks:
      - crypto_network
    restart: unless-stopped

  # Backend Go API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: crypto_exchange_backend
    environment:
      - DB_HOST=database
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - GIN_MODE=release
      - PORT=8080
      - LOG_LEVEL=info
    ports:
      - "8080:8080"
    depends_on:
      - database
    networks:
      - crypto_network
    restart: unless-stopped

  # Frontend React App
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: crypto_exchange_frontend
    environment:
      - REACT_APP_API_URL=http://localhost:8080/api/v1
    ports:
      - "3000:3000"
    depends_on:
      - backend
    networks:
      - crypto_network
    restart: unless-stopped

volumes:
  db_data:

networks:
  crypto_network:
    driver: bridge
EOF
```

## 🚀 6. Запуск приложения

### Сборка и запуск всех сервисов
```bash
cd /home/$USER/crypto-exchange

# Сборка образов
docker compose -f docker-compose.prod.yml build

# Запуск всех сервисов
docker compose -f docker-compose.prod.yml up -d

# Проверка статуса
docker compose -f docker-compose.prod.yml ps
```

### Просмотр логов
```bash
# Все логи
docker compose -f docker-compose.prod.yml logs -f

# Логи конкретного сервиса
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f database
```

## 🔍 7. Проверка работоспособности

### Проверка портов
```bash
# Проверяем, что порты слушаются
sudo netstat -tlnp | grep -E ':3000|:8080|:3306'

# Или с помощью ss
ss -tulpn | grep -E ':3000|:8080|:3306'
```

### Проверка API
```bash
# Health check
curl http://localhost:8080/health

# Получение валют
curl http://localhost:8080/api/v1/currencies

# Проверка frontend
curl http://localhost:3000
```

### Проверка базы данных
```bash
# Подключение к контейнеру MariaDB
docker exec -it crypto_exchange_db mysql -u exchange_user -p crypto_exchange

# Проверка таблиц
SHOW TABLES;
SELECT * FROM currencies;
exit
```

## 🌐 8. Настройка доменного имени (опционально)

### Установка Nginx как reverse proxy
```bash
sudo apt install nginx -y

# Создаем конфигурацию сайта
sudo cat > /etc/nginx/sites-available/crypto-exchange << 'EOF'
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Активируем конфигурацию
sudo ln -s /etc/nginx/sites-available/crypto-exchange /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Установка SSL сертификата (Let's Encrypt)
```bash
# Установка Certbot
sudo apt install snapd
sudo snap install core; sudo snap refresh core
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot

# Получение SSL сертификата
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Автообновление сертификата
sudo crontab -e
# Добавьте строку:
# 0 12 * * * /usr/bin/certbot renew --quiet
```

## 📊 9. Мониторинг и обслуживание

### Создание скриптов для управления
```bash
# Скрипт для перезапуска
cat > ~/restart-exchange.sh << 'EOF'
#!/bin/bash
cd /home/$USER/crypto-exchange
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d
echo "Crypto Exchange перезапущен"
EOF

chmod +x ~/restart-exchange.sh
```

### Мониторинг ресурсов
```bash
# Установка htop для мониторинга
sudo apt install htop

# Мониторинг Docker контейнеров
docker stats

# Проверка места на диске
df -h

# Просмотр использования памяти
free -h
```

### Резервное копирование
```bash
# Создание backup скрипта
cat > ~/backup-exchange.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/home/$USER/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup базы данных
docker exec crypto_exchange_db mysqldump -u exchange_user -p$DB_PASSWORD crypto_exchange > $BACKUP_DIR/db_backup_$DATE.sql

# Backup конфигурации
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz /home/$USER/crypto-exchange/.env /home/$USER/crypto-exchange/docker-compose.prod.yml

echo "Backup создан: $BACKUP_DIR"
EOF

chmod +x ~/backup-exchange.sh

# Автоматический backup каждый день в 2:00
(crontab -l 2>/dev/null; echo "0 2 * * * /home/$USER/backup-exchange.sh") | crontab -
```

## 🔧 10. Управление сервисами

### Основные команды
```bash
cd /home/$USER/crypto-exchange

# Остановка всех сервисов
docker compose -f docker-compose.prod.yml down

# Запуск
docker compose -f docker-compose.prod.yml up -d

# Перезапуск конкретного сервиса
docker compose -f docker-compose.prod.yml restart backend

# Просмотр логов
docker compose -f docker-compose.prod.yml logs -f --tail=100

# Обновление образов
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d

# Очистка неиспользуемых образов
docker system prune -a
```

## 🆘 11. Решение проблем

### Проблемы с портами
```bash
# Проверка занятых портов
sudo lsof -i :3000
sudo lsof -i :8080

# Освобождение порта
sudo kill -9 $(sudo lsof -t -i:8080)
```

### Проблемы с Docker
```bash
# Перезапуск Docker
sudo systemctl restart docker

# Очистка Docker
docker system prune -a

# Проверка логов Docker
sudo journalctl -u docker.service
```

### Проблемы с базой данных
```bash
# Полная пересборка базы данных
docker compose -f docker-compose.prod.yml down -v
docker compose -f docker-compose.prod.yml up -d
```

## ✅ Финальная проверка

После завершения установки проверьте:

1. **Frontend доступен:** http://your-server-ip:3000
2. **API работает:** http://your-server-ip:8080/health
3. **База данных:** Контейнер запущен и принимает соединения
4. **Логи:** Нет критических ошибок в логах

```bash
# Итоговая проверка всех сервисов
docker compose -f docker-compose.prod.yml ps
curl -s http://localhost:8080/health | jq
curl -s http://localhost:3000 | grep -o '<title>.*</title>'
```

## 🎉 Готово!

Ваш криптообменник успешно развернут на Ubuntu 22.04!

**Адреса сервисов:**
- **Frontend:** http://your-server-ip:3000
- **Backend API:** http://your-server-ip:8080
- **Документация API:** http://your-server-ip:8080/health

**Следующие шаги:**
1. Настройте доменное имя
2. Установите SSL сертификат
3. Настройте мониторинг
4. Создайте резервные копии