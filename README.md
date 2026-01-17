# Jerky-vault Backend

## General Overview
Jerky-vault Backend is a REST API server written in Go that provides backend functionality for managing jerky production recipes, ingredients, products, orders, and clients. The project uses a modern technology stack and follows best development practices.

## Technology Stack
- **Programming language**: Go 1.23
- **Web framework**: Gin
- **Database**: PostgreSQL (using GORM and pgx driver)
- **API documentation**: Swagger
- **Authentication**: JWT (24-hour token expiration)
- **Containerization**: Docker
- **Security**: Rate limiting, input validation, secure JWT handling

## Project Structure
```
.
├── controllers/       # HTTP request handlers
├── database/         # Database configuration and migrations
├── docs/            # Swagger documentation
├── middleware/      # Custom middleware (JWT, rate limiting)
├── models/          # Data models
├── routes/          # API routing
├── utils/           # Utility functions
├── constants/       # Application constants
├── main.go          # Application entry point
├── Dockerfile       # Docker configuration
└── docker-compose.yaml # Docker Compose configuration
```

## Main Dependencies
- `github.com/gin-gonic/gin` - Web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `gorm.io/gorm` - ORM for database access
- `github.com/dgrijalva/jwt-go` - JWT token handling
- `github.com/swaggo/swag` - Swagger documentation generation
- `github.com/joho/godotenv` - Environment variable loader

## Configuration
The project uses the following environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- `FRONT_URL` - Frontend application URL for CORS
- `JWT_SECRET` - Secret key for JWT token signing (min 16 characters)

Environment variables can be defined:
1. Directly in the system
2. In a `.env` file (which must not be committed to version control)

## API Documentation
API documentation is available via Swagger UI at: `http://localhost:8080/swagger/index.html`

## Features

### 🔐 Security
- **JWT Authentication**: Token-based authentication with 24-hour expiration
- **Rate Limiting**: 60 requests/minute globally, 10 requests/minute for auth endpoints
- **Input Validation**: Comprehensive validation for all input data
- **JWT Secret Validation**: Server validates JWT_SECRET on startup (min 16 characters)
- **Signing Method Protection**: Guards against "none algorithm" attacks

### 🚀 Performance
- **Database Indexes**: 18 optimized indexes for fast queries
- **Auto-migration**: Automatic database schema updates on startup
- **Soft Deletes**: Data is marked as deleted, not physically removed

### ✅ Data Integrity
- **Input Validation**: All requests validated before processing
- **Transaction Support**: Multi-step operations use database transactions
- **User Isolation**: All data automatically filtered by user ID

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Authenticate and receive JWT token

### Recipes
- `GET /api/recipes` - Get all recipes
- `GET /api/recipes/:id` - Get recipe by ID
- `POST /api/recipes` - Create new recipe
- `DELETE /api/recipes/:id` - Delete recipe

### Ingredients
- `GET /api/ingredients` - Get all ingredients
- `GET /api/ingredients/check` - Check if ingredient exists by name
- `POST /api/ingredients` - Create new ingredient

### Recipe Ingredients
- `POST /api/recipes/:id/ingredients` - Add ingredient to recipe
- `DELETE /api/recipes/:id/ingredients/:ingredient_id` - Remove ingredient from recipe

