# Руководство по тестированию Auth App

## Запуск приложения

### Вариант 1: Docker Compose (рекомендуется)
```bash
docker-compose up -d
```

### Вариант 2: Локально
```bash
./start.sh
```

## Тестирование функциональности

### 1. Регистрация нового пользователя

1. Откройте http://localhost:3000
2. Вы будете перенаправлены на страницу входа
3. Нажмите "Sign up" внизу формы
4. Заполните форму регистрации:
   - Username: testuser
   - Email: test@example.com
   - Password: password123
   - First Name: Test (необязательно)
   - Last Name: User (необязательно)
5. Нажмите "Sign Up"
6. После успешной регистрации вы будете перенаправлены в личный кабинет

### 2. Личный кабинет

После входа вы увидите:
- Информацию о профиле
- Кнопку "Edit Profile" для редактирования
- Статистику (дата регистрации, последнее обновление)
- Кнопку "Logout" для выхода

### 3. Редактирование профиля

1. Нажмите "Edit Profile"
2. Измените любые поля (username, имя, фамилию)
3. Нажмите "Save Changes"
4. Увидите сообщение об успешном обновлении

### 4. Выход и повторный вход

1. Нажмите "Logout"
2. Вы будете перенаправлены на страницу входа
3. Введите email и пароль
4. Нажмите "Sign In"
5. Вы снова окажетесь в личном кабинете

## API тестирование (curl)

### Регистрация
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "apiuser",
    "email": "api@example.com",
    "password": "password123",
    "first_name": "API",
    "last_name": "User"
  }'
```

### Вход
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "api@example.com",
    "password": "password123"
  }'
```

### Получение профиля (требуется токен)
```bash
TOKEN="ваш_jwt_токен_из_ответа_login"

curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Обновление профиля
```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newusername",
    "first_name": "Updated",
    "last_name": "Name"
  }'
```

## Проверка безопасности

1. **Попытка доступа без авторизации**:
   - Попробуйте открыть http://localhost:3000/dashboard напрямую
   - Вы должны быть перенаправлены на страницу входа

2. **Невалидные данные**:
   - Попробуйте зарегистрироваться с коротким паролем (< 6 символов)
   - Попробуйте войти с неверным паролем
   - Попробуйте зарегистрировать пользователя с существующим email

3. **JWT токены**:
   - Токены истекают через 24 часа
   - При выходе токен удаляется из localStorage

## Остановка приложения

### Docker Compose
```bash
docker-compose down
```

### Локальный запуск
Нажмите Ctrl+C в терминале, где запущен start.sh