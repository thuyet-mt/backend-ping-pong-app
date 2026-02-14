# CDN + Nginx Avatar Serving Architecture

## Overview

This document describes the production architecture for serving user avatar images using:

- Nginx (static file server)
- Cloudflare CDN
- Go backend API
- PostgreSQL database

The goal is to ensure:

- Fast image delivery via CDN
- Clean separation of responsibilities
- Scalable and maintainable backend architecture

---

## High-Level Architecture

```
Flutter App
     ↓
Cloudflare CDN (cdn.example.com)
     ↓
Nginx Static Server (Docker)
     ↓
Local File Storage (/opt/static/files)

Go Backend API (api.example.com)
     ↓
PostgreSQL
```

Key idea:

- Backend API never serves image files
- Nginx serves static files
- Cloudflare caches images globally

---

## File Storage Layout

Static files are stored on the host machine:

```
/opt/static/
└── files/
    └── uploads/
        ├── avatar1.jpg
        ├── avatar2.png
        └── ...
```

Database stores only relative paths:

```
files/uploads/<filename>
```

Example:

```
files/uploads/f213f84309fa8d9f.jpeg
```

---

## Database Design

Table: `players`

Relevant column:

```
avatar_url TEXT
```

Important rule:

- Database stores relative file path
- Database does NOT store full CDN URL

This avoids vendor lock-in and keeps data portable.

---

## Repository Layer Responsibilities

Repository interacts only with PostgreSQL.

It must:

- Return raw data from database
- Not contain CDN or HTTP logic

Example model:

```go
type PlayerListItem struct {
    FullName   string
    BirthYear  *int
    AvatarPath string
}
```

Repository returns `AvatarPath`, not a full URL.

---

## Service Layer Responsibilities

Service layer converts internal data into API responses.

It builds the public CDN URL.

Environment variable:

```
CDN_BASE_URL=https://cdn.example.com
```

Helper function:

```go
func BuildCDNURL(path string) string {
    if path == "" {
        return ""
    }
    return strings.TrimRight(os.Getenv("CDN_BASE_URL"), "/") +
        "/" +
        strings.TrimLeft(path, "/")
}
```

Response DTO:

```go
type PlayerListResponse struct {
    FullName  string `json:"full_name"`
    BirthYear *int   `json:"birth_year"`
    AvatarURL string `json:"avatar_url"`
}
```

Mapping example:

```go
AvatarURL: BuildCDNURL(p.AvatarPath)
```

---

## Handler Layer Responsibilities

Handlers:

- Call service methods
- Return JSON responses
- Do not contain business logic

Example:

```go
c.JSON(http.StatusOK, players)
```

---

## Nginx Configuration

Nginx serves files from mounted volume:

```
/opt/static/files → /files
```

Key configuration:

- Long cache headers
- Immutable assets
- Read-only static serving

Example headers:

```
Cache-Control: public, max-age=2592000, immutable
```

---

## Cloudflare CDN Behavior

Cloudflare:

- Caches images at edge locations
- Reduces latency
- Offloads traffic from origin server

Expected response header after caching:

```
cf-cache-status: HIT
```

---

## API Response Format

Example API response:

```json
{
  "full_name": "Nguyen Van A",
  "avatar_url": "https://cdn.example.com/files/uploads/xxx.png"
}
```

Frontend loads images directly from CDN.

Backend is not involved in image delivery.

---

## Production Best Practices

1. Store relative paths in database
2. Build public URLs in service layer
3. Keep repository infrastructure-agnostic
4. Use environment variables for CDN base URL
5. Enable long-term caching for static assets
6. Separate API domain and CDN domain

---

## Summary

This architecture ensures:

- Clean separation of concerns
- High performance image delivery
- Easy CDN migration in the future
- Maintainable backend codebase

Backend focuses on business logic.
CDN and Nginx handle static delivery.

