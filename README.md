# E-Commerce API

A RESTful e-commerce API built with Go, Gin framework, MongoDB, and Stripe for payment processing.

## Table of Contents

- [Overview](#overview)
- [Fixes Applied](#fixes-applied)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [Authentication](#authentication)
  - [Products](#products)
  - [Cart](#cart)
  - [Checkout & Orders](#checkout--orders)
  - [Admin](#admin)
  - [Webhooks](#webhooks)
- [Implementation Details](#implementation-details)

---

## Overview

This API provides a complete e-commerce backend with the following features:

- **User Authentication**: JWT-based authentication with secure password hashing
- **Product Management**: CRUD operations for products (admin-only for create/update/delete)
- **Shopping Cart**: Add, remove, and view cart items
- **Checkout**: Stripe integration for secure payment processing
- **Order Management**: Order creation and history tracking
- **Role-Based Access Control**: User and Admin roles with middleware protection

---

## Fixes Applied

The following issues were identified and fixed in the codebase:

### 1. **Empty Implementation Files**
| File | Issue | Fix |
|------|-------|-----|
| `internal/utils/jwt.go` | Empty file | Implemented JWT token generation and validation using `golang-jwt/jwt/v5` |
| `internal/utils/password.go` | Empty file | Implemented password hashing using `bcrypt` |
| `internal/controller/product.controller.go` | Empty file | Implemented full CRUD operations for products |
| `internal/controller/cart.controller.go` | Empty file | Implemented cart management endpoints |
| `internal/controller/order.controller.go` | Empty file | Implemented checkout and order endpoints |
| `internal/services/payment.service.go` | Empty file | Implemented Stripe checkout session creation and webhook handling |
| `config/env.go` | Empty file | Removed (unnecessary, functionality in `config.go`) |

### 2. **Incorrect Import Paths**
- **Issue**: All files used `ecommerce/internal/...` instead of the module path
- **Fix**: Updated all imports to use `github.com/iyanuloluwa-Miracle/ecommerce-api/...`

### 3. **Wrong Stripe Version**
- **Issue**: `config/stripe.go` imported `github.com/stripe/stripe-go/v75` but `go.mod` had `v76`
- **Fix**: Updated import to `github.com/stripe/stripe-go/v76`

### 4. **Routes File Missing Package Declaration**
- **Issue**: `routes.go` had no package declaration or imports
- **Fix**: Added proper package, imports, and enhanced with health check endpoint

### 5. **Order Service Field Name Mismatch**
- **Issue**: Used `Total` and `Status` instead of `TotalAmount` and `PaymentStatus`
- **Fix**: Corrected field names to match the Order model

### 6. **Cart Repository Issues**
- **Issue**: No upsert option for saving carts, no error handling
- **Fix**: Added upsert logic and proper error returns

### 7. **Missing JWT Library**
- **Issue**: No JWT library in `go.mod`
- **Fix**: Added `github.com/golang-jwt/jwt/v5` as direct dependency

### 8. **Missing JSON Tags on Models**
- **Issue**: Product, Cart, and Order models lacked JSON tags
- **Fix**: Added JSON tags for proper API serialization; marked password field with `json:"-"`

### 9. **Incorrect Environment Configuration**
- **Issue**: `.env` had PostgreSQL configuration but API uses MongoDB
- **Fix**: Updated with correct MongoDB configuration and all required variables

### 10. **Go Version**
- **Issue**: `go.mod` specified non-existent Go version `1.25.6`
- **Fix**: Updated to `go 1.21`

---

## Architecture

```
ecommerce-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   ├── config.go                # Environment loading
│   ├── database.go              # MongoDB connection
│   └── stripe.go                # Stripe initialization
├── internal/
│   ├── controller/              # HTTP handlers
│   │   ├── auth.controller.go
│   │   ├── cart.controller.go
│   │   ├── order.controller.go
│   │   └── product.controller.go
│   ├── middleware/              # HTTP middleware
│   │   ├── admin.middleware.go
│   │   └── jwt.middleware.go
│   ├── models/                  # Data models
│   │   ├── cart.go
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── repositories/            # Database operations
│   │   ├── cart.repo.go
│   │   ├── order.repo.go
│   │   ├── product.repo.go
│   │   └── user.repo.go
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   ├── services/                # Business logic
│   │   ├── auth.service.go
│   │   ├── cart.service.go
│   │   ├── order.service.go
│   │   ├── payment.service.go
│   │   └── product.service.go
│   └── utils/                   # Utility functions
│       ├── jwt.go
│       └── password.go
├── .env                         # Environment variables
├── go.mod                       # Go module file
└── go.sum                       # Dependency checksums
```

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB 4.4 or higher
- Stripe account (for payment processing)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/iyanuloluwa-Miracle/ecommerce-api.git
   cd ecommerce-api
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start MongoDB**
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   
   # Or start your local MongoDB instance
   ```

5. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:8080`

---

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | Database name | `ecommerce` |
| `JWT_SECRET` | Secret key for JWT signing | `your-secret-key` |
| `STRIPE_SECRET_KEY` | Stripe secret API key | `sk_test_...` |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | `whsec_...` |
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin framework mode | `debug` or `release` |
| `FRONTEND_URL` | Frontend URL for Stripe redirects | `http://localhost:3000` |

---

## API Endpoints

### Health Check

#### Check API Status
```
GET /api/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "API is running"
}
```

---

### Authentication

#### Register User
```
POST /api/auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "message": "user created successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "user already exists"
}
```

---

#### Login
```
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "invalid credentials"
}
```

---

### Products

#### Get All Products
```
GET /api/products
```

**Response (200 OK):**
```json
{
  "products": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Wireless Headphones",
      "description": "High-quality wireless headphones",
      "price": 99.99,
      "stock": 50,
      "created_at": "2026-01-22T10:00:00Z"
    }
  ]
}
```

---

#### Get Product by ID
```
GET /api/products/:id
```

**Response (200 OK):**
```json
{
  "product": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones",
    "price": 99.99,
    "stock": 50,
    "created_at": "2026-01-22T10:00:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "product not found"
}
```

---

### Cart

> **Note:** All cart endpoints require authentication. Include the JWT token in the Authorization header:
> ```
> Authorization: Bearer <token>
> ```

#### Add to Cart
```
POST /api/cart/add
```

**Request Body:**
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "quantity": 2
}
```

**Response (200 OK):**
```json
{
  "message": "item added to cart"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "insufficient stock"
}
```

---

#### Remove from Cart
```
DELETE /api/cart/remove
```

**Request Body:**
```json
{
  "product_id": "507f1f77bcf86cd799439011"
}
```

**Response (200 OK):**
```json
{
  "message": "item removed from cart"
}
```

---

#### View Cart
```
GET /api/cart
```

**Response (200 OK):**
```json
{
  "cart": {
    "id": "507f1f77bcf86cd799439012",
    "user_id": "507f1f77bcf86cd799439010",
    "items": [
      {
        "product_id": "507f1f77bcf86cd799439011",
        "quantity": 2
      }
    ]
  }
}
```

---

### Checkout & Orders

#### Checkout (Create Payment Session)
```
POST /api/checkout
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "checkout session created",
  "checkout_url": "https://checkout.stripe.com/c/pay/cs_test_...",
  "session_id": "cs_test_..."
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "cart is empty"
}
```

---

#### Get User Orders
```
GET /api/orders
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "orders": [
    {
      "id": "507f1f77bcf86cd799439013",
      "user_id": "507f1f77bcf86cd799439010",
      "items": [
        {
          "product_id": "507f1f77bcf86cd799439011",
          "quantity": 2
        }
      ],
      "total_amount": 199.98,
      "payment_status": "paid",
      "stripe_session": "cs_test_...",
      "created_at": "2026-01-22T11:00:00Z"
    }
  ]
}
```

---

### Admin

> **Note:** All admin endpoints require authentication with ADMIN role.

#### Create Product
```
POST /api/admin/products
```

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 49.99,
  "stock": 100
}
```

**Response (201 Created):**
```json
{
  "message": "product created",
  "product": {
    "id": "507f1f77bcf86cd799439014",
    "name": "New Product",
    "description": "Product description",
    "price": 49.99,
    "stock": 100,
    "created_at": "2026-01-22T12:00:00Z"
  }
}
```

---

#### Update Product
```
PUT /api/admin/products/:id
```

**Request Body:**
```json
{
  "name": "Updated Product Name",
  "price": 59.99
}
```

**Response (200 OK):**
```json
{
  "message": "product updated"
}
```

---

#### Delete Product
```
DELETE /api/admin/products/:id
```

**Response (200 OK):**
```json
{
  "message": "product deleted"
}
```

---

### Webhooks

#### Stripe Webhook
```
POST /api/webhook/stripe
```

This endpoint handles Stripe webhook events for payment processing. Configure your Stripe dashboard to send events to this URL.

**Handled Events:**
- `checkout.session.completed` - Updates order status to "paid" and decrements product stock
- `checkout.session.expired` - Updates order status to "expired"

---

## Implementation Details

### Step-by-Step Flow

#### 1. User Registration Flow
```
Client → POST /api/auth/register
       ↓
auth.controller.go (validates request)
       ↓
auth.service.go (checks for existing user, hashes password)
       ↓
user.repo.go (saves user to MongoDB)
       ↓
Response: { "message": "user created successfully" }
```

#### 2. User Login Flow
```
Client → POST /api/auth/login
       ↓
auth.controller.go (validates request)
       ↓
auth.service.go (verifies password, generates JWT)
       ↓
jwt.go (creates signed JWT with user_id and role)
       ↓
Response: { "token": "eyJ..." }
```

#### 3. Add to Cart Flow
```
Client → POST /api/cart/add (with JWT)
       ↓
jwt.middleware.go (validates JWT, extracts user_id)
       ↓
cart.controller.go (parses request)
       ↓
cart.service.go (validates product, checks stock)
       ↓
cart.repo.go (upserts cart in MongoDB)
       ↓
Response: { "message": "item added to cart" }
```

#### 4. Checkout Flow
```
Client → POST /api/checkout (with JWT)
       ↓
jwt.middleware.go (validates JWT)
       ↓
order.controller.go (gets user cart)
       ↓
payment.service.go:
  ├─ Validates all products in cart
  ├─ Calculates total amount
  ├─ Creates Stripe Checkout Session
  └─ Creates pending order in database
       ↓
Response: { "checkout_url": "https://checkout.stripe.com/..." }
       ↓
User completes payment on Stripe
       ↓
Stripe → POST /api/webhook/stripe
       ↓
payment.service.go:
  ├─ Verifies webhook signature
  ├─ Updates order status to "paid"
  ├─ Decrements product stock
  └─ Clears user's cart
```

### Security Features

1. **Password Hashing**: Uses bcrypt with default cost factor
2. **JWT Authentication**: Tokens expire after 24 hours
3. **Bearer Token Support**: Accepts both "Bearer <token>" and plain token formats
4. **Stripe Webhook Verification**: Validates webhook signatures to prevent fraud
5. **Role-Based Access Control**: Admin endpoints protected by middleware

### Database Collections

| Collection | Description |
|------------|-------------|
| `users` | User accounts with hashed passwords |
| `products` | Product catalog |
| `carts` | User shopping carts |
| `orders` | Order history and payment status |

---

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "error message here"
}
```

### Common HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

---

## License

MIT License
