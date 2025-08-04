# E-commerce API Documentation

This document provides a comprehensive overview of all available API endpoints for the e-commerce application.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT-based authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Middleware Types

- **Public**: No authentication required
- **User**: Requires valid user authentication (`UserMiddleware`)
- **Admin**: Requires admin privileges (`AdminMiddleware`)
- **JWT**: Requires JWT token authentication (`JWTMiddleware`)

---

## Authentication Routes

### Public Authentication
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| POST | `/register` | Register a new user | None |
| POST | `/login` | User login | None |

### Protected Authentication
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| POST | `/auth/logout` | User logout | JWT |
| POST | `/auth/reset-password` | Reset user password | JWT |

---

## Product Routes

### User Accessible (Read-only)
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| GET | `/products` | Get all products | User |
| GET | `/products/:id` | Get product by ID | User |

### Admin Only (Full CRUD)
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| POST | `/products` | Create new product | Admin |
| PUT | `/products/:id` | Update product | Admin |
| DELETE | `/products/:id` | Delete product | Admin |
| PATCH | `/products/:id/stock` | Update product stock | Admin |

---

## Cart Routes

### Cart Items Management
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| POST | `/cart-items` | Add item to cart | User |
| GET | `/cart-items/:id` | Get specific cart item | User |
| PUT | `/cart-items/:id` | Update cart item | User |
| DELETE | `/cart-items/:id` | Remove item from cart | User |

### Cart Operations
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| GET | `/carts/:cartId/items` | Get all items in a cart | User |

---

## Order Routes

### Order Management
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| POST | `/api/v1/orders` | Create a new order | TBD* |
| GET | `/api/v1/orders` | Get all orders (with pagination/filtering) | TBD* |
| GET | `/api/v1/orders/:id` | Get order by ID | TBD* |
| GET | `/api/v1/orders/number/:orderNumber` | Get order by order number | TBD* |
| PUT | `/api/v1/orders/:id` | Update order | TBD* |
| PATCH | `/api/v1/orders/:id/status` | Update order status | TBD* |
| DELETE | `/api/v1/orders/:id` | Delete order | TBD* |

### User-Specific Orders
| Method | Endpoint | Description | Middleware |
|--------|----------|-------------|------------|
| GET | `/api/v1/users/:userID/orders` | Get all orders for a user | TBD* |

*Note: Order routes don't specify middleware in the provided code - consider adding appropriate authentication

---

## Request/Response Examples

### User Registration
```bash
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

### User Login
```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Get All Products
```bash
GET /products
Authorization: Bearer <jwt-token>
```

### Create Product (Admin)
```bash
POST /products
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description",
  "price": 29.99,
  "stock": 100,
  "category": "electronics"
}
```

### Add Item to Cart
```bash
POST /cart-items
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "productId": "product-uuid",
  "quantity": 2,
  "cartId": "cart-uuid"
}
```

### Create Order
```bash
POST /api/v1/orders
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "items": [
    {
      "productId": "product-uuid",
      "quantity": 2,
      "price": 29.99
    }
  ],
  "shippingAddress": {
    "street": "123 Main St",
    "city": "New York",
    "zipCode": "10001"
  }
}
```

---

## Error Responses

The API returns standard HTTP status codes:

- **200**: Success
- **201**: Created
- **400**: Bad Request
- **401**: Unauthorized
- **403**: Forbidden
- **404**: Not Found
- **500**: Internal Server Error

Error response format:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

---

## Rate Limiting

API rate limiting may be implemented. Check response headers:
- `X-RateLimit-Limit`: Request limit per time window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the rate limit resets

---

## Pagination

For endpoints that return lists (like products, orders), pagination parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `page` | Page number | 1 |
| `limit` | Items per page | 10 |
| `sort` | Sort field | `created_at` |
| `order` | Sort order (asc/desc) | `desc` |

Example:
```
GET /products?page=2&limit=20&sort=name&order=asc
```

---

## Development Notes

1. **Missing Middleware**: Order routes don't specify authentication middleware - consider adding appropriate protection
2. **Inconsistent Routing**: Some routes use `/api/v1` prefix while others don't - consider standardizing
3. **CORS**: Ensure CORS is properly configured for frontend integration
4. **Validation**: Implement request validation for all POST/PUT endpoints
5. **Logging**: Add request/response logging for debugging and monitoring

---

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set up environment variables
4. Run the application: `go run main.go`
5. The API will be available at `http://localhost:8080`

For questions or issues, please refer to the project documentation or contact the de