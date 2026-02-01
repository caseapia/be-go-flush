# go-flush API Documentation

- [go-flush API Documentation](#go-flush-api-documentation)
  - [User Routes](#user-routes)
    - [1. Get Users List](#1-get-users-list)
    - [2. Get User by ID](#2-get-user-by-id)
    - [3. Ban User](#3-ban-user)
    - [4. Unban User](#4-unban-user)
    - [5. Create User](#5-create-user)
    - [6. Delete User](#6-delete-user)
    - [7. Restore User](#7-restore-user)
    - [8. Set User Status](#8-set-user-status)
  - [Log Routes](#log-routes)
    - [1. Get Logs](#1-get-logs)

## User Routes

### 1. Get Users List
- **Route:** `/users/`
- **Method:** GET
- **Body:** —
- **Response Example:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "isDeleted": false,
    "isBanned": false,
    "status": 0,
    "createdAt": "2026-02-01T10:00:00Z",
    "updatedAt": "2026-02-01T10:00:00Z"
  }
]
````

### 2. Get User by ID

* **Route:** `/user/:id`
* **Method:** GET
* **Body:** —
* **Response Example:**

```json
{
  "id": 1,
  "name": "John Doe",
  "isDeleted": false,
  "isBanned": false,
  "status": 0,
  "createdAt": "2026-02-01T10:00:00Z",
  "updatedAt": "2026-02-01T10:00:00Z"
}
```

### 3. Ban User

* **Route:** `/user/admin/:id/ban`
* **Method:** PUT
* **Body:** 

```json
{
  "reason": "Violation of rules"
}
```

* **Response Example:**

```json
{
  "id": 1,
  "name": "John Doe",
  "isBanned": true,
  "banReason": "Violation of rules",
  "updatedAt": "2026-02-01T10:05:00Z"
}
```

### 4. Unban User

* **Route:** `/user/admin/:id/unban`
* **Method:** DELETE
* **Body:** —
* **Response Example:**

```json
{
  "id": 1,
  "name": "John Doe",
  "isBanned": false,
  "updatedAt": "2026-02-01T10:10:00Z"
}
```

### 5. Create User

* **Route:** `/user/admin/create/`
* **Method:** PUT
* **Body Example:**

```json
{
  "name": "Jane Doe",
}
```

* **Response Example:**

```json
{
  "id": 2,
  "name": "Jane Doe",
  "status": 0,
  "createdAt": "2026-02-01T10:15:00Z",
  "updatedAt": "2026-02-01T10:15:00Z"
}
```

### 6. Delete User

* **Route:** `/user/admin/:id/delete`
* **Method:** DELETE
* **Body:** —
* **Response Example:**

```json
{
  "id": 2,
  "isDeleted": true,
  "updatedAt": "2026-02-01T10:20:00Z"
}
```

### 7. Restore User

* **Route:** `/user/admin/:id/restore`
* **Method:** POST
* **Body:** —
* **Response Example:**

```json
{
  "id": 2,
  "isDeleted": false,
  "updatedAt": "2026-02-01T10:25:00Z"
}
```

### 8. Set User Status

* **Route:** `/user/admin/:id/setStatus`
* **Method:** PATCH
* **Body Example:**

```json
{
  "status": 1
}
```

* **Response Example:**

```json
{
  "id": 1,
  "name": "John Doe",
  "status": 1,
  "updatedAt": "2026-02-01T10:30:00Z"
}
```

---

## Log Routes

### 1. Get Logs

* **Route:** `/logs/`
* **Method:** GET
* **Body:** —
* **Response Example:**

```json
[
  {
    "id": 1,
    "adminId": 1,
    "userId": 2,
    "action": "ban",
    "additionalInfo": "Violation of rules",
    "createdAt": "2026-02-01T10:05:00Z"
  }
]
```
