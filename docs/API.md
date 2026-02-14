# API Reference

Base URL: `http://localhost:8000/api/v1`

Swagger UI (requires auth): `http://localhost:8000/swagger/index.html`

## Authentication

Protected endpoints require a valid JWT.
Only these endpoints are public:

- `POST /register`
- `POST /login`

Supported by backend:

- `Authorization: Bearer <token>`
- `Authorization: <token>`
- `task_tracker_token` cookie

In current frontend implementation, protected calls use the `Authorization` header.

## User Endpoints

### `POST /register`

Creates a user.

Request body:

```json
{
  "firstName": "Alice",
  "lastName": "Smith",
  "email": "alice@example.com",
  "password": "secret123"
}
```

Success:

- `201 Created`

Validation:

- `email` must be valid
- `password` length `3..130`
- duplicate email returns `400`

### `POST /login`

Authenticates user and returns JWT.

Request body:

```json
{
  "email": "alice@example.com",
  "password": "secret123"
}
```

Success response (`200 OK`):

```json
{
  "token": "<jwt>"
}
```

### `GET /profile` (protected)

Returns current user profile.

### `PUT /profile` (protected)

Updates first and last name.

Request body:

```json
{
  "firstName": "Alice",
  "lastName": "Smith"
}
```

### `PUT /profile/password` (protected)

Changes password.

Request body:

```json
{
  "currentPassword": "old-secret",
  "newPassword": "new-secret"
}
```

### `GET /users/lookup` (protected)

Returns user lookup list for assignment UI.

### `GET /users/tasks` (protected)

Returns all users with their current assigned tasks (`todo`, `in_progress`).

## Goals Endpoints

### `GET /goals` (protected)

Returns all goals with nested tasks for authenticated users.

Success response (`200 OK`):

```json
[
  {
    "id": 1,
    "title": "Launch MVP",
    "description": "Ship first release",
    "ownerId": 1,
    "ownerName": "Alice Smith",
    "createdAt": "2026-02-13T10:00:00Z",
    "tasks": []
  }
]
```

### `POST /goals` (protected)

Creates a goal owned by current user.

Request body:

```json
{
  "title": "Launch MVP",
  "description": "Ship first release"
}
```

Success: `201 Created`

Validation:

- `title` length `3..255`
- `description` is optional, max length `2000`

### `PUT /goals/{goalID}` (protected)

Updates goal title and description. Only goal owner can update.

Request body:

```json
{
  "title": "Updated goal title",
  "description": "Updated goal description"
}
```

Success: `200 OK`

### `DELETE /goals/{goalID}` (protected)

Deletes a goal owned by requester (and nested tasks via cascade).

Success: `204 No Content`

## Tasks Endpoints

### `POST /goals/{goalID}/tasks` (protected)

Creates a task under goal. Only goal owner can create.

Request body:

```json
{
  "title": "Create onboarding flow",
  "description": "Draft screens and copy",
  "assigneeId": 2
}
```

Notes:

- `assigneeId` is optional and may be `null`
- `description` is optional
- returns `403` when requester does not own the goal

### `GET /goals/{goalID}/tasks` (protected)

Returns one goal object with nested tasks. Any authenticated user can view.

### `GET /tasks/assigned` (protected)

Returns tasks assigned to current user.

### `PUT /tasks/{taskID}` (protected)

Updates task fields. Goal owner permission required.

Request body:

```json
{
  "goalId": 1,
  "title": "Updated task title",
  "description": "Updated task description",
  "status": "in_progress",
  "assigneeId": 2
}
```

Success: `200 OK`

### `PUT /tasks/{taskID}/assign` (protected)

Assigns or unassigns task. Only goal owner can assign.

Request body:

```json
{
  "assigneeId": 2
}
```

To unassign, send:

```json
{
  "assigneeId": null
}
```

### `DELETE /tasks/{taskID}` (protected)

Deletes task under a goal owned by requester.

Success: `204 No Content`

## Error Shape

Error responses are JSON and include an error message in `statusMessage` when proxied through Nuxt routes.

Typical status codes:

- `400`: validation or malformed payload
- `401`: missing authorization header in Nuxt proxy
- `403`: invalid token or permission denied
- `500`: unexpected server/database failure
