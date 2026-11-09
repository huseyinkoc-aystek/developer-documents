# API Gateway Route Catalog

**Audience:** Third-party integrators, DevOps engineers, API consumers  
**Owner:** Backend Team  
**Last Updated:** 2025-01-20

---

## Overview

This document catalogs all public API endpoints exposed through the API Gateway for third-party integrations. All routes require OAuth2 Client Credentials authentication and are subject to rate limiting and daily quota enforcement.

---

## Authentication

All requests must include a valid OAuth2 access token:

```http
Authorization: Bearer <access_token>
```

### Obtaining an Access Token

```http
POST /oauth/token
Content-Type: application/json

{
  "grant_type": "client_credentials",
  "client_id": "your_client_id",
  "client_secret": "your_client_secret"
}
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

---

## Rate Limits and Quotas

All API Gateway routes enforce:

- **RPS Limit:** 10 requests per second per API client
- **Daily Quota:** 100,000 requests per day per API client

When limits are exceeded, the API returns `HTTP 429` with retry information:

```json
{
  "error": "rate_limit_exceeded",
  "retry_after": 1
}
```

Or for quota exhaustion:

```json
{
  "error": "daily_quota_exceeded",
  "quota_resets_at": "2025-01-21T00:00:00Z"
}
```

---

## Routes

### Products

#### 1. List Products (Public)

**Endpoint:** `GET /products`  
**Required Scopes:** `read:products`  
**Description:** Retrieve a paginated list of public products.

**Query Parameters:**

| Parameter | Type   | Required | Description                  |
|-----------|--------|----------|------------------------------|
| page      | int    | No       | Page number (default: 1)      |
| limit     | int    | No       | Items per page (default: 50)  |
| category  | string | No       | Filter by category slug       |

**Example Request:**

```http
GET /products?page=1&limit=20&category=electronics
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Example Response (200 OK):**

```json
{
  "data": [
    {
      "id": "67890abcdef1234567890abc",
      "sku": "PROD-001",
      "name": "Wireless Mouse",
      "slug": "wireless-mouse",
      "price": 29.99,
      "currency": "USD",
      "stock": 150,
      "category": {
        "id": "12345abcdef67890abcdef12",
        "name": "Electronics",
        "slug": "electronics"
      },
      "created_at": "2024-11-15T10:30:00Z",
      "updated_at": "2024-12-01T14:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "pages": 8
  }
}
```

---

#### 2. Get Product Details

**Endpoint:** `GET /products/:id`  
**Required Scopes:** `read:products`  
**Description:** Retrieve detailed information about a specific product.

**Path Parameters:**

| Parameter | Type   | Description         |
|-----------|--------|---------------------|
| id        | string | Product ObjectID    |

**Example Request:**

```http
GET /products/67890abcdef1234567890abc
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Example Response (200 OK):**

```json
{
  "id": "67890abcdef1234567890abc",
  "sku": "PROD-001",
  "name": "Wireless Mouse",
  "description": "Ergonomic wireless mouse with 2.4GHz connection",
  "price": 29.99,
  "currency": "USD",
  "stock": 150,
  "images": [
    "https://cdn.example.com/products/prod001_1.jpg",
    "https://cdn.example.com/products/prod001_2.jpg"
  ],
  "attributes": {
    "color": "Black",
    "weight": "85g",
    "battery": "AA"
  },
  "category": {
    "id": "12345abcdef67890abcdef12",
    "name": "Electronics",
    "slug": "electronics"
  },
  "created_at": "2024-11-15T10:30:00Z",
  "updated_at": "2024-12-01T14:20:00Z"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "product_not_found",
  "message": "Product with ID '67890abcdef1234567890abc' does not exist"
}
```

---

### Orders

#### 3. List Orders

**Endpoint:** `GET /orders`  
**Required Scopes:** `read:orders`  
**Description:** Retrieve a paginated list of orders for the authenticated API client's tenant.

**Query Parameters:**

| Parameter | Type   | Required | Description                            |
|-----------|--------|----------|----------------------------------------|
| page      | int    | No       | Page number (default: 1)                |
| limit     | int    | No       | Items per page (default: 50)            |
| status    | string | No       | Filter by order status (pending, paid, shipped, completed, cancelled) |
| from_date | string | No       | ISO 8601 date filter (created_at >= from_date) |
| to_date   | string | No       | ISO 8601 date filter (created_at <= to_date)   |

**Example Request:**

```http
GET /orders?page=1&limit=10&status=paid&from_date=2025-01-01T00:00:00Z
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Example Response (200 OK):**

