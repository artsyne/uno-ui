# Server API Specification

**Version**: 1.0.0
**Status**: Active
**Last Updated**: 2025-12-18

## 1. Overview

This document serves as the Source of Truth for the iDiary.ai Hebe Server API. It describes the RESTful endpoints, request/response structures, and authentication requirements.

## 2. Global Standards

### 2.1 Response Format

All API responses follow a unified generic structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "request_id": "req_123456"
}
```

**Error Response**:

```json
{
  "success": false,
  "code": 4001001,
  "message": "Validation failed",
  "details": {
    "field": "error description"
  },
  "request_id": "req_123456"
}
```

### 2.2 Authentication

- **Method**: Cookie-based (`ud` cookie).
- **Header**: `X-CSRF-Token` required for state-changing methods (POST, PUT, DELETE).

## 3. API Endpoints

### 3.1 User Module (`/u`)

| Method | Endpoint | Description | Auth |
|:-------|:---------|:------------|:-----|
| POST   | `/register` | Register new user | No |
| POST   | `/login` | User login | No |
| POST   | `/check-username` | Check username availability | No |
| POST   | `/check-email` | Check email availability | No |
| POST   | `/u/settings/profile` | Update user profile | Yes |
| POST   | `/u/settings/password` | Change password | Yes |

### 3.2 Diary Module (`/diary`)

| Method | Endpoint | Description | Auth |
|:-------|:---------|:------------|:-----|
| POST   | `/diary` | Create new diary | Yes |
| PUT    | `/diary/:id` | Update diary | Yes |
| GET    | `/api/diary/like/status` | Get like status | No |
| POST   | `/api/diary/like` | Toggle like | Yes |

### 3.3 Post Module (`/post`)

| Method | Endpoint | Description | Auth |
|:-------|:---------|:------------|:-----|
| POST   | `/post` | Create new post | Yes |
| PUT    | `/post/:id` | Update post | Yes |
| DELETE | `/post/:id` | Delete post | Yes |
| POST   | `/api/post/stats` | Get post statistics | No |
| POST   | `/api/post/like` | Toggle like | Yes |

### 3.4 Comment Module (`/comment`)

| Method | Endpoint | Description | Auth |
|:-------|:---------|:------------|:-----|
| GET    | `/comment/list` | Get comments list | No |
| POST   | `/comment/create` | Create comment | Yes |
| POST   | `/comment/delete` | Delete comment | Yes |

## 4. Data Models

### 4.1 Diary

```go
type Diary struct {
    DiaryID       int64  `json:"diary_id"`
    Title         string `json:"title"`
    Remark        string `json:"remark"`
    ViewSetting   int8   `json:"view_setting"`   // 1:Public, 2:Private
    SocialSetting int8   `json:"social_setting"` // 0:Off, 1:On
}
```

### 4.2 Entry

```go
type Entry struct {
    EntryID int64  `json:"entry_id"`
    DiaryID int64  `json:"diary_id"`
    Title   string `json:"title"`
    Body    string `json:"body"`
    Summary string `json:"summary"`
}
```