### Products
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get product by ID
- `POST /api/products` - Create new product
- `PUT /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product

### Prices
- `GET /api/prices` - Get ingredient prices
- `POST /api/prices` - Add new price for ingredient

### Dashboard
- `GET /api/dashboard` - Get dashboard statistics
- `GET /api/dashboard/profit` - Get profit analysis

### Clients
- `GET /api/clients` - Get all clients
- `GET /api/clients/:id` - Get client by ID
- `POST /api/clients` - Create new client
- `PUT /api/clients/:id` - Update client
- `DELETE /api/clients/:id` - Delete client

### Orders
- `GET /api/orders` - Get all orders (sorted by creation date)
- `GET /api/orders/:id` - Get order by ID
- `POST /api/orders` - Create new order
- `PUT /api/orders/:id` - Update order
- `PUT /api/orders/:id/status` - Update order status
- `DELETE /api/orders/:id` - Delete order

### Packages
- `GET /api/packages` - Get all packages
- `POST /api/packages` - Create new package

### Profile
- `POST /api/profile/change-password` - Change user password

## Authentication

All API endpoints except `/api/auth/register` and `/api/auth/login` require authentication using a JWT Bearer token.

**Header format:** `Authorization: Bearer <your_token>`

**Login Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-01-18T14:30:00Z"
}
```

**Token expiration:** 24 hours

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Global limit:** 60 requests per minute per user/IP
- **Authentication endpoints:** 10 requests per minute per IP
- **Headers:** `X-RateLimit-Limit` included in responses

**Rate Limit Response (429 Too Many Requests):**
```json
{
  "error": "Rate limit exceeded: maximum 60 requests per minute"
}
```

## Error Responses

All errors are returned in the following format:

```json
{
  "error": "Error description"
}
```

### Common Error Codes:
- `400` - Bad Request (invalid request data)
- `401` - Unauthorized (authentication required or invalid token)
- `404` - Not Found (resource not found)
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error (server error)

## Input Validation

The API validates all input data:

**Orders:**
- `quantity` must be greater than 0
- `price` cannot be negative
- `cost_price` cannot be negative
- At least one item is required

**Products:**
- `name` is required (min 1 character)
- `price` is required and cannot be negative
- `cost` cannot be negative
- `package_id` is required

**Clients:**
- `name` is required (min 1 character)
- `surname` is required (min 1 character)

**Recipes, Ingredients, Packages:**
- `name` is required (min 1 character)

## Running the Project

### Local Run
1. Install Go 1.23 or higher
2. Create a `.env` file with the required environment variables
3. Run the application: `go run main.go`

### Running with Docker
1. Build and run with Docker Compose: `docker-compose up --build`
2. View logs: `docker-compose logs -f`

## Development

### Code Quality
1. Code must follow Go standards
2. All comments and documentation must be in English
3. All new endpoints must be documented via Swagger
4. Database changes must use GORM auto-migration
5. All user data must be filtered by `user_id`

### Generating Swagger Documentation
```bash
swag init -g main.go
```

## Security Best Practices

1. **Always validate input data** - Use binding tags and custom validation
2. **Filter by user_id** - All queries must filter by the authenticated user's ID
3. **Use transactions** - For multi-step database operations
4. **Never log sensitive data** - Don't log passwords, tokens, or personal info
5. **Validate JWT_SECRET** - Minimum 16 characters for security
6. **Handle errors gracefully** - Return generic error messages to clients

## Database Optimization

The application automatically creates 18 indexes on startup for optimal query performance:
- Orders: `user_id`, `client_id`, `status`, `created_at`, composite indexes
- Products: `user_id`, `package_id`
- And more...

Indexes are created automatically - no manual setup required.

## Limitations

- The server runs on port 8080
- JWT tokens expire after 24 hours
- Rate limiting is in-memory (resets on server restart)
- Soft delete is enabled on all models

## Recent Improvements

✅ **Security Enhancements:**
- JWT token expiration validation
- Signing method verification
- JWT_SECRET validation on startup
- Secure error messages (no user enumeration)

✅ **Performance:**
- 18 database indexes for fast queries
- Automatic index creation on startup

✅ **Data Integrity:**
- Comprehensive input validation
- Transaction support for multi-step operations
- Positive quantity validation
- Non-negative price validation

✅ **API Protection:**
- Rate limiting (60/10 requests per minute)
- User-based and IP-based limits
- Automatic cleanup of old rate limit data

✅ **Documentation:**
- All code comments in English
- Updated API documentation
- Swagger documentation

## License

[Add your license information here]