```json
{
  "data": [
    {
      "id": "abc123def456789012345678",
      "order_number": "ORD-2025-0001",
      "status": "paid",
      "customer": {
        "id": "user789012345678901234567",
        "name": "John Doe",
        "email": "john.doe@example.com"
      },
      "items": [
        {
          "product_id": "67890abcdef1234567890abc",
          "sku": "PROD-001",
          "name": "Wireless Mouse",
          "quantity": 2,
          "unit_price": 29.99,
          "subtotal": 59.98
        }
      ],
      "subtotal": 59.98,
      "tax": 5.40,
      "shipping": 10.00,
      "total": 75.38,
      "currency": "USD",
      "created_at": "2025-01-05T14:22:10Z",
      "updated_at": "2025-01-05T14:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 78,
    "pages": 8
  }
}
```

---

#### 4. Get Order Details

**Endpoint:** `GET /orders/:id`  
**Required Scopes:** `read:orders`  
**Description:** Retrieve detailed information about a specific order.

**Path Parameters:**

| Parameter | Type   | Description      |
|-----------|--------|------------------|
| id        | string | Order ObjectID   |

**Example Request:**

```http
GET /orders/abc123def456789012345678
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Example Response (200 OK):**

```json
{
  "id": "abc123def456789012345678",
  "order_number": "ORD-2025-0001",
  "status": "paid",
  "customer": {
    "id": "user789012345678901234567",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1-555-0123"
  },
  "shipping_address": {
    "street": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94102",
    "country": "US"
  },
  "items": [
    {
      "product_id": "67890abcdef1234567890abc",
      "sku": "PROD-001",
      "name": "Wireless Mouse",
      "quantity": 2,
      "unit_price": 29.99,
      "subtotal": 59.98
    }
  ],
  "subtotal": 59.98,
  "tax": 5.40,
  "shipping": 10.00,
  "total": 75.38,
  "currency": "USD",
  "payment_method": "credit_card",
  "payment_status": "paid",
  "created_at": "2025-01-05T14:22:10Z",
  "updated_at": "2025-01-05T14:30:00Z"
}
```

**Error Response (403 Forbidden):**

```json
{
  "error": "access_denied",
  "message": "Order does not belong to your tenant"
}
```

---

#### 5. Create Order

**Endpoint:** `POST /orders`  
**Required Scopes:** `write:orders`  
**Description:** Create a new order on behalf of a customer.

**Request Body:**

```json
{
  "customer_email": "jane.smith@example.com",
  "items": [
    {
      "product_id": "67890abcdef1234567890abc",
      "quantity": 1
    }
  ],
  "shipping_address": {
    "street": "456 Oak Ave",
    "city": "Los Angeles",
    "state": "CA",
    "postal_code": "90001",
    "country": "US"
  },
  "payment_method": "credit_card"
}
```

**Example Request:**

```http
POST /orders
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "customer_email": "jane.smith@example.com",
  "items": [
    {
      "product_id": "67890abcdef1234567890abc",
      "quantity": 1
    }
  ],
  "shipping_address": {
    "street": "456 Oak Ave",
    "city": "Los Angeles",
    "state": "CA",
    "postal_code": "90001",
    "country": "US"
  },
  "payment_method": "credit_card"
}
```

**Example Response (201 Created):**

```json
{
  "id": "def456789012345678901234",
  "order_number": "ORD-2025-0002",
  "status": "pending",
  "customer": {
    "id": "user234567890123456789012",
    "email": "jane.smith@example.com"
  },
  "items": [
    {
      "product_id": "67890abcdef1234567890abc",
      "sku": "PROD-001",
      "name": "Wireless Mouse",
      "quantity": 1,
      "unit_price": 29.99,
      "subtotal": 29.99
    }
  ],
  "subtotal": 29.99,
  "tax": 2.70,
  "shipping": 10.00,
  "total": 42.69,
  "currency": "USD",
  "created_at": "2025-01-20T09:15:30Z"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "validation_failed",
  "details": {
    "items": "At least one item is required",
    "shipping_address.postal_code": "Invalid postal code format"
  }
}
```

**Error Response (409 Conflict):**

```json
{
  "error": "insufficient_stock",
  "message": "Product 'Wireless Mouse' has only 5 units in stock, requested 10"
}
```

---

## Route Authorization Matrix

| Endpoint           | Method | Required Scopes     | Route Permission               |
|--------------------|--------|---------------------|--------------------------------|
| `/products`        | GET    | `read:products`     | `api_gateway.products.list`    |
| `/products/:id`    | GET    | `read:products`     | `api_gateway.products.get`     |
| `/orders`          | GET    | `read:orders`       | `api_gateway.orders.list`      |
| `/orders/:id`      | GET    | `read:orders`       | `api_gateway.orders.get`       |
| `/orders`          | POST   | `write:orders`      | `api_gateway.orders.create`    |

**Middleware Enforcement Order:**

1. **Audit Trail Logging** - Records request start
2. **Authentication** - Validates OAuth2 JWT
3. **Rate Limiting** - Enforces RPS + daily quota
4. **Route Authorization** - Validates route permission + scopes
5. **Handler Execution** - Processes request
6. **Audit Trail Completion** - Records response status and latency

---

## Error Codes

| HTTP Status | Error Code                    | Description                                             |
|-------------|-------------------------------|---------------------------------------------------------|
| 400         | `validation_failed`           | Request body validation failed                          |
| 401         | `missing_authorization`       | Authorization header missing                            |
| 401         | `invalid_token`               | JWT validation failed or token expired                  |
| 401         | `invalid_token_claims`        | JWT claims malformed or missing required fields         |
| 403         | `token_type_not_allowed`      | Token type is not `api_client`                          |
| 403         | `route_not_permitted`         | API client not authorized for this route                |
| 403         | `insufficient_scope`          | Required OAuth2 scope missing                           |
| 403         | `access_denied`               | Resource does not belong to client's tenant             |
| 404         | `product_not_found`           | Product does not exist                                  |
| 404         | `order_not_found`             | Order does not exist                                    |
| 409         | `insufficient_stock`          | Product stock insufficient for order                    |
| 429         | `rate_limit_exceeded`         | RPS limit or daily quota exceeded                       |
| 500         | `internal_server_error`       | Unexpected server error                                 |
| 503         | `api_gateway_not_configured`  | API Gateway services unavailable                        |

---

## Metrics and Observability

All API Gateway requests are instrumented with Prometheus metrics:

### Metrics Exported

1. **`api_gateway_requests_total`**  
   Counter by tenant, user, method, route, status

2. **`api_gateway_request_duration_seconds`**  
   Histogram for latency tracking by tenant, user, method, route

3. **`api_gateway_auth_failures_total`**  
   Counter by failure reason

4. **`api_gateway_rate_limit_exceeded_total`**  
   Counter by tenant, user, reason (rps_exceeded, quota_exceeded)

5. **`api_gateway_route_denied_total`**  
   Counter by tenant, user, reason, method, route

### Audit Logs

Every API Gateway request is logged to the `api_audit_logs` MongoDB collection:

**Schema:**

```json
{
  "_id": "ObjectId",
  "tenant_id": "ObjectId",
  "site_id": "ObjectId",
  "api_user_id": "ObjectId",
  "method": "GET",
  "path": "/products",
  "query_params": "?page=1&limit=20",
  "request_headers": {"Authorization": "Bearer <redacted>"},
  "response_status": 200,
  "request_size": 0,
  "response_size": 3421,
  "latency_ms": 45,
  "error_message": "",
  "client_ip": "203.0.113.42",
  "user_agent": "MyIntegration/1.0",
  "created_at": "2025-01-20T10:30:45Z"
}
```

**Indexes:**

- `{ tenant_id: 1, created_at: -1 }` - For tenant-scoped audit queries
- `{ api_user_id: 1, created_at: -1 }` - For user-specific audit trail

---

## SDK and Code Examples

### cURL Example

```bash
# Obtain access token
TOKEN=$(curl -X POST https://api.example.com/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"grant_type":"client_credentials","client_id":"your_client_id","client_secret":"your_secret"}' \
  | jq -r '.access_token')

