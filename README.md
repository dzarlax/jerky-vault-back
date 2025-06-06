# Jerky-vault Backend

## Общее описание
Jerky-vault Backend - это REST API сервер, написанный на Go, который предоставляет backend функциональность для проекта Jerky-vault. Проект использует современный стек технологий и следует лучшим практикам разработки.

## Технический стек
- **Язык программирования**: Go 1.23
- **Веб-фреймворк**: Gin
- **База данных**: PostgreSQL (с использованием GORM и pgx)
- **Документация API**: Swagger
- **Аутентификация**: JWT
- **Контейнеризация**: Docker

## Структура проекта
```
.
├── controllers/     # Обработчики HTTP запросов
├── database/       # Конфигурация и миграции базы данных
├── docs/          # Swagger документация
├── models/        # Модели данных
├── routes/        # Маршрутизация API
├── middleware/    # Промежуточное ПО
├── utils/         # Вспомогательные функции
├── main.go        # Точка входа приложения
├── Dockerfile     # Конфигурация Docker
└── docker-compose.yml.example # Пример конфигурации Docker Compose
```

## Основные зависимости
- `github.com/gin-gonic/gin` - Веб-фреймворк
- `github.com/gin-contrib/cors` - Middleware для CORS
- `github.com/jackc/pgx/v5` - Драйвер PostgreSQL
- `gorm.io/gorm` - ORM для работы с базой данных
- `github.com/dgrijalva/jwt-go` - Работа с JWT токенами
- `github.com/swaggo/swag` - Генерация Swagger документации
- `github.com/joho/godotenv` - Загрузка переменных окружения

## Конфигурация
Проект использует следующие переменные окружения:
- `DATABASE_URL` - URL подключения к базе данных
- `FRONT_URL` - URL фронтенд приложения

Переменные окружения могут быть определены:
1. Напрямую в системе
2. В файле `.env` (который не должен быть в системе контроля версий)

## API Документация
API документация доступна через Swagger UI по адресу: `http://localhost:8080/swagger/*`

## API Documentation (English)

### Authentication
All API endpoints except `/auth/login` and `/auth/register` require authentication using JWT Bearer token.
Include the token in the Authorization header: `Authorization: Bearer <your_token>`

### Endpoints

#### Authentication
- `POST /auth/register`
  - Register a new user
  - Request body:
    ```json
    {
      "email": "string",
      "password": "string",
      "username": "string"
    }
    ```
  - Response: JWT token and user data

- `POST /auth/login`
  - Authenticate user
  - Request body:
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
  - Response: JWT token and user data

#### User Management
- `GET /user/profile`
  - Get current user profile
  - Requires authentication
  - Response: User profile data

- `PUT /user/profile`
  - Update user profile
  - Requires authentication
  - Request body:
    ```json
    {
      "username": "string",
      "email": "string"
    }
    ```
  - Response: Updated user profile

#### Jerky Management
- `GET /jerky`
  - Get list of all jerky items
  - Requires authentication
  - Query parameters:
    - `page`: Page number (default: 1)
    - `limit`: Items per page (default: 10)
  - Response: Paginated list of jerky items

- `POST /jerky`
  - Create new jerky item
  - Requires authentication
  - Request body:
    ```json
    {
      "name": "string",
      "description": "string",
      "type": "string",
      "weight": "number",
      "price": "number",
      "expiry_date": "string (ISO date)",
      "storage_location": "string"
    }
    ```
  - Response: Created jerky item

- `GET /jerky/{id}`
  - Get specific jerky item
  - Requires authentication
  - Response: Jerky item details

- `PUT /jerky/{id}`
  - Update jerky item
  - Requires authentication
  - Request body: Same as POST /jerky
  - Response: Updated jerky item

- `DELETE /jerky/{id}`
  - Delete jerky item
  - Requires authentication
  - Response: Success message

#### Storage Management
- `GET /storage`
  - Get list of storage locations
  - Requires authentication
  - Response: List of storage locations

- `POST /storage`
  - Create new storage location
  - Requires authentication
  - Request body:
    ```json
    {
      "name": "string",
      "description": "string",
      "temperature": "number",
      "humidity": "number"
    }
    ```
  - Response: Created storage location

### Response Formats

#### Success Response
```json
{
  "status": "success",
  "data": {
    // Response data
  },
  "message": "Optional success message"
}
```

#### Error Response
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

### Common Error Codes
- `UNAUTHORIZED`: Authentication required or invalid token
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `VALIDATION_ERROR`: Invalid request data
- `INTERNAL_ERROR`: Server error

### Rate Limiting
- API requests are limited to 100 requests per minute per IP
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Maximum requests per minute
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time until rate limit resets

### Data Types
- All dates are in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)
- Numbers are represented as floats
- Boolean values are true/false
- Strings are UTF-8 encoded

### Best Practices for Mobile Integration
1. Implement token refresh mechanism
2. Cache responses when appropriate
3. Handle offline mode
4. Implement proper error handling
5. Use pagination for large data sets
6. Implement proper retry logic for failed requests
7. Monitor rate limits
8. Implement proper data validation
9. Use proper security measures for token storage
10. Implement proper logging for debugging

## Запуск проекта

### Локальный запуск
1. Установите Go 1.23 или выше
2. Скопируйте `docker-compose.yml.example` в `docker-compose.yml` и настройте его
3. Создайте файл `.env` с необходимыми переменными окружения
4. Запустите базу данных: `docker-compose up -d db`
5. Запустите приложение: `go run main.go`

### Запуск через Docker
1. Соберите образ: `docker build -t jerky-vault-back .`
2. Запустите контейнер: `docker-compose up`

## Безопасность
- Используется JWT для аутентификации
- Настроен CORS для защиты от межсайтовых запросов
- Поддерживается только HTTPS в продакшене
- Чувствительные данные хранятся в переменных окружения

## Разработка
1. Код должен соответствовать стандартам Go
2. Все новые эндпоинты должны быть документированы через Swagger
3. Изменения в базе данных должны быть отражены в миграциях
4. Тесты должны быть написаны для новой функциональности

## Мониторинг и логирование
- Используется стандартный логгер Go
- Логируются все критические ошибки
- Доступны метрики через Swagger UI

## Деплой
Проект может быть развернут:
1. Нативно на сервере
2. В Docker контейнере
3. В облачной платформе (например, AWS, GCP, Azure)

## Ограничения
- Сервер работает на порту 8080
- Поддерживаются только определенные методы HTTP (GET, POST, PUT, DELETE, OPTIONS)
- CORS настроен только для определенных доменов 