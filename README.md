# SSO Service

gRPC-сервис аутентификации и авторизации для онлайн-магазина, написанный на Go с использованием принципов Clean Architecture.

---

## Возможности

- Регистрация пользователей
- Аутентификация пользователей
- Выдача JWT access token
- Проверка прав администратора
- Хранение данных в PostgreSQL
- gRPC API
- Graceful Shutdown
- Интеграционные тесты
- Автоматическое применение миграций
- Docker-инфраструктура

---

## Технологический стек

- Go
- gRPC
- PostgreSQL
- Docker
- golang-migrate
- JWT
- bcrypt
- slog
- Testify
- GoFakeIt
- Task

---

## Архитектура проекта

Проект построен в соответствии с принципами Clean Architecture.

```text
.
├── cmd/                    # Точки входа приложения
│   └── sso/
│
├── config/                # Конфигурационные файлы
│
├── internal/
│   ├── app/               # Инициализация приложения
│   ├── auth/
│   │   ├── delivery/      # gRPC handlers
│   │   ├── domain/        # Доменные сущности
│   │   ├── repository/    # Работа с БД
│   │   ├── usecase/       # Бизнес-логика
│   │   └── deps/          # Интерфейсы зависимостей
│   │
│   └── infrastructure/    # Конфигурация и инфраструктурный код
│
├── migrations/            # SQL-миграции
├── pkg/                   # Переиспользуемые пакеты
├── tests/                 # Интеграционные тесты
├── docker-compose.yaml
├── Taskfile.yml
└── README.md
```

---

## gRPC API

### Register

Регистрирует нового пользователя.

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse);
```

Пример запроса:

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

---

### Login

Аутентифицирует пользователя и возвращает JWT-токен.

```protobuf
rpc Login(LoginRequest) returns (LoginResponse);
```

Пример запроса:

```json
{
  "email": "user@example.com",
  "password": "password",
  "app_id": 1
}
```

---

### IsAdmin

Проверяет наличие административных прав.

```protobuf
rpc IsAdmin(IsAdminRequest) returns (IsAdminResponse);
```

---

## Запуск проекта



---

###  Клонировать репозиторий

```bash
git clone https://github.com/barnigator/sso.git
cd sso
```

---

###  Запустить приложение

```bash
task run
```


---

###  Запуск тестов

```bash
task test
```





