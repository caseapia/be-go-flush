# go-flush API Documentation

> All successful responses are wrapped in `{ "response": <data> }`.
> All error responses follow `{ "error": "<message>", "code": <status_code> }`.

---

- [go-flush API Documentation](#go-flush-api-documentation)
  - [Public Routes](#public-routes)
    - [Auth](#auth-public)
    - [Developer](#developer-public)
  - [Private Routes](#private-routes)
    - [Auth](#auth-private)
    - [User](#user-routes)
    - [Admin — User Management](#admin--user-management)
    - [Support — User](#support--user)
    - [Badges](#badges-routes)
    - [Changelog](#changelog-routes)
    - [Developer](#developer-private)
    - [Invite Codes](#invite-codes-routes)
    - [Logs](#logs-routes)
    - [Notifications](#notifications-routes)
    - [Ranks](#ranks-routes)
    - [Tickets](#tickets-routes)

---

## Public Routes

> Base path: `/api/public`
> No authentication required.

### Auth (Public)

#### 1. Register

* **Route:** `POST /api/public/auth/register`
* **Body:**

```json
{
  "login": "JohnDoe",
  "email": "john@example.com",
  "password": "securepassword",
  "invite_code": "ABC123"
}
```

* **Response** (`201`):

```json
{
  "response": {
    "user": { /* User object */ },
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

---

#### 2. Login

* **Route:** `POST /api/public/auth/login`
* **Body:**

```json
{
  "login": "JohnDoe",
  "password": "securepassword"
}
```

* **Response** (`200`):

```json
{
  "response": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

> Sets `auth_token` and `refresh_token` cookies.

---

### Developer (Public)

#### 3. Ping

* **Route:** `GET /api/public/ping`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "status": "pong"
  }
}
```

---

## Private Routes

> Base path: `/api/private`
> Requires `Authorization` header with a valid access token.
> Middleware automatically updates `last_login` and loads user rank.

---

### Auth (Private)

#### 4. Refresh Token

* **Route:** `POST /api/private/auth/refresh`
* **Body:** — (uses `refresh_token` cookie)
* **Response** (`200`):

```json
{
  "response": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

---

#### 5. Discord OAuth Redirect

* **Route:** `GET /api/private/auth/discord`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "url": "https://discord.com/oauth2/authorize?...",
    "state": "random_state_string"
  }
}
```

---

#### 6. Discord Callback (Link)

* **Route:** `POST /api/private/auth/discord/callback`
* **Body:**

```json
{
  "code": "oauth_code",
  "state": "state_param",
  "saved_state": "saved_state_from_cookie"
}
```

* **Response** (`200`):

```json
{
  "response": {
    "status": "success",
    "user": { /* User object */ }
  }
}
```

---

#### 7. Discord Unlink

* **Route:** `DELETE /api/private/auth/discord/unlink`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "status": "success",
    "user": { /* User object */ }
  }
}
```

---

#### 8. Logout

* **Route:** `DELETE /api/private/auth/logout`
* **Body:** —
* **Response** (`200`): status boolean

---

### User Routes

#### 9. Get Own Account

* **Route:** `GET /api/private/user/account`
* **Body:** —
* **Response** (`200`): full `User` object of the authenticated user.

---

#### 10. Get User by ID

* **Route:** `GET /api/private/user/:id`
* **Body:** —
* **Response** (`200`): `User` object.

> Only staff members with `ADMIN` flag can view other users' profiles.

---

#### 11. Change User Name

* **Route:** `PATCH /api/private/user/edit/name/:id`
* **Body:**

```json
{
  "new_name": "NewName"
}
```

* **Response** (`200`): updated `User` object.

> Staff members can change names of other users.

---

#### 12. Change User Email

* **Route:** `PATCH /api/private/user/edit/email/:id`
* **Body:**

```json
{
  "new_email": "new@example.com"
}
```

* **Response** (`200`): updated `User` object.

---

#### 13. Change User Password

* **Route:** `PATCH /api/private/user/edit/password/:id`
* **Body:**

```json
{
  "old_password": "currentPassword",
  "new_password": "newSecurePassword"
}
```

* **Response** (`200`): updated `User` object.

> `old_password` can be `null` if a staff member is resetting password for another user.

---

#### 14. Get User Sessions

* **Route:** `GET /api/private/user/sessions/:id`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "sessions": [
      {
        "id": "session-uuid",
        "user_agent": "Mozilla/5.0...",
        "ip": "192.168.1.1",
        "is_revoked": false,
        "expires_at": "2026-02-08T10:00:00Z",
        "created_at": "2026-02-01T10:00:00Z"
      }
    ]
  }
}
```

> Only staff with `DEV` (developer rank) or `SENIOR` (staff rank) flags can view sessions of other users.

---

#### 15. Terminate Session

* **Route:** `DELETE /api/private/user/session/terminate`
* **Body:**

```json
{
  "session_id": "session-uuid-to-terminate"
}
```

* **Response** (`200`):

```json
{
  "response": {
    "state": true
  }
}
```

---

#### 16. Terminate All Sessions

* **Route:** `DELETE /api/private/user/session/terminate/all`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "state": true
  }
}
```

---

### Admin — User Management

> Base path: `/api/private/admin/user`
> All endpoints require authentication and specific rank flags.

#### 17. Search All Users

* **Route:** `GET /api/private/admin/user/all`
* **Required Flag:** `ADMIN`
* **Body:** —
* **Response** (`200`): array of `User` objects.

---

#### 18. Get User Private Data

* **Route:** `GET /api/private/admin/user/:id`
* **Required Flag:** `ADMIN`
* **Body:** —
* **Response** (`200`): `User` object with sensitive data (`register_ip`, `last_ip`).

---

#### 19. Create User

* **Route:** `PUT /api/private/admin/user/create`
* **Required Flag:** `SENIOR`
* **Body:**

```json
{
  "name": "JaneDoe",
  "email": "jane@example.com",
  "password": "initialPassword"
}
```

* **Response** (`201`): created `User` object.

---

#### 20. Ban User

* **Route:** `PATCH /api/private/admin/user/ban/:id`
* **Required Flag:** `ADMIN`
* **Body:**

```json
{
  "unban_date": "2026-03-01T00:00:00Z",
  "reason": "Violation of rules"
}
```

* **Response** (`200`): `Ban` object.

---

#### 21. Unban User

* **Route:** `DELETE /api/private/admin/user/unban/:id`
* **Required Flag:** `ADMIN`
* **Body:** —
* **Response** (`200`): unban result.

---

#### 22. Delete User

* **Route:** `DELETE /api/private/admin/user/delete/:id`
* **Required Flag:** `SENIOR`
* **Body:** —
* **Response** (`200`): deletion result.

> First call — soft delete (sets status to `Deleted`). Second call — hard delete.

---

#### 23. Restore User

* **Route:** `PUT /api/private/admin/user/restore/:id`
* **Required Flag:** `MANAGER`
* **Body:** —
* **Response** (`200`): restore result.

> Only works on soft-deleted accounts.

---

#### 24. Set Staff Rank

* **Route:** `PATCH /api/private/admin/user/rank/staff/:id`
* **Required Flag:** `STAFFMANAGEMENT`
* **Body:**

```json
{
  "status": 2
}
```

* **Response** (`200`): updated `User` object.

---

#### 25. Set Developer Rank

* **Route:** `PATCH /api/private/admin/user/rank/developer/:id`
* **Required Flag:** `MANAGER`
* **Body:**

```json
{
  "status": 2
}
```

* **Response** (`200`): updated `User` object.

---

#### 26. Change User Status

* **Route:** `PATCH /api/private/admin/user/edit/status/:id`
* **Required Flag:** `LEAD`
* **Body:**

```json
{
  "new_status": 1
}
```

* **Response** (`200`): updated `User` object.

> Status enum: `0` = NotVerified, `1` = Active, `2` = Disabled, `3` = Deleted, `4` = RequiresPasswordChange.

---

#### 27. Set Donate Points

* **Route:** `PATCH /api/private/admin/user/edit/donate/:id`
* **Required Flag:** `MANAGER`
* **Body:**

```json
{
  "points": 100
}
```

* **Response** (`200`):

```json
{
  "response": {
    "user": { /* User object */ },
    "newPoints": 100
  }
}
```

---

#### 28. Edit User Flags

* **Route:** `PATCH /api/private/admin/user/editflag/:id`
* **Required Flag:** `STAFFMANAGEMENT`
* **Body:**

```json
{
  "flags": ["DEV", "ADMIN"]
}
```

* **Response** (`200`): updated `User` object.

---

#### 29. Edit User Badges

* **Route:** `PATCH /api/private/admin/user/editbadges/:id`
* **Required Flag:** `SENIOR`
* **Body:**

```json
{
  "badges": [1, 2, 3]
}
```

* **Response** (`200`): updated `User` object.

---

#### 30. Reset User Sensitive Data

* **Route:** `DELETE /api/private/admin/user/reset/:id`
* **Required Flag:** `SENIOR`
* **Body:** —
* **Response** (`200`): updated `User` object (IPs and last seen cleared).

---

#### 31. Force Unlink Discord

* **Route:** `DELETE /api/private/admin/user/discord/forceunlink/:id`
* **Required Flag:** `SENIOR`
* **Body:** —
* **Response** (`200`): updated `User` object.

---

### Support — User

#### 32. Populate Ban List

* **Route:** `GET /api/private/support/user/bans/populate`
* **Required Flag:** `STAFF`
* **Body:** —
* **Response** (`200`): array of `Ban` objects.

---

### Badges Routes

> Base path: `/api/private/admin/badges`

#### 33. Populate All Badges

* **Route:** `GET /api/private/admin/badges/populate`
* **Required Flag:** `LEAD`
* **Body:** —
* **Response** (`200`): array of `BadgeAdminResponse` objects.

```json
[
  {
    "id": 1,
    "name": "Early Tester",
    "description": "Awarded to early testers",
    "conditions": "manual",
    "color": "#ff5500",
    "icon_name": "star",
    "created_at": "2026-02-01T10:00:00Z",
    "updated_at": "2026-02-01T10:00:00Z",
    "created_by": { "id": 1, "name": "Admin" }
  }
]
```

---

#### 34. Create Badge

* **Route:** `POST /api/private/admin/badges/create`
* **Required Flag:** `LEAD`
* **Body:**

```json
{
  "name": "Early Tester",
  "description": "Awarded to early testers",
  "conditions": "manual",
  "color": "#ff5500",
  "icon_name": "star"
}
```

* **Response** (`201`): created `Badge` object.

---

#### 35. Edit Badge

* **Route:** `PATCH /api/private/admin/badges/edit/:id`
* **Required Flag:** `LEAD`
* **Body:**

```json
{
  "name": "Updated Name",
  "description": "Updated description",
  "conditions": "automatic",
  "color": "#00ff00",
  "icon_name": "check"
}
```

* **Response** (`200`): updated `Badge` object.

---

#### 36. Delete Badge

* **Route:** `DELETE /api/private/admin/badges/delete/:id`
* **Required Flag:** `LEAD`
* **Body:** —
* **Response** (`200`): boolean.

---

### Changelog Routes

> Base path: `/api/private/changelog`

#### 37. Populate Changelogs

* **Route:** `GET /api/private/changelog/populate`
* **Body:** —
* **Response** (`200`): array of `Changelog` objects.

```json
[
  {
    "id": 1,
    "created_at": "2026-02-01T10:00:00Z",
    "creator": { "id": 1, "name": "Admin" },
    "title": "v1.0.0 Release",
    "content": [
      { "type": 0, "content": "Added new feature X" },
      { "type": 1, "content": "Fixed bug Y" }
    ],
    "version": "1.0.0",
    "is_staff": false
  }
]
```

> Staff changelogs are appended if user has `STAFF` flag on their account or rank.

---

#### 38. Create Changelog

* **Route:** `POST /api/private/changelog/create`
* **Required Flag:** `DEV`
* **Body:**

```json
{
  "title": "v1.1.0 Release",
  "content": [
    { "type": 0, "content": "New feature" }
  ],
  "version": "1.1.0",
  "is_staff": false
}
```

* **Response** (`201`): created `Changelog` object.

---

#### 39. Delete Changelog

* **Route:** `DELETE /api/private/changelog/delete/:id`
* **Required Flag:** `DEV`
* **Body:** —
* **Response** (`200`): deletion result.

---

### Developer (Private)

> Base path: `/api/private/developer`

#### 40. Ping

* **Route:** `GET /api/private/developer/ping`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "status": "pong"
  }
}
```

---

#### 41. Populate Server Info

* **Route:** `GET /api/private/developer/server/populate`
* **Required Flag:** `DEV`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "timestamp": "2026-03-14T12:00:00Z",
    "system": {
      "cpu": 25.5,
      "ram": 60.2,
      "ram_gb": 3.8,
      "uptime": "72h30m"
    }
  }
}
```

---

#### 42. Populate Debug Stacktrace

* **Route:** `GET /api/private/developer/stacktrace`
* **Required Flag:** `DEV`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "stack": "goroutine 1 [running]:..."
  }
}
```

---

#### 43. Populate Services

* **Route:** `GET /api/private/developer/service/populate`
* **Required Flag:** `DEV`
* **Body:** —
* **Response** (`200`): services list.

---

#### 44. Service Interaction

* **Route:** `PATCH /api/private/developer/service/interaction`
* **Required Flag:** `SENIORDEV`
* **Body:**

```json
{
  "name": "invites",
  "action": "enable"
}
```

> `action` values: `enable`, `disable`, `status`

* **Response** (`200`):

```json
{
  "response": {
    "service": "invites",
    "state": true
  }
}
```

---

### Invite Codes Routes

> Base path: `/api/private/admin/invite`

#### 45. Get Invite Codes List

* **Route:** `GET /api/private/admin/invite/list`
* **Required Flag:** `STAFF`
* **Body:** —
* **Response** (`200`): array of `Invite` objects.

```json
[
  {
    "id": 1,
    "code": "ABC123XYZ",
    "created_by": 1,
    "used": false,
    "used_by": null,
    "created_at": "2026-02-01T10:00:00Z",
    "creator": { /* User object */ },
    "user": null
  }
]
```

---

#### 46. Create Invite Code

* **Route:** `POST /api/private/admin/invite/create`
* **Required Flag:** `STAFF`
* **Body:** —
* **Response** (`201`): created `Invite` object.

---

#### 47. Delete Invite Code

* **Route:** `DELETE /api/private/admin/invite/delete/:id`
* **Required Flag:** `LEAD`
* **Body:** —
* **Response** (`200`):

```json
{
  "response": {
    "status": "success"
  }
}
```

---

### Logs Routes

> Base path: `/api/private/admin/logs`

#### 48. Search Admin Logs

* **Route:** `POST /api/private/admin/logs/populate`
* **Required Flag:** `ADMIN`
* **Body:**

```json
{
  "type": "staffCommon",
  "date_start": "2026-01-01",
  "date_end": "2026-03-01",
  "keywords": "banned"
}
```

> `type` values: `staffCommon`, `staffPunish`, `adminTickets`, `adminAuth`.
> `keywords` is optional.

* **Response** (`200`):

```json
{
  "response": {
    "data": [ /* log entries */ ],
    "limit": 100
  }
}
```

---

#### 49. Search Staff Logs

* **Route:** `POST /api/private/admin/logs/staff/populate`
* **Required Flag:** `STAFFMANAGEMENT`
* **Body:** same as [Search Admin Logs](#48-search-admin-logs).
* **Response** (`200`): same structure.

---

### Notifications Routes

> Base path: `/api/private/notifications`

#### 50. Populate Own Notifications

* **Route:** `GET /api/private/notifications/populate`
* **Body:** —
* **Response** (`200`): array of `Notification` objects.

```json
[
  {
    "id": 1,
    "created_at": "2026-02-01T10:00:00Z",
    "type": "information",
    "title": "Welcome",
    "sender_id": null,
    "user_id": 1,
    "text": "Welcome to the platform!",
    "is_readed": false,
    "sender": null,
    "user": { "id": 1, "name": "JohnDoe" }
  }
]
```

> `type` values: `information`, `error`, `success`.

---

#### 51. Read Notifications

* **Route:** `POST /api/private/notifications/read`
* **Body:** —
* **Response** (`200`): updated notifications.

---

#### 52. Clear All Notifications

* **Route:** `DELETE /api/private/notifications/clear`
* **Body:** —
* **Response** (`200`): result.

---

#### 53. Remove Own Notification

* **Route:** `DELETE /api/private/notifications/remove`
* **Body:**

```json
{
  "id": 5
}
```

* **Response** (`200`): boolean.

---

#### 54. Send Notification (Admin)

* **Route:** `POST /api/private/admin/notifications/send/:id`
* **Required Flag:** `ADMIN`
* **Body:**

```json
{
  "type": "information",
  "title": "System Notice",
  "text": "Your account has been reviewed.",
  "sender_id": 1,
  "user_id": 2
}
```

> `:id` — target user ID.

* **Response** (`200`):

```json
{
  "response": {
    "status": "success"
  }
}
```

---

#### 55. Populate User Notifications (Admin)

* **Route:** `GET /api/private/admin/notifications/populate/:id`
* **Required Flag:** `ADMIN`
* **Body:** —
* **Response** (`200`): array of `Notification` objects for the specified user.

---

#### 56. Remove Notification (Admin)

* **Route:** `DELETE /api/private/admin/notifications/remove/:id`
* **Required Flag:** `SENIOR`
* **Body:**

```json
{
  "id": 5
}
```

> `:id` in URL — target user ID. `id` in body — notification ID.

* **Response** (`200`): boolean.

---

### Ranks Routes

> Base path: `/api/private/admin/ranks`

#### 57. Get Ranks List

* **Route:** `GET /api/private/admin/ranks/`
* **Body:** —
* **Response** (`200`): array of `Rank` objects.

```json
[
  {
    "id": 1,
    "name": "Tester",
    "color": "#660000",
    "flags": ["REGISTERAPPLICATIONS", "TICKETS"]
  }
]
```

---

#### 58. Create Rank

* **Route:** `POST /api/private/admin/ranks/create`
* **Required Flag:** `STAFFMANAGEMENT`
* **Body:**

```json
{
  "name": "Moderator",
  "color": "#00ff00",
  "flags": ["ADMIN", "TICKETS"]
}
```

> Flags are automatically uppercased.

* **Response** (`201`): created `Rank` object.

---

#### 59. Edit Rank

* **Route:** `PATCH /api/private/admin/ranks/edit/:id`
* **Required Flag:** `MANAGER`
* **Body:**

```json
{
  "name": "Updated Rank",
  "color": "#ff0000",
  "flags": ["ADMIN", "SENIOR"]
}
```

* **Response** (`200`): updated `Rank` object.

---

#### 60. Delete Rank

* **Route:** `DELETE /api/private/admin/ranks/delete/:id`
* **Required Flag:** `MANAGER`
* **Body:** —
* **Response** (`200`): boolean.

---

### Tickets Routes

> Base path: `/api/private/tickets`

#### 61. Create Ticket

* **Route:** `POST /api/private/tickets/create`
* **Body:**

```json
{
  "title": "Bug Report",
  "category": "bugs",
  "message": "I found a bug in..."
}
```

* **Response** (`201`): created `Ticket` object with messages.

---

#### 62. Populate Ticket

* **Route:** `GET /api/private/tickets/populate/:id`
* **Body:** —
* **Response** (`200`): `TicketResponse` object.

```json
{
  "response": {
    "ticket": {
      "id": 1,
      "created_at": "2026-02-01T10:00:00Z",
      "status": "Waiting for staff",
      "prioriy": "Medium",
      "author_id": 1,
      "title": "Bug Report",
      "category": "bugs",
      "updated_at": "2026-02-01T10:00:00Z",
      "author": { "id": 1, "name": "JohnDoe" },
      "handler": null
    },
    "messages": [
      {
        "id": 1,
        "ticket_id": 1,
        "created_at": "2026-02-01T10:00:00Z",
        "content": "I found a bug in...",
        "author": { "id": 1, "name": "JohnDoe" }
      }
    ],
    "actions": []
  }
}
```

> Only staff members can populate tickets created by other users.

---

#### 63. Get My Tickets

* **Route:** `GET /api/private/tickets/mytickets`
* **Body:** —
* **Response** (`200`): array of `Ticket` objects for the authenticated user.

---

#### 64. Send Ticket Message

* **Route:** `POST /api/private/tickets/send`
* **Body:**

```json
{
  "ticket": { "id": 1 },
  "author_id": 1,
  "content": "Here is more information..."
}
```

* **Response** (`201`): created `TicketMessage` object.

---

#### 65. Close Ticket

* **Route:** `PATCH /api/private/tickets/close/:id`
* **Body:** —
* **Response** (`200`): updated `Ticket` object.

---

#### 66. Search All Tickets (Admin)

* **Route:** `GET /api/private/admin/tickets/populate`
* **Required Flag:** `STAFF`
* **Body:** —
* **Response** (`200`): array of `Ticket` objects.

---

#### 67. Assign Ticket (Admin)

* **Route:** `POST /api/private/admin/tickets/assign/:id`
* **Required Flag:** `STAFF`
* **Body:** —
* **Response** (`200`): updated `Ticket` object with handler assigned.

---

#### 68. Change Ticket Category (Admin)

* **Route:** `PATCH /api/private/admin/tickets/edit/category/:id`
* **Required Flag:** `STAFF`
* **Body:**

```json
{
  "new_category": "feature-request"
}
```

* **Response** (`200`): updated `Ticket` object.

---

## Models Reference

### User

```json
{
  "id": 1,
  "name": "JohnDoe",
  "created_at": "2026-02-01T10:00:00Z",
  "updated_at": "2026-02-01T10:00:00Z",
  "last_login": "2026-03-14T12:00:00Z",
  "badges": [],
  "status": 1,
  "donate_points": 0,
  "staff_rank": 1,
  "developer_rank": 1,
  "staff_flags": null,
  "active_ban": null,
  "discord_name": null,
  "discord_id": null,
  "email": "john@example.com",
  "invited_by": { "id": 2, "name": "Admin" }
}
```

### Ban

```json
{
  "id": 1,
  "date": "2026-02-01T10:00:00Z",
  "expiration_date": "2026-03-01T00:00:00Z",
  "reason": "Violation of rules",
  "status": 0,
  "admin": { "id": 1, "name": "Admin" },
  "user": { "id": 2, "name": "JohnDoe" }
}
```

> Ban status: `0` = Active, `1` = Removed, `2` = Expired.

### Rank

```json
{
  "id": 1,
  "name": "Tester",
  "color": "#660000",
  "flags": ["REGISTERAPPLICATIONS", "TICKETS"]
}
```

### Notification

```json
{
  "id": 1,
  "created_at": "2026-02-01T10:00:00Z",
  "type": "information",
  "title": "Welcome",
  "sender_id": null,
  "user_id": 1,
  "text": "Welcome!",
  "is_readed": false,
  "sender": null,
  "user": { "id": 1, "name": "JohnDoe" }
}
```

### Badge

```json
{
  "id": 1,
  "name": "Early Tester",
  "description": "Awarded to early testers",
  "conditions": "manual",
  "color": "#ff5500",
  "icon_name": "star"
}
```

### Invite

```json
{
  "id": 1,
  "code": "ABC123XYZ",
  "created_by": 1,
  "used": false,
  "used_by": null,
  "created_at": "2026-02-01T10:00:00Z",
  "creator": { /* User object */ },
  "user": null
}
```

### Session

```json
{
  "id": "session-uuid",
  "user_agent": "Mozilla/5.0...",
  "ip": "192.168.1.1",
  "is_revoked": false,
  "expires_at": "2026-02-08T10:00:00Z",
  "created_at": "2026-02-01T10:00:00Z"
}
```
