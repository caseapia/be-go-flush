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
    - [8. Set Staff Rank](#8-set-staff-rank)
    - [9. Set Developer Rank](#9-set-developer-rank)
    - [10. Edit user flags](#10-edit-user-flags)
    - [11. Get ranks list](#11-get-ranks-list)
    - [12. Delete rank](#12-delete-rank)
  - [Log Routes](#log-routes)
    - [1. Get Logs](#1-get-logs)

## User Routes

### 1. Get Users List
- **Route:** `/api/users/`
- **Method:** <span>GET</span>
- **Body:** —
- **Response Example:**
```json
[
  {
    "id": 1, // uint64
    "name": "John Doe", // string
    "isBanned": false, // bool
    "isDeleted": false, // bool
    "staffRank": 0, // int
    "developerRank": 0, // int
    "flags": null, // []string
    "createdAt": "2026-02-01T10:00:00Z", // time.Time
    "updatedAt": "2026-02-01T10:00:00Z" // time.Time
  }
]
````

### 2. Get User by ID

* **Route:** `/api/user/:id`
* **Method:** <span>GET</span>
* **Body:** —
* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "isBanned": false, // bool
  "isDeleted": false, // bool
  "staffRank": 0, // int
  "developerRank": 0, // int
  "flags": null, // []string
  "createdAt": "2026-02-01T10:00:00Z", // time.Time
  "updatedAt": "2026-02-01T10:00:00Z" // time.Time
}
```

### 3. Ban User

* **Route:** `/api/admin/user/ban/:id`
* **Method:** <span style="color: #ff9a1f">PUT</span>
* **Body:** 

```json
{
  "duration": 15, // int
  "reason": "Violation of rules" // string
}
```

* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "isBanned": true, // bool
  "banReason": "Violation of rules", // *string
  "isDeleted": false, // bool
  "staffRank": 0, // int
  "developerRank": 0, // int
  "flags": null, // []string
  "createdAt": "2026-02-01T10:00:00Z", // time.Time
  "updatedAt": "2026-02-01T10:00:00Z" // time.Time
}
```

### 4. Unban User

* **Route:** `/api/admin/user/unban/:id`
* **Method:** <span style="color: #ff5631">DELETE</span>
* **Body:** —
* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "isBanned": false, // bool
  "isDeleted": false, // bool
  "staffRank": 0, // int
  "developerRank": 0, // int
  "flags": null, // []string
  "createdAt": "2026-02-01T10:00:00Z", // time.Time
  "updatedAt": "2026-02-01T10:00:00Z" // time.Time
}
```

### 5. Create User

* **Route:** `/api/admin/user/create`
* **Method:** <span style="color: #ff9a1f">PUT</span>
* **Body Example:**

```json
{
  "name": "Jane Doe", // string
}
```

* **Response Example:**

```json
{
  "id": 2, // uint64
  "name": "Jane Doe", // string
  "isBanned": false, // bool
  "isDeleted": false, // bool
  "staffRank": 0, // int
  "developerRank": 0, // int
  "flags": null, // []string
  "createdAt": "2026-02-01T10:00:00Z", // time.Time
  "updatedAt": "2026-02-01T10:00:00Z" // time.Time
}
```

### 6. Delete User
> #### The first time you use this method, the account will be soft-deleted, and the second time, it will be finally deleted.

* **Route:** `/api/admin/user/delete/:id`
* **Method:** <span style="color: #ff5631">DELETE</span>
* **Body:** —
* **Response Example:**

```json
{
  "id": 2, // int
  "isDeleted": true, // bool
  "updatedAt": "2026-02-01T10:20:00Z" // time.Time
}
```

### 7. Restore User
> #### This method works only with soft-deleted account

* **Route:** `/api/admin/user/restore/:id`
* **Method:** <span style="color: #7ecf2b">POST</span>
* **Body:** —
* **Response Example:**

```json
{
  "id": 2, // int
  "isDeleted": false, // bool
  "updatedAt": "2026-02-01T10:20:00Z" // time.Time
}
```

### 8. Set Staff Rank

* **Route:** `/api/admin/user/rank/staff/:id`
* **Method:** <span style="color: #f0e137">PATCH</span>
* **Body Example:**

```json
{
  "status": 1 // int
}
```

* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "staffRank": 1, // int
  "updatedAt": "2026-02-01T10:30:00Z" // time.Time
}
```

### 9. Set Developer Rank

* **Route:** `/api/admin/user/rank/developer/:id`
* **Method:** <span style="color: #f0e137">PATCH</span>
* **Body Example:**

```json
{
  "status": 1 // int
}
```

* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "developerRank": 1, // int
  "updatedAt": "2026-02-01T10:30:00Z" // time.Time
}
```

### 10. Edit user flags
* **Route:** `/api/admin/user/flags/edit/:id`
* **Method:** <span style="color: #f0e137">PATCH</span>
* **Body Example:**

```json
{
  "flags": ["DEV"] // []string
}
```

* **Response Example:**

```json
{
  "id": 1, // uint64
  "name": "John Doe", // string
  "isBanned": false, // bool
  "isDeleted": false, // bool
  "staffRank": 0, // int
  "developerRank": 0, // int
  "flags": ["DEV"], // []string
  "createdAt": "2026-02-01T10:00:00Z", // time.Time
  "updatedAt": "2026-02-01T10:00:00Z" // time.Time
}
```

### 11. Get ranks list
* **Route:** `/api/admin/rank/list`
* **Method:** GET
* **Body Example:** —
* **Response Example:**

```json
[
	{
		"id": 0, // int64
		"name": "None", // string
		"color": "#ffffff", // string
		"flags": null, // []string
	},
	{
		"id": 1, // int64
		"name": "Tester", // string
		"color": "#660000", // string
		"flags": [
			"REGISTERAPPLICATIONS",
			"TICKETS",
		], // []string
	},
]
```

### 12. Delete rank
* **Route:** `/api/admin/rank/delete/:id`
* **Method:** <span style="color: #ff5631">DELETE</span>
* **Body Example:** —
* **Response Example:**

```true // bool```

---

## Log Routes

### 1. Get Logs

* **Route:** `/logs/:type`
* **Method:** GET
* **Body:** —
* **Response Example:**

```json
[
  {
    "id": 1, // int
    "date": "2026-02-01T10:00:00Z", // time.Time
    "adminName": "John Doe", // string
    "adminId": 0, // uint64
    "userName": "Jane Doe", // string
    "userId": 1, // uint64
    "additionalInfo": "Reason: Violation of rules", // string
    "action": "has banned", // string
  }
]
```
