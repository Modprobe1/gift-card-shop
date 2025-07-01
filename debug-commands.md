# Команды для диагностики проблем

## 1. Проверить статус контейнеров
```bash
docker compose ps
```

## 2. Проверить логи фронтенда
```bash
docker compose logs frontend
```

## 3. Проверить логи бэкенда
```bash
docker compose logs backend
```

## 4. Проверить логи базы данных
```bash
docker compose logs database
```

## 5. Проверить доступность портов
```bash
netstat -tlnp | grep :3000
netstat -tlnp | grep :8080
```

## 6. Проверить работу API
```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/rates
```

## 7. Проверить содержимое контейнера фронтенда
```bash
docker compose exec frontend ls -la /usr/share/nginx/html/
```

## 8. Проверить конфигурацию Nginx
```bash
docker compose exec frontend cat /etc/nginx/conf.d/default.conf
```

## 9. Перезапустить все сервисы
```bash
docker compose down
docker compose up -d
```

## 10. Проверить переменные окружения
```bash
docker compose exec frontend env
docker compose exec backend env
```