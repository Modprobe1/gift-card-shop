# Инструкции для загрузки на GitHub

## Вариант 1: Полная замена содержимого репозитория (рекомендуется)

Выполните эти команды в терминале из папки `/workspace/auth-app`:

```bash
# Добавляем удаленный репозиторий
git remote add origin https://github.com/Modprobe1/gift-card-shop.git

# Принудительно загружаем, заменяя все содержимое
git push -f origin main
```

## Вариант 2: Если хотите сохранить историю коммитов

```bash
# Клонируем существующий репозиторий
cd /workspace
git clone https://github.com/Modprobe1/gift-card-shop.git temp-repo
cd temp-repo

# Удаляем все файлы кроме .git
find . -not -path "./.git/*" -not -name ".git" -delete

# Копируем новые файлы
cp -r /workspace/auth-app/* .
cp /workspace/auth-app/.gitignore .

# Коммитим изменения
git add -A
git commit -m "Replace project with auth application"
git push origin main
```

## После загрузки

1. Проверьте репозиторий: https://github.com/Modprobe1/gift-card-shop
2. Обновите описание репозитория на GitHub
3. Добавьте топики: `react`, `golang`, `mongodb`, `authentication`

## Важные замечания

- Убедитесь, что у вас есть права на запись в репозиторий
- При использовании варианта 1 вся история предыдущих коммитов будет удалена
- Не забудьте обновить файл `.env` в backend с вашими настройками перед деплоем