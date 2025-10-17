# User Activity Service

## 🧩 Stack
- **Language:** Go 1.25+
- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM / Query Builder:** SQLX + go-sqlbuilder
- **Validation:** go-playground/validator
- **UUIDs:** google/uuid
- **CORS:** gin-contrib/cors
- **Containerization:** Docker + Docker Compose

---

## ⚙️ Environment Variables
`BACKGROUND_WORKER_INTERVAL` – how often background jobs run (e.g. 1m)  
`BACKGROUND_WORKER_PERIOD_DURATION` – time range processed each cycle (e.g. 10m)  
`DATABASE_HOST` – database host (e.g. localhost)  
`DATABASE_PORT` – database port (default 5432)  
`DATABASE_USERNAME` – PostgreSQL user (e.g. admin)  
`DATABASE_PASSWORD` – PostgreSQL password (e.g. admin)  
`DATABASE_DB_NAME` – PostgreSQL database name (e.g. activityservice)  
`APP_PORT` – HTTP port the service listens on (default 8080)  
`LOG_LEVEL` – logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`)(default `INFO`)  
`CORS_ORIGINS` – allowed origins for CORS, separated by comma (e.g. http://localhost)

---

## ️⚙️ Docker build

```bash
docker build . \
 -t image-tag 
```

---

## 📦 Dependencies
- [React Client Demo](https://github.com/nejkit/react-client-demo)
- [Postgres](https://hub.docker.com/_/postgres)


---

## 📡 Example Requests

### 🧍‍♂️ User Management

#### POST `/api/v1/users`
Registers a new user.

**Request Body**
```json
{
  "name": "John Doe"
}
```

**Response Body**
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### GET `/api/v1/users/550e8400-e29b-41d4-a716-446655440000`
Returns user details by ID.

**Response Body**
```json
{
  "username": "John Doe"
}
```

---

#### DELETE `/api/v1/users/550e8400-e29b-41d4-a716-446655440000`
Deletes a user by ID.

Response: ```200 OK```

---

### ⚙️ User Activities

#### POST `/api/v1/activities`
Saves a new activity for user.

**Request Body**
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "action": "login",
  "metadata": {
    "ip": "127.0.0.1",
    "device": "desktop"
  }
}
```

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "action": "logout",
  "metadata": {
    "page": "/home"
  }
}
```

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "action": "uploadPhoto"
}
```

Response: ```200 OK```

---

#### GET `/api/v1/activities/events/group`
**Query Parameters**  
- fromDate (optional) - Start date filter in ISO format (e.g. 2025-10-01T00:00:00Z)
- toDate   (optional) - End date filter in ISO format (e.g. 2025-10-17T23:59:59Z)

**Response Body**
```json
{
  "userActivities": {
    "550e8400-e29b-41d4-a716-446655440000": [
      {
        "actionDate": "2025-10-17T12:00:00Z",
        "action": "login",
        "metadata": { "ip": "127.0.0.1" }
      }
    ]
  }
}
```

---

#### GET `/api/v1/activities/statistics/group`
**Query Parameters**
- fromDate (optional) - Start date filter in ISO format (e.g. 2025-10-01T00:00:00Z)
- toDate   (optional) - End date filter in ISO format (e.g. 2025-10-17T23:59:59Z)

**Response Body**
```json
{
  "userActivities": {
    "550e8400-e29b-41d4-a716-446655440000": [
      {
        "fromDate": "2025-10-17T00:00:00Z",
        "toDate": "2025-10-17T10:00:00Z",
        "actionsCount": 5
      }
    ]
  }
}
```

---

#### GET `/api/v1/activities/events`
**Query Parameters**
- userId (required) - User ID in UUID format (e.g. 550e8400-e29b-41d4-a716-446655440000)
- fromDate (optional) - Start date filter in ISO format (e.g. 2025-10-01T00:00:00Z)
- toDate   (optional) - End date filter in ISO format (e.g. 2025-10-17T23:59:59Z)

**Response Body**
```json
{
  "activities": [
    {
      "actionDate": "2025-10-17T12:00:00Z",
      "action": "login",
      "metadata": { "ip": "127.0.0.1" }
    }
  ]
}
```

---

#### GET `/api/v1/activities/statistics`
**Query Parameters**
- userId (required) - User ID in UUID format (e.g. 550e8400-e29b-41d4-a716-446655440000)
- fromDate (optional) - Start date filter in ISO format (e.g. 2025-10-01T00:00:00Z)
- toDate   (optional) - End date filter in ISO format (e.g. 2025-10-17T23:59:59Z)

**Response Body**
```json
{
  "statistics": [
    {
      "fromDate": "2025-10-17T00:00:00Z",
      "toDate": "2025-10-17T10:00:00Z",
      "actionsCount": 5
    }
  ]
}
```