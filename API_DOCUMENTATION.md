# Jerky Vault Backend API Documentation

## Overview

Jerky Vault Backend API - это RESTful API для управления рецептами, ингредиентами, продукцией, заказами и клиентами в производстве вяленого мяса.

**Base URL:** `http://localhost:8080`

**Swagger UI:** `http://localhost:8080/swagger/index.html`

## Authentication

API использует JWT (JSON Web Token) для аутентификации. После успешного входа в систему вы получите токен, который необходимо включать в заголовок `Authorization` для всех защищенных эндпоинтов.

**Header format:** `Authorization: Bearer <your_jwt_token>`

---

## API Endpoints

### 🔐 Authentication

#### POST `/api/auth/register`
Регистрация нового пользователя.

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (201):**
```json
{
  "id": 1,
  "username": "test_user",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Пользователь уже существует
- `400` - Неверные данные запроса

---

#### POST `/api/auth/login`
Аутентификация пользователя.

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `401` - Неверные учетные данные
- `400` - Неверные данные запроса

---

### 🍳 Recipes

#### GET `/api/recipes`
Получение списка рецептов пользователя.

**Query Parameters:**
- `recipe_id` (optional) - Фильтрация по ID рецепта
- `ingredient_id` (optional) - Фильтрация по ID ингредиента

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Beef Jerky Original",
    "user_id": 1,
    "total_cost": 250.50,
    "recipe_ingredients": [
      {
        "id": 1,
        "recipe_id": 1,
        "ingredient_id": 1,
        "quantity": 1000,
        "unit": "g",
        "calculated_cost": 180.00,
        "ingredient": {
          "id": 1,
          "name": "Говядина",
          "type": "Мясо",
          "prices": [
            {
              "id": 1,
              "price": 450.00,
              "quantity": 1,
              "unit": "kg",
              "date": "2024-01-01T00:00:00Z"
            }
          ]
        }
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### GET `/api/recipes/{id}`
Получение рецепта по ID.

**Path Parameters:**
- `id` - ID рецепта

**Response (200):** Аналогично GET `/api/recipes`, но один объект

**Errors:**
- `404` - Рецепт не найден
- `400` - Неверный ID рецепта

---

#### POST `/api/recipes`
Создание нового рецепта.

**Request Body:**
```json
{
  "name": "New Recipe Name"
}
```

**Response (201):**
```json
{
  "id": 2,
  "name": "New Recipe Name",
  "user_id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

#### DELETE `/api/recipes/{id}`
Удаление рецепта.

**Path Parameters:**
- `id` - ID рецепта

**Response (200):**
```json
{
  "message": "Recipe deleted successfully"
}
```

**Errors:**
- `404` - Рецепт не найден

---

### 🥗 Ingredients

#### GET `/api/ingredients`
Получение списка ингредиентов.

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Говядина",
    "type": "Мясо",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### POST `/api/ingredients`
Создание нового ингредиента.

**Request Body:**
```json
{
  "name": "Новый ингредиент",
  "type": "Тип ингредиента"
}
```

**Response (201):**
```json
{
  "id": 2,
  "name": "Новый ингредиент",
  "type": "Тип ингредиента",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Ингредиент с таким именем уже существует

---

#### GET `/api/ingredients/check`
Проверка существования ингредиента по имени.

**Query Parameters:**
- `name` (required) - Имя ингредиента для проверки

**Response (200):**
```json
{
  "exists": true,
  "ingredient": {
    "id": 1,
    "name": "Говядина",
    "type": "Мясо"
  }
}
```

или

```json
{
  "exists": false
}
```

---

### 🔗 Recipe Ingredients

#### POST `/api/recipes/{id}/ingredients`
Добавление ингредиента к рецепту.

**Path Parameters:**
- `id` - ID рецепта

**Request Body:**
```json
{
  "ingredient_id": 1,
  "quantity": 1000,
  "unit": "g"
}
```

**Response (201):**
```json
{
  "id": 1,
  "recipe_id": 1,
  "ingredient_id": 1,
  "quantity": 1000,
  "unit": "g",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

#### DELETE `/api/recipes/{id}/ingredients/{ingredient_id}`
Удаление ингредиента из рецепта.

**Path Parameters:**
- `id` - ID рецепта
- `ingredient_id` - ID ингредиента

**Response (200):**
```json
{
  "message": "Ingredient removed from recipe successfully"
}
```

---

### 📦 Products

#### GET `/api/products`
Получение списка продуктов пользователя.

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Beef Jerky Original",
    "description": "Classic beef jerky",
    "price": 300.00,
    "cost": 180.00,
    "image": "beef_jerky.jpg",
    "package_id": 1,
    "user_id": 1,
    "package": {
      "id": 1,
      "name": "Стандартная упаковка",
      "user_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "options": [
      {
        "id": 1,
        "product_id": 1,
        "recipe_id": 1,
        "user_id": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### GET `/api/products/{id}`
Получение продукта по ID.

**Path Parameters:**
- `id` - ID продукта

**Response (200):** Аналогично GET `/api/products`, но один объект

**Errors:**
- `404` - Продукт не найден
- `400` - Неверный ID продукта

---

#### POST `/api/products`
Создание нового продукта.

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description"
}
```

**Response (201):**
```json
{
  "id": 2,
  "name": "New Product",
  "description": "Product description",
  "user_id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

#### PUT `/api/products/{id}`
Обновление продукта.

**Path Parameters:**
- `id` - ID продукта

**Request Body:**
```json
{
  "name": "Updated Product",
  "description": "Updated description"
}
```

**Response (200):**
```json
{
  "id": 1,
  "name": "Updated Product",
  "description": "Updated description",
  "user_id": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

#### DELETE `/api/products/{id}`
Удаление продукта.

**Path Parameters:**
- `id` - ID продукта

**Response (200):**
```json
{
  "message": "Product deleted successfully"
}
```

---

### 💰 Prices

#### GET `/api/prices`
Получение списка цен на ингредиенты.

**Query Parameters:**
- `ingredient_id` (optional) - Фильтрация по ID ингредиента

**Response (200):**
```json
[
  {
    "id": 1,
    "ingredient_id": 1,
    "price": 450.00,
    "quantity": 1,
    "unit": "kg",
    "date": "2024-01-01T00:00:00Z",
    "user_id": 1,
    "ingredient": {
      "id": 1,
      "name": "Говядина",
      "type": "Мясо"
    }
  }
]
```

---

#### POST `/api/prices`
Добавление новой цены для ингредиента.

**Request Body:**
```json
{
  "ingredient_id": 1,
  "price": 500.00,
  "quantity": 1,
  "unit": "kg",
  "date": "2024-01-01T00:00:00Z"
}
```

**Response (201):**
```json
{
  "id": 2,
  "ingredient_id": 1,
  "price": 500.00,
  "quantity": 1,
  "unit": "kg",
  "date": "2024-01-01T00:00:00Z",
  "user_id": 1
}
```

---

### 👥 Clients

#### GET `/api/clients`
Получение списка клиентов.

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@example.com",
    "phone": "+7-999-123-45-67",
    "address": "Москва, ул. Примерная, 1",
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### GET `/api/clients/{id}`
Получение клиента по ID.

**Path Parameters:**
- `id` - ID клиента

**Response (200):** Аналогично GET `/api/clients`, но один объект

---

#### POST `/api/clients`
Создание нового клиента.

**Request Body:**
```json
{
  "name": "Новый Клиент",
  "email": "client@example.com",
  "phone": "+7-999-000-00-00",
  "address": "Адрес клиента"
}
```

**Response (201):**
```json
{
  "id": 2,
  "name": "Новый Клиент",
  "email": "client@example.com",
  "phone": "+7-999-000-00-00",
  "address": "Адрес клиента",
  "user_id": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

#### PUT `/api/clients/{id}`
Обновление клиента.

**Path Parameters:**
- `id` - ID клиента

**Request Body:**
```json
{
  "name": "Обновленное имя",
  "email": "updated@example.com",
  "phone": "+7-999-111-11-11",
  "address": "Новый адрес"
}
```

**Response (200):** Обновленный объект клиента

---

#### DELETE `/api/clients/{id}`
Удаление клиента.

**Path Parameters:**
- `id` - ID клиента

**Response (200):**
```json
{
  "message": "Client deleted successfully"
}
```

---

### 📋 Orders

#### GET `/api/orders`
Получение списка заказов.

**Response (200):**
```json
[
  {
    "id": 1,
    "client_id": 1,
    "status": "pending",
    "user_id": 1,
    "client": {
      "id": 1,
      "name": "Иван Иванов",
      "email": "ivan@example.com",
      "phone": "+7-999-123-45-67",
      "address": "Москва, ул. Примерная, 1"
    },
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 1,
        "quantity": 5,
        "price": 300.00,
        "cost_price": 180.00,
        "product": {
          "id": 1,
          "name": "Beef Jerky Original",
          "description": "Classic beef jerky",
          "price": 300.00,
          "cost": 180.00,
          "package": {
            "id": 1,
            "name": "Стандартная упаковка"
          }
        }
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### GET `/api/orders/{id}`
Получение заказа по ID.

**Path Parameters:**
- `id` - ID заказа

**Response (200):** Аналогично GET `/api/orders`, но один объект

---

#### POST `/api/orders`
Создание нового заказа.

**Request Body:**
```json
{
  "client_id": 1,
  "status": "pending",
  "items": [
    {
      "product_id": 1,
      "quantity": 5,
      "price": 300.00,
      "cost_price": 180.00
    },
    {
      "product_id": 2,
      "quantity": 3,
      "price": 250.00
    }
  ]
}
```

**Примечания:**
- `cost_price` - необязательное поле. Если не указано, автоматически используется себестоимость из продукта
- `status` - необязательное поле. По умолчанию устанавливается "pending"

**Response (201):**
```json
{
  "id": 1,
  "client_id": 1,
  "status": "pending",
  "user_id": 1,
  "client": {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@example.com",
    "phone": "+7-999-123-45-67",
    "address": "Москва, ул. Примерная, 1"
  },
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "quantity": 5,
      "price": 300.00,
      "cost_price": 180.00,
      "product": {
        "id": 1,
        "name": "Beef Jerky Original",
        "description": "Classic beef jerky",
        "price": 300.00,
        "cost": 180.00,
        "package": {
          "id": 1,
          "name": "Стандартная упаковка"
        }
      }
    },
    {
      "id": 2,
      "order_id": 1,
      "product_id": 2,
      "quantity": 3,
      "price": 250.00,
      "cost_price": 150.00,
      "product": {
        "id": 2,
        "name": "Beef Jerky Spicy",
        "description": "Spicy beef jerky",
        "price": 250.00,
        "cost": 150.00
      }
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Неверные данные запроса
- `400` - Клиент не найден или не принадлежит пользователю
- `400` - Продукт не найден или не принадлежит пользователю

---

#### PUT `/api/orders/{id}`
Обновление заказа.

**Path Parameters:**
- `id` - ID заказа

**Request Body:** Аналогично POST `/api/orders`

**Response (200):** Обновленный объект заказа

---

#### PUT `/api/orders/{id}/status`
Обновление статуса заказа.

**Path Parameters:**
- `id` - ID заказа

**Request Body:**
```json
{
  "status": "completed"
}
```

**Response (200):**
```json
{
  "message": "Order status updated successfully"
}
```

---

#### DELETE `/api/orders/{id}`
Удаление заказа.

**Path Parameters:**
- `id` - ID заказа

**Response (200):**
```json
{
  "message": "Order deleted successfully"
}
```

---

### 📦 Packages

#### GET `/api/packages`
Получение списка упаковок.

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Стандартная упаковка",
    "weight": 100,
    "price": 50.00,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### POST `/api/packages`
Создание новой упаковки.

**Request Body:**
```json
{
  "name": "Новая упаковка",
  "weight": 150,
  "price": 75.00
}
```

**Response (201):** Созданный объект упаковки

---

### 📊 Dashboard

#### GET `/api/dashboard`
Получение данных для дашборда.

**Response (200):**
```json
{
  "total_recipes": 15,
  "total_products": 8,
  "total_orders": 42,
  "pending_orders": 5,
  "recent_orders": [
    {
      "id": 1,
      "client_name": "Иван Иванов",
      "total_amount": 1500.00,
      "status": "pending",
      "order_date": "2024-01-01T00:00:00Z"
    }
  ],
  "order_type_distribution": [
    {
      "type": "pending",
      "count": 5
    },
    {
      "type": "completed",
      "count": 37
    }
  ]
}
```

---

### 👤 Profile

#### POST `/api/profile/change-password`
Изменение пароля пользователя.

**Request Body:**
```json
{
  "current_password": "current_password",
  "new_password": "new_password"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Errors:**
- `400` - Неверный текущий пароль
- `400` - Неверные данные запроса

---

## Error Responses

Все ошибки возвращаются в следующем формате:

```json
{
  "error": "Error description"
}
```

### Common Error Codes:
- `400` - Bad Request (неверные данные запроса)
- `401` - Unauthorized (требуется аутентификация)
- `403` - Forbidden (недостаточно прав)
- `404` - Not Found (ресурс не найден)
- `500` - Internal Server Error (внутренняя ошибка сервера)

---

## Data Models

### User
```json
{
  "id": "uint",
  "username": "string",
  "password": "string (hashed)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Recipe
```json
{
  "id": "uint",
  "name": "string",
  "user_id": "uint",
  "total_cost": "float64 (calculated)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Ingredient
```json
{
  "id": "uint",
  "name": "string (unique)",
  "type": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Price
```json
{
  "id": "uint",
  "ingredient_id": "uint",
  "price": "float64",
  "quantity": "float64",
  "unit": "string",
  "date": "timestamp",
  "user_id": "uint"
}
```

### Client
```json
{
  "id": "uint",
  "name": "string",
  "email": "string",
  "phone": "string",
  "address": "string",
  "user_id": "uint",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Product
```json
{
  "id": "uint",
  "name": "string",
  "description": "string",
  "price": "float64",
  "cost": "float64",
  "image": "string",
  "package_id": "uint",
  "user_id": "uint",
  "package": "Package object",
  "options": "ProductOption array",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Package
```json
{
  "id": "uint",
  "name": "string",
  "user_id": "uint",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### ProductOption
```json
{
  "id": "uint",
  "product_id": "uint",
  "recipe_id": "uint",
  "user_id": "uint",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Order
```json
{
  "id": "uint",
  "client_id": "uint",
  "status": "string",
  "total_amount": "float64",
  "order_date": "timestamp",
  "delivery_date": "timestamp",
  "user_id": "uint",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## Notes

1. **Аутентификация**: Все эндпоинты кроме `/api/auth/*` требуют JWT токен
2. **Фильтрация по пользователю**: Все данные автоматически фильтруются по ID текущего пользователя
3. **Расчет стоимости**: Стоимость рецептов рассчитывается автоматически на основе последних цен ингредиентов
4. **Уникальность ингредиентов**: Система предотвращает создание дубликатов ингредиентов
5. **Soft Delete**: Модели используют soft delete (записи помечаются как удаленные, но не удаляются физически) 