# List products
curl -X GET "https://api.example.com/products?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"

# Create order
curl -X POST https://api.example.com/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_email": "test@example.com",
    "items": [{"product_id": "67890abcdef1234567890abc", "quantity": 1}],
    "shipping_address": {"street": "123 Main St", "city": "SF", "state": "CA", "postal_code": "94102", "country": "US"},
    "payment_method": "credit_card"
  }'
```

### JavaScript/Node.js Example

```javascript
const axios = require('axios');

const API_BASE = 'https://api.example.com';
const CLIENT_ID = 'your_client_id';
const CLIENT_SECRET = 'your_client_secret';

// Obtain access token
async function getAccessToken() {
  const response = await axios.post(`${API_BASE}/oauth/token`, {
    grant_type: 'client_credentials',
    client_id: CLIENT_ID,
    client_secret: CLIENT_SECRET,
  });
  return response.data.access_token;
}

// List products
async function listProducts(token, page = 1, limit = 20) {
  const response = await axios.get(`${API_BASE}/products`, {
    headers: { Authorization: `Bearer ${token}` },
    params: { page, limit },
  });
  return response.data;
}

// Create order
async function createOrder(token, orderData) {
  const response = await axios.post(`${API_BASE}/orders`, orderData, {
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  return response.data;
}

// Usage
(async () => {
  const token = await getAccessToken();
  const products = await listProducts(token, 1, 10);
  console.log('Products:', products);

  const order = await createOrder(token, {
    customer_email: 'test@example.com',
    items: [{ product_id: '67890abcdef1234567890abc', quantity: 1 }],
    shipping_address: {
      street: '123 Main St',
      city: 'San Francisco',
      state: 'CA',
      postal_code: '94102',
      country: 'US',
    },
    payment_method: 'credit_card',
  });
  console.log('Order created:', order);
})();
```

### Python Example

```python
import requests
import json

API_BASE = 'https://api.example.com'
CLIENT_ID = 'your_client_id'
CLIENT_SECRET = 'your_client_secret'

# Obtain access token
def get_access_token():
    response = requests.post(f'{API_BASE}/oauth/token', json={
        'grant_type': 'client_credentials',
        'client_id': CLIENT_ID,
        'client_secret': CLIENT_SECRET,
    })
    return response.json()['access_token']

# List products
def list_products(token, page=1, limit=20):
    headers = {'Authorization': f'Bearer {token}'}
    params = {'page': page, 'limit': limit}
    response = requests.get(f'{API_BASE}/products', headers=headers, params=params)
    return response.json()

# Create order
def create_order(token, order_data):
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json',
    }
    response = requests.post(f'{API_BASE}/orders', headers=headers, json=order_data)
    return response.json()

# Usage
if __name__ == '__main__':
    token = get_access_token()
    products = list_products(token, page=1, limit=10)
    print('Products:', json.dumps(products, indent=2))

    order = create_order(token, {
        'customer_email': 'test@example.com',
        'items': [{'product_id': '67890abcdef1234567890abc', 'quantity': 1}],
        'shipping_address': {
            'street': '123 Main St',
            'city': 'San Francisco',
            'state': 'CA',
            'postal_code': '94102',
            'country': 'US',
        },
        'payment_method': 'credit_card',
    })
    print('Order created:', json.dumps(order, indent=2))
```

---

## Testing and Sandbox

### Sandbox Environment

A sandbox environment is available for testing integrations:

**Base URL:** `https://sandbox-api.example.com`

**Test Credentials:**

- Client ID: `sandbox_client_test123`
- Client Secret: `sandbox_secret_test456`

### Rate Limits (Sandbox)

- RPS Limit: 5 requests/second
- Daily Quota: 10,000 requests/day

### Test Data

The sandbox comes preloaded with sample products and orders for testing.

---

## Support and Feedback

For API issues, rate limit increases, or new route requests:

- **Email:** api-support@example.com
- **Documentation:** https://docs.example.com/api-gateway
- **Status Page:** https://status.example.com

---

## Change Log

| Date       | Version | Changes                                           |
|------------|---------|---------------------------------------------------|
| 2025-01-20 | 1.0.0   | Initial release with Products and Orders routes   |

---

**Next Steps:**

1. Review route catalog with product team for completeness
2. Add integration tests for multi-tenant isolation scenarios
3. Deploy to staging environment for partner testing
4. Prepare Phase 4: Frontend UI for API client management
