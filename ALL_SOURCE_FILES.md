# Все файлы проекта Auth App

Ниже приведены все файлы проекта. Создавайте их в GitHub в соответствующих папках.

## Структура проекта:
```
auth-app/
├── backend/
│   ├── config/
│   │   └── database.go
│   ├── controllers/
│   │   └── auth.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   └── user.go
│   ├── utils/
│   │   └── jwt.go
│   ├── main.go
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── components/
│   │   │   └── PrivateRoute.tsx
│   │   ├── contexts/
│   │   │   └── AuthContext.tsx
│   │   ├── pages/
│   │   │   ├── Login.tsx
│   │   │   ├── Register.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Auth.css
│   │   │   └── Dashboard.css
│   │   ├── services/
│   │   │   └── api.ts
│   │   ├── App.tsx
│   │   ├── App.css
│   │   └── index.tsx
│   ├── package.json
│   └── tsconfig.json
├── docker-compose.yml
├── README.md
└── .gitignore
```

Для каждого файла я покажу путь и содержимое ниже.