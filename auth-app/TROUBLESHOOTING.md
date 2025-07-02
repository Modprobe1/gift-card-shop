# Устранение проблем Auth App

## Проблема: Агент/контейнер постоянно удаляется или перезапускается

### Возможные причины и решения:

1. **Несоответствие версий Go**
   - Проверьте версию в `backend/go.mod` - должна быть `go 1.22`
   - Проверьте версию в `backend/Dockerfile` - должна быть `golang:1.22-alpine`
   - Убедитесь, что нет строки `toolchain` в `go.mod`

2. **Проблемы с зависимостями**
   ```bash
   cd backend
   go mod tidy
   go mod download
   ```

3. **Проблемы с Docker кешем**
   - Используйте скрипт `rebuild.sh` для полной пересборки:
   ```bash
   ./rebuild.sh
   ```

4. **Конфликты портов**
   - Убедитесь, что порты 3000, 8080 и 27017 не заняты:
   ```bash
   netstat -tulpn | grep -E '3000|8080|27017'
   ```

5. **Проблемы с памятью**
   - Проверьте логи контейнера:
   ```bash
   docker logs auth_backend
   ```

## Команды для диагностики

```bash
# Проверить статус контейнеров
docker-compose ps

# Посмотреть логи backend
docker logs auth_backend -f

# Проверить использование ресурсов
docker stats

# Полная пересборка
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Исправленные проблемы

- ✅ Версия Go изменена с несуществующей 1.24.2 на стабильную 1.22
- ✅ Удалена строка `toolchain` из go.mod
- ✅ Синхронизированы версии Go в go.mod и Dockerfile
- ✅ Удалены volumes из backend сервиса для избежания конфликтов
- ✅ Добавлен .dockerignore для оптимизации сборки