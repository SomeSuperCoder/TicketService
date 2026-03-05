# Admin Users API Documentation

## Overview
Admin endpoints for managing CRM users (employees and ROI).

Base path: `/api/v1/admin/users`

## Endpoints

### 1. List Users
**GET** `/api/v1/admin/users`

List all CRM users with optional filtering.

**Query Parameters:**
- `role` (optional): Filter by role - `admin`, `org`, or `executor`
- `status` (optional): Filter by status - `active` or `blocked`
- `email` (optional): Search by email (partial match, case-insensitive)
- `limit` (optional): Number of results per page (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```json
{
  "users": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "role": "executor",
      "status": "active",
      "department_id": 1,
      "department_name": "IT Department",
      "first_name": "John",
      "last_name": "Doe",
      "middle_name": "Smith",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 42
}
```

### 2. Create User
**POST** `/api/v1/admin/users`

Create a new CRM user (employee or ROI).

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "role": "executor",
  "department_id": 1,
  "first_name": "Jane",
  "last_name": "Smith",
  "middle_name": "Marie"
}
```

**Required Fields:**
- `email`: Valid email address (max 255 chars)
- `role`: One of `admin`, `org`, or `executor`

**Optional Fields:**
- `department_id`: Department ID (must exist in departments table)
- `first_name`: First name (max 100 chars)
- `last_name`: Last name (max 100 chars)
- `middle_name`: Middle name (max 100 chars)

**Response:**
```json
{
  "id": "uuid",
  "email": "newuser@example.com",
  "role": "executor",
  "status": "active",
  "department_id": 1,
  "first_name": "Jane",
  "last_name": "Smith",
  "middle_name": "Marie",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 3. Update User
**PATCH** `/api/v1/admin/users/{id}`

Update user permissions or block/unblock user.

**Path Parameters:**
- `id`: User UUID

**Request Body (all fields optional):**
```json
{
  "role": "admin",
  "status": "blocked",
  "department_id": 2,
  "first_name": "Jane",
  "last_name": "Doe",
  "middle_name": "Marie"
}
```

**Fields:**
- `role`: Change user role (`admin`, `org`, or `executor`)
- `status`: Block or activate user (`active` or `blocked`)
- `department_id`: Change department
- `first_name`, `last_name`, `middle_name`: Update name fields

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "role": "admin",
  "status": "blocked",
  "department_id": 2,
  "first_name": "Jane",
  "last_name": "Doe",
  "middle_name": "Marie",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 4. Delete User
**DELETE** `/api/v1/admin/users/{id}`

Permanently delete user and revoke CRM access.

**Path Parameters:**
- `id`: User UUID

**Response:**
```json
{
  "id": "uuid",
  "email": "deleted@example.com",
  "role": "executor",
  "status": "active",
  "department_id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "Smith",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## User Roles

- **admin**: Full system access, can manage all users and settings
- **org**: Organization representative (ROI), can view and manage tickets
- **executor**: Department employee, can work on assigned tickets

## User Status

- **active**: User can access the CRM system
- **blocked**: User is blocked and cannot access the system

## Notes

- Email addresses must be unique
- Blocked users are filtered out from the `GetUser` query
- All timestamps are in UTC
- Department ID must reference an existing department
